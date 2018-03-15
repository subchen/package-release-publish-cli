#!/bin/bash -e

CWD=$(cd $(dirname $0); pwd)

COMPONENTS="
    go-build
    sha256sum-files
    bintray-upload
    github-release-upload
"

for app in $COMPONENTS; do
    cd $CWD/$app && make release
done

github_release_upload=$(ls $CWD/github-release-upload/_releases/github-release-upload-*-linux-amd64)

if [ -n "$TARVIS_TAG" ]; then
    echo "Uploading into github release: $TARVIS_TAG ..."

    for app in $COMPONENTS; do
        "$github_release_upload" "$CWD/$app/_releases"
    done

fi
