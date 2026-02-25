#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${script_dir}"

if [[ ! -f "./diragent" ]]; then
  echo "[DirAgent][ERROR] diragent not found in current folder." >&2
  echo "Please keep uninstall.sh next to diragent." >&2
  exit 1
fi

chmod +x ./diragent ./scripts/uninstall.sh

echo "[DirAgent] Uninstalling context menu and local assets/config..."
./scripts/uninstall.sh ./diragent
./diragent uninstall-assets --remove-config >/dev/null 2>&1 || true

echo "[DirAgent] Uninstall completed."
