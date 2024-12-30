package presentation

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/KathurimaKimathi/extract-helm/pkg/helmextractor/infrastructure"
	"github.com/KathurimaKimathi/extract-helm/pkg/helmextractor/infrastructure/services/helm"
	"github.com/KathurimaKimathi/extract-helm/pkg/helmextractor/presentation/rest"
	helmUsecase "github.com/KathurimaKimathi/extract-helm/pkg/helmextractor/usecases/helmer"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// AllowedHeaders is a list of CORS allowed headers service
var AllowedHeaders = []string{
	"Accept",
	"Accept-Charset",
	"Accept-Language",
	"Accept-Encoding",
	"Origin",
	"Host",
	"User-Agent",
	"Content-Length",
	"Content-Type",
}

// PrepareServer sets up the HTTP server
func PrepareServer(ctx context.Context, port int) {
	// start up the router
	ginEngine := gin.Default()

	err := StartGinRouter(ctx, ginEngine)
	if err != nil {
		msg := fmt.Sprintf("Could not start the router: %v", err)
		log.Panic(msg)
	}

	addr := fmt.Sprintf(":%v", port)

	if err := ginEngine.Run(addr); err != nil {
		log.Panic(err)
	}
}

// StartGinRouter sets up the GIN router
func StartGinRouter(ctx context.Context, engine *gin.Engine) error {
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowWildcard:    true,
		AllowMethods:     []string{http.MethodPut, http.MethodPatch, http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length", "Link"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowWebSockets:  true,
		AllowAllOrigins:  true,
	}))

	serviceHelm := helm.NewHelmSvc(&http.Client{
		Timeout: 15 * time.Second,
	})

	infrastructure := infrastructure.NewInfrastructureInteractor(serviceHelm)

	usecases := helmUsecase.NewHelm(infrastructure)

	handlers := rest.NewPresentationHandlers(*usecases)

	v1 := engine.Group("/api/v1/")

	analyzeChart := v1.Group("analyze-helm-chart")
	{
		analyzeChart.GET("", handlers.HelmChartImageExtractor)
	}

	return nil
}
