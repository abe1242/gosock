BINS=build/gsoc_windows_x86_64.exe build/gsoc_linux_x86_64 build/gsoc_linux_arm64

all: $(BINS)

build/gsoc_windows_x86_64.exe:
	GOOS=windows GOARCH=amd64 go build -tags netgo -ldflags "-s" -o gsoc.exe
	@zip $@.zip gsoc.exe
	@rm gsoc.exe

build/gsoc_linux_x86_64:
	GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags "-s" -o gsoc
	@tar czf $@.tar.gz gsoc
	@rm gsoc

build/gsoc_linux_arm64:
	GOOS=linux GOARCH=arm64 go build -tags netgo -ldflags "-s" -o gsoc
	@tar czf $@.tar.gz gsoc
	@rm gsoc

clean:
	rm build/*
