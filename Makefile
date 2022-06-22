build:
	go build -o build/sdunetd

run:
	go run .

build-all:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags=all="-s -w" -o build/sdunetd-linux-arm
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags=all="-s -w" -o build/sdunetd-linux-arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags=all="-s -w" -o build/sdunetd-linux-386
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags=all="-s -w" -o build/sdunetd-linux-amd64
	CGO_ENABLED=0 GOMIPS=softfloat GOOS=linux GOARCH=mips go build -ldflags=all="-s -w" -o build/sdunetd-linux-mips-softfloat
	CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags=all="-s -w" -o build/sdunetd-linux-mips64
	CGO_ENABLED=0 GOMIPS=softfloat GOOS=linux GOARCH=mipsle go build -ldflags=all="-s -w" -o build/sdunetd-linux-mipsle-softfloat
	CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags=all="-s -w" -o build/sdunetd-linux-mips64le
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags=all="-s -w" -o build/sdunetd-linux-arm
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags=all="-s -w" -o build/sdunetd-linux-arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags=all="-s -w" -o build/sdunetd-windows-amd64.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags=all="-s -w" -o build/sdunetd-windows-386.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags=all="-s -w" -o build/sdunetd-windows-arm64.exe
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags=all="-s -w" -o build/sdunetd-darwin-amd64
	CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags=all="-s -w" -o build/sdunetd-freebsd-amd64

upx: build-all
	upx --best --ultra-brute build/*

clean:
	rm -r build/

all: build-all test
