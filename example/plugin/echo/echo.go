package main

import (
	"context"
	"fmt"
	pluginServer "github.com/faycheng/gob/plugin/server"
)

func Echo(c context.Context, args interface{}) error {
	fmt.Println(args)
	return nil
}

func main() {
	server := pluginServer.NewPluginServer("test.echo")
	server.Register("echo", Echo)
	err := server.Serve()
	if err != nil {
		panic(err)
	}
}
