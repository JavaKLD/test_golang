package middleware

import (
	pb "dolittle2/proto"
)

func MaskSensitiveFields(req interface{}) interface{} {
	switch r := req.(type) {
	case *pb.CreateScheduleRequest:
		return struct {
			AidName   string
			AidPerDay uint64
			Duration  int64
			UserID    string
		}{
			AidName:   r.AidName,
			AidPerDay: r.AidPerDay,
			Duration:  r.Duration,
			UserID:    "***",
		}
	case *pb.GetUserScheduleRequest:
		return struct {
			ID string
		}{
			ID: "***",
		}
	case *pb.GetScheduleRequest:
		return struct {
			UserID     string
			ScheduleID uint64
		}{
			UserID:     "***",
			ScheduleID: r.ScheduleId,
		}
	case *pb.GetNextTakingsRequest:
		return struct {
			UserID string
		}{
			UserID: "***",
		}
	default:
		return ""
	}
}
