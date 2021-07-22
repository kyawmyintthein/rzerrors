clean:
	rm proto/*.go

generate:
	/usr/bin/buf-Linux-x86_64 generate --path ./proto/errorpb/errors.proto