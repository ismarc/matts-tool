#!/usr/bin/env bash

go test -coverprofile=coverage.out -v \
  ./internal/... \
  ./cmd/...
