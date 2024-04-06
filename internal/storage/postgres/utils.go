package postgres

import "time"

func doWithTrials(cmd func() error, trials int) error {
	var err error
	ticker := time.Tick(time.Second)

	for i := 0; i < trials; i += 1 {
		if err = cmd(); err == nil {
			return nil
		}

		<-ticker
	}

	// All attempts failed
	return err
}
