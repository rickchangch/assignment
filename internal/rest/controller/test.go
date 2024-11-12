package controller

import (
	"assignment-pe/internal/cx"
	"assignment-pe/internal/errs"
	"assignment-pe/internal/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TestController struct {
}

func NewTestController() *TestController {
	return &TestController{}
}

func (ctrl *TestController) Isolation(c *gin.Context) {
	level := c.Query("level") // RR,RC
	action := c.Query("action")
	sleepSec, _ := strconv.Atoi(c.Query("sleep"))

	ctx := c.Request.Context()
	tx := cx.GetTx(ctx)
	logger := cx.GetLogger(ctx).WithField("txid", util.GenXid())

	var err error
	if level == "RR" {
		_, err = tx.ExecContext(ctx, "BEGIN TRANSACTION ISOLATION LEVEL REPEATABLE READ")
	} else {
		_, err = tx.ExecContext(ctx, "BEGIN TRANSACTION ISOLATION LEVEL READ COMMITTED")
	}
	if err != nil {
		_ = c.Error(errs.ErrInvalidArgument.Rewrap(err))
		return
	}

	switch action {
	case "UPSERT":
		args := make([]any, 0, 5)
		stmt := `
		INSERT INTO user_campaign_tasks (
			user_id, campaign_task_id, campaign_id, amount, points
		) VALUES (?,?,?,?,?)
		ON CONFLICT (user_id, campaign_task_id) 
		DO UPDATE SET
			points = user_campaign_tasks.points + 1;
	`
		args = append(args,
			"testUser",
			"testCampaignTask",
			"testCampaign",
			1,
			1,
		)

		_, err = tx.ExecContext(ctx, tx.Rebind(stmt), args...)
		if err != nil {
			_ = c.Error(errs.ErrInternal.Rewrap(err))
			return
		}

		if sleepSec > 0 {
			logger.Debug("executed")
			time.Sleep(10 * time.Second)
			logger.Debug("wakeup")
		}

	case "GET":
		stmt := `
			SELECT user_id, campaign_task_id, campaign_id, amount, points
			FROM user_campaign_tasks;
		`
		type dbUserCampaignTask struct {
			UserID         string  `db:"user_id"`
			CampaignTaskID string  `db:"campaign_task_id"`
			CampaignID     string  `db:"campaign_id"`
			Amount         float64 `db:"amount"`
			Points         int64   `db:"points"`
		}
		var dbTasks []dbUserCampaignTask
		err = tx.SelectContext(ctx, &dbTasks, stmt)
		if err != nil {
			_ = c.Error(errs.ErrInternal.Rewrap(err))
			return
		}

		logger.Debug("first time GET: ", dbTasks)

		time.Sleep(time.Duration(sleepSec) * time.Second)

		err = tx.SelectContext(ctx, &dbTasks, stmt)
		if err != nil {
			_ = c.Error(errs.ErrInternal.Rewrap(err))
			return
		}
		logger.Debug("second time GET: ", dbTasks)
	}

	c.Status(http.StatusOK)
}
