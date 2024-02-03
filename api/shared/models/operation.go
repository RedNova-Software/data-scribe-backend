package models

// Used to store ongoing and complete operations
type Operation struct {
	OperationID string
	Completed   bool
	DeleteAt    int64
}
