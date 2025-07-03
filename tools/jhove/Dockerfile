FROM golang:alpine3.21 AS build
ARG BORG_VERSION=${BORG_VERSION}
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o jhove_api -ldflags "-X main.version=${BORG_VERSION}" ./cmd

FROM alpine:3.21 AS extract
RUN apk add --no-cache openjdk17
COPY third_party /third_party
WORKDIR /third_party
# install JHOVE with jar, enter sequence of characters required for the installation wizard
RUN sh -c '/bin/echo -e "1\njhove\nO\n1\nY\n1" | java -jar jhove_1.30.1.jar'

FROM alpine:3.21 AS prod
RUN apk add --no-cache openjdk17
WORKDIR /borg/tools/jhove
COPY --from=build /build/jhove_api .
COPY --from=extract /third_party/jhove ./jhove
CMD ["./jhove_api"]