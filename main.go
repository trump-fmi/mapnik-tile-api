// Copyright (C) {2017}  {Florian Barth florianbarth@gmx.de, Matthias Wagner matzew.mail@gmail.com, Marc Schubert marcschubert1@gmx.de}
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// The osm_label_server delivers an REST endpoint which clients that
// display openstreetmap data can use to obtain labels which can
// be displayed on top of the tiles from the tile-server. This enables
// clients to rotate the map without rotating the labels. It should be
// used as follows
package main

import (
	"encoding/json"
	"flag"
	"github.com/Terry-Mao/goconf"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

var renderdConfigPath string

func main() {
	var pPort int
	flag.IntVar(&pPort, "port", 8081, "Port where the server is reachable")
	flag.StringVar(&renderdConfigPath, "path", "/etc/renderd.conf", "Path to renderd.conf used to parse urls of tiles")
	flag.Parse()

	// Flag validation
	if pPort <= 0 || pPort > 65535 {
		log.Printf("Port not in allowed range. Cannot start with that configuration. Please use a free port out of [1, 65535].")
		return
	}

	log.Printf("Socket startup at :%d/tileEndpoints... ", pPort)
	mainRouter := mux.NewRouter()
	mainRouter.HandleFunc("/tileEndpoints", getTileEndpoints)

	// http timeout 15 s
	srv := &http.Server{
		Handler:      mainRouter,
		Addr:         ":" + strconv.Itoa(pPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

// method called when request for the tile endpoints URI is received
// parses the tile endpoints and returns them as JSON
func getTileEndpoints(w http.ResponseWriter, r *http.Request) {
	tileEndpoints, err := parseEndpoints(renderdConfigPath)
	if err != nil {
		log.Printf("Error during parsing: %s ", err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(tileEndpoints)
}

type tileEndpoint struct {
	Name        string `json:"name"`
	Uri         string `json:"uri"`
	Description string `json:"description"`
}

// parseEndpoints takes a path to a renderd config file and reads the
// endpoints out of it
func parseEndpoints(file string) ([]tileEndpoint, error) {
	var endpoints = make([]tileEndpoint, 0)
	conf := goconf.New()
	conf.Spliter = "="
	conf.Comment = ";"
	if err := conf.Parse(file); err != nil {
		return nil, err
	}

	for _, sectionName := range conf.Sections() {
		section := conf.Get(sectionName)
		uri, err := section.String("URI")
		if err != nil {
			continue
		}
		description, _ := section.String("DESCRIPTION")
		endpoints = append(endpoints, tileEndpoint{sectionName, uri, description})
	}

	return endpoints, nil
}
