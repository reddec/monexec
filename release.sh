#!/usr/bin/env bash
version=$(git describe --tags | cut -c1-7)
snapcraft list-revisions fluent-amqp | grep $version  | awk '{print $1}' | xargs -n 1 -i snapcraft release fluent-amqp '{}' beta,candidate,stable
goreleaser --rm-dist