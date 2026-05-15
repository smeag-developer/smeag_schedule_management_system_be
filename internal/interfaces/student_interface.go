package interfaces

import (
	models "nxt_match_event_manager_api/internal/models/student"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type StudentRepositoryInterface interface {
	CreateStudent(ctx *gin.Context, req *models.CreateStudentRequest) (bson.ObjectID, error)
	UpdateStudent(ctx *gin.Context, id bson.ObjectID, req *models.UpdateStudentRequest) error
	DeleteStudent(ctx *gin.Context, id bson.ObjectID) error
	GetStudentsByLimit(ctx *gin.Context, cursor bson.ObjectID, limit int) (*[]models.Student, error)
}

type StudentServiceInterface interface {
	CreateStudentService(ctx *gin.Context, req *models.CreateStudentRequest) (bson.ObjectID, error)
	UpdateStudentService(ctx *gin.Context, id bson.ObjectID, req *models.UpdateStudentRequest) error
	DeleteStudentService(ctx *gin.Context, id bson.ObjectID) error
	GetStudentsByLimitService(ctx *gin.Context, cursor bson.ObjectID, limit int) (*[]models.Student, error)
}
