BINS=build/gso_windows_x86_64.exe build/gso_linux_x86_64 build/gso_linux_arm64

all: $(BINS)

build/gso_windows_x86_64.exe:
	GOOS=windows GOARCH=amd64 go build -tags netgo -ldflags "-s" -o gso.exe
	@zip $@.zip gso.exe
	@rm gso.exe

build/gso_linux_x86_64:
	GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags "-s" -o gso
	@tar czf $@.tar.gz gso
	@rm gso

build/gso_linux_arm64:
	GOOS=linux GOARCH=arm64 go build -tags netgo -ldflags "-s" -o gso
	@tar czf $@.tar.gz gso
	@rm gso

clean:
	rm build/*
