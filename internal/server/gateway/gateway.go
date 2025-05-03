package gateway

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-kit/log"
)

type RegisterEndpointFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error

type Config struct {
	Endpoints []RegisterEndpointFunc
	Addr      string
	Dev       bool
	Logger    log.Logger
}

const (
	// GrpcMaxSize is the maxium size message for grpc.
	GrpcMaxSize = 524288000
)

func NewProxyGateway(cfg *Config) (*http.Server, error) {
	httpHandler := mux.NewRouter()
	muxcfg := muxconfig{
		logger: cfg.Logger,
	}

	proxyMux := runtime.NewServeMux(
		runtime.WithErrorHandler(muxcfg.proxyGatewayErrorEncoder),
		runtime.WithForwardResponseOption(muxcfg.responseEncoder),
	)
	proxyDialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(GrpcMaxSize),
			grpc.MaxCallSendMsgSize(GrpcMaxSize),
		),
	}

	// Register our gRPC endpoint with proxy
	for _, register := range cfg.Endpoints {
		if rerr := register(context.Background(), proxyMux, cfg.Addr, proxyDialOpts); rerr != nil {
			return nil, errors.WithStack(rerr)
		}
	}

	httpHandler.PathPrefix("/").Handler(proxyMux)
	server := &http.Server{
		Addr: cfg.Addr,
		BaseContext: func(l net.Listener) context.Context {
			return context.Background()
		},
		Handler:           httpHandler,
		ReadHeaderTimeout: time.Second * 60,
		ReadTimeout:       time.Second * 300,
		WriteTimeout:      time.Second * 300,
	}
	return server, nil
}
