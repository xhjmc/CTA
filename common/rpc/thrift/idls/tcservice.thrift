namespace go tcservice

struct BranchRegisterRequest {
    1: i32 BranchType (go.tag = "json:\"branch_type\"")
    2: string Xid (go.tag = "json:\"xid\"")
    3: string ResourceId (go.tag = "json:\"reosurce_id\"")
}

struct BranchRegisterResponse {
    1: i64 BranchId (go.tag = "json:\"branch_id\"")
    2: string Error (go.tag = "json:\"error\"")
}

struct BranchReportRequest {
    1: i32 BranchType (go.tag = "json:\"branch_type\"")
    2: string Xid (go.tag = "json:\"xid\"")
    3: i64 BranchId (go.tag = "json:\"branch_id\"")
    4: i32 BranchStatus (go.tag = "json:\"branch_status\"")
}

struct BranchReportResponse {
    1: string Error (go.tag = "json:\"error\"")
}

struct GlobalLockRequest {
    1: i32 BranchType (go.tag = "json:\"branch_type\"")
    2: string Xid (go.tag = "json:\"xid\"")
    3: string ResourceId (go.tag = "json:\"reosurce_id\"")
    4: string LockKeys (go.tag = "json:\"lock_keys\"")
}

struct GlobalLockResponse {
    1: string Error (go.tag = "json:\"error\"")
}

struct TransactionBeginRequest {
}

struct TransactionBeginResponse {
    1: string Xid (go.tag = "json:\"xid\"")
    2: string Error (go.tag = "json:\"error\"")
}

struct TransactionRequest {
    1: string Xid (go.tag = "json:\"xid\"")
}

struct TransactionResponse {
    1: i32 TransactionStatus (go.tag = "json:\"transaction_status\"")
    2: string Error (go.tag = "json:\"error\"")
}

service TransactionCoordinatorService {
    string Ping(1: string req)
    BranchRegisterResponse BranchRegister(1: BranchRegisterRequest req)
    BranchReportResponse BranchReport(1: BranchReportRequest req)
    GlobalLockResponse GlobalLock(1: GlobalLockRequest req)
    TransactionBeginResponse TransactionBegin(1: TransactionBeginRequest req)
    TransactionResponse TransactionCommit(1: TransactionRequest req)
    TransactionResponse TransactionRollback(1: TransactionRequest req)
    TransactionResponse GetTransactionStatus(1: TransactionRequest req)
}

// thrift --gen go idls/tcservice.thrift