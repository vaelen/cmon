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
	"os"
	"sync"
	"time"
)

const ChannelBufferSize int = 100

func Monitor(ifNames []string, filter string, printInterval time.Duration, verbose bool) {

	wg := sync.WaitGroup{}

	// Start the sniffers

	sources := make([]SnifferConfig, 0)

	for _, ifName := range ifNames {
		c := SnifferConfig{
			IfName:   ifName,
			Filter:   filter,
			Request:  make(chan bool),
			Response: make(chan *Stats),
			Event:    make(chan Event, ChannelBufferSize),
			Verbose:  verbose,
		}

		sources = append(sources, c)

		log.Printf("Starting Sniffer for Interface %s, Filter: %s\n", c.IfName, c.Filter)

		wg.Add(1)
		go func(c SnifferConfig) {
			defer wg.Done()
			packetSniffer(c)
		}(c)

		wg.Add(1)
		go func(c SnifferConfig) {
			defer wg.Done()
			counter(c)
		}(c)
	}

	// Start the printer

	wg.Add(1)
	go func(s []SnifferConfig, p time.Duration) {
		defer wg.Done()
		printer(PrinterConfig{
			Logger:        log.New(os.Stdout, "", 0),
			Sources:       s,
			PrintInterval: p,
		})
	}(sources, printInterval)

	wg.Wait()

}
