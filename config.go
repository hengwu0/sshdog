package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type Debugger bool

var dbg Debugger = false

func (d Debugger) Debug(format string, args ...interface{}) {
	if d {
		msg := fmt.Sprintf(format, args...)
		fmt.Fprintf(os.Stderr, "[DEBUG] %s\n", msg)
	}
}

type config struct {
	dir          string
	usrpasswd    bool
	setuid       bool
	shouldDaemon bool
}

// Lookup the port number
func (c *config) getPort() int16 {
	if len(os.Args) > 1 {
		if port, err := strconv.Atoi(os.Args[1]); err != nil {
			dbg.Debug("Error parsing %s as port: %v", os.Args[1], err)
		} else {
			return int16(port)
		}
	}
	if portData, err := c.getBytes("port"); err == nil {
		portData := strings.TrimSpace(string(portData))
		if port, err := strconv.Atoi(portData); err != nil {
			dbg.Debug("Error parsing %s as port: %v", portData, err)
		} else {
			return int16(port)
		}
	}
	return 1022 // default
}

// Just check if a file exists
func (c *config) fileExists(name string) bool {
	_, err := c.getBytes(name)
	return err == nil
}

func (c *config) getPath(name string) string {
	return path.Join(c.dir, name)
}

// in debug?
func (c *config) beDebug() bool {
	return c.fileExists("debug")
}

func (c *config) getBytes(name string) ([]byte, error) {
	file, err := os.Open(path.Join(c.dir, name))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (c *config) findDir(name string) bool {
	if f, err := os.Stat(path.Join(c.dir, name)); err != nil {
		return false
	} else {
		if f.IsDir() {
			return true
		}
	}
	return false
}

func (c *config) changePWD() {
	user, err := user.Current()
	if nil == err {
		if p, err := filepath.Abs(c.dir); err == nil {
			c.dir = p
		} else {
			dbg.Debug("Can't get abs filepath! Change dir canceled.")
			return
		}
		if nil == os.Chdir(user.HomeDir) {
			dbg.Debug("PWD: %s", user.HomeDir)
			return
		}
		dbg.Debug("Can't chdir to %s", user.HomeDir)
		return
	}
	dbg.Debug("Can't find HomeDir.")
}

func mustFindConfig(name string) *config {
	c := &config{dir: name}
	if !c.findDir(".") {
		fmt.Fprintf(os.Stderr, "Must create an `config` dir! use -h show help.\n")
		os.Exit(-1)
		return nil
	}
	if c.beDebug() {
		dbg = true
	}
	c.usrpasswd = c.fileExists("passwd")
	c.shouldDaemon = !c.fileExists("nodaemon")
	c.setuid = c.fileExists("setuid")
	return c
}
