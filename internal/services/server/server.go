package server

import (
	"context"
	"fmt"
	"net"
	"server-events/internal/apiendpoints"
	"server-events/internal/apitransports"
	userapisvc "server-events/internal/services/usersvc"
	pbapiv1 "server-events/pkg/genproto/pb"
	"time"

	"server-events/internal/server/gateway" // Ensure this path is correct and the gateway package exists

	"github.com/soheilhy/cmux"

	"github.com/go-kit/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const defaultListenTime = 5 * time.Second

type APIServices struct {
	UserSvc   userapisvc.UserSvc
	Endpoints []gateway.RegisterEndpointFunc
}

type GrpcServer struct {
	server     *grpc.Server
	logger     log.Logger
	addr       string
	endpoints  []gateway.RegisterEndpointFunc
	serverType string
}

func GetDialOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(gateway.GrpcMaxSize),
			grpc.MaxCallSendMsgSize(gateway.GrpcMaxSize),
		),
	}
}

func DefaultAPIProxyEndpoints() []gateway.RegisterEndpointFunc {
	return []gateway.RegisterEndpointFunc{
		pbapiv1.RegisterUserSvcHandlerFromEndpoint,
	}
}

func NewAPIServer(svcs *APIServices, logger log.Logger) *GrpcServer {
	grpcServer := grpc.NewServer(

		grpc.MaxRecvMsgSize(gateway.GrpcMaxSize),
		grpc.MaxSendMsgSize(gateway.GrpcMaxSize),
	)

	//create gokit endpoints
	userEndpoints := apiendpoints.NewEndpoints(svcs.UserSvc)
	//userapisvc.RegisterUserSvcServer(grpcServer, userapisvc.NewUserAPIServer(userEndpoints, logger))
	pbapiv1.RegisterUserSvcServer(grpcServer, apitransports.NewUserAPIServer(userEndpoints, logger))
	reflection.Register(grpcServer)

	return &GrpcServer{
		server:     grpcServer,
		logger:     logger,
		addr:       fmt.Sprintf("0.0.0.0:%d", 8080),
		endpoints:  svcs.Endpoints,
		serverType: "grpc",
	}
}

func (s *GrpcServer) Start() error {
	var (
		listener net.Listener
		err      error
	)

	for i := 0; i < 3; i++ {
		if listener, err = net.Listen("tcp", s.addr); err != nil {
			s.logger.Log("msg", "failed to listen", "err", err)
			time.Sleep(defaultListenTime)
			continue
		}
		break
	}
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.addr, err)
	}

	// lis, err := gateway.NewListener(s.addr)
	// if err != nil {
	// 	return err
	// }

	// s.logger.Log("msg", "starting grpc server", "addr", s.addr)

	// go func() {
	// 	if err := s.server.Serve(lis); err != nil {
	// 		s.logger.Log("msg", "grpc server error", "err", err)
	// 	}
	// }()

	return s.serve(listener)

}

func (s *GrpcServer) serve(listener net.Listener) error {
	g, _ := errgroup.WithContext(context.Background())

	m := cmux.New(listener)

	grpcL := m.Match(cmux.HTTP2())
	httpL := m.Match(cmux.Any())

	httpGateway, err := gateway.NewProxyGateway(&gateway.Config{
		Logger: s.logger, Addr: s.addr,
		Endpoints: s.endpoints,
	})

	if err != nil {
		s.logger.Log(
			"error while serving proxy during serve - gateway disabled, err: %s", err,
		)
	} else {
		g.Go(func() error {
			s.logger.Log(
				"started http proxy server")

			return errors.WithStack(httpGateway.Serve(httpL))
		})
	}

	g.Go(func() error {
		s.logger.Log("started gRPC %s server", s.serverType)
		return errors.WithStack(s.server.Serve(grpcL))
	})

	g.Go(func() error {
		s.logger.Log("multiplexer serving on -> [%v]", listener.Addr().String())

		serveErr := m.Serve()

		s.logger.Log(
			"multiplexer error while serving, err: %s", serveErr,
		)

		return errors.WithStack(serveErr)
	})

	return g.Wait()
}
