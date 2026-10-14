#!/bin/sh
# Tải binary mới nhất từ GitHub Releases và chạy.
# Dùng: curl -fsSL https://raw.githubusercontent.com/dominhhieu1405/Simple-Auction-system/main/scripts/run.sh | sh
set -eu

REPO="dominhhieu1405/Simple-Auction-system"
VERSION="${VERSION:-latest}"
DIR="${INSTALL_DIR:-simple-auction}"

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
  linux|darwin) ;;
  *) echo "Unsupported OS: $os" >&2; exit 1 ;;
esac

arch=$(uname -m)
case "$arch" in
  x86_64|amd64) arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) echo "Unsupported architecture: $arch" >&2; exit 1 ;;
esac

name="auction-server-$os-$arch"
if [ "$VERSION" = "latest" ]; then
  url="https://github.com/$REPO/releases/latest/download/$name"
else
  url="https://github.com/$REPO/releases/download/$VERSION/$name"
fi

mkdir -p "$DIR"
cd "$DIR"

echo "Downloading $url"
curl -fL "$url" -o auction-server
chmod +x auction-server

if [ ! -f .env ]; then
  secret=$(head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \n')
  printf 'JWT_SECRET=%s\nDB_PATH=data/auction.db\nSERVER_PORT=8080\n' "$secret" > .env
  echo "Created .env with random JWT_SECRET"
fi

echo "Starting server in $(pwd)"
exec ./auction-server