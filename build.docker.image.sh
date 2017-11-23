#!/bin/bash
image_folder="docker"
app="idp"
image_name="reg.igenetech.com/webapp/$app"
timestamp=`date +%Y%m%d_%H%M`
if [ -z "$1" ]
then
    version="test"
else
    version=$1
fi
image_tag=$image_name:${version}_${timestamp}
image_file=$image_folder/$app.${version}_${timestamp}
echo ">>> building app: $version"
go build
echo ">>> building images: $image_tag"
docker build \
    -f Dockerfile.prod \
    -t $image_tag \
    --build-arg GIT_COMMIT=$(git log -1 --format=%h) \
    .
echo ">>> saving images: $image_file"
docker save -o $image_file $image_tag
