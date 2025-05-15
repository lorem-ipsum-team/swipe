package http

import (
	"log/slog"
	"net/http"
	"time"

	matches_uc "github.com/lorem-ipsum-team/swipe/internal/usecase/matches"
	swipes_uc "github.com/lorem-ipsum-team/swipe/internal/usecase/swipes"
)

type Server struct {
	Server    *http.Server
	matchesUC matches_uc.UseCase
	swipesUC  swipes_uc.UseCase
	log       *slog.Logger
}

func New(
	log *slog.Logger,
	addr string,
	matchesUC matches_uc.UseCase,
	swipesUC swipes_uc.UseCase,
) Server {
	log = log.WithGroup("http_server")
	serv := Server{
		Server: &http.Server{ //nolint:exhaustruct
			Addr:              addr,
			ReadHeaderTimeout: time.Second / 2,
		},
		matchesUC: matchesUC,
		swipesUC:  swipesUC,
		log:       log,
	}

	log.Info("register handlers")
	serv.registerHandlers()

	return serv
}

func (s Server) registerHandlers() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /swipes", s.handleGetSwipes)
	mux.HandleFunc("POST /swipes", s.handleCreateSwipe)

	mux.HandleFunc("GET /matches", s.handleGetMatches)

	mux.HandleFunc("GET /healthy", s.handleHealthy)

	s.Server.Handler = loggingMiddleware(mux, s.log)
}

func loggingMiddleware(next http.Handler, log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(lrw, r)

		log.DebugContext(r.Context(), "request", slog.Group(
			"request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", lrw.statusCode),
			slog.Duration("dur", time.Since(start)),
			slog.String("remote_ip", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		))
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
