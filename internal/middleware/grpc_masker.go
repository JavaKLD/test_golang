package middleware

import (
	pb "dolittle2/proto"
)

func MaskSensitiveFields(req interface{}) interface{} {
	switch r := req.(type) {
	case *pb.CreateScheduleRequest:
		return struct {
			Aid_name    string
			Aid_per_day uint64
			Duration    int64
			UserID      string
		}{
			Aid_name:    r.AidName,
			Aid_per_day: r.AidPerDay,
			Duration:    r.Duration,
			UserID:      "***",
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
