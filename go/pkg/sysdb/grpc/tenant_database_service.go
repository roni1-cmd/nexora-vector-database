package grpc

import (
	"context"
	"errors"

	"github.com/chroma-core/chroma/go/pkg/grpcutils"
	"github.com/pingcap/log"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/chroma-core/chroma/go/pkg/common"
	"github.com/chroma-core/chroma/go/pkg/proto/coordinatorpb"
	"github.com/chroma-core/chroma/go/pkg/sysdb/coordinator/model"
)

func (s *Server) CreateDatabase(ctx context.Context, req *coordinatorpb.CreateDatabaseRequest) (*coordinatorpb.CreateDatabaseResponse, error) {
	res := &coordinatorpb.CreateDatabaseResponse{}
	createDatabase := &model.CreateDatabase{
		ID:     req.GetId(),
		Name:   req.GetName(),
		Tenant: req.GetTenant(),
	}
	_, err := s.coordinator.CreateDatabase(ctx, createDatabase)
	if err != nil {
		log.Error("error CreateDatabase", zap.String("request", req.String()), zap.Error(err))
		if errors.Is(err, common.ErrDatabaseUniqueConstraintViolation) {
			return res, grpcutils.BuildAlreadyExistsGrpcError(err.Error())
		}
		return res, grpcutils.BuildInternalGrpcError(err.Error())
	}
	log.Info("CreateDatabase success", zap.String("request", req.String()))
	return res, nil
}

func (s *Server) GetDatabase(ctx context.Context, req *coordinatorpb.GetDatabaseRequest) (*coordinatorpb.GetDatabaseResponse, error) {
	res := &coordinatorpb.GetDatabaseResponse{}
	getDatabase := &model.GetDatabase{
		Name:   req.GetName(),
		Tenant: req.GetTenant(),
	}
	database, err := s.coordinator.GetDatabase(ctx, getDatabase)
	if err != nil {
		log.Error("error GetDatabase", zap.String("request", req.String()), zap.Error(err))
		if err == common.ErrDatabaseNotFound || err == common.ErrTenantNotFound {
			return res, grpcutils.BuildNotFoundGrpcError(err.Error())
		}
		return res, grpcutils.BuildInternalGrpcError(err.Error())
	}
	res.Database = &coordinatorpb.Database{
		Id:     database.ID,
		Name:   database.Name,
		Tenant: database.Tenant,
	}
	return res, nil
}

func (s *Server) ListDatabases(ctx context.Context, req *coordinatorpb.ListDatabasesRequest) (*coordinatorpb.ListDatabasesResponse, error) {
	res := &coordinatorpb.ListDatabasesResponse{}
	listDatabases := &model.ListDatabases{
		Limit:  req.Limit,
		Offset: req.Offset,
		Tenant: req.GetTenant(),
	}
	databases, err := s.coordinator.ListDatabases(ctx, listDatabases)
	if err != nil {
		log.Error("error ListDatabases", zap.String("request", req.String()), zap.Error(err))
		if err == common.ErrTenantNotFound {
			return res, grpcutils.BuildNotFoundGrpcError(err.Error())
		}
		return res, grpcutils.BuildInternalGrpcError(err.Error())
	}
	for _, database := range databases {
		res.Databases = append(res.Databases, &coordinatorpb.Database{
			Id:     database.ID,
			Name:   database.Name,
			Tenant: database.Tenant,
		})
	}
	return res, nil
}

func (s *Server) DeleteDatabase(ctx context.Context, req *coordinatorpb.DeleteDatabaseRequest) (*coordinatorpb.DeleteDatabaseResponse, error) {
	deleteDatabase := &model.DeleteDatabase{
		Name:   req.GetName(),
		Tenant: req.GetTenant(),
	}
	err := s.coordinator.DeleteDatabase(ctx, deleteDatabase)
	if err != nil {
		log.Error("error DeleteDatabase", zap.String("request", req.String()), zap.Error(err))
		if errors.Is(err, common.ErrDatabaseNotFound) {
			return nil, grpcutils.BuildNotFoundGrpcError(err.Error())
		}
		return nil, grpcutils.BuildInternalGrpcError(err.Error())
	}
	return &coordinatorpb.DeleteDatabaseResponse{}, nil
}

