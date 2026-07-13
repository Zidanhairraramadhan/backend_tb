.PHONY: dev build swagger clean

# ── Development ──
# Jalankan server dengan hot-reload menggunakan Air
# Install Air: go install github.com/air-verse/air@latest
dev:
	air

# ── Build ──
# Build binary untuk production
build:
	@echo "🔨 Building MusicLink Backend..."
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o musiclink-server .
	@echo "✅ Build complete: ./musiclink-server"

# ── Build (Windows) ──
build-win:
	@echo "🔨 Building for Windows..."
	go build -o musiclink-backend.exe .
	@echo "✅ Build complete: ./musiclink-backend.exe"

# ── Swagger ──
# Generate/regenerate Swagger documentation
# Install Swag: go install github.com/swaggo/swag/cmd/swag@latest
swagger:
	@echo "📝 Generating Swagger documentation..."
	swag init --parseDependency --parseInternal
	@echo "✅ Swagger docs updated in ./docs/"

# ── Clean ──
# Hapus file build
clean:
	rm -f musiclink-server musiclink-backend.exe
	@echo "🧹 Cleaned build artifacts"

# ── Tidy ──
# Update go.sum dan hapus dependency yang tidak dipakai
tidy:
	go mod tidy
	@echo "✅ go.mod and go.sum updated"

# ── Run (tanpa hot-reload) ──
run:
	go run main.go
