# Cross-Platform Deployment Guide

This guide covers building and deploying the Card Game API across multiple platforms and architectures.

## Supported Platforms

The Card Game API supports cross-compilation for the following platforms:

| Platform | Architecture | Binary Suffix | Notes |
|----------|-------------|---------------|--------|
| **Linux** | amd64 (x86_64) | `linux-amd64` | Most common server architecture |
| **Linux** | arm64 (aarch64) | `linux-arm64` | ARM servers, Raspberry Pi 4+ |
| **Windows** | amd64 (x86_64) | `windows-amd64.exe` | Windows 10/11, Server 2019+ |
| **Windows** | arm64 (aarch64) | `windows-arm64.exe` | Windows on ARM |
| **macOS** | amd64 (Intel) | `darwin-amd64` | Intel-based Macs |
| **macOS** | arm64 (Apple Silicon) | `darwin-arm64` | M1/M2/M3 Macs |

## Build Methods

### Method 1: Using Makefile (Recommended)

The Makefile provides convenient targets for cross-compilation:

```bash
# Build for all supported platforms
make build-all

# Build for specific platforms
make build-linux          # Linux (amd64 + arm64)
make build-windows         # Windows (amd64 + arm64)  
make build-darwin          # macOS (amd64 + arm64)

# Build for specific architecture
make build-linux-amd64     # Linux x86_64 only
make build-darwin-arm64    # macOS Apple Silicon only
make build-windows-amd64   # Windows x86_64 only

# Clean build artifacts
make clean

# View help
make help
```

Built binaries will be placed in the `dist/` directory:
```
dist/
├── cardgame-api-linux-amd64
├── cardgame-api-linux-arm64
├── cardgame-api-windows-amd64.exe
├── cardgame-api-windows-arm64.exe
├── cardgame-api-darwin-amd64
└── cardgame-api-darwin-arm64
```

### Method 2: Using Build Script

The build script provides more flexibility and detailed output:

```bash
# Make script executable (first time only)
chmod +x scripts/build.sh

# Build all platforms
./scripts/build.sh all

# Build specific platform
./scripts/build.sh linux amd64
./scripts/build.sh darwin arm64
./scripts/build.sh windows amd64

# Clean builds
./scripts/build.sh clean

# List built binaries
./scripts/build.sh list

# Show help
./scripts/build.sh help
```

### Method 3: Manual Cross-Compilation

For advanced users or CI/CD integration:

```bash
# Set build variables
VERSION=$(git describe --tags --always --dirty)
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_HASH=$(git rev-parse --short HEAD)
LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitHash=${COMMIT_HASH} -w -s"

# Build for Linux amd64
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "${LDFLAGS}" -o dist/cardgame-api-linux-amd64 .

# Build for macOS arm64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "${LDFLAGS}" -o dist/cardgame-api-darwin-arm64 .

# Build for Windows amd64
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "${LDFLAGS}" -o dist/cardgame-api-windows-amd64.exe .
```

## Docker Multi-Architecture Builds

### Using Docker Buildx (Recommended)

The project includes a Docker build script for multi-architecture images:

```bash
# Make script executable (first time only)
chmod +x scripts/docker-build.sh

# Build for multiple architectures (default: linux/amd64,linux/arm64)
./scripts/docker-build.sh

# Build and push to registry
./scripts/docker-build.sh --push --registry your-registry.com/username

# Build with custom tags
./scripts/docker-build.sh -t latest -t v1.0.0 --push

# Build for specific platforms only
./scripts/docker-build.sh --platforms linux/amd64

# Show help
./scripts/docker-build.sh --help
```

### Manual Docker Buildx

```bash
# Create and use buildx instance
docker buildx create --name multiarch --driver docker-container --use
docker buildx inspect --bootstrap

# Build for multiple architectures
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --build-arg VERSION=v1.0.0 \
  --build-arg BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ') \
  --build-arg COMMIT_HASH=$(git rev-parse --short HEAD) \
  -t cardgame-api:latest \
  --push \
  .
```

## Deployment Strategies

### 1. Direct Binary Deployment

Download the appropriate binary for your platform:

```bash
# Linux x86_64
wget https://github.com/yourusername/cardgame-api/releases/latest/download/cardgame-api-linux-amd64
chmod +x cardgame-api-linux-amd64
./cardgame-api-linux-amd64

# macOS Apple Silicon
curl -L -o cardgame-api-darwin-arm64 https://github.com/yourusername/cardgame-api/releases/latest/download/cardgame-api-darwin-arm64
chmod +x cardgame-api-darwin-arm64
./cardgame-api-darwin-arm64

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/yourusername/cardgame-api/releases/latest/download/cardgame-api-windows-amd64.exe" -OutFile "cardgame-api.exe"
.\cardgame-api.exe
```

### 2. Container Deployment

#### Docker

```bash
# Pull and run (multi-arch image automatically selects correct architecture)
docker run -d \
  --name cardgame-api \
  -p 8080:8080 \
  -e GIN_MODE=release \
  -e LOG_LEVEL=INFO \
  --restart unless-stopped \
  cardgame-api:latest
```

#### Docker Compose

```yaml
version: '3.8'
services:
  cardgame-api:
    image: cardgame-api:latest
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=release
      - LOG_LEVEL=INFO
      - TRUSTED_PROXIES=127.0.0.1
    restart: unless-stopped
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp
```

#### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cardgame-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: cardgame-api
  template:
    metadata:
      labels:
        app: cardgame-api
    spec:
      containers:
      - name: cardgame-api
        image: cardgame-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: GIN_MODE
          value: "release"
        - name: LOG_LEVEL
          value: "INFO"
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
        securityContext:
          allowPrivilegeEscalation: false
          runAsNonRoot: true
          runAsUser: 65532
          readOnlyRootFilesystem: true
        readinessProbe:
          httpGet:
            path: /hello
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        livenessProbe:
          httpGet:
            path: /hello
            port: 8080
          initialDelaySeconds: 15
          periodSeconds: 20
---
apiVersion: v1
kind: Service
metadata:
  name: cardgame-api-service
spec:
  selector:
    app: cardgame-api
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: LoadBalancer
```

### 3. Cloud Platform Deployment

#### AWS Lambda (with custom runtime)

```bash
# Build for Lambda (Linux amd64)
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap .
zip lambda-deployment.zip bootstrap

# Deploy using AWS CLI
aws lambda create-function \
  --function-name cardgame-api \
  --runtime provided.al2 \
  --role arn:aws:iam::your-account:role/lambda-execution-role \
  --handler bootstrap \
  --zip-file fileb://lambda-deployment.zip
```

#### Google Cloud Run

```bash
# Build and deploy (automatic multi-arch)
gcloud run deploy cardgame-api \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

#### Azure Container Instances

```bash
# Deploy multi-arch container
az container create \
  --resource-group myResourceGroup \
  --name cardgame-api \
  --image cardgame-api:latest \
  --ports 8080 \
  --environment-variables GIN_MODE=release LOG_LEVEL=INFO
```

## Build Information API

The application exposes build information via the `/version` endpoint:

```bash
curl http://localhost:8080/version
```

Response:
```json
{
  "version": "v1.2.3",
  "build_time": "2024-01-15_14:30:45",
  "commit_hash": "a1b2c3d",
  "go_version": "go1.24.4",
  "go_os": "linux",
  "go_arch": "amd64"
}
```

## Performance Considerations

### Binary Size Optimization

The build process uses several flags to minimize binary size:

- `-ldflags="-w -s"`: Strip debug info and symbol table
- `-trimpath`: Remove local path information
- `CGO_ENABLED=0`: Disable CGO for static linking

Typical binary sizes:
- **Linux/macOS**: ~15-20 MB
- **Windows**: ~16-21 MB

### Runtime Performance

- **Memory Usage**: ~20-50 MB at startup
- **CPU Usage**: Low idle, scales with request volume
- **Startup Time**: < 1 second on most platforms
- **Response Time**: < 10ms for most endpoints

## Troubleshooting

### Common Build Issues

1. **Permission Denied**
   ```bash
   chmod +x scripts/build.sh
   chmod +x scripts/docker-build.sh
   ```

2. **Missing Git Information**
   ```bash
   # Ensure you're in a git repository
   git init
   git add .
   git commit -m "Initial commit"
   ```

3. **Docker Buildx Not Available**
   ```bash
   # Install Docker Desktop or enable buildx
   docker buildx install
   ```

### Platform-Specific Issues

#### Windows
- Use PowerShell or WSL2 for best experience
- Ensure Windows version supports your target architecture
- Some antivirus software may flag Go binaries

#### macOS
- Apple Silicon Macs can run both arm64 and amd64 binaries
- Code signing may be required for distribution
- Gatekeeper may block unsigned binaries

#### Linux
- Ensure target system has compatible libc version
- Use static linking (`CGO_ENABLED=0`) for maximum compatibility
- Consider minimal base images for containers

## Security Considerations

### Binary Signing

For production deployments, consider signing binaries:

```bash
# macOS
codesign -s "Developer ID Application: Your Name" cardgame-api-darwin-arm64

# Windows
signtool sign /f certificate.p12 /p password cardgame-api-windows-amd64.exe
```

### Container Security

- Images run as non-root user (UID 65532)
- Read-only root filesystem support
- Minimal attack surface with distroless/alpine base
- Security scanning integrated in CI/CD

### Supply Chain Security

- Dependency verification with `go mod verify`
- Reproducible builds with fixed Go version
- SBOM generation for compliance
- Vulnerability scanning in CI/CD pipeline

## Monitoring and Observability

### Metrics

Prometheus metrics available at `/metrics`:
- HTTP request metrics
- Application-specific metrics
- Runtime metrics (Go runtime)

### Logging

Structured JSON logging with configurable levels:
```bash
# Set log level
export LOG_LEVEL=DEBUG

# Configure log format
export LOG_FORMAT=json
```

### Health Checks

- Liveness: `GET /hello`
- Readiness: `GET /hello`  
- Metrics: `GET /metrics`
- Version: `GET /version`

## Automated Releases

The project includes GitHub Actions workflows for:

1. **Continuous Integration**: Test on multiple Go versions
2. **Cross-Platform Builds**: Build for all supported platforms
3. **Container Builds**: Multi-arch Docker images
4. **Automated Releases**: Tagged releases with binaries
5. **Security Scanning**: Vulnerability detection

Releases are automatically created when tags are pushed:

```bash
git tag v1.2.3
git push origin v1.2.3
```

This triggers builds for all platforms and creates a GitHub release with downloadable binaries.