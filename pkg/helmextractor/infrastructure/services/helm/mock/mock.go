package mock

import (
	"context"

	"github.com/KathurimaKimathi/extract-helm/pkg/helmextractor/application/dto"
)

type HelmMock struct {
	MockProcessHelmChartExtractionFn func(ctx context.Context, url string) (*dto.Response, error)
}

func NewHelmMock() *HelmMock {
	return &HelmMock{
		MockProcessHelmChartExtractionFn: func(ctx context.Context, url string) (*dto.Response, error) {
			return &dto.Response{
				Images: []dto.ImageDetails{
					{
						Name:   "image1",
						Size:   10,
						Layers: 2,
					},
					{
						Name:   "image1",
						Size:   10,
						Layers: 2,
					},
				},
			}, nil
		},
	}
}

// ProcessHelmChartExtraction mocks the implementation of extracting helm charts
func (m *HelmMock) ProcessHelmChartExtraction(ctx context.Context, url string) (*dto.Response, error) {
	return m.MockProcessHelmChartExtractionFn(ctx, url)
}
