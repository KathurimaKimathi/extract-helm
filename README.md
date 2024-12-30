# Helm Chart Image Analyzer API

## Overview
The **Helm Chart Image Analyzer API** is a simple REST API that accepts a Helm chart URL, extracts the container images defined in the chart's `values.yaml` file, and retrieves details about these images, such as their size and number of layers. This tool is particularly useful for auditing container images used in Helm charts.

## Features
- Parses Helm charts to identify container images.
- Retrieves image metadata, including:
  - Image name
  - Size (in bytes)
  - Number of layers
- Provides a RESTful API endpoint for integration with other tools.

## Prerequisites
- Go (1.19 or later)
- Docker installed and running on the host machine
- Internet access to pull images from Docker registries

## Installation
### Method 1
1. Clone this repository:
   ```bash
   git clone git@github.com:KathurimaKimathi/extract-helm.git
   cd extract-helm
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the application
   ```bash
   go run server.go
   ```
### Method 2
1. Clone this repository:
   ```bash
   git clone git@github.com:KathurimaKimathi/extract-helm.git
   cd extract-helm
   ```
2. Run the docker image
   ```bash
   docker run --name helmextractor -d --network host --restart unless-stopped --env-file /path/to/env-file/file.sh -v /var/run/docker.sock:/var/run/docker.sock -v $(pwd):/app helmextractor
   ```

## API Endpoint

### `GET /analyze-helm-chart`

#### Request
- **URL Parameter**:
  - `url` (string): URL to the `.tgz` file of a Helm chart.
- **Headers**:
  - `Content-Type`: `application/json`

#### Example Request
```bash
curl -X POST "http://127.0.0.1:[YOUR-PORT]/api/v1/analyze-helm-chart?helm-chart-url=https://github.com/helm/examples/releases/download/hello-world-0.1.0/hello-world-0.1.0.tgz"
```

#### Response
- **Status**: 200 OK
- **Content-Type**: `application/json`

##### Response Body
```json
{
  "images": [
    {
      "image": "nginx:latest",
      "size": 22284874,
      "layers": 3
    }
  ]
}
```

#### Error Responses
- **400 Bad Request**: Invalid or missing parameters.
- **500 Internal Server Error**: `"error": "failed to create gzip reader: gzip: invalid header"` since no valid zip file was found.

## How It Works
1. **Download and Extract**:
   - Downloads the Helm chart from the provided URL.
   - Extracts the contents of the `.tgz` archive to a temporary directory.
2. **Parse `values.yaml`**:
   - Identifies the `values.yaml` file within the extracted chart directory.
   - Parses the file to extract container image information.
3. **Fetch Image Details**:
   - Pulls the images using the Docker client.
   - Inspects each image to determine its size and number of layers.
4. **Return Response**:
   - Sends the image details as a JSON response.

## License
This project is licensed under the MIT License. See the `LICENSE` file for details.

## Contact
For questions or support, please contact Kathurima at [kathurimakimathi415@gmail.com].

