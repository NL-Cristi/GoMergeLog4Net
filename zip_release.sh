#!/usr/bin/env bash
# ---------------------------------------------------------------------------
# zip_release.sh
# Creates a .zip archive for every sub-folder inside Release/Binaries
# The resulting <folder>.zip files are placed in Release/ZIPS.
# ---------------------------------------------------------------------------
# Example result:
#   Release/
#     Binaries/
#       linux_amd64/GoMergeLog4Net
#       windows_amd64/GoMergeLog4Net.exe
#     ZIPS/
#       linux_amd64.zip
#       windows_amd64.zip
#
# Usage:
#   bash ./zip_release.sh
# ---------------------------------------------------------------------------

set -euo pipefail

# Resolve repository root (directory where this script lives)
SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

BINARIES_DIR="${SCRIPT_DIR}/Release/Binaries"
ZIPS_DIR="${SCRIPT_DIR}/Release/ZIPS"

if [[ ! -d "${BINARIES_DIR}" ]]; then
  echo "Error: Binaries directory not found: ${BINARIES_DIR}" >&2
  exit 1
fi

mkdir -p "${ZIPS_DIR}"

for dir in "${BINARIES_DIR}"/*/; do
  [[ -d "${dir}" ]] || continue # skip if no directory
  folder_name="$(basename "${dir}")"
  dest_zip="${ZIPS_DIR}/${folder_name}.zip"
  rm -f "${dest_zip}"

  # Find first regular file (the binary) inside the folder
  binary_file="$(find "${dir}" -maxdepth 1 -type f -print -quit)"
  if [[ -z "${binary_file}" ]]; then
    echo "Warning: No binary found in ${dir}. Skipping."
    continue
  fi

  echo "Zipping $(basename "${binary_file}") from ${folder_name} -> ${dest_zip}"
  # -j = junk the path, so the file is at root of zip
  (cd "$(dirname "${binary_file}")" && zip -q -j "${dest_zip}" "$(basename "${binary_file}")")

done

echo -e "\nDone. Archives created in ${ZIPS_DIR}/" 