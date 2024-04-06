package models

import (
	"github.com/google/uuid"
	"time"
)

type Runnable struct {
	Sid uuid.UUID `json:"sid"`

	Sources []string       `json:"sources"`
	Status  RunnableStatus `json:"status"`

	Info RunnableInfo `json:"info"`
}

func (r *Runnable) UpdateStatus(status RunnableStatus) {
	r.Status = status

	switch status {
	case StatusInProgress:
		r.Info.StartedTime = time.Now()
	case StatusDone, StatusRejected:
		r.Info.ExitTime = time.Now()
	}
}

func (r *Runnable) AddExitCode(code int) {
	r.Info.ExitCode = code
}

type RunnableInfo struct {
	ScheduledTime time.Time `json:"sched_time,omitempty"`
	StartedTime   time.Time `json:"start_time,omitempty"`
	ExitTime      time.Time `json:"exit_time,omitempty"`

	ExitCode int      `json:"exit_code"`
	Output   []string `json:"output,omitempty"`
}

func NewRunnable(sources []string) *Runnable {
	// Internally calls getTime which always returns nil for error. Handling can be ignored
	sid, _ := uuid.NewUUID()
	return &Runnable{
		Sid:     sid,
		Sources: sources,
		Status:  StatusScheduled,
		Info: RunnableInfo{
			ScheduledTime: time.Now(),
			Output:        make([]string, 0),
		},
	}
}
