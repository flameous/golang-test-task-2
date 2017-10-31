package main

import (
	"github.com/flameous/golang-test-task-2"
	"log"
	"flag"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	path := flag.String("path", "./", "path to configs")
	flag.Parse()

	e := task.NewtElasticClient(*path)
	server := task.NewServer(e)
	server.Run()
}
