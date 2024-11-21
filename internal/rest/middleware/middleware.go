package middleware

import (
	"assignment-pe/internal/cx"
	"assignment-pe/internal/errs"
	"assignment-pe/internal/log"
	"assignment-pe/internal/rest/middleware/ratelimiter"
	"errors"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Middleware interface {
	Logger() gin.HandlerFunc
	Metrics() gin.HandlerFunc
	Trace() gin.HandlerFunc
	AccessLog() gin.HandlerFunc
	ErrorHandler() gin.HandlerFunc
	PanicRecovery() gin.HandlerFunc
	Auth() gin.HandlerFunc
	Tx() gin.HandlerFunc
	RateLimit() gin.HandlerFunc
}

type middleware struct {
	logger      log.AppLogger
	pgdb        *sqlx.DB
	ratelimiter ratelimiter.Ratelimiter
}

func NewMiddleware(
	logger log.AppLogger,
	pgdb *sqlx.DB,
	ratelimiter ratelimiter.Ratelimiter,
) Middleware {
	return &middleware{
		logger:      logger,
		pgdb:        pgdb,
		ratelimiter: ratelimiter,
	}
}

func (m *middleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := cx.SetLogger(c.Request.Context(), m.logger)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// TODO
func (m *middleware) Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// TODO
func (m *middleware) Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// TODO
func (m *middleware) AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		cx.GetLogger(c.Request.Context()).
			WithFields(log.Fields{
				"url":           c.Request.URL.String(),
				"method":        c.Request.Method,
				"response_time": time.Since(startTime).Milliseconds(),
			}).
			Info("access log")
	}
}

func (m *middleware) ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		logger := cx.GetLogger(c.Request.Context())
		if err := c.Errors.Last(); err != nil {
			if appErr, ok := err.Err.(errs.AppError); ok {
				logger.WithAppError(appErr).Error("app error occurred")
				c.AbortWithStatusJSON(appErr.StatusCode, appErr)
			} else {
				logger.WithError(err).Error("unexpected error occurred")
				c.AbortWithStatusJSON(http.StatusInternalServerError, errs.ErrInternal)
			}
		}
	}
}

func (m *middleware) PanicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger := cx.GetLogger(c.Request.Context())
				logger.WithFields(log.Fields{
					"error": err,
					"stack": string(debug.Stack()),
				}).Error("panic recovered")

				_ = c.Error(errs.ErrInternal.Rewrap(errors.New("panic")))
				c.Abort()
			}
		}()

		c.Next()
	}
}

// TODO
func (m *middleware) Tx() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := cx.GetLogger(c.Request.Context())

		tx, err := m.pgdb.Beginx()
		if err != nil {
			_ = c.Error(errs.ErrInternal.Rewrap(err))
			c.Abort()
			return
		}

		// Rollback when error happened
		defer func() {
			if err := recover(); err != nil || c.Errors.Last() != nil {
				logger.Debug("tx rollback triggered")

				if err := tx.Rollback(); err != nil {
					logger.WithError(err).
						WithField("fatal", true).
						Error("tx rollback failed")
				}

				// Rethrow to PanicRecovery layer
				if err != nil {
					panic(err)
				}
			}
		}()

		logger.Debug("tx started")
		c.Request = c.Request.WithContext(cx.SetTx(c.Request.Context(), tx))

		c.Next()

		// Commit
		if c.Errors.Last() == nil {
			if err := tx.Commit(); err != nil {
				logger.WithError(err).
					WithField("fatal", true).
					Error("tx commit failed")

				panic(err)
			}

			logger.Debug("tx committed")
		}
	}
}

// TODO: Implement JWT
func (m *middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("userID")
		ctx := cx.SetUserID(c.Request.Context(), userID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func (m *middleware) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		userID := cx.GetUserID(ctx)
		logger := cx.GetLogger(ctx)

		allow, retryAfter, err := m.ratelimiter.AllowByTokenBucket(ctx, userID, 1)
		if err != nil {
			_ = c.Error(errs.ErrInternal.Rewrap(err))
			c.Abort()
		}

		logger.Debug(userID, allow, retryAfter)

		if !allow {
			_ = c.Error(errs.ErrTooManyRequests.New())
			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.Abort()
		}

		c.Next()
	}
}
