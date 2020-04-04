package rmmodel

type BranchStatus int32

const (
	Unknown BranchStatus = iota // unused
	Registered
	PhaseOne_Done
	PhaseOne_Failed
	PhaseOne_Timeout // unused
	PhaseTwo_CommitDone
	PhaseTwo_CommitFailed_Retryable   // unused
	PhaseTwo_CommitFailed_Unretryable // unused
	PhaseTwo_RollbackDone
	PhaseTwo_RollbackFailed_Retryable
	PhaseTwo_RollbackFailed_Unretryable // unused
)

func (s BranchStatus) String() string {
	switch s {
	case Unknown:
		return "Unknown"
	case Registered:
		return "Registered"
	case PhaseOne_Done:
		return "PhaseOne_Done"
	case PhaseOne_Failed:
		return "PhaseOne_Failed"
	case PhaseOne_Timeout:
		return "PhaseOne_Timeout"
	case PhaseTwo_CommitDone:
		return "PhaseTwo_CommitDone"
	case PhaseTwo_CommitFailed_Retryable:
		return "PhaseTwo_CommitFailed_Retryable"
	case PhaseTwo_CommitFailed_Unretryable:
		return "PhaseTwo_CommitFailed_Unretryable"
	case PhaseTwo_RollbackDone:
		return "PhaseTwo_RollbackDone"
	case PhaseTwo_RollbackFailed_Retryable:
		return "PhaseTwo_RollbackFailed_Retryable"
	case PhaseTwo_RollbackFailed_Unretryable:
		return "PhaseTwo_RollbackFailed_Unretryable"
	default:
		return ""
	}
}
