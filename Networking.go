/*******************************************************
 * Copyright (C) 2018 Palanjyan Zhorzhik
 *
 * This file is part of WS-Discovery project.
 *
 * WS-Discovery can be copied and/or distributed without the express
 * permission of Palanjyan Zhorzhik
 *******************************************************/

package WS_Discovery

import (
	"fmt"
	"net"
	"golang.org/x/net/ipv4"
	"time"
	"log"
)

const bufSize  = 8192

func sendUDPMulticast (msg string, interfaceName string) []string {
	var result []string
	data := []byte(msg)
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		fmt.Println(err)
	}
	group := net.IPv4(239, 255, 255, 250)

	c, err := net.ListenPacket("udp4", "0.0.0.0:1024")
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()

	p := ipv4.NewPacketConn(c)
	if err := p.JoinGroup(iface, &net.UDPAddr{IP: group}); err != nil {
		fmt.Println(err)
	}

	dst := &net.UDPAddr{IP: group, Port: 3702}
	for _, ifi := range []*net.Interface{iface} {
		if err := p.SetMulticastInterface(ifi); err != nil {
			fmt.Println(err)
		}
		p.SetMulticastTTL(2)
		if _, err := p.WriteTo(data, nil, dst); err != nil {
			fmt.Println(err)
		}
	}

	if err := p.SetReadDeadline(time.Now().Add(time.Second * 1)); err != nil {
		log.Fatal(err)
	}

	for {
		b := make([]byte, bufSize)
		n, _, _, err := p.ReadFrom(b)
		if err != nil {
			fmt.Println(err)
			break
		}
		result = append(result, string(b[0:n]))
	}
	return result
}