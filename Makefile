# Build all executables.
all:
	go build -o . ./...

clean:
	$(RM) godotenvvault
