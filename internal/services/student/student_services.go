package services

import (
	"nxt_match_event_manager_api/internal/interfaces"
	models "nxt_match_event_manager_api/internal/models/student"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type StudentService struct {
	repo interfaces.StudentRepositoryInterface
}

func NewStudentService(repo interfaces.StudentRepositoryInterface) *StudentService {
	return &StudentService{repo: repo}
}

func (s *StudentService) CreateStudentService(ctx *gin.Context, req *models.CreateStudentRequest) (bson.ObjectID, error) {
	return s.repo.CreateStudent(ctx, req)
}

func (s *StudentService) UpdateStudentService(ctx *gin.Context, id bson.ObjectID, req *models.UpdateStudentRequest) error {
	return s.repo.UpdateStudent(ctx, id, req)
}

func (s *StudentService) DeleteStudentService(ctx *gin.Context, id bson.ObjectID) error {
	return s.repo.DeleteStudent(ctx, id)
}

func (s *StudentService) GetStudentsByLimitService(ctx *gin.Context, cursor bson.ObjectID, limit int) (*[]models.Student, error) {
	return s.repo.GetStudentsByLimit(ctx, cursor, limit)
}