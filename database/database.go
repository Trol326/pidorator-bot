package database

import "context"

type Database interface {
	NewConnection(context.Context) error
	Disconnect(context.Context)
}
