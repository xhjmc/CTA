namespace go rmservice

struct ResourceRequest {
    1: i32 BranchType (go.tag = "json:\"branch_type\"")
    2: string Xid (go.tag = "json:\"xid\"")
    3: i64 BranchId (go.tag = "json:\"branch_id\"")
    4: string ResourceId (go.tag = "json:\"reosurce_id\"")
}

struct ResourceResponse {
    1: i32 BranchStatus (go.tag = "json:\"branch_status\"")
    2: string Error (go.tag = "json:\"error\"")
}

service ResourceManagerBaseService {
    string Ping(1: string req)
    ResourceResponse BranchCommit(1: ResourceRequest req)
    ResourceResponse BranchRollback(1: ResourceRequest req)
}

// thrift --gen go idls/rmservice.thrift