package loggers

import (
	"fmt"
	"log/slog"
	"net/http"
	cc "smeag_sms_be/internal/constants"
	common "smeag_sms_be/internal/models/common"
	"smeag_sms_be/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

/**
* Logs Definition
* @Maintainer : nxt_match dev
* @Version : 1.0.0
* @Description: Define the logging statement here
**/

func StatusOK(ctx *gin.Context, data any) {
	utils.LoggerHandler(
		slog.LevelInfo,
		fmt.Sprintf("%v", data),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusOK,
	)

	ctx.JSON(http.StatusOK, data)
}

func StatusNoContent(ctx *gin.Context, err any) {
	utils.LoggerHandler(
		slog.LevelError,
		fmt.Sprintf("%v", err),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusNoContent,
	)

	ctx.JSON(http.StatusNoContent, gin.H{"error": err})
}

func StatusUnauthorizedError(ctx *gin.Context, err any) {
	utils.LoggerHandler(
		slog.LevelError,
		fmt.Sprintf("%v", err),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusUnauthorized,
	)

	ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
}

func StatusInternalServerError(ctx *gin.Context, err any) {
	utils.LoggerHandler(
		slog.LevelError,
		fmt.Sprintf("%v", err),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusInternalServerError,
	)

	ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
}

func StatusBadRequestError(ctx *gin.Context, err any) {
	utils.LoggerHandler(
		slog.LevelError,
		fmt.Sprintf("%v", err),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusBadRequest,
	)

	ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
}

func StatusBadGateway(ctx *gin.Context, err any) {
	utils.LoggerHandler(
		slog.LevelError,
		fmt.Sprintf("%v", err),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusBadGateway,
	)

	ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
}

func StatusNotFound(ctx *gin.Context, err any) {
	utils.LoggerHandler(
		slog.LevelError,
		fmt.Sprintf("%v", err),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusNotFound,
	)

	ctx.JSON(http.StatusNotFound, gin.H{"error": err})
}

// ErrorHandler captures errors and returns a consistent JSON error response
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Step1: Process the request first.

		// Step2: Check if any errors were added to the context
		if len(c.Errors) > 0 {
			// Step3: Use the last error
			err := c.Errors.Last().Err

			errCode := http.StatusInternalServerError
			// Step4: Respond with a generic error message
			c.JSON(errCode, gin.H{
				"errCode": errCode,
				"message": err.Error(),
			})
		}
		// Any other steps if no errors are found
	}
}

// Not Severe Error
func StatusWarnNotFound(ctx *gin.Context, msg string) {
	utils.LoggerHandler(
		slog.LevelWarn,
		fmt.Sprintf("%v", msg),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		http.StatusNotFound,
	)
	ctx.JSON(http.StatusUnauthorized, gin.H{"warn": msg})
}

func GetCommonError(ctx *gin.Context, err string, errCode int) {

	r := &common.ErrorResponse{
		ErrorCode: http.StatusInternalServerError,
		TimeStamp: time.Now(),
		ErrorMsg:  err,
	}

	switch errCode {
	case http.StatusNoContent:
		StatusNoContent(ctx, r)
	case http.StatusBadRequest:
		StatusBadRequestError(ctx, r)
	case http.StatusInternalServerError:
		StatusInternalServerError(ctx, r)
	case http.StatusUnauthorized:
		StatusUnauthorizedError(ctx, r)
	case http.StatusBadGateway:
		StatusBadGateway(ctx, r)
	case http.StatusNotFound:
		StatusNotFound(ctx, r)
	default:
		slog.Error("unable to find error code")
	}
}

func SuccessResponse(ctx *gin.Context, data any) {
	StatusOK(ctx, &common.SuccessResponse{
		SuccessID:  bson.NewObjectID(),
		Status:     cc.SUCCESS,
		HttpCode:   http.StatusOK,
		ResponseAt: time.Now(),
		Body:       data,
	})
}
