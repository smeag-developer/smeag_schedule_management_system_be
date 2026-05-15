package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type StudentCourse struct {
	Name     string `bson:"name" json:"name" validate:"required"`
	Duration int    `bson:"duration" json:"duration" validate:"required,min=1"`
}

type SATeacher struct {
	TeacherId string `bson:"teacher_id" json:"teacher_id"`
	Name      string `bson:"name" json:"name"`
}

type Student struct {
	Id               bson.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	StudentIdNumber  string          `bson:"student_id_number" json:"student_id_number"`
	FullName         string          `bson:"fullname" json:"fullname"`
	Age              int             `bson:"age" json:"age"`
	Nationality      string          `bson:"nationality" json:"nationality"`
	Sex              string          `bson:"sex" json:"sex"`
	EnglishName      string          `bson:"english_name" json:"english_name"`
	Level            string          `bson:"level" json:"level"`
	DurationOfCourse int             `bson:"duration_of_course" json:"duration_of_course"`
	StartDate        time.Time       `bson:"start_date" json:"start_date"`
	EndDate          time.Time       `bson:"end_date" json:"end_date"`
	Courses          []StudentCourse `bson:"courses" json:"courses"`
	SATeacher        SATeacher       `bson:"sa_teacher" json:"sa_teacher"`
	CreatedAt        time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time       `bson:"updated_at" json:"updated_at"`
}

type CreateStudentRequest struct {
	StudentIdNumber  string          `json:"student_id_number" validate:"required"`
	FullName         string          `json:"fullname" validate:"required"`
	Age              int             `json:"age" validate:"required,min=1,max=120"`
	Nationality      string          `json:"nationality" validate:"required"`
	Sex              string          `json:"sex" validate:"required,oneof=male female other"`
	EnglishName      string          `json:"english_name" validate:"required"`
	Level            string          `json:"level" validate:"required"`
	DurationOfCourse int             `json:"duration_of_course" validate:"required,min=1"`
	StartDate        time.Time       `json:"start_date" validate:"required"`
	EndDate          time.Time       `json:"end_date" validate:"required"`
	Courses          []StudentCourse `json:"courses" validate:"required,min=1,dive"`
	SATeacher        SATeacher       `json:"sa_teacher"`
}

type UpdateStudentRequest struct {
	FullName         string          `json:"fullname" bson:"fullname,omitempty"`
	Age              int             `json:"age" bson:"age,omitempty" validate:"omitempty,min=1,max=120"`
	Nationality      string          `json:"nationality" bson:"nationality,omitempty"`
	Sex              string          `json:"sex" bson:"sex,omitempty" validate:"omitempty,oneof=male female other"`
	EnglishName      string          `json:"english_name" bson:"english_name,omitempty"`
	Level            string          `json:"level" bson:"level,omitempty"`
	DurationOfCourse int             `json:"duration_of_course" bson:"duration_of_course,omitempty" validate:"omitempty,min=1"`
	StartDate        *time.Time      `json:"start_date" bson:"start_date,omitempty"`
	EndDate          *time.Time      `json:"end_date" bson:"end_date,omitempty"`
	Courses          []StudentCourse `json:"courses" bson:"courses,omitempty" validate:"omitempty,min=1,dive"`
	SATeacher        *SATeacher      `json:"sa_teacher" bson:"sa_teacher,omitempty"`
}

type StudentPaginationResponse struct {
	Students   interface{} `json:"students"`
	NextCursor string      `json:"nextCursor"`
	HasMore    bool        `json:"hasMore"`
}
