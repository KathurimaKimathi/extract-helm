package infrastructure

import (
	"context"

	"github.com/KathurimaKimathi/extract-helm/pkg/helmextractor/application/dto"
)

type ServiceHelm interface {
	ProcessHelmChartExtraction(ctx context.Context, url string) (*dto.Response, error)
}

// Infrastructure ...
type Infrastructure struct {
	Helm ServiceHelm
}

// NewInfrastructureInteractor sets up a new constructor for the any project associated services, data stores etc.
func NewInfrastructureInteractor(
	helm ServiceHelm,
) Infrastructure {
	return Infrastructure{
		Helm: helm,
	}
}
