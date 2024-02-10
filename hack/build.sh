#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

bin_dir="${base_dir}/bin"

GOOS="windows"
GOARCH="${GOARCH:-$(go env GOARCH)}"

GO_LD_FLAGS="${GO_LD_FLAGS:-"-s"}"

output_name="${bin_dir}/wsl2-ssh-pageant.exe"

pushd "${base_dir}" >/dev/null

echo "Building $(basename "${output_name}")"
GOOS="${GOOS}" GOARCH="${GOARCH}" go build -ldflags="${GO_LD_FLAGS}" -o "${output_name}" ./cmd/...
