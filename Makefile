build:
	go build -o spot-look-back

build-arm:
	GOARM=7 GOARCH=arm GOOS=linux go build

migrate:
	cd migrations
	goose postgres "$(db)" up
	cd ..
