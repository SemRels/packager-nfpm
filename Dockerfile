# SPDX-License-Identifier: Apache-2.0
# SPDX-FileCopyrightText: 2026 The packager-nfpm Authors

FROM golang:1.25-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/plugin ./cmd/plugin && \
    CGO_ENABLED=0 GOBIN=/out go install github.com/goreleaser/nfpm/v2/cmd/nfpm@v2.43.0

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/plugin /usr/local/bin/plugin
COPY --from=build /out/nfpm /usr/local/bin/nfpm
ENTRYPOINT ["/usr/local/bin/plugin"]
