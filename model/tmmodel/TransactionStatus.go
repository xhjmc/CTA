package tmmodel

type TransactionStatus int32

const (
	Unknown TransactionStatus = iota
	Begin
	Committing
	CommitDone
	CommitFailed
	Rollbacking
	RollbackDone
	RollbackFailed
)

func (s TransactionStatus) String() string {
	switch s {
	case Unknown:
		return "Unknown"
	case Begin:
		return "Begin"
	case Committing:
		return "Committing"
	case CommitDone:
		return "CommitDone"
	case CommitFailed:
		return "CommitFailed"
	case Rollbacking:
		return "Rollbacking"
	case RollbackDone:
		return "RollbackDone"
	case RollbackFailed:
		return "RollbackFailed"
	default:
		return ""
	}
}
