package grpc

import (
	"accounts/internal/api"
	"accounts/internal/server/service"
	"context"
	"errors"
	"net"
	"sync"
	"unsafe"

	"google.golang.org/grpc"
)

type accountsServiceServer struct {
	service service.AccountsService
}

func (srv *accountsServiceServer) GetAmount(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	amount, err := srv.service.GetAmount(ctx, req.BalanceId)
	if err != nil {
		return nil, err
	}

	return &api.GetResponse{
		BalanceId: req.GetBalanceId(),
		Amount:    amount,
	}, nil
}

func (srv *accountsServiceServer) AddAmount(ctx context.Context, req *api.AddRequest) (*api.AddResponse, error) {
	if err := srv.service.AddAmount(ctx, req.BalanceId, req.Value); err != nil {
		return nil, err
	}

	return &api.AddResponse{}, nil
}

type statisticsServiceServer struct {
	service service.StatisticsService
}

func (srv *statisticsServiceServer) Reset(context.Context, *api.Empty) (*api.Empty, error) {
	srv.service.Reset()

	return &api.Empty{}, nil
}

type Server struct {
	accountsServiceServer   *accountsServiceServer
	statisticsServiceServer *statisticsServiceServer

	address string

	mu         sync.Mutex
	grpcServer *grpc.Server

	unaryInt  []grpc.UnaryServerInterceptor
	streamInt []grpc.StreamServerInterceptor
}

func NewServer(address string, accountsSvc service.AccountsService, statisticsSvc service.StatisticsService) *Server {
	return &Server{
		accountsServiceServer:   &accountsServiceServer{service: accountsSvc},
		statisticsServiceServer: &statisticsServiceServer{service: statisticsSvc},
		address:                 address,
	}
}

func (srv *Server) WithUnaryInterceptors(interceptors ...UnaryServerInterceptor) *Server {
	srv.unaryInt = append(srv.unaryInt, toOriginalGRPCUnaryInterceptors(interceptors)...)

	return srv
}

func (srv *Server) WithStreamInterceptors(interceptors ...StreamServerInterceptor) *Server {
	srv.streamInt = append(srv.streamInt, toOriginalGRPCStreamInterceptors(interceptors)...)

	return srv
}

func (srv *Server) Start() error {
	srv.mu.Lock()

	if srv.grpcServer != nil {
		srv.mu.Unlock()
		return errors.New("server already started")
	}

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(srv.unaryInt...),
		grpc.ChainStreamInterceptor(srv.streamInt...),
	)

	api.RegisterAccountsServiceServer(grpcSrv, srv.accountsServiceServer)
	api.RegisterStatisticsServiceServer(grpcSrv, srv.statisticsServiceServer)

	listener, err := net.Listen("tcp", srv.address)
	if err != nil {
		srv.mu.Unlock()

		return err
	}

	srv.grpcServer = grpcSrv
	srv.mu.Unlock()

	return grpcSrv.Serve(listener)
}

func (srv *Server) GracefulStop() {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.grpcServer != nil {
		srv.grpcServer.GracefulStop()
		srv.grpcServer = nil
	}
}

//UnaryServerInfo - это копия одноименного типа из пакета google.golang.org/grpc
type UnaryServerInfo struct {
	Server     interface{}
	FullMethod string
}

//UnaryHandler - это копия одноименного типа из пакета google.golang.org/grpc
type UnaryHandler func(ctx context.Context, req interface{}) (interface{}, error)

//UnaryServerInterceptor - это копия типа из пакета google.golang.org/grpc
type UnaryServerInterceptor func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error)

//ServerStream - это копия одноименного типа из пакета google.golang.org/grpc
type ServerStream interface {
	SetHeader(map[string][]string) error
	SendHeader(map[string][]string) error
	SetTrailer(map[string][]string)
	Context() context.Context
	SendMsg(m interface{}) error
	RecvMsg(m interface{}) error
}

//StreamServerInfo - это копия одноименного типа из пакета google.golang.org/grpc
type StreamServerInfo struct {
	FullMethod     string
	IsClientStream bool
	IsServerStream bool
}

//StreamHandler - это копия одноименного типа из пакета google.golang.org/grpc
type StreamHandler func(srv interface{}, stream ServerStream) error

//StreamServerInterceptor - это копия одноименного типа из пакета google.golang.org/grpc
type StreamServerInterceptor func(srv interface{}, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error

func toOriginalGRPCUnaryInterceptors(interceprors []UnaryServerInterceptor) []grpc.UnaryServerInterceptor {
	var res []grpc.UnaryServerInterceptor

	for _, f := range interceprors {
		res = append(res, func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return f(ctx, req, (*UnaryServerInfo)(info), UnaryHandler(handler))
		})
	}

	return res
}

func toOriginalGRPCStreamInterceptors(interceprors []StreamServerInterceptor) []grpc.StreamServerInterceptor {
	var res []grpc.StreamServerInterceptor

	for _, f := range interceprors {
		res = append(res, func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			return f(srv, *(*ServerStream)(unsafe.Pointer(&ss)), (*StreamServerInfo)(info), *(*StreamHandler)(unsafe.Pointer(&handler)))
		})
	}

	return res
}
