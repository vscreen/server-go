
build:
	go build

pi:
	GOOS=linux GOARCH=arm go build -tags pi