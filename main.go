package main

import (
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
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"strings"
)

func main() {

	serviceName := "service_ocr"

	var ocrConfig config.OcrConfig
	err := frame.ConfigProcess("", &ocrConfig)
	if err != nil {
		logrus.WithError(err).Fatal("could not process configs")
		return
	}

	ctx, service := frame.NewService(serviceName, frame.Config(&ocrConfig))
	defer service.Stop(ctx)
	log := service.L()

	serviceOptions := []frame.Option{frame.Datastore(ctx)}

	if ocrConfig.DoDatabaseMigrate() {

		service.Init(serviceOptions...)

		err := service.MigrateDatastore(ctx, ocrConfig.GetDatabaseMigrationPath(),
			&models.OcrLog{})

		if err != nil {
			log.Fatalf("main -- Could not migrate successfully because : %+v", err)
		}

		return
	}

	err = service.RegisterForJwt(ctx)
	if err != nil {
		log.WithError(err).Fatal("main -- could not register fo jwt")
	}

	oauth2ServiceHost := ocrConfig.GetOauth2ServiceURI()
	oauth2ServiceURL := fmt.Sprintf("%s/oauth2/token", oauth2ServiceHost)

	audienceList := make([]string, 0)
	oauth2ServiceAudience := ocrConfig.Oauth2ServiceAudience
	if oauth2ServiceAudience != "" {
		audienceList = strings.Split(oauth2ServiceAudience, ",")
	}

	filesCli, err := fapi.NewFilesClient(ctx,
		apis.WithEndpoint(ocrConfig.FilesServiceURI),
		apis.WithTokenEndpoint(oauth2ServiceURL),
		apis.WithTokenUsername(service.JwtClientID()),
		apis.WithTokenPassword(ocrConfig.Oauth2ServiceClientSecret),
		apis.WithAudiences(audienceList...))
	if err != nil {
		log.Fatalf("main -- Could not setup files service : %+v", err)
	}

	jwtAudience := ocrConfig.Oauth2JwtVerifyAudience
	if jwtAudience == "" {
		jwtAudience = serviceName
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpcctxtags.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
			service.UnaryAuthInterceptor(jwtAudience, ocrConfig.Oauth2JwtVerifyIssuer),
		)),
		grpc.StreamInterceptor(service.StreamAuthInterceptor(jwtAudience, ocrConfig.Oauth2JwtVerifyIssuer)),
	)

	implementation := &handlers.OCRServer{
		Service:  service,
		FilesCli: filesCli,
	}
	ocr.RegisterOCRServiceServer(grpcServer, implementation)

	grpcServerOpt := frame.GrpcServer(grpcServer)
	serviceOptions = append(serviceOptions, grpcServerOpt)

	ocrSyncQueueHandler := queue.NewOCRQueueHandler(service)
	ocrSyncQueue := frame.RegisterSubscriber(ocrConfig.QueueOcrSyncName, ocrConfig.QueueOcrSync, 2, ocrSyncQueueHandler)
	ocrSyncQueueP := frame.RegisterPublisher(ocrConfig.QueueOcrSyncName, ocrConfig.QueueOcrSync)
	serviceOptions = append(serviceOptions, ocrSyncQueue, ocrSyncQueueP)

	service.Init(serviceOptions...)

	log.WithField("server http port", ocrConfig.HttpServerPort).
		WithField("server grpc port", ocrConfig.GrpcServerPort).
		Info(" Initiating server operations")

	err = service.Run(ctx, "")
	if err != nil {
		log.Printf("main -- Could not run Server : %v", err)
	}

}
