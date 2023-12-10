FROM golang:alpine AS build
RUN mkdir -p /tmp/build/
WORKDIR /tmp/build
COPY . .
RUN go mod tidy
RUN go build -o /usr/bin/access-service "./services/access/cmd/main.go"
ARG CONFIGPATH=./services/access/cmd/config.yaml
RUN mkdir -p /etc/access-service
COPY ${CONFIGPATH} /etc/access-service/config.yaml

FROM alpine:latest
COPY --from=build /usr/bin/access-service /usr/bin/access-service
COPY --from=build /etc/access-service /etc/access-service
CMD /usr/bin/access-service --config /etc/access-service/config.yaml
