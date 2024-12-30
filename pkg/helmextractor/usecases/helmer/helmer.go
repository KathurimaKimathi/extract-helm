package helmer

import (
	"context"

	"github.com/KathurimaKimathi/extract-helm/pkg/helmextractor/application/dto"
	"github.com/KathurimaKimathi/extract-helm/pkg/helmextractor/infrastructure"
)

type Helm struct {
	infrastructure infrastructure.Infrastructure
}

func NewHelm(infra infrastructure.Infrastructure) *Helm {
	return &Helm{
		infrastructure: infra,
	}
}

// HelmChartImageExtractor contains the chart image extraction logic
func (h *Helm) HelmChartImageExtractor(ctx context.Context, chartLocationURL string) (*dto.Response, error) {
	output, err := h.infrastructure.Helm.ProcessHelmChartExtraction(ctx, chartLocationURL)
	if err != nil {
		return nil, err
	}

	return output, nil
}
