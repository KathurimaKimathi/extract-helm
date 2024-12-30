package rest

import (
	"net/http"

	helmUsecase "github.com/KathurimaKimathi/extract-helm/pkg/helmextractor/usecases/helmer"
	"github.com/gin-gonic/gin"
)

// AcceptedContentTypes is a list of all the accepted content types
var AcceptedContentTypes = []string{"application/json", "application/x-www-form-urlencoded"}

// PresentationHandlersImpl represents the usecase implementation object
type PresentationHandlersImpl struct {
	usecase helmUsecase.Helm
}

// NewPresentationHandlers initializes a new rest handlers usecase
func NewPresentationHandlers(usecase helmUsecase.Helm) *PresentationHandlersImpl {
	return &PresentationHandlersImpl{
		usecase: usecase,
	}
}

// HelmChartImageExtractor handles the GET request to retrieve charts and analyze them
func (h *PresentationHandlersImpl) HelmChartImageExtractor(c *gin.Context) {
	helmChartURL := c.Request.URL.Query().Get("helm-chart-url")

	output, err := h.usecase.HelmChartImageExtractor(c.Request.Context(), helmChartURL)
	if err != nil {
		jsonErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

func jsonErrorResponse(c *gin.Context, statusCode int, err error) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"error": err.Error(),
	})
}
