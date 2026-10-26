#!/usr/bin/env bash
# This script installs the Gine project dependencies and sets up the environment.
version="0.1.5"
must=("curl" "tar")
platform="$(uname -s)"
arch="$(uname -m)"
source_url="https://github.com/sumit-poudel/gine/releases/download/v$version/gine_Linux_"

if ! command -v "${must[0]}" &> /dev/null || ! command -v "${must[1]}" &> /dev/null; then
    echo "Required tools are not installed. Please install curl or tar and try again."
    exit 1
fi

if [ -n "${TMPDIR:-}" ] && [ -d "${TMPDIR}" ]; then
    temp="$(mktemp -d "$TMPDIR/gine-${version}-XXXXXX")"
else
    temp="$(mktemp -d "/tmp/gine-${version}-XXXXXX")"
fi
trap 'rm -rf "$temp"' EXIT

if [ "$platform" = "Linux" ]; then
    platform="linux"
else
    echo "Unsupported OS: $(uname -s)"
    exit 1
fi

case "$platform-$arch" in
    linux-x86_64)
        arch="x86_64"
        ;;
    linux-aarch64 | linux-arm64)
        arch="arm64"
        ;;
    *)
        echo "Unsupported architecture: $arch"
        exit 1
        ;;
esac

if command -v curl >/dev/null 2>&1; then
    curl () {
        command curl -fL "$@"
    }
elif command -v wget >/dev/null 2>&1; then
    curl () {
        command wget -qO- "$@"
    }
else
    echo "Neither curl nor wget is installed. Please install one of them and try again."
    exit 1
fi

echo "Downloading gine $version for ${platform}/${arch}..."
archive="$temp/gine.tar.gz"
if ! curl "$source_url$arch.tar.gz" > "$archive"; then
    echo "Download failed. Check that gine v$version supports ${platform}/${arch}."
    exit 1
fi

tar -xzf "$archive" -C "$temp"
if [ ! -f "$temp/gine" ]; then
    echo "Archive did not contain the gine binary."
    exit 1
fi

mkdir -p "$HOME/.local/bin"
install -m 0755 "$temp/gine" "$HOME/.local/bin/gine"
echo "Installed gine $version to $HOME/.local/bin/gine"

if [ "$(command -v gine)" = "$HOME/.local/bin/gine" ]; then
    echo "Run it with: gine"
else
    echo "To run Gine from your terminal, you must add ~/.local/bin to your PATH"
    echo "Run:"

    case "$SHELL" in
        *fish)
            echo "   fish_add_path -U $HOME/.local/bin"
            ;;
        *zsh)
            echo "   echo 'export PATH=\$HOME/.local/bin:\$PATH' >> ~/.zshrc"
            echo "   source ~/.zshrc"
            ;;
        *)
            echo "   echo 'export PATH=\$HOME/.local/bin:\$PATH' >> ~/.bashrc"
            echo "   source ~/.bashrc"
            ;;
    esac
fi
