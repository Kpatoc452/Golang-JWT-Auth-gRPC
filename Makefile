.PHONNY:
run:
	@go run -gcflags="all=-lang=go1.23" ./cmd/auth --config=./config/local.yaml 