package tcp_packets

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)

type IPstat struct {
	IP            string
	All           int
	Retransmitted int
}

type SequenceConfig struct {
	n         uint32
	ticker    *time.Ticker
	interrupt chan bool
}

const (
	// The same default as tcpdump.
	defaultSnapLen = 262144
	rto            = 1000 * time.Millisecond
)

func Analyze(c chan *IPstat, device string) {
	handle, err := pcap.OpenLive(device, defaultSnapLen, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}

	defer handle.Close()

	packets := gopacket.NewPacketSource(handle, handle.LinkType()).Packets()

	ipMap := make(map[string]*IPstat)
	seqMap := make(map[uint32]*SequenceConfig)

	for pkt := range packets {
		tcpLayer := pkt.Layer(layers.LayerTypeTCP)
		networkLayer := pkt.Layer(layers.LayerTypeIPv4)

		if tcpLayer != nil && networkLayer != nil {
			tcp, ok := tcpLayer.(*layers.TCP)

			// stop ticker and skip
			if seq, ok := seqMap[tcp.Seq]; ok {
				seq.interrupt <- true

				delete(seqMap, tcp.Seq)
			}

			if !ok {
				log.Printf("TCP layer decoding error: SEQ - %d", tcp.Seq)
			}

			ip, ok := networkLayer.(*layers.IPv4)
			if !ok {
				log.Printf("IP layer decoding error: IP - %d", ip.DstIP)
			}

			if ip.DstIP.IsPrivate() {
				continue
			}

			if _, ok := ipMap[ip.DstIP.String()]; !ok {
				ipMap[ip.DstIP.String()] = &IPstat{
					IP: ip.DstIP.String(),
				}
			}

			targetIP := ipMap[ip.DstIP.String()]
			targetIP.All++

			if len(tcp.Payload) > 0 {
				seqMap[getNextSequence(tcp)] = &SequenceConfig{
					n:         tcp.Seq,
					ticker:    time.NewTicker(rto),
					interrupt: make(chan bool, 1),
				}
				targetSeq := seqMap[getNextSequence(tcp)]

				//defer targetSeq.ticker.Stop()

				go func(s *SequenceConfig, i *IPstat, c chan *IPstat) {
					for {
						select {
						case <-s.ticker.C:
							i.Retransmitted++
							c <- targetIP
							return
						case <-s.interrupt:
							return
						}
					}
				}(targetSeq, targetIP, c)
			}
		}
	}
}

// it takes the current sequence number
// (specified usually just in front of the "next sequence number" field)
// and adds the tcp payload length
func getNextSequence(packet *layers.TCP) uint32 {
	return packet.Seq + uint32(len(packet.Payload))
}
