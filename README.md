# Subfinder Service

[subfinder-ui](https://github.com/NRGGIT/subfinder-ui.git)

A Docker service that provides a REST API for the [subfinder](https://github.com/projectdiscovery/subfinder) subdomain enumeration tool.

## Features

- **Asynchronous Processing**: Submit multiple domain search requests simultaneously
- **Configurable**: Control the depth level, sources, and other subfinder options
- **RESTful API**: Simple HTTP API for submitting jobs and retrieving results
- **Containerized**: Easy deployment with Docker
- **Cloud-Ready**: Kubernetes manifests for cloud deployment

## Project Structure

```
/
├── README.md
├── docker-compose.yml
├── .env.example
├── deploy.sh
├── backend/
│   ├── Dockerfile
│   ├── Makefile
│   ├── cmd/
│   ├── internal/
│   └── pkg/
├── frontend/
│   ├── Dockerfile
│   ├── .env.example
│   ├── app.vue
│   └── ...
└── k8s/
    ├── deploy-k8s.sh
    ├── ingress.yaml
    ├── backend/
    │   ├── deployment.yaml
    │   ├── service.yaml
    │   └── configmap.yaml
    └── frontend/
        ├── deployment.yaml
        ├── service.yaml
        └── configmap.yaml
```

## API Endpoints

### Submit a Job

```
POST /subfinder
```

Request body:

```json
{
  "domain": "example.com",
  "config": {
    "max_depth": 2,
    "include_ips": true,
    "sources": ["virustotal", "crtsh"],
    "timeout": 60,
    "rate_limit": 10,
    "include_wildcards": false,
    "exclude_unresolvable": true,
    "exclude_www": true
  }
}
```

Response:

```json
{
  "job_id": "unique-job-id",
  "status": "queued",
  "estimated_completion_time": "2025-03-04T12:45:00Z"
}
```

### Get Job Status/Results

```
GET /subfinder/{job_id}
```

Response:

```json
{
  "job_id": "unique-job-id",
  "domain": "example.com",
  "config": {
    "max_depth": 2,
    "include_ips": true,
    "sources": ["virustotal", "crtsh"],
    "timeout": 60,
    "rate_limit": 10,
    "include_wildcards": false,
    "exclude_unresolvable": true,
    "exclude_www": true
  },
  "status": "completed",
  "created_at": "2025-03-04T12:30:00Z",
  "started_at": "2025-03-04T12:30:05Z",
  "completed_at": "2025-03-04T12:31:10Z",
  "subdomains": [
    "api.example.com",
    "mail.example.com",
    "subdomain.example.com",
    "a.b.example.com"
  ],
  "stats": {
    "total_found": 42,
    "execution_time": "1m5s",
    "sources_used": ["virustotal", "crtsh"]
  }
}
```

### Get Service Status

```
GET /subfinder/status
```

Response:

```json
{
  "status": "ok",
  "jobs": {
    "total": 10,
    "queued": 2,
    "running": 3,
    "completed": 4,
    "failed": 1
  },
  "time": "2025-03-04T12:35:00Z"
}
```

## Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `max_depth` | Maximum depth level for subdomains | 1 |
| `include_ips` | Include IP addresses in results | false |
| `sources` | List of sources to use | all available |
| `timeout` | Timeout in seconds | 60 |
| `rate_limit` | Rate limit for requests (per second) | 10 |
| `include_wildcards` | Include wildcard subdomains | false |
| `exclude_unresolvable` | Exclude subdomains that don't resolve | false |
| `exclude_www` | Exclude subdomains with www prefix | false |

## Deployment Options

### Local Deployment with Docker Compose

The simplest way to deploy the service locally is to use Docker Compose:

```bash
# Create .env file from example
cp .env.example .env

# Build and start the services
docker compose up -d

# View logs
docker compose logs -f

# Stop the services
docker compose down
```

### Using the Deployment Script

For more flexibility, you can use the provided deployment script:

```bash
# Make the script executable (if needed)
chmod +x deploy.sh

# Deploy locally
./deploy.sh

# Deploy with custom backend URL
./deploy.sh --backend-url http://api.example.com

# Deploy to a specific environment
./deploy.sh --env aws
```

### Kubernetes Deployment

For cloud environments, you can use the Kubernetes manifests:

```bash
# Make the script executable (if needed)
chmod +x k8s/deploy-k8s.sh

# Deploy to Kubernetes
cd k8s
./deploy-k8s.sh

# Deploy with custom settings
./deploy-k8s.sh --backend-image my-registry/subfinder-backend:v1.0.0 \
                --frontend-image my-registry/subfinder-frontend:v1.0.0 \
                --host subfinder.mydomain.com \
                --namespace subfinder
```

## Development

### Backend Development

```bash
# Navigate to backend directory
cd backend

# Build the application
make build

# Run the application
make run

# Run tests
make test
```

### Frontend Development

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build
```

## Example Usage

### Submit a Job

```bash
curl -X POST http://localhost:8080/subfinder \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "example.com",
    "config": {
      "max_depth": 2,
      "include_ips": true,
      "exclude_www": true
    }
  }'
```

### Get Job Results

```bash
curl -X GET http://localhost:8080/subfinder/job-id-here
```

### Get Service Status

```bash
curl -X GET http://localhost:8080/subfinder/status
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
