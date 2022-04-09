package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/dao"
	"coolcar/auth/token"
	"coolcar/auth/wechat"
	"io/ioutil"
	"log"
	"os"
	"time"

	server "coolcar/shared/server"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/coolcar?readPreference=primary&ssl=false"))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}

	pkFile, err := os.Open("auth/private.key")
	if err != nil {
		logger.Fatal("cannot open private key", zap.Error(err))
	}
	pkBitys, err := ioutil.ReadAll(pkFile)
	if err != nil {
		logger.Fatal("cannot read private key", zap.Error(err))

	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBitys)
	if err != nil {
		logger.Fatal("cannot parse private key", zap.Error(err))
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:   "auth",
		Addr:   ":8081",
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			authpb.RegisterAuthServiceServer(s, &auth.Service{
				OpenIDResolver: &wechat.Service{
					AppID: "wx320cc87c906b9df2",
					// 不可明文存入
					AppSecret: "8adff52a469ba53ea8b178e036c54bfd",
				},
				Mongo:          dao.NewMongo(mongoClient.Database("coolcar")),
				Logger:         logger,
				TokenExpire:    15 * time.Second,
				TokenGenerator: token.NewJWTTokenGen("coolcar/auth", privKey),
			})
		},
	}))

	// lis, err := net.Listen("tcp", ":8081")
	// if err != nil {
	// 	logger.Fatal("cannot listen", zap.Error(err))
	// }

	// s := grpc.NewServer()
	// authpb.RegisterAuthServiceServer(s, &auth.Service{
	// 	OpenIDResolver: &wechat.Service{
	// 		AppID: "wx320cc87c906b9df2",
	// 		// 不可明文存入
	// 		AppSecret: "8adff52a469ba53ea8b178e036c54bfd",
	// 	},
	// 	Mongo:          dao.NewMongo(mongoClient.Database("coolcar")),
	// 	Logger:         logger,
	// 	TokenExpire:    2 * time.Hour,
	// 	TokenGenerator: token.NewJWTTokenGen("coolcar/auth", privKey),
	// })

	// err = s.Serve(lis)
	// logger.Fatal("cannot server", zap.Error(err))
}
