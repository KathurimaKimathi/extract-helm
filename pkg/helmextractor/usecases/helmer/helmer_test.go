package helmer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/KathurimaKimathi/extract-helm/pkg/helmextractor/application/dto"
)

func TestHelm_HelmChartImageExtractor(t *testing.T) {
	type args struct {
		ctx              context.Context
		chartLocationURL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: Extract helm chart",
			args: args{
				ctx:              context.Background(),
				chartLocationURL: "https://test.docker.com",
			},
			wantErr: false,
		},
		{
			name: "Sad case: Unable to extract helm chart",
			args: args{
				ctx:              context.Background(),
				chartLocationURL: "https://test.docker.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, mocks := setupMocks()

			if tt.name == "Sad case: Unable to extract helm chart" {
				mocks.HelmMock.MockProcessHelmChartExtractionFn = func(ctx context.Context, helmChartURL string) (*dto.Response, error) {
					return nil, fmt.Errorf("unable to extract charts")
				}
			}

			_, err := h.HelmChartImageExtractor(tt.args.ctx, tt.args.chartLocationURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("Helm.HelmChartImageExtractor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
