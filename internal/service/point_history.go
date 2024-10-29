package service

import (
	"assignment-pe/internal/model"
	"assignment-pe/internal/repo"
	"context"
)

// nolint: lll
type PointHistoryService interface {
	BatchInsert(ctx context.Context, batch []model.PointHistory) error
	ListPagination(ctx context.Context, condition model.PointHistorySearchCondition) ([]model.PointHistory, error)
}

type pointHistoryService struct {
	pointHistoryRepo repo.PointHistoryRepo
}

func NewPointHistoryService(
	pointHistoryRepo repo.PointHistoryRepo,
) PointHistoryService {
	return &pointHistoryService{
		pointHistoryRepo: pointHistoryRepo,
	}
}

func (srv *pointHistoryService) BatchInsert(
	ctx context.Context,
	batch []model.PointHistory,
) error {
	err := srv.pointHistoryRepo.BatchInsert(ctx, batch)
	if err != nil {
		return err
	}
	return nil
}

func (srv *pointHistoryService) ListPagination(
	ctx context.Context,
	condition model.PointHistorySearchCondition,
) ([]model.PointHistory, error) {
	pointHistories, err := srv.pointHistoryRepo.ListPagination(ctx, condition)
	if err != nil {
		return nil, err
	}
	return pointHistories, nil
}
