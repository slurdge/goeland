FROM golang:alpine as builder

ARG TARGETARCH
ARG TAG

COPY . /usr/src/goeland
WORKDIR /usr/src/goeland
RUN go get -v ./...
ENV GOOS=linux
ENV GOARCH=$TARGETARCH
ENV CGO_ENABLED=0
RUN go build -o /goeland \
    -ldflags "-X github.com/slurdge/goeland/version.GitCommit=${GIT_COMMIT} -X github.com/slurdge/goeland/version.BuildDate=${BUILD_DATE} -X github.com/slurdge/goeland/internal/goeland/fetch.clientID=${IMGUR_CLIENT_ID}"

FROM alpine
COPY --from=builder --chown=1000 /goeland /goeland
USER 1000
WORKDIR /data
ENTRYPOINT ["/goeland"]
CMD ["daemon"]