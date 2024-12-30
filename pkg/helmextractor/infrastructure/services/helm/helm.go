package helm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/KathurimaKimathi/extract-helm/pkg/helmextractor/application/dto"
	pullOpts "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/mholt/archiver/v3"
	"gopkg.in/yaml.v2"
)

type Helm struct {
	HttpClient *http.Client
}

func NewHelmSvc(httpClient *http.Client) *Helm {
	return &Helm{
		HttpClient: httpClient,
	}
}

// ProcessHelmChartExtraction downloads the chart zip file and extracts in toa temporary directory.
// Then retrieves the container images and checks their details
func (h *Helm) ProcessHelmChartExtraction(ctx context.Context, url string) (*dto.Response, error) {
	chartPath, err := h.downloadAndExtractHelmChart(url)
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(chartPath)

	imageList, err := getContainerImagesFromHelmChart(chartPath)
	if err != nil {
		return nil, err
	}

	imageDetails, err := fetchImageDetails(ctx, imageList)
	if err != nil {
		return nil, err
	}

	return &dto.Response{
		Images: imageDetails,
	}, nil
}

// downloadHelmChart is used to download charts from a given URL, creates a temporary directory for extraction,
// saves the downloaded file as a temporary .tgz file and then extracts the .tgz file into the temporary directory
func (h *Helm) downloadAndExtractHelmChart(url string) (string, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to download Helm chart: %w", err)
	}

	resp, err := h.HttpClient.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	tempDir, err := os.MkdirTemp("", "helm-chart-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	tempFile := filepath.Join(tempDir, "chart.tgz")

	outFile, err := os.Create(tempFile)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}

	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return "", fmt.Errorf("failed to write Helm chart to temp file: %w", err)
	}

	extractor := archiver.NewTarGz()
	if err := extractor.Unarchive(tempFile, tempDir); err != nil {
		return "", fmt.Errorf("failed to extract Helm chart: %w", err)
	}

	return tempDir, nil
}

// getContainerImagesFromHelmChart is used to search for values in the chart directory, locate values.yml,
// access the image and return the corresponding image
func getContainerImagesFromHelmChart(chartDir string) ([]string, error) {
	var valuesFile string

	// search for values.yaml in the directory
	err := filepath.Walk(chartDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, "values.yaml") {
			valuesFile = path
			return io.EOF // halt the search once we find the file
		}

		return nil
	})
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("error searching for values.yaml: %w", err)
	}

	if valuesFile == "" {
		return nil, fmt.Errorf("values.yaml not found in chart directory")
	}

	valuesData, err := os.ReadFile(valuesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read values.yaml: %w", err)
	}

	// parse the values
	var values Values
	if err := yaml.Unmarshal(valuesData, &values); err != nil {
		return nil, fmt.Errorf("failed to parse values.yaml: %w", err)
	}

	image := values.Image.Repository
	if values.Image.Tag != "" {
		image += ":" + values.Image.Tag
	} else {
		image += ":latest"
	}

	return []string{image}, nil
}

// fetchImageDetails is used to fetch the details of a given image
func fetchImageDetails(ctx context.Context, images []string) ([]dto.ImageDetails, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	var (
		details      []dto.ImageDetails
		detailsMutex sync.Mutex
		wg           sync.WaitGroup
		errors       []error
		errorsMutex  sync.Mutex
	)

	for _, image := range images {
		image := image

		wg.Add(1)

		go func() {
			defer wg.Done()

			reader, err := cli.ImagePull(ctx, image, pullOpts.PullOptions{})
			if err != nil {
				errorsMutex.Lock()
				errors = append(errors, fmt.Errorf("failed to pull image %s: %w", image, err))
				errorsMutex.Unlock()

				return
			}

			_, err = io.Copy(io.Discard, reader) // Consume the output to avoid blocking
			if err != nil {
				errorsMutex.Lock()
				errors = append(errors, fmt.Errorf("failed to consume output from image %s: %w", image, err))
				errorsMutex.Unlock()

				return
			}

			imageInfo, _, err := cli.ImageInspectWithRaw(ctx, image)
			if err != nil {
				errorsMutex.Lock()
				errors = append(errors, fmt.Errorf("failed to inspect image %s: %w", image, err))
				errorsMutex.Unlock()

				return
			}

			detailsMutex.Lock()
			details = append(details, dto.ImageDetails{
				Name:   image,
				Size:   imageInfo.Size,
				Layers: len(imageInfo.RootFS.Layers),
			})
			detailsMutex.Unlock()
		}()
	}

	wg.Wait()

	if len(errors) > 0 {
		return nil, fmt.Errorf("one or more errors occurred: %v", errors)
	}

	return details, nil
}
