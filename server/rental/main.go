package main

import (
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/api/trip"
	server "coolcar/shared/server"

	"log"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:              "rental",
		Addr:              ":8082",
		AuthPublicKeyFile: "shared/auth/public.key",
		Logger:            logger,
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				Logger: logger,
			})
		},
	}))

	// // 给一个文件地址，建一个interceptor
	// in, err := auth.Interceptor("shared/auth/public.key")
	// if err != nil {
	// 	logger.Fatal("cannot create auth interceptor", zap.Error(err))
	// }
	// // 加入interceptor
	// s := grpc.NewServer(grpc.UnaryInterceptor(in))
	// rentalpb.RegisterTripServiceServer(s, &trip.Service{
	// 	Logger: logger,
	// })

	// err = s.Serve(lis)
	// logger.Fatal("cannot server", zap.Error(err))
}
