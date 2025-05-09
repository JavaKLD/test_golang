package server

import (
	"context"
	"dolittle2/internal/domain/models"
	"dolittle2/internal/domain/services"
	"dolittle2/internal/utils"
	"dolittle2/pkg/middleware"
	"dolittle2/proto"
	"log/slog"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ScheduleServer struct {
	proto.UnimplementedScheduleServiceServer
	scheduleService scheduleService
}

func (s *ScheduleServer) CreateSchedule(_ context.Context, req *proto.CreateScheduleRequest) (*proto.CreateScheduleResponse, error) {
	schedule := &models.Schedule{
		AidName:   req.AidName,
		AidPerDay: req.AidPerDay,
		UserID:    req.UserId,
		Duration:  req.Duration,
	}

	schedule.CreatedAt = utils.RoundTime(time.Now())

	id, err := s.scheduleService.CreateSchedule(schedule)
	if err != nil {
		if err.Error() == "Запись с таким именем для пользователя уже существует" {
			return nil, status.Errorf(
				codes.AlreadyExists,
				"Запись с таким aid_name для данного пользователя уже существует",
			)
		} else {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"Лекарства принимаются с 8 до 22",
			)
		}
	}

	return &proto.CreateScheduleResponse{
		Id:      id,
		Message: "Запись создана",
	}, nil
}

func (s *ScheduleServer) GetUserSchedule(_ context.Context, req *proto.GetUserScheduleRequest) (*proto.GetUserScheduleResponse, error) {
	userID := req.GetId()
	if userID == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Не может быть равным 0",
		)
	}
	exists, _ := s.scheduleService.CheckUserExists(userID)

	if !exists {
		return nil, status.Error(
			codes.InvalidArgument,
			"Такого пользователя нет ",
		)
	}

	scheduleID, err := s.scheduleService.FindByUserID(userID)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Не удалось получить данные %v:",
			err,
		)
	}

	return &proto.GetUserScheduleResponse{
		Message:   "Успешный ответ с расписанием",
		Schedules: scheduleID,
	}, nil
}

func (s *ScheduleServer) GetSchedule(_ context.Context, req *proto.GetScheduleRequest) (*proto.GetDailyScheduleResponse, error) {
	userID := req.GetUserId()
	scheduleID := req.GetScheduleId()
	if userID == 0 || scheduleID == 0 {
		return nil, status.Errorf(
			codes.Internal,
			"Не указан user_id или schedule_id",
		)
	}

	scheduleTimes, err := s.scheduleService.GetDailySchedule(userID, scheduleID)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Ошибка вывода графика приема лекарств",
		)
	}

	var formattedTimes []string
	for _, t := range scheduleTimes {
		formattedTimes = append(formattedTimes, t.Format("15:04"))
	}

	return &proto.GetDailyScheduleResponse{
		FormattedTimes: formattedTimes,
	}, nil
}

func (s *ScheduleServer) GetNextTakings(_ context.Context, req *proto.GetNextTakingsRequest) (*proto.GetNextTakingsResponse, error) {
	userID := req.GetUserId()

	if userID == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Отсутствует user_id",
		)
	}

	exists, _ := s.scheduleService.CheckUserExists(userID)
	if !exists {
		return nil, status.Error(
			codes.InvalidArgument,
			"Такого пользователя нет ",
		)
	}

	nextTakings, err := s.scheduleService.GetNextTakings(userID)
	if err != nil {
		if err.Error() == "Нет ближайших приемов" {
			return &proto.GetNextTakingsResponse{
				Message: "Нет ближайших приемов",
			}, nil
		}
		return nil, status.Errorf(
			codes.Internal,
			"Ошибка получения ближайших приемов: %v",
			err,
		)
	}

	var schedule []*proto.KeyValuePair
	for k, v := range nextTakings {
		schedule = append(schedule, &proto.KeyValuePair{
			Key:   k,
			Value: strings.Join(v, ", "),
		})
	}

	return &proto.GetNextTakingsResponse{
		Schedule: schedule,
		Message:  "Успешно",
	}, nil
}

func StartGRPCServer(service *services.ScheduleService) {
	logger := slog.Default()

	var lis net.Listener
	var err error

	for {
		lis, err = net.Listen("tcp", "127.0.0.1:50051")
		if err != nil {
			logger.Error("Ошибка подключения gRPC сервера", "error", err)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.LoggingInterceptor(logger)),
	)

	proto.RegisterScheduleServiceServer(
		grpcServer,
		&ScheduleServer{scheduleService: service},
	)

	logger.Info("gRPC сервер на порте 50051")

	err = grpcServer.Serve(lis)
	if err != nil {
		logger.Error("Ошибка запуска gRPC сервера", "error", err)
	}
}
