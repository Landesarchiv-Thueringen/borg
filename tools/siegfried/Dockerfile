FROM golang:alpine3.23 AS build
ARG BORG_VERSION=${BORG_VERSION}
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o siegfried_api -ldflags "-X main.version=${BORG_VERSION}" ./cmd

FROM alpine:3.23 AS extract
RUN apk add --no-cache unzip
COPY third_party /third_party
WORKDIR /third_party
RUN unzip -j siegfried_1-11-2_linux64.zip
RUN rm siegfried_1-11-2_linux64.zip

FROM alpine:3.23 AS prod
WORKDIR /borg/tools/siegfried
COPY --from=build /build/siegfried_api .
COPY --from=extract /third_party/sf third_party/
RUN ./third_party/sf -update
CMD ["./siegfried_api"]