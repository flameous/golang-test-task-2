package main

import (
	"github.com/flameous/golang-test-task-2"
	"log"
	"os"
	"fmt"
	"flag"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	fmt.Println(os.Args)

	path := flag.String("path", "./", "path to configs")
	flag.Parse()

	ctx, client := task.InitElasticClient(*path)
	server := task.NewServer(ctx, client)
	server.Run()
}
