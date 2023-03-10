package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/web-programming-fall-2022/airline-auth/internal/bootstrap"
	"github.com/web-programming-fall-2022/airline-auth/internal/bootstrap/job"
	"github.com/web-programming-fall-2022/airline-auth/internal/cfg"
	"github.com/web-programming-fall-2022/airline-auth/internal/storage"
	"github.com/web-programming-fall-2022/airline-auth/internal/token"
	pb "github.com/web-programming-fall-2022/airline-auth/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"regexp"
)

func RunServer(ctx context.Context, config cfg.Config) job.WithGracefulShutdown {
	serverRunner, err := bootstrap.NewGrpcServerRunner(config.GrpcServerRunnerConfig)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	// Create the gRPC server
	grpcServer := serverRunner.GetGrpcServer()

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	store := storage.NewStorage(&config.MainDB)

	if err := store.Migrate(); err != nil {
		log.Fatal(err)
	}

	tokenManager := token.NewJWTManager(config.JWT.Secret, store, rdb)

	registerServer(
		grpcServer,
		tokenManager,
		store,
		config.JWT.AuthTokenExpire,
		config.JWT.RefreshTokenExpire,
	)

	go func() {
		defer store.Close()
		logrus.Infoln("Starting grpc server...")
		if err := serverRunner.Run(ctx); err != nil {
			logrus.Fatal(err.Error())
		}
	}()
	return serverRunner
}

func registerServer(
	server *grpc.Server,
	tokenManager token.Manager,
	storage *storage.Storage,
	authTokenExpire int64,
	refreshTokenExpire int64,
) {
	pb.RegisterAuthServiceServer(server, NewAuthServiceServer(
		tokenManager,
		storage,
		authTokenExpire,
		refreshTokenExpire,
	))
}

func RunHttpServer(ctx context.Context, config cfg.Config) job.WithGracefulShutdown {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", config.Server.Port), opts); err != nil {
		logrus.Fatal("Failed to start HTTP gateway", err.Error())
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HttpServer.Port),
		Handler: cors(mux),
	}

	logrus.Info("Starting HTTP/REST Gateway...", srv.Addr)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			logrus.Fatal("Failed to start HTTP gateway", err.Error())
		}
	}()
	return srv
}

func allowedOrigin(origin string) bool {
	if viper.GetString("cors") == "*" {
		return true
	}
	if matched, _ := regexp.MatchString(viper.GetString("cors"), origin); matched {
		return true
	}
	return false
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowedOrigin(r.Header.Get("Origin")) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
		}
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
