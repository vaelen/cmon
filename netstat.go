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
	"os/exec"
	//"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Mode int

const (
	SS Mode = iota
	NETSTAT
)

func FindMode() Mode {
	cmd := exec.Command("ss")
	_, err := cmd.Output()
	if err == nil {
		return SS
	}
	cmd = exec.Command("netstat")
	_, err2 := cmd.Output()
	if err2 != nil {
		log.Panicf("Need either ss or netstat.\n%s\n%s", err.Error(), err2.Error())
	}
	return NETSTAT
}

func Connections(mode Mode) []Event {
	var cmd *exec.Cmd
	var minFields, stateField, srcField, dstField int
	
	switch mode {
	case SS:
		cmd = exec.Command("ss", "-tn")
		minFields = 4
		stateField = 0
		srcField = 3
		dstField = 4
	default:
		cmd = exec.Command("netstat", "-tn")
		minFields = 5
		stateField = 5
		srcField = 3
		dstField = 4
	}
	
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Could not retrieve active port list.  Command: %s, Error: %s\n", cmd, err)
	}
	
	conns := make([]Event, 0)
	for _, s := range strings.Split(string(out), "\n") {
		fields := strings.Fields(s)
		if len(fields) > minFields {
			stateString := fields[stateField]
			src := fields[srcField]
			dst := fields[dstField]
		
			srcIP, srcPort, _ := net.SplitHostPort(src)
			dstIP, dstPort, _ := net.SplitHostPort(dst)
			state := OTHER
			if len(stateString) > 4 && stateString[:5] == "ESTAB" {
				state = ESTABLISHED
			}

			sp, _ := strconv.Atoi(srcPort)
			dp, _ := strconv.Atoi(dstPort)
		
			e := Event{
				SrcIP: net.ParseIP(srcIP),
				SrcPort: uint16(sp),
				DstIP: net.ParseIP(dstIP),
				DstPort: uint16(dp),
				State: state,
			}
			conns = append(conns, e)
		}
	}
	return conns
}
