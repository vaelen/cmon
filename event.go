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
	"fmt"
	"net"
)

type State int

const (
	ACTIVE       State = iota
	ESTABLISHED
	OTHER
)

type Event struct {
	SrcIP net.IP
	SrcPort uint16
	DstIP net.IP
	DstPort uint16
	State State
}

func (e *Event) String() string {
	dstIP := ""
	if e.DstIP != nil {
		dstIP = fmt.Sprintf("%s", e.DstIP)
	}
	srcIP := ""
	if e.SrcIP != nil {
		srcIP = fmt.Sprintf("%s", e.SrcIP)
	}
	return fmt.Sprintf("%s:%d -> %s:%d", srcIP, e.SrcPort, dstIP, e.DstPort)
}
