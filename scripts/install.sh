#!/usr/bin/env bash
# Installs the latest linden CLI release for the current OS/arch.
#
#   curl -fsSL https://raw.githubusercontent.com/mylinden-tech/linden-cli/main/scripts/install.sh | bash
#
# Override the install directory with LINDEN_INSTALL_DIR (default: /usr/local/bin,
# falling back to ~/.local/bin if that isn't writable).

set -euo pipefail

REPO="mylinden-tech/linden-cli"
BINARY="linden"

log() { printf '%s\n' "$*" >&2; }
die() { log "error: $*"; exit 1; }

detect_os() {
  case "$(uname -s)" in
    Linux) echo "linux" ;;
    Darwin) echo "darwin" ;;
    *) die "unsupported OS: $(uname -s). Build from source instead: https://github.com/${REPO}#installation" ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    *) die "unsupported architecture: $(uname -m)" ;;
  esac
}

main() {
  command -v curl >/dev/null 2>&1 || die "curl is required"
  command -v tar >/dev/null 2>&1 || die "tar is required"

  local os arch version url tmp_dir
  os="$(detect_os)"
  arch="$(detect_arch)"

  version="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep -m1 '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')"
  [ -n "$version" ] || die "could not determine latest release version"

  url="https://github.com/${REPO}/releases/download/${version}/${BINARY}_${os}_${arch}.tar.gz"
  log "Downloading ${BINARY} ${version} for ${os}/${arch}..."

  tmp_dir="$(mktemp -d)"
  trap 'rm -rf "$tmp_dir"' EXIT

  curl -fsSL "$url" -o "${tmp_dir}/${BINARY}.tar.gz" \
    || die "download failed: ${url}"
  tar -xzf "${tmp_dir}/${BINARY}.tar.gz" -C "$tmp_dir" "$BINARY"

  local install_dir="${LINDEN_INSTALL_DIR:-/usr/local/bin}"
  if [ ! -w "$install_dir" ] 2>/dev/null; then
    install_dir="${HOME}/.local/bin"
    mkdir -p "$install_dir"
  fi

  install -m 0755 "${tmp_dir}/${BINARY}" "${install_dir}/${BINARY}"
  log "Installed ${BINARY} ${version} to ${install_dir}/${BINARY}"

  case ":$PATH:" in
    *":${install_dir}:"*) ;;
    *) log "Note: ${install_dir} is not on your PATH. Add it, e.g.: export PATH=\"${install_dir}:\$PATH\"" ;;
  esac

  "${install_dir}/${BINARY}" --version
}

main "$@"
