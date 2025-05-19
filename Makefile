.DEFAULT_GOAL := push

build:
	go build -o go-flv main.go

clean:
	rm -f go-flv

test:
	go test ./...

run: build
	./go-flv

push:
	git add .
	git commit -m "Update $(shell date +'%Y-%m-%d %H:%M:%S')."
	git pull
	git push origin main