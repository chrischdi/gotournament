#!/usr/bin/env bash

set -e

platforms="darwin/386 darwin/amd64 linux/386 linux/amd64 linux/arm linux/arm64 windows/386 windows/amd64"

name="gotournament"
gofile="cmd/gotournament.go"

mkdir -p bin tar

for platform in ${platforms}
do
  split=(${platform//\// })
  goos=${split[0]}
  goarch=${split[1]}
	echo "# Building for $platform: $goos $goarch"

  # prepare
  ext=""
  if [ "$goos" == "windows" ]; then
    ext=".exe"
  fi
  mkdir -p bin/$goos/$goarch

  # build
  CGO_ENABLED=0 GOOS=$goos GOARCH=$goarch go build -ldflags='-s -w' -v -o bin/$goos/$goarch/$name$ext $gofile
	
	# add tpl
	cp -r tpl bin/$goos/$goarch/tpl
  
  # pack
	if [ "$goos" == "windows" ]; then
		base=$(pwd)
		pushd bin/$goos/$goarch
		zip -r $base/tar/$name-$goos-$goarch.zip ./*
		popd
	else
	  tar cfvz tar/$name-$goos-$goarch.tar.gz -C bin/$goos/$goarch .
	fi
done
