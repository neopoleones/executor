package models

type RunnableStatus string

const (
	StatusScheduled RunnableStatus = "scheduled"
	StatusDone                     = "done"
	StatusRejected                 = "rejected"
)
