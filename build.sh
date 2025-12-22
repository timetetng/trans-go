#!/bin/bash

platforms=("windows/amd64" "darwin/amd64" "darwin/arm64" "linux/amd64")

for platform in "${platforms[@]}"; do
  platform_split=(${platform//\// })
  GOOS=${platform_split[0]}
  GOARCH=${platform_split[1]}

  output_name="trans-go-$GOOS-$GOARCH"
  if [ $GOOS = "windows" ]; then
    output_name+='.exe'
  fi

  echo "正在构建: $platform"
  export GOOS=$GOOS
  export GOARCH=$GOARCH
  go build -o "dist/$output_name"
done
