package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

func NewCreateSubscriptionHandler(log *slog.Logger, subscription *usecase.Subscription) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateSubscriptionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		req.Owner = strings.TrimSpace(req.Owner)
		req.RepoName = strings.TrimSpace(req.RepoName)

		created, err := subscription.Create(r.Context(), req.Owner, req.RepoName)
		if err != nil {
			log.Error("failed to create subscription", "error", err, "owner", req.Owner, "repo", req.RepoName)
			writeGRPCError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(dto.SubscriptionResponse{
			Owner:    created.Owner,
			RepoName: created.RepoName,
		}); err != nil {
			log.Error("failed to write create subscription response", "error", err)
		}
	}
}

func NewDeleteSubscriptionHandler(log *slog.Logger, subscription *usecase.Subscription) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := strings.TrimSpace(r.PathValue("owner"))
		repo := strings.TrimSpace(r.PathValue("repo"))

		if err := subscription.Delete(r.Context(), owner, repo); err != nil {
			log.Error("failed to delete subscription", "error", err, "owner", owner, "repo", repo)
			writeGRPCError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func NewListSubscriptionsHandler(log *slog.Logger, subscription *usecase.Subscription) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subscriptions, err := subscription.List(r.Context())
		if err != nil {
			log.Error("failed to list subscriptions", "error", err)
			writeGRPCError(w, err)
			return
		}

		response := dto.ListSubscriptionsResponse{
			Subscriptions: make([]dto.SubscriptionResponse, 0, len(subscriptions)),
		}
		for _, item := range subscriptions {
			response.Subscriptions = append(response.Subscriptions, dto.SubscriptionResponse{
				Owner:    item.Owner,
				RepoName: item.RepoName,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write list subscriptions response", "error", err)
		}
	}
}

func NewSubscriptionsInfoHandler(log *slog.Logger, subscriptionInfo *usecase.SubscriptionInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repositories, err := subscriptionInfo.Fetch(r.Context())
		if err != nil {
			log.Error("failed to get subscriptions info", "error", err)
			writeGRPCError(w, err)
			return
		}

		response := dto.SubscriptionsInfoResponse{
			Repositories: make([]dto.RepositoryInfoResponse, 0, len(repositories)),
		}
		for _, repo := range repositories {
			response.Repositories = append(response.Repositories, dto.RepositoryInfoResponse{
				Name:         repo.Name,
				Owner:        repo.Owner,
				FullName:     repo.FullName,
				Description:  repo.Description,
				Forks:        repo.Forks,
				Stars:        repo.Stars,
				CreatedAt:    repo.CreatedAt,
				CommitsCount: repo.CommitsCount,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write subscriptions info response", "error", err)
		}
	}
}
