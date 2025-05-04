package gateway

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/log"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

type muxconfig struct {
	logger log.Logger
}

func (m *muxconfig) responseEncoder(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	// Implement the response encoding logic here
	// Example: Write a JSON response
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}

func (m *muxconfig) proxyGatewayErrorEncoder(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
