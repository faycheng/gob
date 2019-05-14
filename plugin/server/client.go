package server

import (
	"context"
	"github.com/faycheng/gob/plugin/server/api"
	"google.golang.org/grpc"
)

type PluginClient struct {
	path   string
	conn   *grpc.ClientConn
	client proto.TaskServiceClient
}

func NewPluginClient(path string) *PluginClient {
	// TODO: parse the unix socket path
	conn, err := grpc.Dial("/")
	client := proto.NewTaskServiceClient(conn)
	if err != nil {
		panic(err)
	}
	return &PluginClient{
		path:   path,
		conn:   conn,
		client: client,
	}
}

func (c *PluginClient) Call(ctx context.Context, method, args string) (err error) {
	req := &proto.CallReq{
		Method: method,
		Args:   args,
	}
	_, err = c.client.Call(ctx, req)
	return
}
