package compliment

type CreateComplimentRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Points      int    `json:"points"`
	ToUserID    int64  `json:"to_user_id"`
}

type MarkAsViewedRequest struct {
	IDs []int64 `json:"ids"`
}
