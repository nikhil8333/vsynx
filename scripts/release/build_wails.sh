#!/usr/bin/env bash
set -euo pipefail

PLATFORM=${1:?"platform (e.g., windows/amd64) required"}
OUTPUT=${2:?"output name required"}
EXT=${3:-}
INSTALLER=${4:-false}

CMD=(wails build -platform "$PLATFORM" -o "${OUTPUT}${EXT}")
if [[ "$INSTALLER" == "true" ]]; then
  CMD=(wails build -nsis -platform "$PLATFORM" -o "${OUTPUT}${EXT}")
fi

echo "Running: ${CMD[*]}"
"${CMD[@]}"
