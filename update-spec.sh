#!/bin/sh

UPSTREAM_URL="https://codeload.github.com/vscreen/spec/zip/master"
OUT_DIR="pb"

cd `dirname $0`

mkdir -p ${OUT_DIR}
curl -o /tmp/spec.zip -L ${UPSTREAM_URL}
unzip -d /tmp /tmp/spec.zip
mv /tmp/spec-master/go/* ${OUT_DIR}/