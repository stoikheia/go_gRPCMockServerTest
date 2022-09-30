package cancelablelistener

import (
	"context"
	"net"
)

func Listen(ctx context.Context, network, address string) (net.Listener, error) {
	lis, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			lis.Close()
		}
	}(ctx)
	return lis, nil
}
