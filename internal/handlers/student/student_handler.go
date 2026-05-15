package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	cc "nxt_match_event_manager_api/internal/constants"
	common "nxt_match_event_manager_api/internal/models/common"
	models "nxt_match_event_manager_api/internal/models/student"
	"nxt_match_event_manager_api/internal/utils"
	loggers "nxt_match_event_manager_api/internal/utils/loggers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *NewStudentHandler) handleCreateStudent(ctx *gin.Context) {
	var req models.CreateStudentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		loggers.StatusBadRequestError(ctx, fmt.Errorf("invalid payload"))
		return
	}

	if err := utils.Validate.Struct(&req); err != nil {
		loggers.StatusBadRequestError(ctx, err)
		return
	}

	id, err := h.studentService.CreateStudentService(ctx, &req)
	if err != nil {
		loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	loggers.StatusOK(ctx, &common.SuccessResponse{
		SuccessID:  bson.NewObjectID(),
		Status:     cc.SUCCESS_CREATE_STUDENT,
		HttpCode:   http.StatusOK,
		ResponseAt: time.Now(),
		Body:       gin.H{"id": id.Hex()},
	})
}

func (h *NewStudentHandler) handleUpdateStudent(ctx *gin.Context) {
	id, err := utils.IDValidation(ctx, "id")
	if err != nil {
		loggers.StatusBadRequestError(ctx, fmt.Errorf("invalid student id"))
		return
	}

	var req models.UpdateStudentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		loggers.StatusBadRequestError(ctx, fmt.Errorf("invalid payload"))
		return
	}

	if err := utils.Validate.Struct(&req); err != nil {
		loggers.StatusBadRequestError(ctx, err)
		return
	}

	if err := h.studentService.UpdateStudentService(ctx, id, &req); err != nil {
		loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	loggers.StatusOK(ctx, &common.SuccessResponse{
		SuccessID:  bson.NewObjectID(),
		Status:     cc.SUCCESS_UPDATE_STUDENT,
		HttpCode:   http.StatusOK,
		ResponseAt: time.Now(),
		Body:       "student updated successfully",
	})
}

func (h *NewStudentHandler) handleDeleteStudent(ctx *gin.Context) {
	id, err := utils.IDValidation(ctx, "id")
	if err != nil {
		loggers.StatusBadRequestError(ctx, fmt.Errorf("invalid student id"))
		return
	}

	if err := h.studentService.DeleteStudentService(ctx, id); err != nil {
		loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	loggers.StatusOK(ctx, &common.SuccessResponse{
		SuccessID:  bson.NewObjectID(),
		Status:     cc.SUCCESS_DELETE_STUDENT,
		HttpCode:   http.StatusOK,
		ResponseAt: time.Now(),
		Body:       "student deleted successfully",
	})
}

func (h *NewStudentHandler) handleGetStudents(ctx *gin.Context) {
	limitStr := ctx.Query(cc.LIMIT)
	cursorStr := ctx.Query(cc.CURSOR)

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		loggers.StatusBadRequestError(ctx, fmt.Errorf("invalid limit: must be a positive integer"))
		return
	}

	var cursor bson.ObjectID
	if cursorStr != "" {
		cursor, err = utils.MustHex(cursorStr)
		if err != nil {
			loggers.StatusBadRequestError(ctx, fmt.Errorf("invalid cursor id"))
			return
		}
	}

	students, err := h.studentService.GetStudentsByLimitService(ctx, cursor, limit)
	if err != nil {
		loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	var nextCursor string
	if len(*students) == limit {
		last := (*students)[len(*students)-1]
		if !bson.ObjectID.IsZero(last.Id) {
			nextCursor = last.Id.Hex()
		}
	}

	loggers.StatusOK(ctx, &common.SuccessResponse{
		SuccessID:  bson.NewObjectID(),
		Status:     cc.SUCCESS_GET_STUDENTS,
		HttpCode:   http.StatusOK,
		ResponseAt: time.Now(),
		Body: &models.StudentPaginationResponse{
			Students:   students,
			NextCursor: nextCursor,
			HasMore:    nextCursor != "",
		},
	})
}
