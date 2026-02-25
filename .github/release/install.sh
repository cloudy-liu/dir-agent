#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${script_dir}"

if [[ ! -f "./diragent" ]]; then
  echo "[DirAgent][ERROR] diragent not found in current folder." >&2
  echo "Please keep install.sh next to diragent." >&2
  exit 1
fi

chmod +x ./diragent ./scripts/install.sh ./scripts/uninstall.sh

echo "[DirAgent] Cleaning previous install (keep existing config)..."
./scripts/uninstall.sh ./diragent >/dev/null 2>&1 || true

echo "[DirAgent] Installing context menu..."
./scripts/install.sh ./diragent

echo "[DirAgent] Install completed."
