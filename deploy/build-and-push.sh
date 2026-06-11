#!/usr/bin/env bash
set -euo pipefail
# Build devfleet images from the current checkout and push them to GHCR.
# Run this on wowbagger after a git sync (e.g. from a webhook or cron).
#
# Environment:
#   REGISTRY        GHCR registry (default: ghcr.io/flatout-works)
#   TAG             Image tag (default: main)
#   GHCR_USERNAME   GHCR username (required for push)
#   GHCR_TOKEN      GHCR personal access token (required for push)

cd "$(dirname "$0")/.."

: "${REGISTRY:=ghcr.io/flatout-works}"
: "${TAG:=main}"

MCP_IMAGE="${REGISTRY}/devfleet-mcp:${TAG}"
RUNNER_BASE_IMAGE="${REGISTRY}/devfleet-runner-base:${TAG}"
RUNNER_IMAGE="${REGISTRY}/devfleet-runner:${TAG}"

if [[ -n "${GHCR_TOKEN:-}" ]]; then
  echo "${GHCR_TOKEN}" | docker login ghcr.io -u "${GHCR_USERNAME:-gokr}" --password-stdin
fi

echo "=== Building MCP image ==="
docker build -t "$MCP_IMAGE" .

echo "=== Building runner base image ==="
docker build -f runner/Dockerfile.devfleet-base -t "$RUNNER_BASE_IMAGE" .

echo "=== Building runner image ==="
docker build -f runner/Dockerfile.devfleet -t "$RUNNER_IMAGE" .

echo "=== Pushing images ==="
docker push "$MCP_IMAGE"
docker push "$RUNNER_BASE_IMAGE"
docker push "$RUNNER_IMAGE"

echo "=== Done ==="
