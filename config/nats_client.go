package config

import (
	nats "github.com/nats-io/nats.go"
)


func NATSClient() (*nats.EncodedConn, error) {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}
	ns, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}
	return ns, nil
}
