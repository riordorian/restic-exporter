package repository

import (
	"github.com/google/uuid"
)

type Request interface{}
type BaseRepositoryInterface interface {
	Delete(id uuid.UUID) (bool, error)
}
