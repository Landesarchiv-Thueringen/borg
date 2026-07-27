FROM golang:alpine3.23 AS build
ARG BORG_VERSION=${BORG_VERSION}
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o verapdf_api -ldflags "-X main.version=${BORG_VERSION}" ./cmd

FROM alpine:3.23 AS extract
RUN apk add --no-cache unzip openjdk25
COPY third_party /third_party
WORKDIR /third_party
RUN unzip verapdf-installer.zip
RUN sh -c '/bin/echo -e "1\n/third_party/verapdf\nO\n1\nN\nY\nN\nN\n1\nN\n" | ./verapdf-greenfield-1.30.2/verapdf-install'

FROM alpine:3.23 AS prod
RUN apk add --no-cache openjdk25
WORKDIR /borg/tools/verapdf
COPY --from=build /build/verapdf_api .
COPY --from=extract /third_party/verapdf ./third_party
CMD ["./verapdf_api"]
