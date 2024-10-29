package model

type Campaign struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	StartedAt     int64  `json:"startedAt"`
	EndedAt       int64  `json:"endedAt"`
	IsDistributed bool   `json:"isDistributed"`
}

type CampaignTask struct {
	ID            string `json:"id"`
	CampaignID    string `json:"campaignID"`
	Pair          string `json:"pair"`
	Points        int64  `json:"points"`
	StartedAt     int64  `json:"startedAt"`
	EndedAt       int64  `json:"endedAt"`
	IsDistributed bool   `json:"isDistributed"`
}

type ActiveCampaignWithTask struct {
	CampaignID     string `json:"campaignID"`
	CampaignTaskID string `json:"campaignTaskID"`
}

type LeaderboardRow struct {
	UserID string `json:"userID"`
	Points int64  `json:"points"`
}
