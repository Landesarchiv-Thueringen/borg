#!/usr/bin/env bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

source "$SCRIPT_DIR/../.env"

tags=(
    gui
    server
    droid
    tika
    jhove
    magika
    verapdf
    odf-validator
    ooxml-validator
)

echo "Build and publish"
for tag in "${tags[@]}"; do
    echo "    $IMAGE_PREFIX/$tag:$IMAGE_VERSION"
done
echo
read -p "Continue? [y/N] "
if [[ $REPLY =~ ^[Yy]$ ]]
then
    echo
    echo Building images...
    (
        cd "$SCRIPT_DIR/.."
        docker compose build --pull
    )
    echo Publishing images...
    for tag in "${tags[@]}"; do
        docker push "$IMAGE_PREFIX/$tag:$IMAGE_VERSION"
    done
fi
