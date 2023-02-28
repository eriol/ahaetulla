.PHONY: build clean

# build ahaetulla
build:
	@go build ./cmd/ahaetulla

# remove artifacs
clean:
	@rm -rf ahaetulla
