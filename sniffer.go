/*
   This file is part of CMon.

   Copyright 2017, Andrew Young <andrew@vaelen.org>

   CMon is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   CMon is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with CMon.  If not, see <http://www.gnu.org/licenses/>.
*/

package cmon

import (
	"log"
	"reflect"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type SnifferConfig struct {
	IfName   string
	Filter   string
	Request  chan bool
	Response chan *Stats
	Event    chan Event
	Verbose  bool
}

func packetSniffer(c SnifferConfig) {
	if handle, err := pcap.OpenLive(c.IfName, 1600, true, pcap.BlockForever); err != nil {
		log.Panicf("Couldn't Open Interface %s, Error: %s\n", c.IfName, err.Error())
	} else if err := handle.SetBPFFilter(c.Filter); err != nil { // optional
		log.Panicf("Couldn't Set Filter on Interface %s, Filter: %s, Error: %s\n", c.IfName, c.Filter, err.Error)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		if c.Verbose {
			log.Printf("Sniffer Started for %s\n", c.IfName)
			defer log.Printf("Sniffer Finished for %s\n", c.IfName)
		}
		for packet := range packetSource.Packets() {
			handlePacket(c, packet) // Do something with a packet here.
		}
	}
}

func handlePacket(c SnifferConfig, packet gopacket.Packet) {
	e := Event{}
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		// TCP Packet
		tcp, _ := tcpLayer.(*layers.TCP)
		e.SrcPort = AsUint16(tcp.SrcPort)
		e.DstPort = AsUint16(tcp.DstPort)
	}

	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		// UDP Packet
		udp, _ := udpLayer.(*layers.UDP)
		e.SrcPort = AsUint16(udp.SrcPort)
		e.DstPort = AsUint16(udp.DstPort)
	}

	if ipv4Layer := packet.Layer(layers.LayerTypeIPv4); ipv4Layer != nil {
		// IPv4 Packet
		ipv4, _ := ipv4Layer.(*layers.IPv4)
		e.SrcIP = ipv4.SrcIP
		e.DstIP = ipv4.DstIP
	}

	if ipv6Layer := packet.Layer(layers.LayerTypeIPv6); ipv6Layer != nil {
		// IPv6 Packet
		ipv6, _ := ipv6Layer.(*layers.IPv6)
		e.SrcIP = ipv6.SrcIP
		e.DstIP = ipv6.DstIP
	}
	//log.Printf("Event Sent on %s: %s\n", s.IfName, e.String());
	c.Event <- e
}

func AsUint16(val interface{}) uint16 {
	ref := reflect.ValueOf(val)
	if ref.Kind() != reflect.Uint16 {
		return 0
	}
	return uint16(ref.Uint())
}
