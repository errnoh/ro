// Copyright 2013 errnoh. All rights reserved.
// Use of this source code is governed by
// MIT License that can be found in the LICENSE file.

// Sample text based Reittiopas application
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/errnoh/ro"
)

var route *string = flag.String("route", "", "Reitti (~/.ro tiedostosta)")
var reverse *bool = flag.Bool("reverse", false, "Paluureitti")
var from *string = flag.String("from", "", "Mistä? (\"Kauppakuja 1, helsinki\" for example")
var to *string = flag.String("to", "", "Mihin?")
var time *string = flag.String("time", "", "Aika (HHMM)")
var date *string = flag.String("date", "", "Päivämäärä (YYYYMMDD)")
var optimize *string = flag.String("optimize", "", "fastest|least_transfers|least_walking")
var limit *int = flag.Int("limit", 3, "Lukumäärä (max 5)")

func main() {
	flag.Parse()

	ro.SetCredentials("ro-term", "roTermPublic")

        if *route != "" {
            var ok bool
            if *from, *to, ok = ro.GetNamedRoute(*route); !ok {
                fmt.Printf("Reittiä '%s' ei löytynyt\n", *route)
                return
            }
        }

        if *reverse {
            *from, *to = *to, *from
        }

	start, ok := ro.GetLocation(*from)
	if !ok {
		fmt.Printf("Osoitetta '%s' ei löytynyt.\n", *from)
		return
	} else if len(start) > 1 {
		fmt.Printf("Löytyi %d vaihtoehtoa haulla '%s', käytetään osoitetta: '%s'\n",
			len(start),
			*from,
			address(start[0]),
		)
	}
	end, ok := ro.GetLocation(*to)
	if !ok {
		fmt.Printf("Osoitetta '%s' ei löytynyt.\n", *to)
		return
	} else if len(end) > 1 {
		fmt.Printf("Löytyi %d vaihtoehtoa haulla '%s', käytetään osoitetta: '%s'\n",
			len(end),
			*to,
			address(end[0]),
		)
	}
	routes := ro.GetRoute((start)[0], (end)[0], *time, *date, *optimize, *limit)
	fmt.Println()
	printRoutes(routes)
}

func address(location ro.Location) string {
	s := fmt.Sprintf("%s", location.Name)
	if location.Details.HouseNumber != 0 {
		s += fmt.Sprintf(" %d", int(location.Details.HouseNumber))
	}
	s += fmt.Sprintf(", %s", location.City)
	return s
}

type Route ro.Route

func (r Route) String() string {
	first := r.Legs[0].Locs[0]
	last := r.Legs[len(r.Legs)-1].Locs[len(r.Legs[len(r.Legs)-1].Locs)-1]
	duration := int(r.Duration)

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("Lähtö:%7s\n", first.DepTime[8:]))
	buffer.WriteString(fmt.Sprintf("Perillä:%5s\n", last.ArrTime[8:]))
	buffer.WriteString(fmt.Sprintf("Kesto:%4d:%02d\n", duration/3600, duration%3600/60))
	for _, v := range r.Legs {
		buffer.WriteString(fmt.Sprintf("%s\n", Legs(v)))
	}
	return buffer.String()
}

func printRoutes(routes ro.RouteResponse) {
	var buffer bytes.Buffer
	for i := 0; i < len(routes); i++ {
		buffer.WriteString(fmt.Sprintf("%s\n", Route(routes[i][0])))
	}
	fmt.Println(buffer.String())
}

type Legs ro.Legs

func (l Legs) String() string {
	line := (ro.Legs)(l).Line()
	start := l.Locs[0]
	end := l.Locs[len(l.Locs)-1]

	s := fmt.Sprintf("%20s %s -%5s - %s %20s", start.Name, start.DepTime[8:], line, end.ArrTime[8:], end.Name)
	return s
}
