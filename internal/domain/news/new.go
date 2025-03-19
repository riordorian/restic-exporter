package news

import (
	"github.com/google/uuid"
	"grpc/internal/domain/tags"
	"grpc/internal/shared"
	"time"
)

const (
	INACTIVE = iota
	ACTIVE
	DRAFT
)

type News struct {
	Id          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Text        string    `json:"text" db:"text"`
	Status      int32     `json:"status" db:"status"`
	Media       []shared.Media
	Attachments []shared.Media
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	AcceptedBy  uuid.UUID `json:"accepted_by" db:"accepted_by"`
	CreatedBy   uuid.UUID `json:"created_by" db:"created_by"`
	Tags        []tags.Tag
}
