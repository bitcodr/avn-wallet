package helper

import (
	"google.golang.org/grpc"
)

func ClientConnection(connURL string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(connURL, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
