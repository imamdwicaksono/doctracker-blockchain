#!/bin/bash

APP_NAME="doc-tracker"

PLATFORMS=("windows/amd64" "linux/amd64" "darwin/amd64" "linux/arm64" "windows/arm64")

mkdir -p builds

for PLATFORM in "${PLATFORMS[@]}"
do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}

    OUTPUT_NAME=$APP_NAME"-"$GOOS"-"$GOARCH
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME+=".exe"
    fi

    echo "🔧 Building for $GOOS/$GOARCH..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -o builds/$OUTPUT_NAME .

    if [ $? -ne 0 ]; then
        echo "❌ Build failed for $GOOS/$GOARCH"
        exit 1
    else
        echo "✅ Built: builds/$OUTPUT_NAME"
    fi
done

echo "🎉 All builds completed!"
