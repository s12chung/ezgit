build:
	go install .
	GOOS=windows GOARCH=amd64 go build -o bin/ezgit.exe .

goimports:
	goimports -w .
