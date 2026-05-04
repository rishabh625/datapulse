#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${BASH_VERSION:-}" ]]; then
  echo "This script must be run with bash."
  echo "Example: bash ./deploy-cloudrun.sh"
  exit 1
fi

PROJECT_ID="${PROJECT_ID:-$(gcloud config get-value project 2>/dev/null)}"
REGION="${REGION:-asia-south1}"
SERVICE_NAME="${SERVICE_NAME:-datapulse-mcp}"
AR_REPO="${AR_REPO:-datapulse}"
IMAGE_TAG="${IMAGE_TAG:-$(date +%Y%m%d-%H%M%S)}"

if [[ -z "${PROJECT_ID}" ]]; then
  echo "PROJECT_ID is empty. Set it, or run: gcloud config set project <project-id>"
  exit 1
fi

IMAGE_URI="${REGION}-docker.pkg.dev/${PROJECT_ID}/${AR_REPO}/mcpserver:${IMAGE_TAG}"

echo "Project: ${PROJECT_ID}"
echo "Region: ${REGION}"
echo "Service: ${SERVICE_NAME}"
echo "Image: ${IMAGE_URI}"

gcloud services enable \
  run.googleapis.com \
  artifactregistry.googleapis.com \
  cloudbuild.googleapis.com

if ! gcloud artifacts repositories describe "${AR_REPO}" --location="${REGION}" >/dev/null 2>&1; then
  gcloud artifacts repositories create "${AR_REPO}" \
    --repository-format=docker \
    --location="${REGION}"
fi

gcloud builds submit --tag "${IMAGE_URI}" .

gcloud run deploy "${SERVICE_NAME}" \
  --image "${IMAGE_URI}" \
  --region "${REGION}" \
  --platform managed \
  --allow-unauthenticated \
  --port 8080

SERVICE_URL="$(gcloud run services describe "${SERVICE_NAME}" --region "${REGION}" --format='value(status.url)')"
echo "Deployed: ${SERVICE_URL}"
