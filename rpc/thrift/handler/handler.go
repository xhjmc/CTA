package handler

import (
	"context"
	"cta/logs"
	"cta/rpc/thrift/gen-go/tc"
)

type TCServiceHandler struct {
}

func (h *TCServiceHandler) Ping(ctx context.Context, req *tc.PingRequest) (*tc.PingResponse, error) {
	resp := tc.NewPingResponse()
	resp.Msg = req.Msg + " pong"
	logs.Infof(resp.Msg)
	return resp, nil
}
