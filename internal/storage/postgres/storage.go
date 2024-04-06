package postgres

import (
	"context"
	"executor/internal/config"
	"executor/internal/models"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type RemoteStorage struct {
	pool *pgxpool.Pool
}

func (r RemoteStorage) GetCommands(ctx context.Context) ([]*models.Runnable, error) {
	// Get command bases
	commands, err := r.getCommandBases(ctx)
	if err != nil {
		return nil, err
	}

	for _, c := range commands {
		if err := r.enrichRunnable(ctx, c); err != nil {
			return nil, err
		}
	}

	return commands, nil
}

func (r RemoteStorage) enrichRunnable(ctx context.Context, rc *models.Runnable) error {
	if err := r.inplaceGetCommandInfo(ctx, rc); err != nil {
		return err
	}

	if err := r.inplaceGetCommandSources(ctx, rc); err != nil {
		return err
	}

	if err := r.inplaceGetCommandOutputLines(ctx, rc); err != nil {
		return err
	}

	return nil
}

func (r RemoteStorage) getCommandBases(ctx context.Context) ([]*models.Runnable, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Commit(ctx)

	commands := make([]*models.Runnable, 0)

	rows, err := tx.Query(ctx, queryGetCommandBases)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		runnable := new(models.Runnable)

		if err := rows.Scan(&runnable.Sid, &runnable.Status); err != nil {
			return nil, err
		}

		commands = append(commands, runnable)
	}

	return commands, nil
}

func (r RemoteStorage) inplaceGetCommandInfo(ctx context.Context, rn *models.Runnable) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Commit(ctx)

	return tx.QueryRow(ctx, queryGetCommandInfo, pgx.NamedArgs{
		"outSid": rn.Sid,
	}).Scan(&rn.Info.ScheduledTime, &rn.Info.StartedTime, &rn.Info.ExitTime, &rn.Info.ExitCode)
}

func (r RemoteStorage) getCommandLines(ctx context.Context, q string, sid uuid.UUID) ([]string, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Commit(ctx)

	lines := make([]string, 0)

	// Get lines
	rows, err := tx.Query(ctx, q, pgx.NamedArgs{
		"outSid": sid.String(),
	})
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ns string

		if err := rows.Scan(&ns); err != nil {
			return nil, err
		}
		lines = append(lines, ns)
	}

	return lines, nil
}

func (r RemoteStorage) inplaceGetCommandSources(ctx context.Context, rn *models.Runnable) error {
	lines, err := r.getCommandLines(ctx, queryGetCommandSourceLines, rn.Sid)
	if err != nil {
		return err
	}

	rn.Sources = lines
	return nil
}

func (r RemoteStorage) inplaceGetCommandOutputLines(ctx context.Context, rn *models.Runnable) error {
	lines, err := r.getCommandLines(ctx, queryGetCommandOutputLines, rn.Sid)
	if err != nil {
		return err
	}

	rn.Info.Output = lines
	return nil
}

