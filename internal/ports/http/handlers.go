package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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

func (s Server) handleCreateSwipe(w http.ResponseWriter, r *http.Request) {
	userID, err := getJWTUserID(r)
	if err != nil {
		s.log.WarnContext(r.Context(), "no auth user", slog.Any("error", err))
		errStatusCode(w, http.StatusUnauthorized)

		return
	}

	req := createSwipeReq{} //nolint:exhaustruct

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		s.log.DebugContext(r.Context(), "could not unmarshal swipe", slog.Any("error", err))

		errStatusCode(w, http.StatusBadRequest)

		return
	}

	err = s.swipesUC.CreateSwipe(r.Context(),
		domain.Swipe{
			Init:       userID,
			Target:     domain.UserID(req.Target),
			InitResp:   &req.Like,
			TargetResp: nil,
		})
	if err != nil {
		s.log.ErrorContext(r.Context(), "failed to create swipe", slog.Any("error", err))
		errStatusCode(w, http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s Server) handleGetSwipes(w http.ResponseWriter, r *http.Request) {
	userID, err := getJWTUserID(r)
	if err != nil {
		s.log.WarnContext(r.Context(), "no auth user", slog.Any("error", err))
		errStatusCode(w, http.StatusUnauthorized)

		return
	}

	pag := pagination(r)

	swipes, err := s.swipesUC.MySwipes(r.Context(), userID, pag)
	if err != nil {
		s.log.ErrorContext(r.Context(), "failed to get swipes", slog.Any("error", err))
		errStatusCode(w, http.StatusInternalServerError)

		return
	}

	ids := make([]uuid.UUID, 0, len(swipes))
	for _, swipe := range swipes {
		ids = append(ids, uuid.UUID(swipe.Init))
	}

	writeJSON(w, ids)
}

func (s Server) handleGetMatches(w http.ResponseWriter, r *http.Request) {
	userID, err := getJWTUserID(r)
	if err != nil {
		s.log.WarnContext(r.Context(), "no auth user", slog.Any("error", err))
		errStatusCode(w, http.StatusUnauthorized)

		return
	}

	pag := pagination(r)

	matches, err := s.matchesUC.Matches(r.Context(), userID, pag)
	if err != nil {
		s.log.ErrorContext(r.Context(), "failed to get matches", slog.Any("error", err))
		errStatusCode(w, http.StatusInternalServerError)

		return
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

func (s Server) handleHealthy(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func pagination(r *http.Request) domain.Pagination {
	pag := domain.Pagination{
		Offset: 0,
		Limit:  25, //nolint:mnd
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err == nil && page >= 0 {
		pag.Offset = page
	}

	limit, err := strconv.Atoi(limitStr)
	if err == nil && limit > 0 {
		pag.Limit = limit
	}

	return pag
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
