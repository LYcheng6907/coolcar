package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	rentalpb "coolcar/rental/api/gen/v1"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{
			EnumsAsInts: true,
			OrigName:    true,
		},
	))

	serverConfig := []struct {
		name         string
		addr         string
		RegisterFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	}{
		{
			name:         "auth",
			addr:         "localhost:8081",
			RegisterFunc: authpb.RegisterAuthServiceHandlerFromEndpoint,
		},
		{
			name:         "rental",
			addr:         "localhost:8082",
			RegisterFunc: rentalpb.RegisterTripServiceHandlerFromEndpoint,
		},
	}

	for _, s := range serverConfig {

		err := s.RegisterFunc(
			c, mux, s.addr, []grpc.DialOption{grpc.WithInsecure()},
		)
		if err != nil {
			logger.Sugar().Fatalf("cannot register  service %s : %v", s.name, err)
		}
	}
	addr := ":8080"
	logger.Sugar().Infof("grpc gateway started at %s", addr)
	logger.Sugar().Fatal(http.ListenAndServe(addr, mux))
}
