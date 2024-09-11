all: cross native

# the CGO_ENABLED=0
# is used to not rely on libc, to use
# in on debian for the host, as it has different libc.

native:
	env CGO_ENABLED=0 go build -o bmc-test-x86
cross:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o bmc-test-armv5
	ln -s -f bmc-test-armv5 bmc-test-cross

clean:
	rm -f bmc-test-cross bmc-test-armv5 bmc-test-x86
