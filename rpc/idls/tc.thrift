namespace go tc

struct PingRequest {
    1: string Msg
}

struct PingResponse {
    1: string Msg
}

service TCService {
   PingResponse Ping(1: PingRequest req)
}