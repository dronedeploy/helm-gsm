.PHONY: clean macos linux

macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/helm-gsm-darwin

linux:
	GOOS=linux GOARCH=amd64 go build -o bin/helm-gsm-linux

clean:
	rm bin/*
