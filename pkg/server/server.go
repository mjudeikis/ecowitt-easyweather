package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mjudeikis/ecowitt-easyweather/pkg/config"
	"github.com/mjudeikis/ecowitt-easyweather/pkg/utils/ratelimiter"
	"github.com/mjudeikis/ecowitt-easyweather/pkg/utils/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	log          *zap.Logger
	server       *http.Server
	router       *mux.Router
	config       *config.Config
	rateLimiters rateLimiterBuckets
}

type rateLimiterBuckets struct {
	generic ratelimiter.Interface
}

func New(
	log *zap.Logger,
	config *config.Config,
) (*Server, error) {
	rlb := rateLimiterBuckets{
		generic: ratelimiter.NewQuantumRateLimiter("generic", 50, time.Second),
	}

	s := &Server{
		log:          log,
		config:       config,
		rateLimiters: rlb,
	}

	s.router = s.setupRouter()

	apiRouter := s.router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/ingest", s.gRateLimit(s.ingest)).Methods("POST")
	apiRouter.HandleFunc("/metrics", s.gRateLimit(promhttp.Handler().ServeHTTP)).Methods("GET")

	s.server = &http.Server{
		Addr: config.ServerURI,
		Handler: handlers.CORS(
			handlers.AllowCredentials(),
			handlers.AllowedHeaders([]string{"Content-Type"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE"}),
			handlers.AllowedOrigins(config.ServerAllowedOrigins),
		)(s),
	}

	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Run(ctx context.Context) error {
	s.log.Info("Starting API Service")
	go func() {
		defer recover.Panic(s.log)

		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		err := s.server.Shutdown(ctx)
		if err != nil {
			s.log.Error("api shutdown error", zap.Error(err))
		}
		s.log.Info("Stopped API Service")
	}()

	s.log.Info("Server will now listen", zap.String("url", s.config.ServerURI))
	return s.server.ListenAndServe()
}

func (s *Server) setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(Panic(s.log))
	r.Use(Gzip())
	r.Use(Log(s.log))

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
	})

	return r
}
