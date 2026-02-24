#!/usr/bin/env bash
set -euo pipefail

resolve_binary() {
  local provided="${1:-}"
  if [[ -n "${provided}" && -x "${provided}" ]]; then
    printf '%s\n' "${provided}"
    return 0
  fi

  local script_dir
  script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  local local_bin="${script_dir}/../diragent"
  if [[ -x "${local_bin}" ]]; then
    printf '%s\n' "${local_bin}"
    return 0
  fi

  if command -v diragent >/dev/null 2>&1; then
    command -v diragent
    return 0
  fi

  printf '%s\n' ""
}

remove_linux_entries() {
  rm -f "${HOME}/.local/share/file-manager/actions/diragent-open-in-codex.desktop"
  rm -f "${HOME}/.local/share/file-manager/actions/diragent-open-in-claude.desktop"
  rm -f "${HOME}/.local/share/kio/servicemenus/dir-agent.desktop"
  echo "Removed DirAgent Linux context actions."
}

remove_macos_entries() {
  rm -rf "${HOME}/Applications/DirAgent/Open in Codex (DirAgent).app"
  rm -rf "${HOME}/Applications/DirAgent/Open in Claude (DirAgent).app"
  rmdir "${HOME}/Applications/DirAgent" 2>/dev/null || true
  echo "Removed DirAgent macOS launcher apps."
}

main() {
  local binary_path
  binary_path="$(resolve_binary "${1:-}")"

  case "$(uname -s)" in
    Linux)
      remove_linux_entries
      ;;
    Darwin)
      remove_macos_entries
      ;;
    *)
      echo "Unsupported OS for uninstall.sh: $(uname -s)" >&2
      exit 1
      ;;
  esac

  if [[ -n "${binary_path}" ]]; then
    "${binary_path}" uninstall-assets >/dev/null || true
  fi

  echo "Uninstall complete."
}

main "$@"
