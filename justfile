# proto-gen -> generate protobuf
proto-gen:
  @buf generate

# swagger-gen -> generate swagger api docs
swagger-gen:
	@echo "unimplemented"

# build -> build application
build:
	@go build -o main ./cmd

# run -> application
run:
	@./main

# dev -> run build then run it
dev: 
	@watchexec -r -c -e go -- just build run