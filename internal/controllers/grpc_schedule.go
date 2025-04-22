package controllers

import (
	"context"
	"dolittle2/internal/models"
	"dolittle2/internal/services"
	"dolittle2/internal/utils"
	pb "dolittle2/proto"
	"log/slog"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ScheduleServer struct {
	pb.UnimplementedScheduleServiceServer
	Service *services.ScheduleService
}

func (s *ScheduleServer) CreateSchedule(ctx context.Context, req *pb.CreateScheduleRequest) (*pb.CreateScheduleResponse, error) {
	schedule := &models.Schedule{
		Aid_name:    req.AidName,
		Aid_per_day: req.AidPerDay,
		UserID:      req.UserId,
		Duration:    req.Duration,
	}

	schedule.Create_at = utils.RoundTime(time.Now())

	id, err := s.Service.CreateSchedule(schedule)
	if err != nil {
		if err.Error() == "Запись с таким именем для пользователя уже существует" {
			return nil, status.Errorf(codes.AlreadyExists, "Запись с таким aid_name для данного пользователя уже существует")
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "Лекарства принимаются с 8 до 22")
		}
	}

	return &pb.CreateScheduleResponse{
		Id:      id,
		Message: "Запись создана",
	}, nil
}

func (s *ScheduleServer) GetUserSchedule(ctx context.Context, req *pb.GetUserScheduleRequest) (*pb.GetUserScheduleResponse, error) {
	userID := req.GetId()
	if userID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Не может быть равным 0")
	}

	scheduleID, err := s.Service.FindByUserID(userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Не удалось получить данные %v:", err)
	}

	return &pb.GetUserScheduleResponse{
		Message:   "Успешный ответ с расписанием",
		Schedules: scheduleID,
	}, nil
}

func (s *ScheduleServer) GetSchedule(ctx context.Context, req *pb.GetScheduleRequest) (*pb.GetDailyScheduleResponse, error) {
	userID := req.GetUserId()
	scheduleID := req.GetScheduleId()
	if userID == 0 || scheduleID == 0 {
		return nil, status.Errorf(codes.Internal, "Не указан user_id или schedule_id")
	}

	scheduleTimes, err := s.Service.GetDailySchedule(userID, scheduleID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Ошибка вывода графика приема лекарств")
	}

	var formattedTimes []string
	for _, t := range scheduleTimes {
		formattedTimes = append(formattedTimes, t.Format("15:04"))
	}

	return &pb.GetDailyScheduleResponse{
		FormattedTimes: formattedTimes,
	}, nil
}

func (s *ScheduleServer) GetNextTakings(ctx context.Context, req *pb.GetNextTakingsRequest) (*pb.GetNextTakingsResponse, error) {
	userID := req.GetUserId()

	if userID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Отсутствует user_id")
	}

	nextTakings, err := s.Service.GetNextTakings(userID)
	if err != nil {
		if err.Error() == "Нет ближайших приемов" {
			return &pb.GetNextTakingsResponse{
				Message: "Нет ближайших приемов",
			}, nil
		}
	}

	var schedule []*pb.KeyValuePair
	for k, v := range nextTakings {
		schedule = append(schedule, &pb.KeyValuePair{
			Key:   k,
			Value: strings.Join(v, ", "),
		})
	}

	return &pb.GetNextTakingsResponse{
		Schedule: schedule,
		Message:  "Успешно",
	}, nil
}

func StartGRPCServer(service *services.ScheduleService) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		slog.Error("Ошибка подключения gRPC сервера")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterScheduleServiceServer(grpcServer, &ScheduleServer{Service: service})

	err = grpcServer.Serve(lis)
	if err != nil {
		slog.Error("Ошибка ")
	}
}