func (r RemoteStorage) GetCommandByID(ctx context.Context, sid uuid.UUID) (*models.Runnable, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Commit(ctx)

	// Fill base
	cmd := &models.Runnable{Sid: sid}

	if err := tx.QueryRow(ctx, queryGetCommandBaseBySid, pgx.NamedArgs{"outSid": sid}).Scan(&cmd.Status); err != nil {
		return nil, err
	}

	if err := r.enrichRunnable(ctx, cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

func (r RemoteStorage) AddCommand(ctx context.Context, sources []string) (*models.Runnable, error) {
	newRunnable := models.NewRunnable(sources)

	// Add base
	if err := r.addCommandBase(ctx, newRunnable.Sid, newRunnable.Status); err != nil {
		return nil, err
	}

	// Add sources
	if err := r.addCommandSources(ctx, newRunnable.Sid, newRunnable.Sources); err != nil {
		return nil, err
	}

	// Add info
	if err := r.addCommandInfo(ctx, newRunnable); err != nil {
		return nil, err
	}

	return newRunnable, nil
}

func (r RemoteStorage) addCommandBase(ctx context.Context, sid uuid.UUID, status models.RunnableStatus) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Commit(ctx)

	_, err = tx.Exec(ctx, queryInsertCommandBase, pgx.NamedArgs{
		"outSid":    sid.String(),
		"outStatus": status,
	})

	return err
}

func (r RemoteStorage) addCommandLines(ctx context.Context, q string, sid uuid.UUID, lines []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Commit(ctx)

	for _, l := range lines {
		if _, err = tx.Exec(ctx, q, pgx.NamedArgs{
			"outSid":  sid.String(),
			"outLine": l,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (r RemoteStorage) addCommandInfo(ctx context.Context, rn *models.Runnable) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Commit(ctx)

	_, err = tx.Exec(ctx, queryInsertCommandInfo, pgx.NamedArgs{
		"outSid":    rn.Sid.String(),
		"schedTime": rn.Info.ScheduledTime,
		"startTime": rn.Info.StartedTime,
		"exitTime":  rn.Info.ExitTime,
		"exitCode":  rn.Info.ExitCode,
	})

	return err
}

func (r RemoteStorage) addCommandSources(ctx context.Context, sid uuid.UUID, sources []string) error {
	return r.addCommandLines(ctx, queryInsertCommandSourceLine, sid, sources)
}

func (r RemoteStorage) addCommandOutput(ctx context.Context, sid uuid.UUID, output []string) error {
	return r.addCommandLines(ctx, queryInsertCommandOutputLine, sid, output)
}

func (r RemoteStorage) UpdateCommandInfo(ctx context.Context, runnable *models.Runnable) error {
	if err := r.updateCommandStatus(ctx, runnable); err != nil {
		return err
	}

	return r.updateCommandInfo(ctx, runnable)
}

func (r RemoteStorage) updateCommandStatus(ctx context.Context, runnable *models.Runnable) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Commit(ctx)

	_, err = tx.Exec(ctx, queryUpdateCommandStatus, pgx.NamedArgs{
		"outStatus": runnable.Status,
		"outSid":    runnable.Sid,
	})

	return err
}

func (r RemoteStorage) updateCommandInfo(ctx context.Context, runnable *models.Runnable) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Commit(ctx)

	_, err = tx.Exec(ctx, queryUpdateCommandInfo, pgx.NamedArgs{
		"outSid":    runnable.Sid,
		"schedTime": runnable.Info.ScheduledTime,
		"startTime": runnable.Info.StartedTime,
		"exitTime":  runnable.Info.ExitTime,
		"exitCode":  runnable.Info.ExitCode,
	})

	return err
}

func (r RemoteStorage) AddCommandOutput(ctx context.Context, sid uuid.UUID, outputLines []string) error {
	return r.addCommandOutput(ctx, sid, outputLines)
}

func (r RemoteStorage) Close(_ context.Context) {
	slog.Info("closing the RemoteStorage(postgres)")

	r.pool.Close()
}

func compileConString(cfg *config.Configuration) string {
	// fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	dbc := cfg.Database

	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		dbc.Username, dbc.Password,
		dbc.Hostname, dbc.Port, dbc.DB,
	)
}

func GetStorage(ctx context.Context, cfg *config.Configuration) (*RemoteStorage, error) {
	var storage RemoteStorage
	cs := compileConString(cfg)

	// Connect to database using the pgxpool: we need a safe concurent access to database from many routines
	// Here I use wrapper doWithTrials, because postgres usually loads longer than Executor service
	err := doWithTrials(func() error {
		pgxC, err := pgxpool.New(ctx, cs)
		if err != nil {
			return err
		}

		storage.pool = pgxC
		return nil
	}, 5)

	if err != nil {
		return nil, err
	}

	// Try to ping: maybe connection was aborted
	if err = storage.pool.Ping(ctx); err != nil {
		// We can't use this postgres instance
		return nil, err
	}

	return &storage, nil
}
