FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

COPY go.mod main.go /build/
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o confrender .

FROM scratch

# Allow us to create a condition to immediately & successfully exit from this image
# For e.g. creating a service in a composefile whose sole purpose is to pull
# the image down, bypassing issues with inheriting pull_policy from the parent
# service

COPY --from=tianon/true:multiarch /true /true


LABEL org.opencontainers.image.title="ConfRender"
LABEL org.opencontainers.image.description="Renders a file using Go templating"
LABEL org.opencontainers.image.authors="Love Tropics <contact@lovetropics.org>"

LABEL org.opencontainers.image.source=https://github.com/LoveTropics/ConfRender
LABEL org.opencontainers.image.license=MPL-2.0

COPY --from=builder /build/confrender /confrender

ENTRYPOINT [ "/confrender" ]
