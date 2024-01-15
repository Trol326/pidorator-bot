package database

import "context"

type Database interface {
	NewConnection(ctx context.Context) error
	Disconnect()
}
