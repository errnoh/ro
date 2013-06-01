// Copyright 2013 errnoh. All rights reserved.
// Use of this source code is governed by
// MIT License that can be found in the LICENSE file.

// Small wrapper for Reittiopas API ( http://developer.reittiopas.fi/pages/en/http-get-interface-version-2.php )
package ro

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)


var (
	api_user, api_pass string

        // You can save routes in a file named .ro inside your home folder.
        // Each row should be formatted as: name,source,destination
        // Escape strings with quotes if needed (if there are any commas).
        // Example row:
        // keskustaan,"ratapihantie 1, helsinki",rautatientori
        routefile = fmt.Sprintf("%s/.ro", os.Getenv("HOME"))
)

func SetCredentials(user, pass string) {
	api_user = user
	api_pass = pass
}

func SetRoutefile(file string) {
        routefile = file
}

func GetLocation(loc string) (LocationResponse, bool) {
	if loc == "" {
		return nil, false
	}

	var resp LocationResponse

	s := fmt.Sprintf("request=geocode&key=%s", url.QueryEscape(loc))
	ok := Get(s, &resp)
	if !ok {
		return nil, false
	}

	return resp, true
}

func GetRoute(from Location, to Location, time string, date string, optimize string, limit int) RouteResponse {
	var (
		resp   RouteResponse
		buffer bytes.Buffer
	)

	buffer.WriteString(fmt.Sprintf("request=route&from=%s&to=%s", from.Coords, to.Coords))

	if time != "" {
		buffer.WriteString(fmt.Sprintf("&time=%s", time))
	}
	if date != "" {
		buffer.WriteString(fmt.Sprintf("&date=%s", date))
	}
	if optimize != "" {
		buffer.WriteString(fmt.Sprintf("&optimize=%s", optimize))
	}
	if limit != 3 && limit >= 0 && limit <= 5 {
		buffer.WriteString(fmt.Sprintf("&show=%d", limit))
	}

	Get(buffer.String(), &resp)

	resp.fixStartEndLocations(from, to)

	return resp
}

// Add route start & end locations to the route information instead of displaying empty fields.
func (r RouteResponse) fixStartEndLocations(from Location, to Location) {
	for i := 0; i < len(r); i++ {
		r[i][0].Legs[0].Locs[0].Name = from.Name
		r[i][0].Legs[len(r[i][0].Legs)-1].Locs[len(r[i][0].Legs[len(r[i][0].Legs)-1].Locs)-1].Name = to.Name
	}
}

func Get(s string, v interface{}) bool {
	url := fmt.Sprintf("http://api.reittiopas.fi/hsl/prod/?user=%s&pass=%s&%s", api_user, api_pass, s)

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return false
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		if string(b) != "" {
			log.Println(string(b))
		}
		return false
	}

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		// These aren't fatal errors most of the time.
		log.Printf(err.Error())
	}

	return true
}

func (l Legs) Line() string {
	switch l.Type {
	case "1", "2", "3", "4":
		return l.Code[2:][:4]
	case "5", "8":
		return l.Code[1:][:4]
	case "6":
		return "metro"
	case "12":
		return l.Code[4:5]
	case "walk":
		return l.Type
	}
	return l.Code
}

func GetNamedRoute(route string) (from, to string, ok bool) {
	file, err := os.Open(routefile)
	if err != nil {
                fmt.Println(err)
		return
	}
	defer file.Close()

	r := csv.NewReader(file)
	for row, err := r.Read(); err == nil; row, err = r.Read() {
		if len(row) != 3 {
			continue
		}
		if row[0] == route {
			from = row[1]
			to = row[2]
			ok = true
			break
		}
	}
	return
}
