FROM golang:alpine AS build
RUN apk add git opusfile-dev opus-dev pkgconfig gcc musl-dev
ADD . /go/src/ContainerStarterBot/
WORKDIR /go/src/ContainerStarterBot
RUN go build .

FROM alpine
WORKDIR /data
COPY --from=build /go/src/ContainerStarterBot/ContainerStarterBot /bin/ContainerStarterBot
ENTRYPOINT [ "ContainerStarterBot" ]