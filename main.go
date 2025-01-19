package main

import (
	"log"
	"os"

	"traffic_analyzer/internal/db"
	"traffic_analyzer/internal/tcp_packets"
	"traffic_analyzer/internal/web"
)

func main() {
	ipChan := make(chan *tcp_packets.IPstat)

	dbClient := db.New()

	go tcp_packets.Analyze(ipChan, os.Args[1])

	srv := web.New()
	go srv.Start()

	for {
		select {
		case ip := <-ipChan:
			if dbClient.Get(ip.IP) != nil {
				err := dbClient.Update(ip.IP, ip.All, ip.Retransmitted)
				if err != nil {
					log.Println(err)
				}
			} else {
				err := dbClient.Insert(ip.IP, ip.All, ip.Retransmitted)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
