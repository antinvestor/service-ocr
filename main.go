package main

import (
	"context"
	"fmt"
	"github.com/antinvestor/apis"
	fapi "github.com/antinvestor/service-files-api"
	ocr "github.com/antinvestor/service-ocr-api"
	"github.com/antinvestor/service-ocr/config"
	"github.com/antinvestor/service-ocr/service/handlers"
	"github.com/antinvestor/service-ocr/service/models"
	"github.com/antinvestor/service-ocr/service/queue"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pitabwire/frame"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
)

func main() {

	serviceName := "service_ocr"

	ctx := context.Background()

	datasource := frame.GetEnv(config.EnvDatabaseUrl, "postgres://ant:@nt@localhost:5423/service_ocr")
	mainDb := frame.Datastore(ctx, datasource, false)

	readOnlydatasource := frame.GetEnv(config.EnvReplicaDatabaseUrl, datasource)
	readDb := frame.Datastore(ctx, readOnlydatasource, true)

	sysService := frame.NewService(serviceName, mainDb, readDb)

	isMigration, err := strconv.ParseBool(frame.GetEnv(config.EnvMigrate, "false"))
	if err != nil {
		isMigration = false
	}

	stdArgs := os.Args[1:]
	if (len(stdArgs) > 0 && stdArgs[0] == "migrate") || isMigration {

		migrationPath := frame.GetEnv(config.EnvMigrationPath, "./migrations/0001")
		err := sysService.MigrateDatastore(ctx, migrationPath, &models.OcrLog{})
		if err != nil {
			log.Fatalf("main -- Could not migrate successfully because : %v", err)
		}

		return
	}

	var serviceOptions []frame.Option

	filesServiceURL := frame.GetEnv(config.EnvFilesServiceUri, "127.0.0.1:7005")

	oauth2ServiceHost := frame.GetEnv(config.EnvOauth2ServiceUri, "")
	oauth2ServiceURL := fmt.Sprintf("%s/oauth2/token", oauth2ServiceHost)
	oauth2ServiceSecret := frame.GetEnv(config.EnvOauth2ServiceClientSecret, "")

	filesCli, err := fapi.NewFilesClient(ctx,
		apis.WithEndpoint(filesServiceURL), apis.WithTokenEndpoint(oauth2ServiceURL),
		apis.WithTokenUsername(serviceName), apis.WithTokenPassword(oauth2ServiceSecret))

	jwtAudience := frame.GetEnv(config.EnvOauth2JwtVerifyAudience, serviceName)
	jwtIssuer := frame.GetEnv(config.EnvOauth2JwtVerifyIssuer, "")

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpcctxtags.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
			frame.UnaryAuthInterceptor(jwtAudience, jwtIssuer),
		)),
		grpc.StreamInterceptor(frame.StreamAuthInterceptor(jwtAudience, jwtIssuer)),
	)

	implementation := &handlers.OCRServer{
		Service:  sysService,
		FilesCli: filesCli,
	}
	ocr.RegisterOCRServiceServer(grpcServer, implementation)

	grpcServerOpt := frame.GrpcServer(grpcServer)
	serviceOptions = append(serviceOptions, grpcServerOpt)

	ocrSyncQueueHandler := queue.NewOCRQueueHandler(sysService)
	ocrSyncQueueURL := frame.GetEnv(config.EnvQueueOcrSync, fmt.Sprintf("mem://%s", config.QueueOcrSyncName))
	ocrSyncQueue := frame.RegisterSubscriber(config.QueueOcrSyncName, ocrSyncQueueURL, 2, ocrSyncQueueHandler)
	ocrSyncQueueP := frame.RegisterPublisher(config.QueueOcrSyncName, ocrSyncQueueURL)
	serviceOptions = append(serviceOptions, ocrSyncQueue, ocrSyncQueueP)

	sysService.Init(serviceOptions...)

	serverPort := frame.GetEnv(config.EnvServerPort, "7012")

	log.Printf(" main -- Initiating server operations on : %s", serverPort)
	err = sysService.Run(ctx, fmt.Sprintf(":%v", serverPort))
	if err != nil {
		log.Printf("main -- Could not run Server : %v", err)
	}

}
