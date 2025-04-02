package cqrs

import (
	"context"
	"fmt"
)

type CommandInterface interface {
	CommandName() string
}

type QueryInterface interface {
	QueryName() string
}

type CommandHandlerInterface interface {
	Handle(ctx context.Context, cmd CommandInterface) error
}

type QueryHandlerInterface interface {
	Handle(ctx context.Context, query QueryInterface) (interface{}, error)
}

type commandHandlerFunc[C CommandInterface] func(ctx context.Context, cmd C) error

func NewCommandHandler[C CommandInterface](handler commandHandlerFunc[C]) CommandHandlerInterface {
	return &commandHandler[C]{handler: handler}
}

type commandHandler[C CommandInterface] struct {
	handler commandHandlerFunc[C]
}

func (h *commandHandler[C]) Handle(ctx context.Context, cmd CommandInterface) error {
	typedCmd, ok := cmd.(C)
	if !ok {
		return fmt.Errorf("invalid command type %T", cmd)
	}
	return h.handler(ctx, typedCmd)
}

type queryHandlerFunc[Q QueryInterface, R any] func(ctx context.Context, query Q) (R, error)

func NewQueryHandler(handler QueryHandlerInterface) QueryHandlerInterface {
	return handler
}

type queryHandler[Q QueryInterface, R any] struct {
	handler queryHandlerFunc[Q, R]
}

func (h *queryHandler[Q, R]) Handle(ctx context.Context, query QueryInterface) (interface{}, error) {
	typedQuery, ok := query.(Q)
	if !ok {
		return nil, fmt.Errorf("invalid query type %T", query)
	}
	return h.handler(ctx, typedQuery)
}

type DispatcherInterface interface {
	RegisterCommand(name string, handler CommandHandlerInterface)
	RegisterQuery(name string, handler QueryHandlerInterface)
	DispatchCommand(ctx context.Context, cmd CommandInterface) error
	DispatchQuery(ctx context.Context, query QueryInterface) (interface{}, error)
}

type Dispatcher struct {
	cmdHandlers map[string]CommandHandlerInterface
	qryHandlers map[string]QueryHandlerInterface
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		cmdHandlers: make(map[string]CommandHandlerInterface),
		qryHandlers: make(map[string]QueryHandlerInterface),
	}
}

func (d *Dispatcher) RegisterCommand(name string, handler CommandHandlerInterface) {
	d.cmdHandlers[name] = handler
}

func (d *Dispatcher) RegisterQuery(name string, handler QueryHandlerInterface) {
	d.qryHandlers[name] = handler
}

func (d *Dispatcher) DispatchCommand(ctx context.Context, cmd CommandInterface) error {
	handler, ok := d.cmdHandlers[cmd.CommandName()]
	if !ok {
		return fmt.Errorf("command handler not found: %s", cmd.CommandName())
	}
	return handler.Handle(ctx, cmd)
}

func (d *Dispatcher) DispatchQuery(ctx context.Context, query QueryInterface) (interface{}, error) {
	handler, ok := d.qryHandlers[query.QueryName()]

	fmt.Println(query.QueryName())
	if !ok {
		return nil, fmt.Errorf("query handler not found: %s", query.QueryName())
	}
	return handler.Handle(ctx, query)
}
