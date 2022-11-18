#!/usr/bin/env bash

set -o errexit
set -o nounset

binaryname=${1}
ostype=${2:-"native"}
version_path="github.com/pavisalavisa/juggler/internal/config.Version"

echo "[+] Build OS type selected: ${ostype}"

ldf_cmp="-s -w -extldflags '-static'"
f_ver="-X ${version_path}=${VERSION:-dev}"

function build_arch() {
    EXTENSION=$1
    GOOS=$2
    GOARCH=$3
    mkdir -p bin
    final_out=./bin/${binaryname}${EXTENSION:-}
    EXTENSION="${EXTENSION}" GOOS="${GOOS}" GOARCH="${GOARCH}" CGO_ENABLED=0 go build -o ${final_out} --ldflags "${ldf_cmp} ${f_ver}" ./cmd/${binaryname}
}

if [ $ostype == 'linux' ]; then
    echo "[*] Building ${binaryname} binary for ${ostype}..."
    build_arch "-linux-amd64" "linux" "amd64"
elif [ $ostype == 'macos' ]; then
    echo "[*] Building ${binaryname} binary for ${ostype}..."
    build_arch "-darwin-arm64" "darwin" "arm64"
elif [ $ostype == 'windows' ]; then
    echo "[*] Building ${binaryname} binary for ${ostype}..."
    build_arch "-windows-amd64.exe" "windows" "amd64"
else
    echo "[*] Building ${binaryname} binary..."
    build_arch "" "" ""
fi
