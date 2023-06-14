package main

import (
	"flag"

	"github.com/natac13/ggcache/cache"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3000", "listen address")
	leaderAddr := flag.String("leaderaddr", "", "leader address")

	flag.Parse()

	opts := ServerOpts{
		ListenAddr: *listenAddr,
		IsLeader:   len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	server := NewServer(opts, cache.New())
	server.Start()
}
