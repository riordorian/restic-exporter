package ports

import (
	"restic-exporter/internal/infrastructure/ports/http"
)

type Services struct {
	HttpServer *http.Server
}
