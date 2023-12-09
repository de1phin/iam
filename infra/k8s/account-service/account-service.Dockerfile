FROM golang:alpine AS build
RUN mkdir -p /tmp/build/
WORKDIR /tmp/build
COPY . .
RUN go mod tidy
RUN go build -o /usr/bin/account-service "./services/account/cmd/account-service/main.go"
ARG CONFIGPATH=./services/account/cmd/account-service/config.yaml
RUN mkdir -p /etc/account-service
COPY ${CONFIGPATH} /etc/account-service/config.yaml

FROM alpine:latest
COPY --from=build /usr/bin/account-service /usr/bin/account-service
COPY --from=build /etc/account-service /etc/account-service
CMD /usr/bin/account-service --config /etc/account-service/config.yaml
