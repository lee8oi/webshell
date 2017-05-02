/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

//
package main

import (
	"strings"
)

type command struct {
	Desc    string
	Handler func([]string) string
}

var commands = make(map[string]command)

func init() {
	commands["echo"] = command{
		Desc: "Echos the input args back to the shell.",
		Handler: func(args []string) (s string) {
			if len(args) > 1 {
				s = strings.Join(args[1:], " ")
			} else {
				s = "echoing nothing!"
			}
			return
		},
	}
	commands["help"] = command{
		Desc: "Returns help information about available commands.",
		Handler: func(args []string) (s string) {
			if len(args) > 0 {
				if len(args) == 1 {
					cmds := ""
					for k := range commands {
						cmds += " " + k
					}
					s = "Available commands:" + cmds
				} else {
					if cmd, ok := commands[args[1]]; ok {
						s = "(" + args[1] + ") " + cmd.Desc
					} else {
						s = "Command not available: " + args[1]
					}
				}
			}
			return
		},
	}
}
