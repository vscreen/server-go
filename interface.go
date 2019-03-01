package main

import (
	"context"

	"github.com/grandcat/zeroconf"
)

func publish(ctx context.Context) error {
	server, err := zeroconf.Register(name, serviceTag, "local.", port, nil, nil)
	if err != nil {
		return err
	}

	go func() {
		defer server.Shutdown()
		<-ctx.Done()
	}()
	return nil
}
