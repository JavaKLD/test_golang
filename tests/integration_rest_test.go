package tests

import (
	"context"
	"dolittle2/internal/domain/models"
	"dolittle2/internal/domain/services"
	pb "dolittle2/proto"
	"encoding/json"
	"google.golang.org/grpc"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"dolittle2/internal/server"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockScheduleRepo struct{}

func (m *mockScheduleRepo) CreateSchedule(schedule *models.Schedule) (uint64, error) {
	return 1, nil
}

func (m *mockScheduleRepo) AidNameExists(aidName string, userID uint64) (bool, error) {
	return false, nil
}

func (m *mockScheduleRepo) UserIdExists(userID uint64) (bool, error) {
	return true, nil
}

func (m *mockScheduleRepo) FindByUserID(userID uint64) ([]uint64, error) {
	return []uint64{1, 2, 3}, nil
}

func (m *mockScheduleRepo) FindSchedule(userID, scheduleID uint64) (*models.Schedule, error) {
	return &models.Schedule{
		ID:          scheduleID,
		UserID:      userID,
		Aid_name:    "парацетамол",
		Aid_per_day: 3,
		Duration:    4,
		Create_at:   time.Now(),
	}, nil
}

func (m *mockScheduleRepo) NextTakings(userID uint64) ([]models.Schedule, error) {
	return []models.Schedule{
		{
			ID:          1,
			UserID:      userID,
			Aid_name:    "парацетамол",
			Aid_per_day: 3,
			Duration:    4,
			Create_at:   time.Now(),
		},
	}, nil
}

func TestCreateSchedule(t *testing.T) {
	e := echo.New()
	repo := &mockScheduleRepo{}
	svc := services.NewService(repo, time.Hour*24)
	controller := server.NewScheduleController(svc)

	e.POST("/schedule", controller.PostSchedule)

	body := `{
			"aid_name": "парацетамол", 
			"aid_per_day": 3,
			"duration": 4,
			"user_id": 12345
	}`

	req := httptest.NewRequest(echo.POST, "/schedule", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "Запись создана")
	t.Log("Resp body", rec.Body.String())
}

func TestGetUserSchedule_Success(t *testing.T) {
	e := echo.New()
	repo := &mockScheduleRepo{}
	svc := services.NewService(repo, time.Hour*24)
	controller := server.NewScheduleController(svc)

	e.GET("/schedules", controller.GetUserSchedule)

	req := httptest.NewRequest(http.MethodGet, "/schedules?user_id=1", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Успешный ответ с расписанием")
	t.Log("Response body:", rec.Body.String())
}

func TestGetSchedule_Success(t *testing.T) {
	e := echo.New()
	repo := &mockScheduleRepo{}
	svc := services.NewService(repo, time.Hour*24)
	controller := server.NewScheduleController(svc)

	e.GET("/schedule", controller.GetSchedule)

	req := httptest.NewRequest(http.MethodGet, "/schedule?user_id=1&schedule_id=1", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		ScheduleID []string `json:"schedule"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp.ScheduleID[0], ":")
	t.Log("Response body:", rec.Body.String())
}

func TestGetNextTakings_Success(t *testing.T) {
	e := echo.New()
	repo := &mockScheduleRepo{}
	svc := services.NewService(repo, time.Hour*24)
	controller := server.NewScheduleController(svc)

	e.GET("/next_takings", controller.GetNextTakings)

	req := httptest.NewRequest(http.MethodGet, "/next_takings?user_id=1", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Schedule map[string][]string `json:"schedule"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Contains(t, rec.Body.String(), "парацетамол")
	takings, ok := resp.Schedule["парацетамол"]
	assert.True(t, ok)

	assert.Greater(t, len(takings), 0, "список приёмов должен содержать хотя бы одно время")
	t.Log("Response body:", rec.Body.String())
}

func TestGRPCCreateSchedule_Success(t *testing.T) {
	repo := &mockScheduleRepo{}
	svc := services.NewService(repo, time.Hour*24)

	go server.StartGRPCServer(svc)

	time.Sleep(2 * time.Second)

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Error("Ошибка создания клиента", err)
	}
	defer conn.Close()

	client := pb.NewScheduleServiceClient(conn)

	ctx := context.Background()
	resp, err := client.CreateSchedule(ctx, &pb.CreateScheduleRequest{
		AidName:   "парацетамол",
		AidPerDay: 3,
		Duration:  4,
		UserId:    12345,
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint64(1), resp.Id)
	assert.Equal(t, "Запись создана", resp.Message)
	t.Log("Response body:", resp.String())
}

func TestGRPCGetUserSchedule_Success(t *testing.T) {
	repo := &mockScheduleRepo{}
	svc := services.NewService(repo, time.Hour*24)

	go server.StartGRPCServer(svc)

	time.Sleep(2 * time.Second)

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Error("Ошибка создания клиента", err)
	}
	defer conn.Close()

	client := pb.NewScheduleServiceClient(conn)

	ctx := context.Background()

	resp, err := client.GetUserSchedule(ctx, &pb.GetUserScheduleRequest{Id: 12345})
	if err != nil {
		t.Error(err)
	}
	if resp == nil {
		t.Error("Ожидалс	sch := resp.Schedules[0]я не nil")
	}

	assert.NotNil(t, resp)
	assert.Equal(t, "Успешный ответ с расписанием", resp.Message)
	assert.Equal(t, []uint64{1, 2, 3}, resp.Schedules)

	t.Log("Response body:", resp.String())
}

func TestGRPCGetSchedule_Success(t *testing.T) {
	repo := &mockScheduleRepo{}
	svc := services.NewService(repo, time.Hour*24)

	go server.StartGRPCServer(svc)

	time.Sleep(2 * time.Second)

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Error("Ошибка создания клиента", err)
	}
	defer conn.Close()

	client := pb.NewScheduleServiceClient(conn)

	ctx := context.Background()

	resp, err := client.GetSchedule(ctx, &pb.GetScheduleRequest{UserId: 12345, ScheduleId: 1})
	if err != nil {
		t.Error(err)
	}
	if resp == nil {
		t.Error("Ожидался не nil")
	}

	assert.NotNil(t, resp)
	assert.Contains(t, []string{""}, resp.FormattedTimes)
}
