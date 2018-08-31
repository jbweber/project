#!/bin/sh

set -ex

VERSION=$(cat version.txt)

docker build --rm --force-rm -t jbweber/project:${VERSION} .
