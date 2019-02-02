
build:
	go build

pi:
	GOOS=linux GOARCH=arm GOARM=7 go build -tags pi
