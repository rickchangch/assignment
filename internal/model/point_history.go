package model

type PointHistory struct {
	ID             string `json:"id"`
	UserID         string `json:"userID"`
	Points         int64  `json:"points"`
	CampaignID     string `json:"campaignID"`
	CampaignTaskID string `json:"campaignTaskID"`
	CreatedAt      int64  `json:"createdAt"`
}

type PointHistorySearchCondition struct {
	UserID string
	Cursor *string
	Size   int
}
