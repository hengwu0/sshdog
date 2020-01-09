// Copyright 2016 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// TODO: High-level file comment.
package main

import (
	"flag"
	"fmt"
	"github.com/hengwu0/sshdog/daemon"
	"github.com/hengwu0/sshdog/proc"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage1: %s\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "write files in `config` dir to make configuration:\n")
	fmt.Fprintf(os.Stderr, "filename:port\n")
	fmt.Fprintf(os.Stderr, "filename:nodaemon\n")
	fmt.Fprintf(os.Stderr, "filename:setuid\n")
	fmt.Fprintf(os.Stderr, "    #term login root with SUID of sshd.\n")
	fmt.Fprintf(os.Stderr, "filename:passwd\n")
	fmt.Fprintf(os.Stderr, "    #format:\n")
	fmt.Fprintf(os.Stderr, "     usr:passwd\n")
	fmt.Fprintf(os.Stderr, "filename:authorized_keys\n")
	fmt.Fprintf(os.Stderr, "    hostkey file names:\n")
	fmt.Fprintf(os.Stderr, "        ssh_host_dsa_key\n")
	fmt.Fprintf(os.Stderr, "        ssh_host_ecdsa_key\n")
	fmt.Fprintf(os.Stderr, "        ssh_host_rsa_key\n")
	fmt.Fprintf(os.Stderr, "        id_rsa\n")
	fmt.Fprintf(os.Stderr, "    if no key file, RandomHostkey() will auto run,\n")
	fmt.Fprintf(os.Stderr, "    the key fingerprint will be changed everytime.\n")
	fmt.Fprintf(os.Stderr, "usage2: %s <-s/e/stop/exit>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "send kill signal to running sshd.\n\n")
	fmt.Fprintf(os.Stderr, "Any question, please contact 'hengwu0 <wu.heng@zte.com.cn>'.\n")
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(2)
}

var conf *config

func main() {
	flagParse()
	conf = mustFindConfig("./config")

	if conf.shouldDaemon {
		if err := daemon.Daemonize(daemonStart, dbg == true); err != nil {
			dbg.Debug("Error daemonizing: %v", err)
		}
	} else {
		waitFunc, _ := daemonStart()
		if waitFunc != nil {
			waitFunc()
		}
	}
}

func flagParse() {
	var exit bool
	flag.BoolVar(&exit, "s", false, "exit sshd.")
	flag.BoolVar(&exit, "stop", false, "exit sshd.")
	flag.BoolVar(&exit, "e", false, "exit sshd.")
	flag.BoolVar(&exit, "exit", false, "exit sshd.")
	flag.Usage = usage
	flag.Parse()
	if f := flag.Args(); len(f) != 0 {
		fmt.Fprintf(os.Stderr, "Can't parse %v\n", f)
		os.Exit(1)
	}

	if exit {
		dbg = true
		proc.SendExitSignal()
		os.Exit(0)
	}
}

// Actually run the implementation of the daemon
func daemonStart() (waitFunc func(), stopFunc func()) {
	proc.SetSignalExit()
	conf.changePWD()
	server := NewServer()

	hasHostKeys := false
	for _, keyName := range keyNames {
		if keyData, err := conf.getBytes(keyName); err == nil {
			dbg.Debug("Adding hostkey file: %s", keyName)
			if err = server.AddHostkey(keyData); err != nil {
				dbg.Debug("Error adding public key: %v", err)
			}
			hasHostKeys = true
		}
	}
	if !hasHostKeys {
		if err := server.RandomHostkey(); err != nil {
			dbg.Debug("Error adding random hostkey: %v", err)
			return
		}
	}

	if authData, err := conf.getBytes("authorized_keys"); err == nil {
		dbg.Debug("Adding authorized_keys.")
		server.AddAuthorizedKeys(authData)
	} else {
		dbg.Debug("No authorized keys found: %v", err)
	}
	server.ListenAndServe(conf.getPort())
	return server.Wait, server.Stop
}
