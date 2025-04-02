package application

import (
	"restic-exporter/internal/application/cqrs"
)

type Services struct {
	Dispatcher cqrs.DispatcherInterface
}
