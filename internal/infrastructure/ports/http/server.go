package http

import (
	"context"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	proadapt "restic-exporter/internal/infrastructure/adapters/prometheus"
	"time"
)

type Server struct {
	ctx        context.Context
	router     RouterInterface
	listenPort int
	srv        *http.Server
}

func GetServer(port int, router RouterInterface) *Server {
	ctx := context.Background()

	return &Server{
		ctx:        ctx,
		listenPort: port,
		router:     router,
	}
}

func (s *Server) Serve() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	//s.router.HandleFunc("/metrics", s.metricsHandler)

	reg := prometheus.NewRegistry()
	reg.MustRegister(
		proadapt.NewResticCollector(),
	)

	s.router.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Debug: Metrics endpoint called")
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{
			Registry: reg,
		}).ServeHTTP(w, r)
	}))

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", s.listenPort),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.router,
	}
	s.srv = srv

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := srv.Shutdown(ctx); err != nil {
		log.Println(err.Error())
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func (s *Server) metricsHandler(w http.ResponseWriter, r *http.Request) {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		proadapt.NewResticCollector(),
	)

	promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})
	fmt.Println(r)
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
