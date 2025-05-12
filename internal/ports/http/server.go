package http

import (
	"net/http"

	matches_uc "github.com/lorem-ipsum-team/swipe/internal/usecase/matches"
	swipes_uc "github.com/lorem-ipsum-team/swipe/internal/usecase/swipes"
)

type Server struct {
	Server    http.Server
	matchesUC matches_uc.UseCase
	swipesUC  swipes_uc.UseCase
}

func (s *Server) RegisterHandlers() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /swipes", s.handleGetSwipes)
	mux.HandleFunc("POST /swipes", s.handleCreateSwipe)

	mux.HandleFunc("GET /matches", s.handleGetMatches)

	mux.HandleFunc("GET /healthy", s.handleHealthy)

	s.Server.Handler = mux
}
