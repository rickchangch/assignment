package model

type SwapInfo struct {
	UserID                  string                   `json:"user_id"`
	Pair                    string                   `json:"pair"`
	Amount                  float64                  `json:"amount"`
	ActiveCampaignWithTasks []ActiveCampaignWithTask `json:"-"`
}
