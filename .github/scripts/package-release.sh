#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: $0 <tag>" >&2
  exit 1
fi

tag="$1"
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
binaries_dir="${repo_root}/dist/binaries"
release_dir="${repo_root}/dist/release"
quickstart_path="${repo_root}/.github/release/README.quickstart.md"

if [[ ! -f "${quickstart_path}" ]]; then
  echo "missing quickstart template: ${quickstart_path}" >&2
  exit 1
fi

mkdir -p "${release_dir}"

shopt -s nullglob
binaries=("${binaries_dir}/diragent_${tag}_"*)
if [[ ${#binaries[@]} -eq 0 ]]; then
  echo "no build artifacts found under ${binaries_dir}" >&2
  exit 1
fi

for binary_path in "${binaries[@]}"; do
  binary_name="$(basename "${binary_path}")"
  os_arch="${binary_name#diragent_${tag}_}"

  package_os=""
  package_arch=""
  output_binary_name="diragent"
  script_names=(install.sh uninstall.sh)

  if [[ "${binary_name}" == *.exe ]]; then
    os_arch="${os_arch%.exe}"
    output_binary_name="diragent.exe"
    script_names=(install.ps1 uninstall.ps1)
  fi

  package_os="${os_arch%%_*}"
  package_arch="${os_arch#*_}"
  if [[ -z "${package_os}" || -z "${package_arch}" || "${package_os}" == "${package_arch}" ]]; then
    echo "cannot parse os/arch from artifact name: ${binary_name}" >&2
    exit 1
  fi

  package_name="diragent_${tag}_${package_os}_${package_arch}"
  stage_root="$(mktemp -d)"
  package_dir="${stage_root}/${package_name}"
  mkdir -p "${package_dir}/scripts"

  cp "${binary_path}" "${package_dir}/${output_binary_name}"
  cp "${quickstart_path}" "${package_dir}/README.quickstart.md"

  for script_name in "${script_names[@]}"; do
    source_script="${repo_root}/scripts/${script_name}"
    if [[ ! -f "${source_script}" ]]; then
      echo "missing script for package ${package_name}: ${source_script}" >&2
      exit 1
    fi
    cp "${source_script}" "${package_dir}/scripts/${script_name}"
  done

  if [[ "${output_binary_name}" != "diragent.exe" ]]; then
    chmod +x "${package_dir}/diragent" "${package_dir}/scripts/install.sh" "${package_dir}/scripts/uninstall.sh"
  fi

  (
    cd "${stage_root}"
    zip -qr "${release_dir}/${package_name}.zip" "${package_name}"
  )
  rm -rf "${stage_root}"
done
