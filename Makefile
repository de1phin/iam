proto = services/account/api

.PHONY: generate, $(proto)

generate: $(proto)

$(proto):
	$(eval PATH_OUT := ./genproto)
	if [[ ! -d '${PATH_OUT}' ]]; then \
		mkdir -p '${PATH_OUT}'; \
	fi; \
	protoc --go_out=${PATH_OUT} --go_opt=paths=source_relative \
		--go-grpc_out=${PATH_OUT} --go-grpc_opt=paths=source_relative \
		$(shell find ./$@ -iname "*.proto")

