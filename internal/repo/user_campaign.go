package repo

import (
	"assignment-pe/internal/cx"
	"assignment-pe/internal/errs"
	"assignment-pe/internal/model"
	"assignment-pe/internal/util"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// nolint: lll
type UserCampaignRepo interface {
	Get(ctx context.Context, userID string, campaignID string) (*model.UserCampaign, error)
	ListTasks(ctx context.Context, userID string, campaignID string) ([]*model.UserCampaignTask, error)
	UpsertAmountToCampaignAndTask(ctx context.Context, swapInfo model.SwapInfo) ([]model.PointHistory, []int64, error)
	GetLeaderboard(ctx context.Context, campaignID string) ([]model.LeaderboardRow, error)
}

type userCampaignRepo struct {
}

func NewUserCampaignRepo() UserCampaignRepo {
	return &userCampaignRepo{}
}

type dbUserCampaign struct {
	UserID      string  `db:"user_id"`
	CampaignID  string  `db:"campaign_id"`
	IsCompleted bool    `db:"is_completed"`
	Amount      float64 `db:"amount"`
	Points      int64   `db:"points"`
}

func (d dbUserCampaign) toEntity() model.UserCampaign {
	return model.UserCampaign(d)
}

func (repo *userCampaignRepo) Get(
	ctx context.Context,
	userID string,
	campaignID string,
) (*model.UserCampaign, error) {
	tx := cx.GetTx(ctx)

	stmt := `
		SELECT user_id, campaign_id, is_completed, amount, points
		FROM user_campaigns
		WHERE user_id = ?
		AND campaign_id = ?
	`
	var userCampaign dbUserCampaign
	err := tx.GetContext(ctx, &userCampaign, tx.Rebind(stmt), userID, campaignID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errs.ErrInternal.Rewrap(err)
	}

	userCampaignEntity := userCampaign.toEntity()
	return &userCampaignEntity, nil
}

type dbUserCampaignTask struct {
	UserID         string  `db:"user_id"`
	CampaignTaskID string  `db:"campaign_task_id"`
	CampaignID     string  `db:"campaign_id"`
	Amount         float64 `db:"amount"`
	Points         int64   `db:"points"`
}

func (d dbUserCampaignTask) toEntity() model.UserCampaignTask {
	return model.UserCampaignTask(d)
}

func (repo *userCampaignRepo) ListTasks(
	ctx context.Context,
	userID string,
	campaignID string,
) ([]*model.UserCampaignTask, error) {
	tx := cx.GetTx(ctx)
	stmt := `
		SELECT user_id, campaign_task_id, campaign_id, amount, points
		FROM user_campaign_tasks
		WHERE user_id = ?
		AND campaign_id = ?
	`

	var userCampaignTasks []dbUserCampaignTask
	err := tx.SelectContext(ctx, &userCampaignTasks, tx.Rebind(stmt), userID, campaignID)
	if err != nil {
		return nil, errs.ErrInternal.Rewrap(err)
	}

	return util.Map(userCampaignTasks, func(in dbUserCampaignTask) *model.UserCampaignTask {
		e := in.toEntity()
		return &e
	}), nil
}

func (repo *userCampaignRepo) UpsertAmountToCampaignAndTask(
	ctx context.Context,
	swapInfo model.SwapInfo,
) ([]model.PointHistory, []int64, error) {
	tx := cx.GetTx(ctx)
	upsertRows := len(swapInfo.ActiveCampaignWithTasks)

	// Upsert campaigns
	stmt := fmt.Sprintf(`
		INSERT INTO user_campaigns (
			user_id, campaign_id, is_completed, amount, points
		) VALUES %v
		ON CONFLICT (user_id, campaign_id) 
		DO UPDATE SET
			is_completed = (
				CASE WHEN user_campaigns.amount + EXCLUDED.amount >= 1000
				THEN TRUE ELSE FALSE END
			),
			amount = user_campaigns.amount + EXCLUDED.amount,
			points = (
				CASE WHEN user_campaigns.amount < 1000
					AND user_campaigns.amount + EXCLUDED.amount >= 1000
				THEN user_campaigns.points + 100 ELSE user_campaigns.points END
			)
		RETURNING user_campaigns.amount, user_campaigns.points;
	`, strings.Repeat("(?, ?, ?, ?, ?)", upsertRows))

	points := 0
	isCompleted := false
	if swapInfo.Amount >= 1000 {
		points = 100
		isCompleted = true
	}

	args := make([]any, 0, upsertRows*5)
	for i := 0; i < upsertRows; i++ {
		args = append(args,
			swapInfo.UserID,
			swapInfo.ActiveCampaignWithTasks[i].CampaignID,
			isCompleted,
			swapInfo.Amount,
			points,
		)
	}

	type returnDO struct {
		Amount float64 `db:"amount"`
		Points int64   `db:"points"`
	}
	var rows []returnDO
	err := tx.SelectContext(ctx, &rows, tx.Rebind(stmt), args...)
	if err != nil {
		return nil, nil, errs.ErrInternal.Rewrap(err)
	}

	// Check if each campaign is the first time completing the onboarding task
	isOnboardingRows := make([]bool, len(rows))
	for i := 0; i < len(rows); i++ {
		isOnboardingRows[i] = rows[i].Amount >= 1000 && rows[i].Amount-swapInfo.Amount < 1000
	}

	// Upsert tasks
	stmt = fmt.Sprintf(`
		INSERT INTO user_campaign_tasks (
			user_id, campaign_task_id, campaign_id, amount, points
		) VALUES %v
		ON CONFLICT (user_id, campaign_task_id) 
		DO UPDATE SET
			amount = user_campaign_tasks.amount + EXCLUDED.amount,
			points = user_campaign_tasks.points + EXCLUDED.points;
	`, strings.Repeat("(?, ?, ?, ?, ?)", upsertRows))

	pointHistories := []model.PointHistory{}
	args = make([]any, 0, upsertRows*5)
	for i := 0; i < upsertRows; i++ {
		args = append(args,
			swapInfo.UserID,
			swapInfo.ActiveCampaignWithTasks[i].CampaignTaskID,
			swapInfo.ActiveCampaignWithTasks[i].CampaignID,
			swapInfo.Amount,
		)
		if isOnboardingRows[i] {
			args = append(args, 100)

			pointHistories = append(pointHistories, model.PointHistory{
				UserID:         swapInfo.UserID,
				Points:         100,
				CampaignID:     swapInfo.ActiveCampaignWithTasks[i].CampaignID,
				CampaignTaskID: swapInfo.ActiveCampaignWithTasks[i].CampaignTaskID,
			})
		} else {
			args = append(args, 0)
		}
	}

	_, err = tx.ExecContext(ctx, tx.Rebind(stmt), args...)
	if err != nil {
		return nil, nil, errs.ErrInternal.Rewrap(err)
	}

	return pointHistories,
		util.Map(rows, func(in returnDO) int64 {
			return in.Points
		}),
		nil
}

// TODO: complete it
func (repo *userCampaignRepo) GetLeaderboard(
	ctx context.Context,
	campaignID string,
) ([]model.LeaderboardRow, error) {
	stmt := `
		SELECT user_id, points
		FROM user_campaigns
		WHERE campaign_id = ?
	`
	_ = stmt
	return []model.LeaderboardRow{
		{UserID: "fakeUser1", Points: 100},
		{UserID: "fakeUser2", Points: 200},
	}, nil
}
