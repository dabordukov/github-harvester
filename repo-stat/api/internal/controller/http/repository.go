package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewRepositoryInfoHandler(log *slog.Logger, repositoryInfo *usecase.RepositoryInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawURL := r.URL.Query().Get("url")
		if rawURL == "" {
			writeError(w, http.StatusBadRequest, "url is required")
			return
		}

		owner, repo, err := parseRepositoryURL(rawURL)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		info, err := repositoryInfo.Fetch(r.Context(), owner, repo)
		if err != nil {
			log.Error("failed to get repository info", "error", err, "owner", owner, "repo", repo)
			writeGRPCError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(dto.RepositoryInfoResponse{
			Name:         info.Name,
			Owner:        info.Owner,
			Description:  info.Description,
			Forks:        info.Forks,
			Stars:        info.Stars,
			CreatedAt:    info.CreatedAt,
			CommitsCount: info.CommitsCount,
		}); err != nil {
			log.Error("failed to write repository response", "error", err)
		}
	}
}

func parseRepositoryURL(URL string) (string, string, error) {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return "", "", errors.New("invalid repository url")
	}

	if parsedURL.Host != "github.com" {
		return "", "", errors.New("invalid repository url")
	}

	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", errors.New("invalid repository url")
	}

	return parts[0], parts[1], nil
}

func writeGRPCError(w http.ResponseWriter, err error) {
	status, ok := status.FromError(err)
	if !ok {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	switch status.Code() {
	case codes.InvalidArgument:
		writeError(w, http.StatusBadRequest, status.Message())
	case codes.NotFound:
		writeError(w, http.StatusNotFound, status.Message())
	case codes.Unauthenticated:
		writeError(w, http.StatusUnauthorized, status.Message())
	case codes.ResourceExhausted:
		writeError(w, http.StatusTooManyRequests, status.Message())
	default:
		writeError(w, http.StatusInternalServerError, status.Message())
	}
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: message})
}
