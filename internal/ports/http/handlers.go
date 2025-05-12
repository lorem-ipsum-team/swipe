package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lorem-ipsum-team/swipe/internal/domain"
)

var (
	errNoAuth        = errors.New("no Authorization header")
	errInvalidClaims = errors.New("invalid token claims")
)

func (s *Server) handleCreateSwipe(w http.ResponseWriter, r *http.Request) {
	userID, err := getJWTUserID(r)
	if err != nil {
		errStatusCode(w, http.StatusUnauthorized)
	}

	req := createSwipeReq{} //nolint:exhaustruct

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errStatusCode(w, http.StatusBadRequest)
	}

	err = s.swipesUC.CreateSwipe(r.Context(), userID, domain.UserID(req.Target), req.Like)
	if err != nil {
		errStatusCode(w, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleGetSwipes(w http.ResponseWriter, r *http.Request) {
	userID, err := getJWTUserID(r)
	if err != nil {
		errStatusCode(w, http.StatusUnauthorized)
	}

	offset, limit := getOffsetLimit(r)

	swipes, err := s.swipesUC.MySwipes(r.Context(), userID, offset, limit)
	if err != nil {
		errStatusCode(w, http.StatusInternalServerError)
	}

	ids := make([]uuid.UUID, 0, len(swipes))
	for _, swipe := range swipes {
		ids = append(ids, uuid.UUID(swipe.Init))
	}

	writeJSON(w, ids)
}

func (s *Server) handleGetMatches(w http.ResponseWriter, r *http.Request) {
	userID, err := getJWTUserID(r)
	if err != nil {
		errStatusCode(w, http.StatusUnauthorized)
	}

	offset, limit := getOffsetLimit(r)

	matches, err := s.matchesUC.Matches(r.Context(), userID, offset, limit)
	if err != nil {
		errStatusCode(w, http.StatusInternalServerError)
	}

	ids := make([]uuid.UUID, 0, len(matches))

	for _, match := range matches {
		if match.Target != userID {
			ids = append(ids, uuid.UUID(match.Target))
		} else {
			ids = append(ids, uuid.UUID(match.Init))
		}
	}

	writeJSON(w, ids)
}

func (s *Server) handleHealthy(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getOffsetLimit(r *http.Request) (int, int) {
	const (
		defaultPage  = 0
		defaultLimit = 20
	)

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		page = defaultPage
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}

	return page, limit
}

func getJWTUserID(r *http.Request) (domain.UserID, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return domain.UserID{}, errNoAuth
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return domain.UserID{}, fmt.Errorf("failed to parse jwt: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return domain.UserID{}, errInvalidClaims
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return domain.UserID{}, fmt.Errorf("failed to get subject from jwt: %w", err)
	}

	id, err := uuid.Parse(sub)

	return domain.UserID(id), err
}

type createSwipeReq struct {
	Target uuid.UUID `json:"targetId"`
	Like   bool      `json:"like"`
}

func errStatusCode(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func writeJSON(w http.ResponseWriter, j any) {
	err := json.NewEncoder(w).Encode(j)
	if err != nil {
		panic(err)
	}
}
