package concurrent

import "context"

type EventWithCtx struct {
	Event interface{}

	Ctx context.Context
}
