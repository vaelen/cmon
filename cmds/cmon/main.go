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

package main

import (
	"os"
	"log"
	"net"
	"strings"
	"time"
	"github.com/vaelen/cmon"
)

const DefaultFilter string = ""
const DefaultPrintInterval = time.Minute * 1

func main() {
	log.SetOutput(os.Stderr)
	
	var ifNames []string
	var filter string
	var printInterval time.Duration
	var verbose bool

	filter = DefaultFilter
	printInterval = DefaultPrintInterval
	
	if len(os.Args) > 1 {
		ifNames = strings.Split(os.Args[1], ",")
		if len(os.Args) > 2 {
			filter = os.Args[2]
		}
	} else {
		interfaces, err := net.Interfaces()
		if err != nil {
			log.Panicf("Couldn't Retrieve Interface List.  Error: %s\n", err.Error())
		}
		ifNames = make([]string, 0)
		for _, i := range interfaces {
			ifNames = append(ifNames, i.Name)
		}
	}

	cmon.Monitor(ifNames, filter, printInterval, verbose)
}
