package handlers

import (
	"nxt_match_event_manager_api/internal/interfaces"

	"github.com/gin-gonic/gin"
)

type NewStudentHandler struct {
	studentService interfaces.StudentServiceInterface
}

func NewStudentHandlerInit(svc interfaces.StudentServiceInterface) *NewStudentHandler {
	return &NewStudentHandler{studentService: svc}
}

func (h *NewStudentHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/student", h.handleCreateStudent)
	router.PUT("/student/:id", h.handleUpdateStudent)
	router.DELETE("/student/:id", h.handleDeleteStudent)
	router.GET("/students", h.handleGetStudents)
}
