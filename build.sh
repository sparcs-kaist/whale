#!/usr/bin/env bash

VERSION=$1

if [[ $# -ne 1 ]] ; then
  echo "Usage: $(basename $0) <VERSION>"
  exit 1
fi

mkdir -pv /tmp/whale-builds

grunt release
docker build -t sparcs-kaist/whale:linux-amd64-${VERSION} -f build/linux/Dockerfile .
docker push sparcs-kaist/whale:linux-amd64-${VERSION}
docker build -t sparcs-kaist/whale:linux-amd64 -f build/linux/Dockerfile .
docker push sparcs-kaist/whale:linux-amd64
rm -rf /tmp/whale-builds/unix && mkdir -pv /tmp/whale-builds/unix/whale
mv dist/* /tmp/whale-builds/unix/whale
cd /tmp/whale-builds/unix
tar cvpfz whale-${VERSION}-linux-amd64.tar.gz whale
mv whale-${VERSION}-linux-amd64.tar.gz /tmp/whale-builds/
cd -

grunt release-arm
docker build -t sparcs-kaist/whale:linux-arm-${VERSION} -f build/linux/Dockerfile .
docker push sparcs-kaist/whale:linux-arm-${VERSION}
docker build -t sparcs-kaist/whale:linux-arm -f build/linux/Dockerfile .
docker push sparcs-kaist/whale:linux-arm
rm -rf /tmp/whale-builds/arm && mkdir -pv /tmp/whale-builds/arm/whale
mv dist/* /tmp/whale-builds/arm/whale
cd /tmp/whale-builds/arm
tar cvpfz whale-${VERSION}-linux-arm.tar.gz whale
mv whale-${VERSION}-linux-arm.tar.gz /tmp/whale-builds/
cd -

grunt release-arm64
docker build -t sparcs-kaist/whale:linux-arm64-${VERSION} -f build/linux/Dockerfile .
docker push sparcs-kaist/whale:linux-arm64-${VERSION}
docker build -t sparcs-kaist/whale:linux-arm64 -f build/linux/Dockerfile .
docker push sparcs-kaist/whale:linux-arm64
rm -rf /tmp/whale-builds/arm64 && mkdir -pv /tmp/whale-builds/arm64/whale
mv dist/* /tmp/whale-builds/arm64/whale
cd /tmp/whale-builds/arm64
tar cvpfz whale-${VERSION}-linux-arm64.tar.gz whale
mv whale-${VERSION}-linux-arm64.tar.gz /tmp/whale-builds/
cd -

grunt release-macos
rm -rf /tmp/whale-builds/darwin && mkdir -pv /tmp/whale-builds/darwin/whale
mv dist/* /tmp/whale-builds/darwin/whale
cd /tmp/whale-builds/darwin
tar cvpfz whale-${VERSION}-darwin-amd64.tar.gz whale
mv whale-${VERSION}-darwin-amd64.tar.gz /tmp/whale-builds/
cd -

grunt release-win
rm -rf /tmp/whale-builds/win && mkdir -pv /tmp/whale-builds/win/whale
cp -r dist/* /tmp/whale-builds/win/whale
cd /tmp/whale-builds/win
tar cvpfz whale-${VERSION}-windows-amd64.tar.gz whale
mv whale-${VERSION}-windows-amd64.tar.gz /tmp/whale-builds/

exit 0
