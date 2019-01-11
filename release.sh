#!/usr/bin/env bash
version=$(git describe --tags | cut -c1-7)
snapcraft list-revisions monexec | grep $version  | awk '{print $1}' | xargs -n 1 -i snapcraft release monexec '{}' beta,candidate,stable
goreleaser --rm-dist