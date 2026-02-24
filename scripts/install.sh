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

  echo "diragent binary not found. Pass path as first argument." >&2
  exit 1
}

install_linux() {
  local binary_path="$1"
  local data_path="$2"

  local actions_dir="${HOME}/.local/share/file-manager/actions"
  local kde_dir="${HOME}/.local/share/kio/servicemenus"
  mkdir -p "${actions_dir}" "${kde_dir}"

  cat > "${actions_dir}/diragent-open-in-codex.desktop" <<EOF
[Desktop Entry]
Type=Action
Name=Open in Codex (DirAgent)
Icon=${data_path}/assets/icons/linux/codex.png
Profiles=profile-zero;

[X-Action-Profile profile-zero]
MimeTypes=inode/directory;application/octet-stream;text/plain;
Exec=${binary_path} launch --tool codex --path "%f"
EOF

  cat > "${actions_dir}/diragent-open-in-claude.desktop" <<EOF
[Desktop Entry]
Type=Action
Name=Open in Claude Code (DirAgent)
Icon=${data_path}/assets/icons/linux/claude.png
Profiles=profile-zero;

[X-Action-Profile profile-zero]
MimeTypes=inode/directory;application/octet-stream;text/plain;
Exec=${binary_path} launch --tool claude --path "%f"
EOF

  cat > "${kde_dir}/dir-agent.desktop" <<EOF
[Desktop Entry]
Type=Service
X-KDE-ServiceTypes=KonqPopupMenu/Plugin
MimeType=inode/directory;application/octet-stream;text/plain;
Actions=OpenInCodex;OpenInClaude;

[Desktop Action OpenInCodex]
Name=Open in Codex (DirAgent)
Icon=${data_path}/assets/icons/linux/codex.png
Exec=${binary_path} launch --tool codex --path "%f"

[Desktop Action OpenInClaude]
Name=Open in Claude Code (DirAgent)
Icon=${data_path}/assets/icons/linux/claude.png
Exec=${binary_path} launch --tool claude --path "%f"
EOF

  echo "Installed DirAgent Linux context actions."
}

write_macos_app_source() {
  local output_file="$1"
  local binary_path="$2"
  local tool="$3"

  cat > "${output_file}" <<EOF
on run argv
  if (count of argv) > 0 then
    set targetPath to item 1 of argv
    do shell script quoted form of "${binary_path}" & " launch --tool ${tool} --path " & quoted form of targetPath & " >/dev/null 2>&1 &"
  end if
end run

on open inputItems
  if (count of inputItems) > 0 then
    set targetPath to POSIX path of (item 1 of inputItems)
    do shell script quoted form of "${binary_path}" & " launch --tool ${tool} --path " & quoted form of targetPath & " >/dev/null 2>&1 &"
  end if
end open
EOF
}

build_macos_icon() {
  local png_path="$1"
  local app_path="$2"
  local tmp_dir
  tmp_dir="$(mktemp -d)"
  local iconset="${tmp_dir}/icon.iconset"
  mkdir -p "${iconset}"

  sips -z 16 16 "${png_path}" --out "${iconset}/icon_16x16.png" >/dev/null
  sips -z 32 32 "${png_path}" --out "${iconset}/icon_16x16@2x.png" >/dev/null
  sips -z 32 32 "${png_path}" --out "${iconset}/icon_32x32.png" >/dev/null
  sips -z 64 64 "${png_path}" --out "${iconset}/icon_32x32@2x.png" >/dev/null
  sips -z 128 128 "${png_path}" --out "${iconset}/icon_128x128.png" >/dev/null
  sips -z 256 256 "${png_path}" --out "${iconset}/icon_128x128@2x.png" >/dev/null
  sips -z 256 256 "${png_path}" --out "${iconset}/icon_256x256.png" >/dev/null
  sips -z 512 512 "${png_path}" --out "${iconset}/icon_256x256@2x.png" >/dev/null
  sips -z 512 512 "${png_path}" --out "${iconset}/icon_512x512.png" >/dev/null
  cp "${png_path}" "${iconset}/icon_512x512@2x.png"

  iconutil -c icns "${iconset}" -o "${app_path}/Contents/Resources/applet.icns"
  rm -rf "${tmp_dir}"
}

install_macos() {
  local binary_path="$1"
  local data_path="$2"

  if ! command -v osacompile >/dev/null 2>&1; then
    echo "osacompile not found; cannot install macOS launcher apps." >&2
    exit 1
  fi

  local app_dir="${HOME}/Applications/DirAgent"
  mkdir -p "${app_dir}"

  local codex_src claude_src
  codex_src="$(mktemp)"
  claude_src="$(mktemp)"

  write_macos_app_source "${codex_src}" "${binary_path}" "codex"
  write_macos_app_source "${claude_src}" "${binary_path}" "claude"

  local codex_app="${app_dir}/Open in Codex (DirAgent).app"
  local claude_app="${app_dir}/Open in Claude Code (DirAgent).app"
  rm -rf "${codex_app}" "${claude_app}"

  osacompile -o "${codex_app}" "${codex_src}" >/dev/null
  osacompile -o "${claude_app}" "${claude_src}" >/dev/null
  rm -f "${codex_src}" "${claude_src}"

  if command -v iconutil >/dev/null 2>&1 && command -v sips >/dev/null 2>&1; then
    build_macos_icon "${data_path}/assets/icons/macos/codex.png" "${codex_app}"
    build_macos_icon "${data_path}/assets/icons/macos/claude.png" "${claude_app}"
  fi

  touch "${codex_app}" "${claude_app}"
  echo "Installed DirAgent macOS apps at ${app_dir}."
  echo "Use Finder: right-click -> Open With -> Open in Codex/Claude Code (DirAgent)."
}

main() {
  local binary_path
  binary_path="$(resolve_binary "${1:-}")"
  "${binary_path}" install-assets >/dev/null
  local data_path
  data_path="$("${binary_path}" path --kind data)"

  case "$(uname -s)" in
    Linux)
      install_linux "${binary_path}" "${data_path}"
      ;;
    Darwin)
      install_macos "${binary_path}" "${data_path}"
      ;;
    *)
      echo "Unsupported OS for install.sh: $(uname -s)" >&2
      exit 1
      ;;
  esac

  echo "Binary: ${binary_path}"
  echo "Data path: ${data_path}"
}

main "$@"
