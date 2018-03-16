#!/bin/bash -e

tag=v0.1.0

os=$(uname -s | tr '[A-Z]' '[a-z]')

curl -fSL https://github.com/subchen/publish-toolset/releases/download/$tag/publish-toolset-$os.zip -o publish-toolset.zip
unzip publish-toolset.zip -d /usr/local/bin/
rm -f publish-toolset.zip
