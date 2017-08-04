build: 
	go build -o ./build/service1 ./service1/cmd/main.go
	go build -o ./build/service2 ./service2/cmd/main.go
	go build -o ./build/cli ./cli/main.go
clean:
	rm -rf ./build
