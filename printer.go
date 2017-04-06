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
	"time"
)

type PrinterConfig struct {
	Logger *log.Logger
	Sources []SnifferConfig
	PrintInterval time.Duration
}

func printer(c PrinterConfig) {
	mode := FindMode()
	lastConns := connectionStrings(Connections(mode))
	for {
		time.Sleep(c.PrintInterval)
		conns := connectionStrings(Connections(mode))
		active := make(map[string]bool)

		c.Logger.Printf("# %s\n", time.Now().Format(time.RFC1123Z))
		c.Logger.Println()
		
		// Print Active Connections
		c.Logger.Println("## Active Connections")
		c.Logger.Printf("| %20s | %60s | %20s | %6s | \n", "Interface", "Source -> Destination", "Count", "State")
		c.Logger.Printf("| -------------------- | ------------------------------------------------------------ | -------------------- | ------ |\n")
		for _, s := range c.Sources {
			s.Request <- true
			stats := <-s.Response
			for k,v := range stats.Counts {
				active[k] = true
				n := ""
				_, last := lastConns[k]
				_, curr := conns[k]
				if !last && curr {
					n = "New"
				} else if !curr {
					n = "Closed"
				} else {
					n = "Active"
				}
				c.Logger.Printf("| %20s | %60s | %20d | %6s |\n", s.IfName, k, v, n)
			}
		}
		c.Logger.Println()

		// Print Inactive Connections
		c.Logger.Println("## Inactive Connections")
		c.Logger.Printf("| %60s | %8s | \n", "Source -> Destination", "State")
		c.Logger.Printf("| ------------------------------------------------------------ | -------- |\n")
		for k, _ := range lastConns {
			_, act  := active[k]
 			if !act {
				n := ""
				_, curr := conns[k]
				if !curr {
					n = "Closed"
				} else {
					n = "Inactive"
				}
				c.Logger.Printf("| %60s | %8s |\n", k, n)
			}
			for k, _ := range conns {
				_, act  := active[k]
				_, last := lastConns[k]
				if !act && !last {
					n := "New"
					c.Logger.Printf("| %60s | %8s |\n", k, n)
				}
			}
		}

		c.Logger.Println()

	}
}

func connectionStrings(c []Event) map[string]bool {
	r := make(map[string]bool)
	for _, e := range c {
		r[e.String()] = true
	}
	return r
}
