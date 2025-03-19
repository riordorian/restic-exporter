package ports

import (
	"grpc/internal/infrastructure/ports/grpc"
	"grpc/internal/infrastructure/ports/http"
)

type Services struct {
	GrpcServer *grpc.NewsServer
	HttpServer *http.Server
}
