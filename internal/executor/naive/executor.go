package naive

import (
	"context"
	"errors"
	"executor/internal/config"
	"executor/internal/executor"
	"executor/internal/models"
	"executor/internal/storage"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"time"
)

const (
	tempLoc = "/tmp"
	tempSfx = "naive_runner"
)

type SystemExecutor struct {
	interpreterPath string

	storage storage.ExecutorStorage
}

func (s *SystemExecutor) getTemp() (*os.File, error) {
	path, err := os.MkdirTemp(tempLoc, tempSfx)
	if err != nil {
		return nil, err
	}

	return os.CreateTemp(path, tempSfx)
}

func (s *SystemExecutor) toTempScript(commands []string) (string, error) {
	bFile, err := s.getTemp()
	if err != nil {
		return "", err
	}

	for _, c := range commands {
		if _, err = bFile.WriteString(c); err != nil {
			return "", err
		}
	}

	return bFile.Name(), nil
}

func (s *SystemExecutor) prepareCommand(fName string, buffer io.Writer) *exec.Cmd {
	return &exec.Cmd{
		Path:      s.interpreterPath,
		Args:      []string{s.interpreterPath, fName}, // argv[0] = runner_path
		WaitDelay: time.Second,
		Stdout:    buffer,
	}
}

func (s *SystemExecutor) Release(ctx context.Context) {
	// Release is a function for other CommandExecutor implementations
	// like remote environments support

	slog.Info("executor released")
}

func (s *SystemExecutor) Run(ctx context.Context, sid uuid.UUID) (*models.Runnable, error) {
	// First of all, we get a runnable from storage
	runnable, err := s.storage.GetCommandByID(ctx, sid)
	if err != nil {
		return nil, err
	}

	// It should be in scheduled state before starting
	if runnable.Status != models.StatusScheduled {
		return nil, executor.ErrNotScheduled
	}

	// Here we copy content of commands to temp file for running using the interpreter
	fName, err := s.toTempScript(runnable.Sources)
	if err != nil {
		return nil, err
	}
	defer os.Remove(fName)

	// Prepare buffers
	// Get actual command for execution
	outBuf := GetTSafeBuf()
	cmd := s.prepareCommand(fName, outBuf)

	runnable.UpdateStatus(models.StatusInProgress)
	if err := s.storage.UpdateCommandInfo(ctx, runnable); err != nil {
		return nil, err
	}

	// Start the command. Start the routine to track output changes.
	// Also check for Rejected status and kill command on need

	go func() {
		rejectTicker := time.NewTicker(time.Second)

	outer:
		for {
			if nl := outBuf.buf.String(); len(nl) > 0 {
				// TODO: add buffered read for adding several lines at once
				if err := s.storage.AddCommandOutput(ctx, runnable.Sid, []string{nl}); err != nil {
					slog.Warn(
						"failed to add command output",
						slog.String("sid", runnable.Sid.String()),
						slog.String("err", err.Error()),
					)
				}

				// Should be called in order to clean read output
				outBuf.Reset()
			}

			// Check for status update once a second
			select {
			case <-rejectTicker.C:
				runnableUpdate, err := s.storage.GetCommandByID(ctx, runnable.Sid)
				if err != nil {
					slog.Error(
						"failed to get runnable status",
						slog.String("sid", runnable.Sid.String()),
						slog.String("err", err.Error()),
					)
				}

				// Force process to exit
				if runnableUpdate.Status == models.StatusRejected {
					_ = cmd.Process.Kill()
				}

			case <-ctx.Done():
				// TODO: check if context check is useful here

				// When context is fired, server is shutting down
				// So we can't handle execution of runnable anymore

				_ = cmd.Process.Kill()
				break outer
			default:
				// Just check for next line
			}
		}

	}()

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	if err = cmd.Wait(); err != nil {
		runnable.UpdateStatus(models.StatusRejected)
	} else {
		runnable.UpdateStatus(models.StatusDone)
	}

	runnable.AddExitCode(cmd.ProcessState.ExitCode())
	if err := s.storage.UpdateCommandInfo(ctx, runnable); err != nil {
		return nil, err
	}

	return runnable, nil
}

func GetExecutor(es storage.ExecutorStorage, cfg *config.Configuration) *SystemExecutor {
	// Check for interpreter
	nip := cfg.Executor.InterpreterPath

	if _, err := os.Stat(nip); errors.Is(err, os.ErrNotExist) {
		// tough way to say something is wrong
		// but service shouldn't be started with incorrect interpreter
		panic(fmt.Sprintf("interpreter: %s - not found", nip))
	}

	return &SystemExecutor{nip, es}
}
