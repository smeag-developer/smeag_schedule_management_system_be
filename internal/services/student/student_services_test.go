package services_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	models "smeag_sms_be/internal/models/student"
	services "smeag_sms_be/internal/services/student"
)

// mockStudentRepo is a manual mock that satisfies StudentRepositoryInterface.
type mockStudentRepo struct {
	createFn func(*gin.Context, *models.CreateStudentRequest) (bson.ObjectID, error)
	updateFn func(*gin.Context, bson.ObjectID, *models.UpdateStudentRequest) error
	deleteFn func(*gin.Context, bson.ObjectID) error
	getFn    func(*gin.Context, bson.ObjectID, int) (*[]models.Student, error)
}

func (m *mockStudentRepo) CreateStudent(ctx *gin.Context, req *models.CreateStudentRequest) (bson.ObjectID, error) {
	return m.createFn(ctx, req)
}
func (m *mockStudentRepo) UpdateStudent(ctx *gin.Context, id bson.ObjectID, req *models.UpdateStudentRequest) error {
	return m.updateFn(ctx, id, req)
}
func (m *mockStudentRepo) DeleteStudent(ctx *gin.Context, id bson.ObjectID) error {
	return m.deleteFn(ctx, id)
}
func (m *mockStudentRepo) GetStudentsByLimit(ctx *gin.Context, cursor bson.ObjectID, limit int) (*[]models.Student, error) {
	return m.getFn(ctx, cursor, limit)
}

func newTestCtx() *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	return ctx
}

func sampleCreateRequest() *models.CreateStudentRequest {
	return &models.CreateStudentRequest{
		StudentIdNumber:  "2024-0001",
		FullName:         "Juan Dela Cruz",
		Age:              20,
		Nationality:      "Filipino",
		Sex:              "male",
		EnglishName:      "John",
		Level:            "intermediate",
		DurationOfCourse: 4,
		StartDate:        time.Now(),
		EndDate:          time.Now().AddDate(0, 1, 0),
		Courses:          []models.StudentCourse{{Name: "ESL", Duration: 4}},
		SATeacher:        models.SATeacher{TeacherId: bson.NewObjectID().Hex(), Name: "Teacher A"},
	}
}

func TestCreateStudentService(t *testing.T) {
	expectedID := bson.NewObjectID()

	tests := []struct {
		name    string
		repoErr error
		wantID  bson.ObjectID
		wantErr bool
	}{
		{name: "success", repoErr: nil, wantID: expectedID, wantErr: false},
		{name: "repo error", repoErr: errors.New("db error"), wantID: bson.NilObjectID, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := services.NewStudentService(&mockStudentRepo{
				createFn: func(_ *gin.Context, _ *models.CreateStudentRequest) (bson.ObjectID, error) {
					if tt.repoErr != nil {
						return bson.NilObjectID, tt.repoErr
					}
					return expectedID, nil
				},
			})

			id, err := svc.CreateStudentService(newTestCtx(), sampleCreateRequest())

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.wantErr && id != tt.wantID {
				t.Errorf("expected id %v, got %v", tt.wantID, id)
			}
		})
	}
}

func TestUpdateStudentService(t *testing.T) {
	id := bson.NewObjectID()

	tests := []struct {
		name    string
		repoErr error
		wantErr bool
	}{
		{name: "success", repoErr: nil, wantErr: false},
		{name: "not found", repoErr: errors.New("student not found"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := services.NewStudentService(&mockStudentRepo{
				updateFn: func(_ *gin.Context, _ bson.ObjectID, _ *models.UpdateStudentRequest) error {
					return tt.repoErr
				},
			})

			err := svc.UpdateStudentService(newTestCtx(), id, &models.UpdateStudentRequest{FullName: "New Name"})

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestDeleteStudentService(t *testing.T) {
	id := bson.NewObjectID()

	tests := []struct {
		name    string
		repoErr error
		wantErr bool
	}{
		{name: "success", repoErr: nil, wantErr: false},
		{name: "not found", repoErr: errors.New("student not found"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := services.NewStudentService(&mockStudentRepo{
				deleteFn: func(_ *gin.Context, _ bson.ObjectID) error {
					return tt.repoErr
				},
			})

			err := svc.DeleteStudentService(newTestCtx(), id)

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetStudentsByLimitService(t *testing.T) {
	students := []models.Student{
		{Id: bson.NewObjectID(), FullName: "Juan Dela Cruz"},
		{Id: bson.NewObjectID(), FullName: "Maria Santos"},
	}

	tests := []struct {
		name    string
		result  *[]models.Student
		repoErr error
		wantLen int
		wantErr bool
	}{
		{name: "returns results", result: &students, repoErr: nil, wantLen: 2, wantErr: false},
		{name: "empty result", result: &[]models.Student{}, repoErr: nil, wantLen: 0, wantErr: false},
		{name: "repo error", result: nil, repoErr: errors.New("db error"), wantLen: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := services.NewStudentService(&mockStudentRepo{
				getFn: func(_ *gin.Context, _ bson.ObjectID, _ int) (*[]models.Student, error) {
					return tt.result, tt.repoErr
				},
			})

			res, err := svc.GetStudentsByLimitService(newTestCtx(), bson.NilObjectID, 10)

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.wantErr && len(*res) != tt.wantLen {
				t.Errorf("expected %d results, got %d", tt.wantLen, len(*res))
			}
		})
	}
}
