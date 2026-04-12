package dto

type CreateSubscriptionRequest struct {
	Owner    string `json:"owner"`
	RepoName string `json:"repo_name"`
}

type SubscriptionResponse struct {
	Owner    string `json:"owner"`
	RepoName string `json:"repo_name"`
}

type ListSubscriptionsResponse struct {
	Subscriptions []SubscriptionResponse `json:"subscriptions"`
}

type SubscriptionsInfoResponse struct {
	Repositories []RepositoryInfoResponse `json:"repositories"`
}
