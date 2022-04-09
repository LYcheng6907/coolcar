package trip

import (
	"context"

	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/api/trip/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Logger *zap.Logger
	Mongo  *dao.Mongo
}

// 创建TripService
func (s *Service) CreateTrip(c context.Context, req *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// 读
func (s *Service) GetTrip(c context.Context, req *rentalpb.GetTripRequest) (*rentalpb.Trip, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (s *Service) GetTrips(c context.Context, req *rentalpb.GetTripsRequest) (*rentalpb.GetTripsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// 更新
func (s *Service) UpdateTrip(c context.Context, req *rentalpb.UpdateTripRequest) (*rentalpb.Trip, error) {
	aid, err := auth.AccountIDFromcontext(c)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "")
	}

	// 可能用户两个终端设备同时读取，并更新导致冲突，使用乐观锁解决冲突
	tid := id.TripID(req.Id)
	tr, err := s.Mongo.GetTrip(c, tid, aid)
	if req.Current != nil {
		tr.Trip.Current = s.calCurrentStatus(tr.Trip, req.Current)

	}

	if req.EndTrip {
		tr.Trip.End = tr.Trip.Current
		tr.Trip.Status = rentalpb.TripStatus_FINISHED
	}
	// 如果有人改变过值，就会更新UpdatedAt，旧的UpdatedAt会被拒绝，导致更新失败
	s.Mongo.UpdateTrip(c, tid, aid, tr.UpdatedAt, tr.Trip)
	return nil, status.Error(codes.Unimplemented, "")
}

// 辅助函数，计算当前位置的状态（位置，公里数，钱等）
func (s *Service) calCurrentStatus(trip *rentalpb.Trip, cur *rentalpb.Location) *rentalpb.LocationStatus {

	return nil
}
