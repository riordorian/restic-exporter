package ports

import (
	"grpc/internal/infrastructure/ports/http"
)

type Services struct {
	HttpServer *http.Server
}
