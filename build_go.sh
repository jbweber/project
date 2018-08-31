#!/bin/sh

set -xe

VERSION=$(cat version.txt)
GITCOMMIT=$(git rev-parse --short HEAD)
GITUNTRACKEDCHANGES=$(git status --porcelain --untracked-files=no)

if [ "x${GITUNTRACKEDCHANGES}" != "x" ]; then
    GITCOMMIT=${GITCOMMIT}-dirty
fi

CGO_ENABLED=0 go test -v ./...
CGO_ENABLED=0 go install -ldflags "-w -X github.com/jbweber/project/internal.GitCommit=${GITCOMMIT} -X github.com/jbweber/project/internal.Version=${VERSION}" ./...
