package helm

import (
	"context"
	"net/http"
	"testing"
)

func TestHelm_ProcessHelmChartExtraction(t *testing.T) {
	type args struct {
		ctx context.Context
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: process helm chart",
			args: args{
				ctx: context.Background(),
				url: "https://github.com/helm/examples/releases/download/hello-world-0.1.0/hello-world-0.1.0.tgz",
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to process helm chart",
			args: args{
				ctx: context.Background(),
				url: "https://github.com/",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{}

			h := NewHelmSvc(httpClient)

			_, err := h.ProcessHelmChartExtraction(tt.args.ctx, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Helm.ProcessHelmChartExtraction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
