package repo

import (
	"assignment-pe/internal/cx"
	"assignment-pe/internal/errs"
	"assignment-pe/internal/model"
	"assignment-pe/internal/util"
	"context"
	"fmt"
	"strings"
	"time"
)

// nolint: lll
type PointHistoryRepo interface {
	BatchInsert(ctx context.Context, batch []model.PointHistory) error
	ListPagination(ctx context.Context, condition model.PointHistorySearchCondition) ([]model.PointHistory, error)
}

type pointHistoryRepo struct {
}

func NewPointHistoryRepo() PointHistoryRepo {
	return &pointHistoryRepo{}
}

type dbPointHistory struct {
	ID             string `db:"id"`
	UserID         string `db:"user_id"`
	Points         int64  `db:"points"`
	CampaignID     string `db:"campaign_id"`
	CampaignTaskID string `db:"campaign_task_id"`
	CreatedAt      int64  `db:"created_at"`
}

func (d dbPointHistory) toEntity() model.PointHistory {
	return model.PointHistory(d)
}

func (repo *pointHistoryRepo) BatchInsert(
	ctx context.Context,
	batch []model.PointHistory,
) error {
	if len(batch) == 0 {
		return nil
	}

	tx := cx.GetTx(ctx)
	insertRows := len(batch)
	stmt := fmt.Sprintf(`
		INSERT INTO point_histories (
			id, user_id, points, campaign_id, campaign_task_id, created_at
		) VALUES %v
	`, strings.Repeat("(?, ?, ?, ?, ?, ?)", insertRows))

	now := time.Now().UnixMilli()
	args := make([]any, 0, insertRows*6)
	for i := 0; i < insertRows; i++ {
		args = append(args,
			util.GenXid(),
			batch[i].UserID,
			batch[i].Points,
			batch[i].CampaignID,
			batch[i].CampaignTaskID,
			now,
		)
	}

	_, err := tx.ExecContext(ctx, tx.Rebind(stmt), args...)
	if err != nil {
		return errs.ErrInternal.Rewrap(err)
	}

	return nil
}

func (repo *pointHistoryRepo) ListPagination(
	ctx context.Context,
	condition model.PointHistorySearchCondition,
) ([]model.PointHistory, error) {
	tx := cx.GetTx(ctx)

	cond := "user_id = ?"
	args := []any{condition.UserID}

	if condition.Cursor != nil {
		cond += " AND id < ?"
		args = append(args, *condition.Cursor)
	}

	args = append(args, condition.Size)

	stmt := fmt.Sprintf(`
		SELECT id, user_id, points, campaign_id, campaign_task_id, created_at
		FROM point_histories
		WHERE %s
		ORDER BY id DESC
		LIMIT ?
	`, cond)

	var pointHistories []dbPointHistory
	err := tx.SelectContext(ctx, &pointHistories, tx.Rebind(stmt), args...)
	if err != nil {
		return nil, errs.ErrInternal.Rewrap(err)
	}

	return util.Map(pointHistories, func(in dbPointHistory) model.PointHistory {
		return in.toEntity()
	}), nil
}
