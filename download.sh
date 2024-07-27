#!/usr/bin/env sh
set -eu

# Requires jq, curl
if ! command -v jq >/dev/null; then
    echo "jq is not installed; it is required by this script"
    exit 1
fi

if ! command -v curl >/dev/null; then
    echo "curl is not installed; it is required by this script"
    exit 1
fi

# Usage: ./install.sh [tag]

tag="${1:-}"
if [ -z "$tag" ]; then
    echo "No tag provided; downloading the latest release"
    tag="$(curl --retry 3 -s "https://api.github.com/repos/stdedos/junit2html/releases/latest" | jq -r '.tag_name')"
fi

version=$(echo "${tag}" | cut -c 2-)
url="https://github.com/stdedos/junit2html/releases/download/${tag}/junit2html_${version}_$(uname)_$(uname -m | sed 's/aarch64/arm64/').tar.gz"
curl --retry 3 -L "${url}" | tar --wildcards -zxvf - "junit2html_*"

chmod +x junit2html_*
