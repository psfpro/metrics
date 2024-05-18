package agent

import (
	"errors"
	"log"
	"net"
)

func detectIPAddress() (net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println("Error getting network interfaces:", err)
		return nil, err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Println("Error getting addresses for interface:", i.Name, err)
			continue
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			// Check if the IP address is IPv4 (you can skip this check if you want both IPv4 and IPv6)
			if ip.To4() != nil {
				log.Printf("Interface: %s, IP: %s\n", i.Name, ip.String())
				return ip, nil
			}
		}
	}
	return nil, errors.New("no IP address found")
}
