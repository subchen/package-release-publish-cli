#!/bin/bash -e

CWD=$(cd $(dirname $0); pwd)

COMPONENTS="
    go-build
    sha256sum-files
    bintray-upload
    github-release-upload
"

echo "Building publish-toolset ..."

for app in $COMPONENTS; do
    cd $CWD/$app && make build
done
    
if [ -z "$TRAVIS_TAG" ]; then
    exit 0
fi

echo "Uploading files into github release: $TARVIS_TAG ..."

#version="${TRAVIS_TAG}"
github_release_upload="$CWD/_releases/linux/github-release-upload"
sha256sum_files="$CWD/_releases/linux/sha256sum-files"

rm -rf   $CWD/_releases/
mkdir -p $CWD/_releases/linux
mkdir -p $CWD/_releases/darwin
mkdir -p $CWD/_releases/windows

for app in $COMPONENTS; do
    mv -f $CWD/$app/_releases/*-linux-*   $CWD/_releases/linux/$app
    mv -f $CWD/$app/_releases/*-darwin-*  $CWD/_releases/darwin/$app
    mv -f $CWD/$app/_releases/*-windows-* $CWD/_releases/windows/$app.exe
done

zip $CWD/_releases/publish-toolset-linux.zip   $CWD/_releases/linux/*
zip $CWD/_releases/publish-toolset-darwin.zip  $CWD/_releases/darwin/*
zip $CWD/_releases/publish-toolset-windows.zip $CWD/_releases/windows/*

"$sha256sum_files" $CWD/_releases/*.zip

"$github_release_upload" "$CWD/$app/_releases"
