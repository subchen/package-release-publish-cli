#!/bin/bash -e

make glide-vc -C bintray-upload
make glide-vc -C github-release-upload
make glide-vc -C go-build
make glide-vc -C homebrew-tap
make glide-vc -C sha256sum-files
