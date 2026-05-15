package handlers

import (
	"smeag_sms_be/internal/interfaces"

	"github.com/gin-gonic/gin"
)

type NewStudentHandler struct {
	studentService interfaces.StudentServiceInterface
}

func NewStudentHandlerInit(svc interfaces.StudentServiceInterface) *NewStudentHandler {
	return &NewStudentHandler{studentService: svc}
}

func (h *NewStudentHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/student/create", h.handleCreateStudent)
	router.PUT("/student/update/:id", h.handleUpdateStudent)
	router.DELETE("/student/delete/:id", h.handleDeleteStudent)
	router.GET("/students/all", h.handleGetStudents)
}
