package main

import "github.com/go_zipkin_example/service1"

func main() {
	server := service1.NewServer("opensesame")
	server.Start()
}
