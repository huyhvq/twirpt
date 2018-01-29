package main

import (
	"log"
	"net/http"
	pb "github.com/huyhvq/twirpt/rpc/haberdasher"
	"github.com/huyhvq/twirpt/internal/haberdasherserver"
)

func main() {
	//hook := statsd.NewStatsdServerHooks(LoggingStatter{os.Stderr})
	server := &haberdasherserver.Server{}
	twirpHandler := pb.NewHaberdasherServer(server, nil)
	log.Fatal(http.ListenAndServe(":8080", twirpHandler))
}
