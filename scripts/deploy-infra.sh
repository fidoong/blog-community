#!/usr/bin/env bash
set -euo pipefail

# Deploy local infrastructure (PostgreSQL, Redis, ES, Kafka, MinIO)
# onto the Orb VM named 'fidoo' (IP: 192.168.139.191).

VM_HOST="192.168.139.191"
VM_USER="${VM_USER:-root}"
COMPOSE_FILE="../docker-compose.infra.yml"
REMOTE_DIR="/opt/blog-infra"

echo "==> Deploying infra to ${VM_USER}@${VM_HOST}"

# Ensure Docker & Compose are available on VM
echo "==> Checking Docker on VM..."
ssh "${VM_USER}@${VM_HOST}" 'docker compose version >/dev/null 2>&1 || (apt-get update && apt-get install -y docker-compose-plugin || true)'

# Upload compose file
echo "==> Uploading ${COMPOSE_FILE}..."
ssh "${VM_USER}@${VM_HOST}" "mkdir -p ${REMOTE_DIR}"
scp "${COMPOSE_FILE}" "${VM_USER}@${VM_HOST}:${REMOTE_DIR}/docker-compose.infra.yml"

# Pull and start services
echo "==> Pulling images and starting services..."
ssh "${VM_USER}@${VM_HOST}" "cd ${REMOTE_DIR} && docker compose -f docker-compose.infra.yml pull && docker compose -f docker-compose.infra.yml up -d"

# Wait a bit for healthchecks
echo "==> Waiting for services to become healthy..."
sleep 5

# Show status
ssh "${VM_USER}@${VM_HOST}" "cd ${REMOTE_DIR} && docker compose -f docker-compose.infra.yml ps"

echo "==> Done. Services should be available at:"
echo "    PostgreSQL:  ${VM_HOST}:5432"
echo "    Redis:       ${VM_HOST}:6379"
echo "    Elasticsearch: ${VM_HOST}:9200"
echo "    Kafka:       ${VM_HOST}:9092"
echo "    MinIO API:   ${VM_HOST}:9000"
echo "    MinIO Console: ${VM_HOST}:9001"
