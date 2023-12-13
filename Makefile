proto = services/account/api services/token/api services/access/api

.PHONY: generate cli $(proto)

generate: $(proto)

$(proto):
	$(eval PATH_OUT := ./genproto/services)
	if [[ ! -d '${PATH_OUT}' ]]; then \
		mkdir -p '${PATH_OUT}'; \
	fi; \
	protoc -I "./pkg/proto" -I "./services" --go_out=${PATH_OUT} --go_opt=paths=source_relative \
		--go-grpc_out=${PATH_OUT} --go-grpc_opt=paths=source_relative \
		--swagger_out=${PATH_OUT} \
		--swagger_opt=logtostderr=true \
        --grpc-gateway_out ${PATH_OUT} --grpc-gateway_opt paths=source_relative \
		$(shell find ./$@ -iname "*.proto")

test:
	go test ./...

cli:
	go build -o ./iamcli ./cli/cmd/main.go
