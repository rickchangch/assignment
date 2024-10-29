package model

type UserCampaign struct {
	UserID      string  `json:"userID"`
	CampaignID  string  `json:"campaignID"`
	IsCompleted bool    `json:"isCompleted"`
	Amount      float64 `json:"amount"`
	Points      int64   `json:"points"`
}

type UserCampaignTask struct {
	UserID         string  `json:"userID"`
	CampaignTaskID string  `json:"campaignTaskID"`
	CampaignID     string  `json:"campaignID"`
	Amount         float64 `json:"amount"`
	Points         int64   `json:"points"`
}
