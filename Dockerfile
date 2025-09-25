# Get the gocmd image
FROM harbor.cyverse.org/de/gocmd:latest AS gocmd

# Build the app
FROM golang:1.24 AS golang
COPY . /go/src/github.com/cyverse-de/batch-exit-handler
WORKDIR /go/src/github.com/cyverse-de/batch-exit-handler
ENV CGO_ENABLED=0
RUN go build --buildvcs=false .

# Create the deployable container image
FROM debian:stable-slim
WORKDIR /app
COPY --from=gocmd /usr/bin/gocmd /bin/gocmd
COPY --from=golang /go/src/github.com/cyverse-de/batch-exit-handler/batch-exit-handler /bin/batch-exit-handler
ENTRYPOINT ["batch-exit-handler"]
