#!/bin/bash -e

CWD=$(cd $(dirname $0); pwd)

VERSION=0.1.0

COMPONENTS="
    go-build
    sha256sum-files
    bintray-upload
    github-release-upload
"

echo "Building publish-toolset ..."

for app in $COMPONENTS; do
    cd $CWD/$app && make build VERSION=$VERSION
done
    
if [ -z "$TRAVIS_TAG" ]; then
    exit 0
fi

echo "Uploading files into github release: $TRAVIS_TAG ..."

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

cd $CWD/_releases/linux/   && zip ../publish-toolset-linux.zip   *
cd $CWD/_releases/darwin/  && zip ../publish-toolset-darwin.zip  *
cd $CWD/_releases/windows/ && zip ../publish-toolset-windows.zip *

"$sha256sum_files" $CWD/_releases/*.zip

"$github_release_upload" $CWD/_releases/publish-toolset-*
