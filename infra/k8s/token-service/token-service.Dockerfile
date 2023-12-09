FROM golang:alpine AS build
RUN mkdir -p /tmp/build/
WORKDIR /tmp/build
COPY . .
RUN go mod tidy
RUN go build -o /usr/bin/token-service "./services/token/cmd/main.go"
ARG CONFIGPATH=./services/token/cmd/config.yaml
RUN mkdir -p /etc/token-service
COPY ${CONFIGPATH} /etc/token-service/config.yaml

FROM alpine:latest
COPY --from=build /usr/bin/token-service /usr/bin/token-service
COPY --from=build /etc/token-service /etc/token-service
CMD /usr/bin/token-service --config /etc/token-service/config.yaml
