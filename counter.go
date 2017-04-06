/*
    This file is part of CMon.

    Copyright 2017, Andrew Young <andrew@vaelen.org>

    CMon is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    Foobar is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with Foobar.  If not, see <http://www.gnu.org/licenses/>.
*/

package cmon

import (
	"log"
)

type Stats struct {
	Counts map[string] uint64
}

func newStats() *Stats {
	return &Stats {
		Counts: make(map[string]uint64),
	}
}

func counter(s SnifferConfig) {
	if s.Verbose {
		log.Printf("Counter Started for %s\n", s.IfName)
		defer log.Printf("Counter Finished for %s\n", s.IfName)
	}
	stats := newStats()
	for {
		select {
		case e := <-s.Event:
			k := e.String()
			if (s.Verbose) {
				log.Printf("Packet Received on %s: %s\n", s.IfName, k);
			}
			i, ok := stats.Counts[k]
			if !ok {
				i = 0
			}
			i++
			stats.Counts[k] = i
		case _ = <-s.Request:
			s.Response <- stats
			stats = newStats()
		}
	}
}
