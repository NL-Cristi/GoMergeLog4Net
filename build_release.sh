#!/usr/bin/env bash
# Build GoMergeLog4Net for several platforms and place each binary in
#   Release/<goos>_<goarch>/GoMergeLog4Net[.exe]

set -euo pipefail

# -------------------------------------------------------
# Paths relative to where the script itself lives
SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP_DIR="${SCRIPT_DIR}"            # main.go is in the same dir
APP_NAME="GoMergeLog4Net"          # binary name
OUT_ROOT="${SCRIPT_DIR}/Release/Binaries"  # match PowerShell script
# -------------------------------------------------------

PLATFORMS=(               # {GOOS}/{GOARCH}
  "windows/amd64"
  "linux/amd64"
  "darwin/amd64"
  "darwin/arm64"
)

# Optional: inject Git tag / hash
VERSION=$(git -C "${SCRIPT_DIR}" describe --tags --always 2>/dev/null || echo "dev")
LDFLAGS="-s -w -X main.version=${VERSION}"

mkdir -p "${OUT_ROOT}"

echo "Building version: ${VERSION}"

for target in "${PLATFORMS[@]}"; do
  IFS=/ read -r GOOS GOARCH <<<"${target}"
  OUT_DIR="${OUT_ROOT}/${GOOS}_${GOARCH}"
  mkdir -p "${OUT_DIR}"

  EXT=""
  [[ ${GOOS} == "windows" ]] && EXT=".exe"

  echo "=> building ${GOOS}/${GOARCH}"
  env CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} \
      go build -C "${APP_DIR}" \
               -ldflags "${LDFLAGS}" \
               -o "${OUT_DIR}/${APP_NAME}${EXT}" \
               .
done

echo -e "\nBinaries are under ${OUT_ROOT}/"