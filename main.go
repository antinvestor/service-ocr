package main

import (
	"context"
	"net/http"

	"buf.build/gen/go/antinvestor/files/connectrpc/go/files/v1/filesv1connect"
	"buf.build/gen/go/antinvestor/ocr/connectrpc/go/ocr/v1/ocrv1connect"
	"connectrpc.com/connect"
	apis "github.com/antinvestor/apis/go/common"
	"github.com/antinvestor/apis/go/files"
	"github.com/antinvestor/service-ocr/config"
	"github.com/antinvestor/service-ocr/service/business"
	"github.com/antinvestor/service-ocr/service/handlers"
	"github.com/antinvestor/service-ocr/service/models"
	"github.com/antinvestor/service-ocr/service/queue"
	"github.com/antinvestor/service-ocr/service/repository"
	"github.com/pitabwire/frame"
	fconfig "github.com/pitabwire/frame/config"
	"github.com/pitabwire/frame/datastore"
	"github.com/pitabwire/frame/security"
	connectInterceptors "github.com/pitabwire/frame/security/interceptors/connect"
	"github.com/pitabwire/frame/security/openid"
	"github.com/pitabwire/util"
)

func main() {
	tmpCtx := context.Background()

	cfg, err := fconfig.LoadWithOIDC[config.OcrConfig](tmpCtx)
	if err != nil {
		util.Log(tmpCtx).WithError(err).Error("could not process configs")
		return
	}

	if cfg.Name() == "" {
		cfg.ServiceName = "service_ocr"
	}

	ctx, svc := frame.NewServiceWithContext(
		tmpCtx,
		frame.WithConfig(&cfg),
		frame.WithRegisterServerOauth2Client(),
		frame.WithDatastore(),
	)
	defer svc.Stop(ctx)

	log := util.Log(ctx)
	dbManager := svc.DatastoreManager()

	if cfg.DoDatabaseMigrate() {
		dbPool := dbManager.GetPool(ctx, datastore.DefaultMigrationPoolName)
		if dbPool == nil {
			log.Fatal("database pool is nil - check DATABASE_URL environment variable")
			return
		}
		err = dbManager.Migrate(ctx, dbPool, cfg.GetDatabaseMigrationPath(),
			&models.OcrLog{})
		if err != nil {
			log.WithError(err).Fatal("could not migrate successfully")
		}
		return
	}

	sm := svc.SecurityManager()

	audienceList := cfg.GetOauth2ServiceAudience()

	filesCli, err := setupFilesClient(ctx, sm, cfg, audienceList)
	if err != nil {
		log.WithError(err).Fatal("could not setup files client")
	}

	dbPool := dbManager.GetPool(ctx, datastore.DefaultPoolName)
	if dbPool == nil {
		log.Fatal("database pool is nil - check DATABASE_URL environment variable")
		return
	}

	queueMan := svc.QueueManager()
	ocrRepo := repository.NewOcrRepository(dbPool)
	ocrBusiness := business.NewOcrBusiness(ctx, svc, filesCli, ocrRepo, queueMan)

	connectHandler := setupConnectServer(ctx, sm, ocrBusiness)

	ocrSyncQueueHandler := queue.NewOCRQueueHandler(ocrRepo)

	serviceOptions := []frame.Option{
		frame.WithHTTPHandler(connectHandler),
		frame.WithRegisterSubscriber(cfg.QueueOcrSyncName, cfg.QueueOcrSync, ocrSyncQueueHandler),
		frame.WithRegisterPublisher(cfg.QueueOcrSyncName, cfg.QueueOcrSync),
	}

	svc.Init(ctx, serviceOptions...)

	serverPort := cfg.Port()
	if serverPort == "" {
		serverPort = ":7012"
	}

	log.With("port", serverPort).Info("initiating server operations")
	err = svc.Run(ctx, serverPort)
	if err != nil {
		log.WithError(err).Error("could not run Server")
	}
}

func setupFilesClient(
	ctx context.Context,
	clHolder security.InternalOauth2ClientHolder,
	cfg config.OcrConfig,
	audiences []string,
) (filesv1connect.FilesServiceClient, error) {
	return files.NewClient(ctx,
		apis.WithEndpoint(cfg.FilesServiceURI),
		apis.WithTokenEndpoint(cfg.GetOauth2TokenEndpoint()),
		apis.WithTokenUsername(clHolder.JwtClientID()),
		apis.WithTokenPassword(clHolder.JwtClientSecret()),
		apis.WithScopes(openid.ConstSystemScopeInternal),
		apis.WithAudiences(audiences...))
}

func setupConnectServer(ctx context.Context, sm security.Manager, ocrBusiness business.OCRBusiness) http.Handler {
	implementation := handlers.NewOCRServer(ocrBusiness)

	defaultInterceptorList, err := connectInterceptors.DefaultList(ctx, sm.GetAuthenticator(ctx))
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("could not create default interceptors")
	}

	_, serverHandler := ocrv1connect.NewOCRServiceHandler(
		implementation, connect.WithInterceptors(defaultInterceptorList...))

	return serverHandler
}
