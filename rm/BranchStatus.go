package rm

type BranchStatus int

const (
	Unknown BranchStatus = iota
	Registered
	PhaseOne_Done
	PhaseOne_Failed
	PhaseOne_Timeout
	PhaseTwo_CommittDone
	PhaseTwo_CommitFailed_Retryable
	PhaseTwo_CommitFailed_Unretryable
	PhaseTwo_RollbackDone
	PhaseTwo_RollbackFailed_Retryable
	PhaseTwo_RollbackFailed_Unretryable
)
