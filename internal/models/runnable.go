package models

import (
	"github.com/google/uuid"
	"time"
)

type Runnable struct {
	Sid uuid.UUID `json:"sid"`

	Source string         `json:"source"`
	Status RunnableStatus `json:"status"`

	Info RunnableInfo `json:"info"`
}

type RunnableInfo struct {
	ScheduledTime time.Time `json:"sched_time"`
	StartedTime   time.Time `json:"start_time"`
	ExitTime      time.Time `json:"exit_time"`

	ExitCode int      `json:"exit_code"`
	Output   []string `json:"output"`
}
