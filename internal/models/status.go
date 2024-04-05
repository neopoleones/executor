package models

type RunnableStatus string

const (
	StatusScheduled  RunnableStatus = "scheduled"
	StatusInProgress                = "running"
	StatusDone                      = "done"

	// StatusRejected may be used after failing input validation?
	// (should I introduce input validation? If we can schedule arbitrary commands there is no need to check for command injections)
	StatusRejected = "rejected"
)
