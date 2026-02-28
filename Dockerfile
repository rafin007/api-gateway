# =============================================================================
# STAGE 1: builder
# =============================================================================
# Use the official Go Alpine image as the build environment.
# Alpine keeps the builder small and reduces layer cache size compared to
# the full Debian-based golang image (~300 MB vs ~800 MB).
# Pin to a specific patch version for fully reproducible builds.
FROM golang:1.25.6-alpine AS builder

# Install only the packages we genuinely need at build time:
#  - git         : required by `go mod download` for VCS-based module fetches
#  - ca-certificates : needed for HTTPS calls to the Go module proxy
RUN apk add --no-cache git ca-certificates

# Set the working directory for all subsequent build instructions.
WORKDIR /build

# ── Dependency caching ────────────────────────────────────────────────────────
# Copy dependency manifests BEFORE the source code so that Docker can cache
# the `go mod download` layer. As long as go.mod / go.sum are unchanged, this
# layer is reused on every subsequent build, dramatically reducing build time.
# The gateway/ prefix matches the subdirectory layout in this repository.
COPY go.mod go.sum ./

# Download and verify all declared module dependencies.
# `go mod verify` checks that on-disk content matches the checksums in go.sum,
# catching any supply-chain tampering early.
RUN go mod download && go mod verify

# ── Source copy ───────────────────────────────────────────────────────────────
# Copy the full application source into the builder. This step is intentionally
# placed after dependency download to avoid busting the module cache layer on
# every code change.
COPY . .

# ── Compile ───────────────────────────────────────────────────────────────────
# Build a fully static, minimal binary.
#
#  CGO_ENABLED=0     : Disable cgo — produces a static binary with zero C
#                      dependencies, compatible with the distroless final image.
#  GOOS / GOARCH     : Explicitly target linux/amd64 regardless of the host OS.
#  -ldflags="-s -w"  : Strip the symbol table (-s) and DWARF debug info (-w),
#                      reducing binary size by ~30% with no runtime impact.
#  -trimpath         : Remove local filesystem paths embedded in the binary,
#                      improving reproducibility and avoiding accidental secret leakage.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /build/api \
    ./cmd/api/main.go


# =============================================================================
# STAGE 2: final (production) image
# =============================================================================
# Google Distroless contains only the application and its minimal runtime
# dependencies — no shell, no package manager, no libc, no debug tools.
# This dramatically reduces:
#   • Attack surface  (no shell means no RCE via shell injection)
#   • CVE exposure    (fewer installed packages = fewer vulnerabilities)
#   • Image size      (~5 MB vs ~250 MB for a full Alpine image)
#
# The ":nonroot" tag pre-configures the image to run as uid/gid 65532
# ("nonroot") — following the principle of least privilege without any
# extra USER instruction.
FROM gcr.io/distroless/static-debian12:nonroot

# Copy the CA certificate bundle from the builder so the application can
# establish outbound TLS connections (e.g., to external APIs or cloud services).
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy ONLY the compiled binary from the builder stage. Nothing else — no
# source code, no build toolchain, no .env files — ever lands in this image.
COPY --from=builder /build/api /app/api

# Set the working directory. In production, all configuration should be
# supplied as environment variables (DB_HOST, SIGNING_SECRET, etc.) rather
# than a .env file. Viper's AutomaticEnv() will pick them up automatically,
# overriding anything that would have been in the .env file.
WORKDIR /app

# Document the port the HTTP server listens on.
# Actual port binding is controlled at `docker run -p` or in your compose/k8s config.
EXPOSE 8080

# Use ENTRYPOINT (not CMD) so the binary is the container's sole process and
# cannot be accidentally overridden. The required APP_MODE env var ("dev" | "prod")
# must be set at runtime, e.g.:
#   docker run -e APP_MODE=prod -e DB_HOST=... <image>
ENTRYPOINT ["/app/api"]