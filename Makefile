.PHONY: generate

proto = services/account/proto

generate: $(proto)

$(proto):
	if [[ ! -d './pkg/gen/$@' ]]; then \
		mkdir -p './pkg/gen/$@'; \
	fi; \
	protoc --go_out=./pkg/gen/$@ --go_opt=paths=source_relative \
		--go-grpc_out=./pkg/gen/$@ --go-grpc_opt=paths=source_relative \
		$(shell find ./pkg/$@ -iname "*.proto")
		
	

