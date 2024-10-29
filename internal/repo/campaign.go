package repo

import (
	"assignment-pe/internal/cx"
	"assignment-pe/internal/errs"
	"assignment-pe/internal/model"
	"assignment-pe/internal/util"
	"context"
	"fmt"
	"time"
)

// nolint: lll
type CampaignRepo interface {
	ListActiveCampaignWithTask(ctx context.Context, swapInfo model.SwapInfo) ([]model.ActiveCampaignWithTask, error)
}

type campaignRepo struct {
}

func NewCampaignRepo() CampaignRepo {
	return &campaignRepo{}
}

type dbCampaign struct {
	ID            string `db:"id"`
	Name          string `db:"name"`
	StartedAt     int64  `db:"started_at"`
	EndedAt       int64  `db:"ended_at"`
	IsDistributed bool   `db:"is_distributed"`
}

func (d dbCampaign) toEntity() model.Campaign {
	return model.Campaign(d)
}

type dbCampaignTask struct {
	ID            string `db:"id"`
	CampaignID    string `db:"campaign_id"`
	Pair          string `db:"pair"`
	Points        int64  `db:"points"`
	StartedAt     int64  `db:"started_at"`
	EndedAt       int64  `db:"ended_at"`
	IsDistributed bool   `db:"is_distributed"`
}

func (d dbCampaignTask) toEntity() model.CampaignTask {
	return model.CampaignTask(d)
}

type dbActiveCampaignWithTask struct {
	CampaignID     string `db:"campaign_id"`
	CampaignTaskID string `db:"campaign_task_id"`
}

func (d dbActiveCampaignWithTask) toEntity() model.ActiveCampaignWithTask {
	return model.ActiveCampaignWithTask(d)
}

func (repo *campaignRepo) ListActiveCampaignWithTask(
	ctx context.Context,
	swapInfo model.SwapInfo,
) ([]model.ActiveCampaignWithTask, error) {
	tx := cx.GetTx(ctx)
	now := time.Now().UnixMilli()
	stmt := fmt.Sprintf(`
		SELECT c.id AS campaign_id, ct.id AS campaign_task_id
		FROM campaigns c
		LEFT OUTER JOIN campaign_tasks ct ON c.id = ct.campaign_id
		WHERE %d BETWEEN ct.started_at AND ct.ended_at
		AND ct.pair = ?
	`, now)

	var result []dbActiveCampaignWithTask
	err := tx.SelectContext(ctx, &result, tx.Rebind(stmt), swapInfo.Pair)
	if err != nil {
		return nil, errs.ErrInternal.Rewrap(err)
	}

	return util.Map(result, func(in dbActiveCampaignWithTask) model.ActiveCampaignWithTask {
		return in.toEntity()
	}), nil
}
