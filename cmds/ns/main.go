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

package main

import (
	"log"
	"github.com/vaelen/cmon"
)

func main() {
	mode := cmon.FindMode()
	conns := cmon.Connections(mode)
	for _, e := range conns {
		log.Printf("%d %s\n", e.State, e.String())
	}
}
