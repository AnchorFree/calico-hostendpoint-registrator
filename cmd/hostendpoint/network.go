package main

import (
	"net"
)

func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

// GetNetworkConfig returns the name of the default network interface of
// the node, and its' primary IP address with the CIDR mask.
func GetNetworkConfig() (networkInterface, primaryIP string) {

	ifaces, _ := net.Interfaces()
	outboundIP, err := getOutboundIP()
	if err == nil {
		for _, i := range ifaces {
			addrs, _ := i.Addrs()
			for _, addr := range addrs {
				var ip net.IP
				var nt *net.IPNet
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
					nt = v
				case *net.IPAddr:
					ip = v.IP
				}
				if ip.Equal(outboundIP) {
					networkInterface = i.Name
					primaryIP = nt.String()
				}
			}
		}
	}
	return networkInterface, primaryIP

}
