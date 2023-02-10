# build freakble
build:
    @cd cmd/freakble && go build && mv freakble ../..

# remove artifacs
clean:
    @rm -rf freakble
