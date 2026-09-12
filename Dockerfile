FROM golang:1.27-alpine AS builder

WORKDIR /build

COPY ./src/ /build/
RUN go build -v -o confrender .

FROM scratch

LABEL org.opencontainers.image.title="ConfRender"
LABEL org.opencontainers.image.description="Renders a file using Go templating"
LABEL org.opencontainers.image.authors="Love Tropics <contact@lovetropics.org>"

LABEL org.opencontainers.image.source=https://github.com/LoveTropics/ConfRender
LABEL org.opencontainers.image.license=MPL-2.0

COPY --from=builder /build/confrender /confrender

ENTRYPOINT [ "/confrender" ]