func (s *Server) CreateTenant(ctx context.Context, req *coordinatorpb.CreateTenantRequest) (*coordinatorpb.CreateTenantResponse, error) {
	res := &coordinatorpb.CreateTenantResponse{}
	createTenant := &model.CreateTenant{
		Name: req.GetName(),
	}
	_, err := s.coordinator.CreateTenant(ctx, createTenant)
	if err != nil {
		log.Error("error CreateTenant", zap.String("request", req.String()), zap.Error(err))
		if err == common.ErrTenantUniqueConstraintViolation {
			return res, grpcutils.BuildAlreadyExistsGrpcError(err.Error())
		}
		return res, grpcutils.BuildInternalGrpcError(err.Error())
	}
	log.Info("CreateTenant success", zap.String("request", req.String()))
	return res, nil
}

func (s *Server) GetTenant(ctx context.Context, req *coordinatorpb.GetTenantRequest) (*coordinatorpb.GetTenantResponse, error) {
	res := &coordinatorpb.GetTenantResponse{}
	getTenant := &model.GetTenant{
		Name: req.GetName(),
	}
	tenant, err := s.coordinator.GetTenant(ctx, getTenant)
	if err != nil {
		log.Error("error GetTenant", zap.String("request", req.String()), zap.Error(err))
		if err == common.ErrTenantNotFound {
			return res, grpcutils.BuildNotFoundGrpcError(err.Error())
		}
		return res, grpcutils.BuildInternalGrpcError(err.Error())
	}
	res.Tenant = &coordinatorpb.Tenant{
		Name:         tenant.Name,
		ResourceName: tenant.ResourceName,
	}
	return res, nil
}

func (s *Server) SetLastCompactionTimeForTenant(ctx context.Context, req *coordinatorpb.SetLastCompactionTimeForTenantRequest) (*emptypb.Empty, error) {
	err := s.coordinator.SetTenantLastCompactionTime(ctx, req.TenantLastCompactionTime.TenantId, req.TenantLastCompactionTime.LastCompactionTime)
	if err != nil {
		log.Error("error SetTenantLastCompactionTime", zap.String("request", req.String()), zap.Error(err))
		return nil, grpcutils.BuildInternalGrpcError(err.Error())
	}
	log.Info("SetLastCompactionTimeForTenant success", zap.String("request", req.String()))
	return &emptypb.Empty{}, nil
}

func (s *Server) SetTenantResourceName(ctx context.Context, req *coordinatorpb.SetTenantResourceNameRequest) (*coordinatorpb.SetTenantResourceNameResponse, error) {
	err := s.coordinator.SetTenantResourceName(ctx, req.Id, req.ResourceName)
	if err != nil {
		log.Error("error SetTenantResourceName", zap.String("request", req.String()), zap.Error(err))
		return nil, grpcutils.BuildInternalGrpcError(err.Error())
	}
	return &coordinatorpb.SetTenantResourceNameResponse{}, nil
}

func (s *Server) GetLastCompactionTimeForTenant(ctx context.Context, req *coordinatorpb.GetLastCompactionTimeForTenantRequest) (*coordinatorpb.GetLastCompactionTimeForTenantResponse, error) {
	res := &coordinatorpb.GetLastCompactionTimeForTenantResponse{}
	tenantIDs := req.TenantId
	tenants, err := s.coordinator.GetTenantsLastCompactionTime(ctx, tenantIDs)
	if err != nil {
		log.Error("error GetLastCompactionTimeForTenant", zap.Any("tenantIDs", tenantIDs), zap.Error(err))
		return nil, grpcutils.BuildInternalGrpcError(err.Error())
	}
	for _, tenant := range tenants {
		res.TenantLastCompactionTime = append(res.TenantLastCompactionTime, &coordinatorpb.TenantLastCompactionTime{
			TenantId:           tenant.ID,
			LastCompactionTime: tenant.LastCompactionTime,
		})
	}
	return res, nil
}

func (s *Server) FinishDatabaseDeletion(ctx context.Context, req *coordinatorpb.FinishDatabaseDeletionRequest) (*coordinatorpb.FinishDatabaseDeletionResponse, error) {
	res, err := s.coordinator.FinishDatabaseDeletion(ctx, req)
	if err != nil {
		return nil, grpcutils.BuildInternalGrpcError(err.Error())
	}
	return res, nil
}
