#!/usr/bin/env bash
# check_api_shape.sh — Ensure proto public services and RPCs are indexed in docs.
#
# Usage:
#   cd packages/proto
#   ./scripts/check_api_shape.sh
#
# Exits 0 if all proto services and RPCs are listed in docs/README.md.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROTO_DIR="${SCRIPT_DIR}/.."
DOCS_FILE="${SCRIPT_DIR}/../../../docs/README.md"
if [ ! -f "$DOCS_FILE" ]; then
  echo "ERROR: docs index not found: $DOCS_FILE"
  exit 1
fi

service_count=0
rpc_count=0
missing_items=""

for proto_file in $(find "$PROTO_DIR" -name '*.proto' -print | sort); do
  package_name=$(sed -n 's/^package[[:space:]]\{1,\}\([^;]*\);/\1/p' "$proto_file" | head -n 1)
  rel_path=${proto_file#"$PROTO_DIR"/}

  while IFS= read -r line; do
    if [[ "$line" =~ ^[[:space:]]*service[[:space:]]+([A-Za-z0-9_]+) ]]; then
      service_count=$((service_count + 1))
      service_name="${BASH_REMATCH[1]}"
      fq_service="$service_name"
      if [ -n "$package_name" ]; then
        fq_service="${package_name}.${service_name}"
      fi

      if ! grep -Fq "\`$fq_service\`" "$DOCS_FILE"; then
        missing_items="${missing_items}service ${fq_service} (${rel_path})
"
      fi
    fi

    if [[ "$line" =~ ^[[:space:]]*rpc[[:space:]]+([A-Za-z0-9_]+) ]]; then
      rpc_count=$((rpc_count + 1))
      rpc_name="${BASH_REMATCH[1]}"

      if ! grep -Fq "\`$rpc_name\`" "$DOCS_FILE"; then
        missing_items="${missing_items}rpc ${rpc_name} (${rel_path})
"
      fi
    fi
  done < "$proto_file"
done

if [ "$service_count" -eq 0 ]; then
  echo "ERROR: no proto services found under $PROTO_DIR"
  exit 1
fi

if [ -n "$missing_items" ]; then
  echo "Proto services or RPCs missing from docs/README.md:"
  printf "%s" "$missing_items" | sed '/^$/d; s/^/  - /'
  exit 1
fi

echo "OK: docs/README.md indexes $service_count proto services and $rpc_count RPCs"
