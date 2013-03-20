// Copyright 2013 errnoh. All rights reserved.
// Use of this source code is governed by
// MIT License that can be found in the LICENSE file.

package ro

type LocationResponse []Location

type Location struct {
	LocType     string
	LocTypeID   float64
	Name        string
	MatchedName string
	Lang        string
	City        string
	Coords      string
	Details     *Details
}

type Details struct {
	HouseNumber float64
	// more?
}

type RouteResponse [][]Route

type Route struct {
	Length   float64
	Duration float64
	Legs     []Legs
}

type Legs struct {
	Length   float64
	Duration float64
	Type     string
	Code     string
	Locs     []Stop
	Shape    string
}

type Stop struct {
	Coord       Coord
	ArrTime     string
	DepTime     string
	Name        string
	Code        float64
	ShortCode   string
	StopAddress string
}

type Coord struct {
	X, Y float64
}
