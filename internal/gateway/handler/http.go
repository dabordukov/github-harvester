package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github-harvester/internal/gateway/service"

	"google.golang.org/grpc/metadata"
)

type RepoResponse struct {
	Name         string `json:"name"`
	Owner        string `json:"owner"`
	Description  string `json:"description"`
	Forks        int    `json:"forks"`
	Stars        int    `json:"stars"`
	CreatedAt    string `json:"created_at"`
	CommitsCount int    `json:"commits_count"`
}

type RepoHandler struct {
	svc *service.HarvesterService
}

func NewRepoHandler(svc *service.HarvesterService) *RepoHandler {
	return &RepoHandler{svc: svc}
}

// GetRepo godoc
// @Summary      Get GitHub Repository Stats
// @Description  Returns stars, forks, and commit count for a public repo
// @Tags         repositories
// @Accept       json
// @Produce      json
// @Param        owner   path      string  true  "Repository Owner"
// @Param        repo    path      string  true  "Repository Name"
// @Success      200     {object}  model.RepositoryModel
// @Failure      404     {object}  string
// @Failure      500     {object}  string
// @Router       /repo/{owner}/{repo} [get]
func (h *RepoHandler) GetRepo(w http.ResponseWriter, r *http.Request) {
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")

	token := r.Header.Get("Authorization")
	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(r.Context(), md)

	data, err := h.svc.GetRepositoryInfo(ctx, owner, repo)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, service.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, service.ErrUnauthorized):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := RepoResponse{
		Name:         data.Name,
		Owner:        data.Owner,
		Description:  data.Description,
		Forks:        int(data.Forks),
		Stars:        int(data.Stars),
		CreatedAt:    data.CreatedAt,
		CommitsCount: int(data.CommitsCount),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode json response: %v", err)
	}
}
