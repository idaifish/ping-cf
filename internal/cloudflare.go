//go:generate curl -o ips-v4.txt https://www.cloudflare.com/ips-v4
package internal

import (
	_ "embed"
	"encoding/binary"
	"log"
	"math/rand"
	"net"
	"strings"
)

//go:embed ips-v4.txt
var CIDRs string

// GetIP obtains all the CloudFlare IP addresses.
func GetIP() chan net.IP {
	ipChan := make(chan net.IP)

	go func(chan net.IP) {
		defer close(ipChan)

		cidrs := strings.Split(CIDRs, "\n")
		_ipChan := make(chan net.IP)
		var total int
		for _, cidr := range cidrs {
			_, ipnet, err := net.ParseCIDR(cidr)
			if err != nil {
				log.Fatalln(err)
			}

			startIP := binary.BigEndian.Uint32(ipnet.IP)
			mask := binary.BigEndian.Uint32(ipnet.Mask)
			endIP := (startIP & mask) | (mask ^ 0xffffffff)
			length := endIP - startIP
			total += int(length)

			go func(start, l uint32, ipChan chan<- net.IP) {
				for _, v := range rand.Perm(int(l)) {
					ip := make(net.IP, 4)
					binary.BigEndian.PutUint32(ip, start+uint32(v))
					ipChan <- ip
				}
			}(startIP, length, _ipChan)
		}

		for i := 0; i < total; i++ {
			ipChan <- <-_ipChan
		}
	}(ipChan)

	return ipChan
}
