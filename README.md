# AWS S3 Presigned URL Service

A Go-based REST API service that generates presigned URLs for AWS S3 operations, allowing secure temporary access to S3 objects without exposing AWS credentials.

## Table of Contents

- [Quick Start](#quick-start)
- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Configuration](#configuration)
- [Deployment](#deployment)
  - [Local Development](#local-development)
  - [Docker Deployment](#docker-deployment)
  - [Production Deployment](#production-deployment)
- [API Reference](#api-reference)
  - [Endpoints](#endpoints)
  - [Response Schemas](#response-schemas)
  - [Error Handling](#error-handling)
- [Usage Examples](#usage-examples)
- [Testing](#testing)
- [Project Structure](#project-structure)
- [Dependencies](#dependencies)
- [Security Considerations](#security-considerations)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)

## Quick Start

Get up and running in 5 minutes:

1. **Clone the repository:**
```bash
git clone <repository-url>
cd presigned_url
```

2. **Set up environment:**
```bash
cp .env.example .env
# Edit .env with your AWS credentials:
# AWS_REGION=us-east-1
# AWS_ACCESS_KEY_ID=your_access_key
# AWS_SECRET_ACCESS_KEY=your_secret_key
```

3. **Run the service:**
```bash
go run cmd/main.go
```

4. **Test it works:**
```bash
curl "http://localhost:8080/health"
curl "http://localhost:8080/s3/presigned-url?bucket=your-bucket&key=test-file.txt"
```

**Or use Docker:**
```bash
docker-compose up -d
```

## Overview

This service provides a simple HTTP API to generate presigned URLs for AWS S3 operations. It supports both PUT (upload) and GET (download) operations with configurable expiration times.

## Features

- **Presigned PUT URLs**: Generate secure URLs for uploading files to S3
- **Presigned GET URLs**: Generate secure URLs for downloading files from S3
- **Transaction Tracking**: Each request generates a unique transaction ID for tracking
- **Environment-based Configuration**: AWS credentials and configuration via environment variables
- **RESTful API**: Clean REST endpoints using Gin web framework
- **Error Handling**: Comprehensive error responses with transaction details

## Architecture

The service follows a modular architecture pattern:

```
├── cmd/                    # Application entry point
├── infrastructure/         # External service integrations
│   └── aws/               # AWS S3 configuration and operations
├── modules/               # Business logic modules
│   └── presigned/         # Presigned URL module
│       └── presigned/
│           ├── controller/ # HTTP request handlers
│           └── route/     # API route definitions
└── shared/                # Shared utilities
```

## Installation

### Prerequisites

- Go 1.24 or higher
- AWS account with S3 access
- AWS credentials configured

### Steps

1. Clone the repository:
```bash
git clone <repository-url>
cd presigned_url
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o main cmd/main.go
```

## Configuration

Create a `.env` file in the root directory with the following variables:

```env
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_access_key_id
AWS_SECRET_ACCESS_KEY=your_secret_access_key
AWS_SESSION_TOKEN=your_session_token  # Optional, for temporary credentials
```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `AWS_REGION` | AWS region for S3 operations | Yes |
| `AWS_ACCESS_KEY_ID` | AWS access key ID | Yes |
| `AWS_SECRET_ACCESS_KEY` | AWS secret access key | Yes |
| `AWS_SESSION_TOKEN` | AWS session token (for temporary credentials) | No |

## Deployment

### Local Development

1. **Clone and setup the project:**
```bash
git clone <repository-url>
cd presigned_url
go mod download
```

2. **Create environment file:**
```bash
cp .env.example .env
# Edit .env with your AWS credentials
```

3. **Run the application:**
```bash
# Development mode
go run cmd/main.go

# Or build and run
go build -o main cmd/main.go
./main
```

4. **Test the service:**
```bash
curl "http://localhost:8080/health"
curl "http://localhost:8080/s3/presigned-url?bucket=test-bucket&key=test-file.txt"
```

### Docker Deployment

#### Using Docker directly:

1. **Build the Docker image:**
```bash
docker build -t presigned-url-service .
```

2. **Run with Docker:**
```bash
docker run -d \
  --name presigned-url-service \
  -p 8080:8080 \
  -e AWS_REGION=us-east-1 \
  -e AWS_ACCESS_KEY_ID=your_key \
  -e AWS_SECRET_ACCESS_KEY=your_secret \
  presigned-url-service
```

#### Using Docker Compose:

1. **Set up environment:**
```bash
cp .env.example .env
# Edit .env with your credentials
```

2. **Run with docker-compose:**
```bash
docker-compose up -d
```

The docker-compose.yml includes:
- Health checks
- Restart policies
- Environment variable configuration
- Network isolation

### Production Deployment

#### AWS ECS/Fargate

1. **Build and push to ECR:**
```bash
# Login to ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin your-account.dkr.ecr.us-east-1.amazonaws.com

# Build and push
docker build -t presigned-url-service .
docker tag presigned-url-service:latest your-account.dkr.ecr.us-east-1.amazonaws.com/presigned-url-service:latest
docker push your-account.dkr.ecr.us-east-1.amazonaws.com/presigned-url-service:latest
```

2. **Create ECS task definition:**
```json
{
  "family": "presigned-url-service",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::account:role/presigned-url-task-role",
  "containerDefinitions": [
    {
      "name": "presigned-url-service",
      "image": "your-account.dkr.ecr.region.amazonaws.com/presigned-url-service:latest",
      "portMappings": [{"containerPort": 8080, "protocol": "tcp"}],
      "environment": [{"name": "AWS_REGION", "value": "us-east-1"}],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/presigned-url-service",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

#### Kubernetes

1. **Create deployment manifests:**

**Deployment:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: presigned-url-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: presigned-url-service
  template:
    metadata:
      labels:
        app: presigned-url-service
    spec:
      containers:
      - name: presigned-url-service
        image: your-registry/presigned-url-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: AWS_REGION
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: aws-region
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
```

**Service:**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: presigned-url-service
spec:
  selector:
    app: presigned-url-service
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: LoadBalancer
```

2. **Deploy to Kubernetes:**
```bash
kubectl apply -f k8s/
kubectl get pods -l app=presigned-url-service
```

#### Production Best Practices

- **Use HTTPS**: Always enable SSL/TLS in production
- **IAM Roles**: Use IAM roles instead of access keys when possible
- **Load Balancing**: Place behind ALB/NLB for high availability
- **Monitoring**: Implement logging, metrics, and health checks
- **Auto-scaling**: Configure horizontal pod autoscaling
- **Security**: Use VPC, security groups, and network policies

## API Reference

### Base Information

- **Base URL**: `http://localhost:8080`
- **Protocol**: HTTP/HTTPS
- **Content-Type**: `application/json`
- **Response Format**: JSON
- **Authentication**: AWS credentials via environment variables

### Endpoints

#### 1. Health Check

**Endpoint:** `GET /health`

**Description:** Check service health status.

**Response:**
```json
{
  "status": "healthy",
  "service": "AWS S3 Presigned URL Service"
}
```

#### 2. Generate Presigned Upload URL

**Endpoint:** `GET /s3/presigned-url`

**Description:** Generates a presigned URL for uploading files to S3 (15-minute expiration).

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `bucket` | string | Yes | S3 bucket name where the file will be uploaded |
| `key` | string | Yes | S3 object key (path/filename) for the file |

**Example Request:**
```bash
curl "http://localhost:8080/s3/presigned-url?bucket=my-storage-bucket&key=documents/report.pdf"
```

**Success Response (200 OK):**
```json
{
  "sucesso": 0,
  "message": "Presigned URL generated successfully",
  "url": "https://my-storage-bucket.s3.us-east-1.amazonaws.com/documents/report.pdf?AWSAccessKeyId=AKIA...&Expires=1643723400&Signature=...",
  "transaction": {
    "transactionID": "550e8400-e29b-41d4-a716-446655440000",
    "timestamp": "2025-07-28 15:30:45"
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": 1,
  "message": "bucket and key parameters are required",
  "transaction": {
    "transactionID": "550e8400-e29b-41d4-a716-446655440000",
    "timestamp": "2025-07-28 15:30:45"
  }
}
```

#### 3. Generate Presigned Download URL

**Endpoint:** `GET /s3/presigned-get`

**Description:** Generates a presigned URL for downloading files from S3 (15-minute expiration).

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `bucket` | string | Yes | S3 bucket name where the file is stored |
| `key` | string | Yes | S3 object key (path/filename) of the file |

**Example Request:**
```bash
curl "http://localhost:8080/s3/presigned-get?bucket=my-storage-bucket&key=documents/report.pdf"
```

**Success Response (200 OK):**
```json
{
  "sucesso": 0,
  "message": "Presigned GET URL generated successfully",
  "url": "https://my-storage-bucket.s3.us-east-1.amazonaws.com/documents/report.pdf?AWSAccessKeyId=AKIA...&Expires=1643723400&Signature=...",
  "transaction": {
    "transactionID": "550e8400-e29b-41d4-a716-446655440000",
    "timestamp": "2025-07-28 15:30:45"
  }
}
```

### Response Schemas

#### Success Response Schema
```json
{
  "sucesso": {
    "type": "integer",
    "description": "Success indicator (0 = success)",
    "example": 0
  },
  "message": {
    "type": "string",
    "description": "Success message"
  },
  "url": {
    "type": "string",
    "description": "Generated presigned URL"
  },
  "transaction": {
    "type": "object",
    "properties": {
      "transactionID": {
        "type": "string",
        "description": "Unique transaction identifier (UUID)"
      },
      "timestamp": {
        "type": "string",
        "description": "Transaction timestamp (YYYY-MM-DD HH:MM:SS)"
      }
    }
  }
}
```

#### Error Response Schema
```json
{
  "error": {
    "type": "integer",
    "description": "Error indicator (1 = error)"
  },
  "message": {
    "type": "string",
    "description": "Error message"
  },
  "transaction": {
    "type": "object",
    "description": "Transaction details (same as success)"
  }
}
```

### Error Handling

| HTTP Status | Error Code | Description | Resolution |
|-------------|------------|-------------|------------|
| 400 | 1 | Missing required parameters | Provide both `bucket` and `key` parameters |
| 500 | N/A | AWS operation failed | Check AWS credentials and permissions |
| 500 | N/A | Session creation failed | Verify AWS configuration and network connectivity |
  "sucesso": 0,
  "message": "Presigned GET URL generated successfully",
  "url": "https://bucket-name.s3.region.amazonaws.com/file-key?AWSAccessKeyId=...",
  "transaction": {
    "transactionID": "uuid-string",
    "timestamp": "2025-07-28 15:30:45"
  }
}
```

## Usage Examples

### Starting the Server

```bash
# Development mode
go run cmd/main.go

# Production mode (using built binary)
./main
```

The server will start on port 8080.

### Generating Upload URL

```bash
curl -X GET "http://localhost:8080/s3/presigned-url?bucket=my-bucket&key=uploads/file.jpg"
```

### Generating Download URL

```bash
curl -X GET "http://localhost:8080/s3/presigned-get?bucket=my-bucket&key=uploads/file.jpg"
```

### Using the Presigned URL

Once you receive a presigned URL, you can use it directly:

**For Upload (PUT):**
```bash
curl -X PUT "PRESIGNED_URL" \
  -H "Content-Type: image/jpeg" \
  --data-binary @local-file.jpg
```

**For Download (GET):**
```bash
curl "PRESIGNED_URL" -o downloaded-file.jpg
```

## Error Handling

### Error Response Format

```json
{
  "error": 1,
  "message": "Error description",
  "transaction": {
    "transactionID": "uuid-string",
    "timestamp": "2025-07-28 15:30:45"
  }
}
```

### Common Error Codes

| HTTP Status | Error Code | Description |
|-------------|------------|-------------|
| 400 | 1 | Bad Request - Missing required parameters |
| 500 | N/A | Internal Server Error - AWS operation failed |

### Response Status Codes

- `sucesso: 0` - Success
- `error: 1` - Error occurred

## Usage Examples

### Starting the Server

```bash
# Development mode
go run cmd/main.go

# Production mode (using built binary)
go build -o main cmd/main.go && ./main
```

The server will start on port 8080.

### Using the API

#### Generate Upload URL
```bash
curl -X GET "http://localhost:8080/s3/presigned-url?bucket=my-bucket&key=uploads/file.jpg"
```

#### Generate Download URL
```bash
curl -X GET "http://localhost:8080/s3/presigned-get?bucket=my-bucket&key=uploads/file.jpg"
```

### Using Presigned URLs

#### Upload a file:
```bash
# Get the presigned upload URL first
UPLOAD_URL=$(curl -s "http://localhost:8080/s3/presigned-url?bucket=my-bucket&key=uploads/document.pdf" | jq -r '.url')

# Upload the file using the presigned URL
curl -X PUT "$UPLOAD_URL" \
  -H "Content-Type: application/pdf" \
  --data-binary @document.pdf
```

#### Download a file:
```bash
# Get the presigned download URL
DOWNLOAD_URL=$(curl -s "http://localhost:8080/s3/presigned-get?bucket=my-bucket&key=uploads/document.pdf" | jq -r '.url')

# Download the file
curl "$DOWNLOAD_URL" -o downloaded-document.pdf
```

### Important Notes

- **Upload Requirements**: Use HTTP PUT method with appropriate Content-Type header
- **Download Requirements**: Use HTTP GET method
- **Expiration**: All URLs expire after 15 minutes
- **File Size**: Limited by S3 maximum object size (5TB)

## Testing

### API Testing

#### Using curl:
```bash
# Test health endpoint
curl "http://localhost:8080/health"

# Test upload URL generation
curl "http://localhost:8080/s3/presigned-url?bucket=test-bucket&key=test-file.txt"

# Test download URL generation
curl "http://localhost:8080/s3/presigned-get?bucket=test-bucket&key=test-file.txt"
```

#### Using Postman:
1. Create GET requests to the endpoints
2. Add query parameters for `bucket` and `key`
3. Verify response format
4. Test returned presigned URLs

### Unit Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./modules/presigned/presigned/controller/
```

## Project Structure

```
presigned_url/
├── cmd/
│   └── main.go                 # Application entry point
├── infrastructure/
│   └── aws/
│       └── config.go           # AWS S3 configuration and operations
├── modules/
│   └── presigned/
│       └── presigned/
│           ├── controller/
│           │   └── controller.go # HTTP request handlers
│           └── route/
│               └── router.go     # API route definitions
├── shared/
│   └── uuid.go                 # UUID utility functions
├── go.mod                      # Go module dependencies
├── go.sum                      # Dependency checksums
├── Dockerfile                  # Docker configuration
├── docker-compose.yml          # Docker Compose configuration
├── .env.example                # Environment variables template
└── README.md                   # This documentation
```

### Key Components

1. **Main Application** (`cmd/main.go`): Entry point that starts the Gin server
2. **AWS Infrastructure** (`infrastructure/aws/config.go`): AWS S3 client configuration and presigned URL generation
3. **Controllers** (`modules/presigned/presigned/controller/controller.go`): HTTP request handlers
4. **Routes** (`modules/presigned/presigned/route/router.go`): API endpoint definitions
5. **Shared Utilities** (`shared/uuid.go`): Common utility functions

## Dependencies

- **[Gin](https://github.com/gin-gonic/gin)**: HTTP web framework
- **[AWS SDK for Go](https://github.com/aws/aws-sdk-go)**: AWS service integration
- **[Google UUID](https://github.com/google/uuid)**: UUID generation
- **[GoDotEnv](https://github.com/joho/godotenv)**: Environment variable loading

## Security Considerations

### Application Security
1. **URL Expiration**: Presigned URLs expire after 15 minutes for security
2. **Environment Variables**: Store AWS credentials securely using environment variables
3. **HTTPS**: Use HTTPS in production to encrypt API communications
4. **Input Validation**: Validate bucket names and object keys
5. **Rate Limiting**: Consider implementing rate limiting for production

### AWS Security
1. **IAM Permissions**: Use principle of least privilege
2. **IAM Roles**: Use IAM roles instead of access keys when possible
3. **VPC**: Deploy in VPC for network isolation
4. **CloudTrail**: Enable logging for audit trails
5. **S3 Bucket Policies**: Configure appropriate bucket permissions

### Required IAM Permissions
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject"
      ],
      "Resource": "arn:aws:s3:::your-bucket-name/*"
    }
  ]
}
```

## Troubleshooting

### Common Issues

#### 1. AWS Credentials Not Found
**Error:** `NoCredentialProviders` or authentication failures

**Solutions:**
- Verify environment variables are set correctly
- Check IAM permissions for S3 operations
- Ensure AWS region is configured properly
- For ECS/EKS, verify task/pod has proper IAM role

#### 2. Presigned URL Generation Fails
**Error:** `failed to generate presigned URL`

**Solutions:**
- Check S3 bucket exists and is accessible
- Verify bucket is in the same region as configured
- Ensure IAM permissions include `s3:PutObject` or `s3:GetObject`
- Check network connectivity to AWS

#### 3. Service Not Responding
**Error:** Connection refused or timeouts

**Solutions:**
- Verify service is running: `curl http://localhost:8080/health`
- Check port 8080 is accessible
- Review application logs for errors
- Ensure no firewall blocking traffic

#### 4. Upload/Download Failures
**Error:** Presigned URL returns 403 or other errors

**Solutions:**
- Verify URL hasn't expired (15-minute limit)
- Check Content-Type header for uploads
- Ensure using correct HTTP method (PUT for upload, GET for download)
- Verify file permissions and size limits

### Debug Logging

Enable debug mode:
```bash
export GIN_MODE=debug
go run cmd/main.go
```

Check application logs:
```bash
# Docker
docker logs presigned-url-service

# Kubernetes
kubectl logs -l app=presigned-url-service

# Local development
# Logs will appear in terminal where service is running
```

### Performance Considerations

1. **Connection Pooling**: AWS SDK handles connection pooling automatically
2. **Timeouts**: Configure appropriate timeouts for AWS operations
3. **Memory Usage**: Monitor memory usage under load
4. **Concurrent Requests**: Service handles concurrent requests via Gin
5. **Rate Limits**: Consider AWS API rate limits for high-volume usage

## Contributing

We welcome contributions!

### Quick Contributing Guide

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes** and add tests
4. **Commit your changes**: `git commit -m 'Add amazing feature'`
5. **Push to the branch**: `git push origin feature/amazing-feature`
6. **Open a Pull Request**

### Development Standards

- Follow Go best practices and formatting (`go fmt`)
- Add tests for new functionality
- Update documentation for any API changes
- Use meaningful commit messages
- Ensure all tests pass before submitting PR


---

**Need Help?** 
- 📖 Check the troubleshooting section above
- 🐛 Report bugs by creating an issue
- 💡 Request features by creating an issue
- 💬 Ask questions in discussions

