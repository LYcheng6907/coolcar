package main

import (
	"context"
	trippb "coolcar/proto/gen/go"
	trip "coolcar/tripservice"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(log.Lshortfile)
	go startGRPCGateWay()
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		// 输完log程序退出
		// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	trippb.RegisterTripServiceServer(s, &trip.Service{})
	// Fatal is equivalent to Print() followed by a call to os.Exit(1).
	log.Fatal(s.Serve(lis))

}

func startGRPCGateWay() {
	c := context.Background()
	c, cancel := context.WithCancel(c)
	// 服务完成后断开连接
	defer cancel()

	// 连接注册在 runtime.NewServeMux()
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{
			EnumsAsInts: true, // 枚举类型返回具体值，而不是字符串
			OrigName:    true, // 驼峰命名变为小写和下滑线组合（前端小程序命名一般采用下划线组合）
		},
	))
	err := trippb.RegisterTripServiceHandlerFromEndpoint(
		c,
		mux,
		"localhost:8081", // grpc的端口
		[]grpc.DialOption{grpc.WithInsecure()},
	)
	if err != nil {
		log.Fatalf("cannot start grpc gateway: %v", err)
	}
	// 建立与grpc的连接
	err = http.ListenAndServe(":8080", mux) // http 端口
	if err != nil {
		log.Fatalf("cannot listen and serve: %v", err)
	}
}
