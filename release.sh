#!/bin/bash
set -e

VERSION="${1:-v1.0.0}"
APP="LongDarkBackup"
OUT="dist"
SRC="./app/"
LDFLAGS="-s -w"
LINUX_IMAGE="ldb-linux-builder"
DOCS=(LICENSE NOTICE README.md README_RU.md)

rm -rf "$OUT"
mkdir -p "$OUT"

pack_zip() {
    local name=$1 dir="${OUT}/$1"
    for f in "${DOCS[@]}"; do cp "$f" "$dir/"; done
    (cd "$OUT" && zip -rq "${name}.zip" "${name}")
    rm -rf "$dir"
    echo "  -> ${OUT}/${name}.zip"
}

build() {
    local os=$1 arch=$2 suffix=$3 ext=$4 cgo=$5
    local name="${APP}-${suffix}"
    echo "Building ${name}..."
    mkdir -p "${OUT}/${name}"
    CGO_ENABLED="$cgo" GOOS="$os" GOARCH="$arch" \
        go build -ldflags="$LDFLAGS" -o "${OUT}/${name}/${APP}${ext}" "$SRC"
    pack_zip "$name"
}

build_macos_dmg() {
    local arch=$1 suffix=$2
    local name="${APP}-${suffix}"
    local appdir="${OUT}/${APP}.app"
    local dmg="${OUT}/${name}.dmg"
    local staging="${OUT}/_dmg_${suffix}"
    echo "Building ${name}.dmg..."

    # Build binary
    CGO_ENABLED=1 GOOS=darwin GOARCH="$arch" \
        go build -ldflags="$LDFLAGS" -o "${OUT}/${APP}" "$SRC"

    # Ad-hoc sign binary
    codesign --sign - --force "${OUT}/${APP}"

    # Create .app bundle
    rm -rf "$appdir"
    mkdir -p "${appdir}/Contents/MacOS"
    mkdir -p "${appdir}/Contents/Resources"
    cp "${OUT}/${APP}" "${appdir}/Contents/MacOS/${APP}"
    cp pkg/macos/icon.icns "${appdir}/Contents/Resources/icon.icns"
    sed "s/VERSION/${VERSION#v}/g" pkg/macos/Info.plist > "${appdir}/Contents/Info.plist"

    # Ad-hoc sign .app bundle
    codesign --sign - --force "${appdir}"

    # Create DMG
    rm -rf "$staging"
    mkdir -p "$staging"
    cp -R "$appdir" "$staging/"
    ln -s /Applications "$staging/Applications"

    hdiutil create -volname "${APP}" \
        -srcfolder "$staging" \
        -ov -format UDZO \
        "$dmg" > /dev/null

    rm -rf "$appdir" "$staging" "${OUT}/${APP}"
    echo "  -> ${dmg}"
}

build_linux() {
    local arch=$1 suffix=$2
    local name="${APP}-${suffix}"
    echo "Building ${name} (docker)..."
    mkdir -p "${OUT}/${name}"

    docker build -q -t "$LINUX_IMAGE" -f Dockerfile.linux .
    docker run --rm \
        -v "$PWD":/src -w /src \
        -e GOFLAGS=-buildvcs=false \
        "$LINUX_IMAGE" \
        go build -ldflags="$LDFLAGS" -o "dist/${name}/${APP}" "$SRC"

    pack_zip "$name"
}

# ── Windows (CGO off) ──
build windows 386   "windows-x32" ".exe" 0
build windows amd64 "windows-x64" ".exe" 0

# ── macOS .dmg (ad-hoc signed .app bundle) ──
build_macos_dmg amd64 "macos-intel"
build_macos_dmg arm64 "macos-arm"

# ── Linux (Docker with GTK) ──
build_linux amd64 "linux-x64"

echo ""
echo "Done. Release files in ${OUT}/:"
ls -lh "$OUT"/${APP}-*

# Create GitHub release
if [ "$2" = "--publish" ]; then
    echo ""
    echo "Pushing to origin..."
    git push origin main

    echo "Creating tag ${VERSION}..."
    git tag "$VERSION"
    git push origin "$VERSION"

    echo "Creating GitHub release ${VERSION}..."
    gh release create "$VERSION" \
        "$OUT"/${APP}-windows-x32.zip \
        "$OUT"/${APP}-windows-x64.zip \
        "$OUT"/${APP}-macos-intel.dmg \
        "$OUT"/${APP}-macos-arm.dmg \
        "$OUT"/${APP}-linux-x64.zip \
        --title "$VERSION" \
        --notes "$(cat <<EOF
## ${APP} ${VERSION}

### Downloads

| Platform | File |
|----------|------|
| Windows x32 | ${APP}-windows-x32.zip |
| Windows x64 | ${APP}-windows-x64.zip |
| macOS Intel | ${APP}-macos-intel.dmg |
| macOS Apple Silicon | ${APP}-macos-arm.dmg |
| Linux x64 | ${APP}-linux-x64.zip |

### macOS

Open the .dmg and drag **${APP}** to **Applications**.
If macOS blocks it: **System Settings → Privacy & Security → Open Anyway**.

### Windows / Linux

Unzip and run. No installation required.
EOF
)"
    echo "Release ${VERSION} published!"
fi
