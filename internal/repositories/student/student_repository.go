package repository

import (
	"time"

	models "smeag_sms_be/internal/models/student"
	"smeag_sms_be/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *StudentRepository) CreateStudent(ctx *gin.Context, req *models.CreateStudentRequest) (bson.ObjectID, error) {
	doc := models.Student{
		Id:               bson.NewObjectID(),
		StudentIdNumber:  req.StudentIdNumber,
		FullName:         req.FullName,
		Age:              req.Age,
		Nationality:      req.Nationality,
		Sex:              req.Sex,
		EnglishName:      req.EnglishName,
		Level:            req.Level,
		DurationOfCourse: req.DurationOfCourse,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
		Courses:          req.Courses,
		SATeacher:        req.SATeacher,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return insertDoc(ctx, r.StudentsCol, doc)
}

func (r *StudentRepository) UpdateStudent(ctx *gin.Context, id bson.ObjectID, req *models.UpdateStudentRequest) error {
	fields, err := utils.StructToBsonM(req)
	if err != nil {
		return err
	}

	fields["updated_at"] = time.Now()

	return updateDoc(ctx, r.StudentsCol, id, fields)
}

func (r *StudentRepository) DeleteStudent(ctx *gin.Context, id bson.ObjectID) error {
	return deleteDoc(ctx, r.StudentsCol, id)
}

func (r *StudentRepository) GetStudentsByLimit(ctx *gin.Context, cursor bson.ObjectID, limit int) (*[]models.Student, error) {
	filter := bson.D{}
	if cursor != bson.NilObjectID {
		filter = bson.D{{Key: "_id", Value: bson.M{"$gt": cursor}}}
	}

	return findWithPagination[models.Student](ctx, r.StudentsCol, filter, int64(limit))
}
