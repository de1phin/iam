proto = services/token/api

.PHONY: generate, $(proto)

generate: $(proto)

$(proto):
	$(eval PATH_OUT := ./genproto)
	if [[ ! -d '${PATH_OUT}' ]]; then \
		mkdir -p '${PATH_OUT}'; \
	fi; \
	protoc -I ./services \
	    --go_out=${PATH_OUT} --go_opt=paths=source_relative \
		--go-grpc_out=${PATH_OUT} --go-grpc_opt=paths=source_relative \
		--swagger_out=${PATH_OUT} \
		--swagger_opt=logtostderr=true \
       	--grpc-gateway_out ${PATH_OUT} --grpc-gateway_opt paths=source_relative \
		$(shell find ./$@ -iname "*.proto")

test:
	go test ./...

gen-token-service:
	protoc -I ./services \
       --go_out ./genproto/services --go_opt paths=source_relative \
       --go-grpc_out ./genproto/services --go-grpc_opt paths=source_relative \
       --swagger_out=./genproto/services \
	   --swagger_opt=logtostderr=true \
       --grpc-gateway_out ./genproto/services --grpc-gateway_opt paths=source_relative \
       ./services/token/api/token-service.proto