/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

//
package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	flag.Parse()
	cfg := loadConfig(*cfgPath)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(cfg.PubDir))))
	http.HandleFunc("/ajax", ajaxHandler)
	go func() {
		err := http.ListenAndServeTLS(":"+cfg.HTTPSPort, cfg.CertPem, cfg.KeyPem, nil)
		if err != nil {
			log.Fatal("ListenAndServeTLS:", err)
		}
	}()
	err := http.ListenAndServe(":"+cfg.HTTPPort, http.RedirectHandler("https://"+cfg.Domain+":"+cfg.HTTPSPort, 301))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var cfgPath = flag.String("config", "config.json", "path to config file (in JSON format)")

// config type contains the necessary server configuration strings.
type config struct {
	HTTPPort, HTTPSPort,
	Domain, PubDir, CertPem, KeyPem string
}

// ajaxHandler processes the input text sent via ajax.
func ajaxHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	text, args := "", getArgs(body)
	if len(args) > 0 {
		if cmd, exists := commands[strings.ToLower(args[0])]; exists {
			text = cmd.Handler(args)
		} else {
			text = (args[0] + ": Command not found.")
		}
	}
	io.WriteString(w, text)
}

// loadConfig loads configuration values from file.
func loadConfig(path string) (c config) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// getArgs splits a slice of bytes into a slice of string arguments.
// Anything in '', "", or `` are consider a single argument.
func getArgs(b []byte) (s []string) {
	re := regexp.MustCompile("`([\\S\\s]*)`|('([\\S \\t\\r]*)'|\"([\\S ]*)\"|\\S+)")
	args := re.FindAllSubmatch(b, -1)
	for _, val := range args {
		s = append(s, string(val[0]))
	}
	return
}
