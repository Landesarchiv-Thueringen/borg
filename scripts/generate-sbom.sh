#!/usr/bin/env bash

# Go to repository root dir
cd -- "$( dirname -- "${BASH_SOURCE[0]}" )/.."

tags=(
    gui
    server
    droid
    siegfried
    tika
    jhove
    magika
    mediainfo
    verapdf
    odf-validator
    ooxml-validator
)

mkdir -p docs/sbom
docker compose -f compose.yml build --pull
for tag in "${tags[@]}"; do
    docker scout sbom --format list "localhost/borg/$tag" > "docs/sbom/sbom_borg_$tag.txt"
    docker scout sbom --format cyclonedx "localhost/borg/$tag" > "docs/sbom/sbom_borg_$tag.cdx"
done