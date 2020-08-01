bin/helm-gsm-darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/helm-gsm-darwin

bin/helm-gsm-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/helm-gsm-linux
