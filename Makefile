proto = services/account/api

.PHONY: generate, $(proto)

generate: $(proto)

$(proto):
	if [[ ! -d './genproto/$@' ]]; then \
		mkdir -p './genproto/$@'; \
	fi; \
	protoc --go_out=./genproto/$@ --go_opt=paths=source_relative \
		--go-grpc_out=./genproto/$@ --go-grpc_opt=paths=source_relative \
		$(shell find ./$@ -iname "*.proto")

