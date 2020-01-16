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
	shouldDaemon bool
}

// Lookup the port number
func (c *config) getPort() int16 {
	if portData, err := c.getBytes("port"); err == nil {
		portData := strings.TrimSpace(string(portData))
		if port, err := strconv.Atoi(portData); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing %s as port: %v, use 1022 for default.", portData, err)
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
	pwd, _ := os.Getwd()
	if pwd == "" {
		pwd = "Unknown"
	}
	dbg.Debug("changePWD failed. Now PWD: %s", pwd)
}

func (c *config) getPasswdUpdateMsg() {
	if c.usrpasswd != c.fileExists("passwd") {
		dbg.Debug("You need to restart %s to let passwd-file take effect!!!", os.Args[0])
	}
}

func mustFindConfig(name string) *config {
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't find exe path, use cwd instead.\n")
		name = filepath.Join("./", name)
	} else {
		if exe2, err := filepath.EvalSymlinks(exe); err != nil {
			fmt.Fprintf(os.Stderr, "Can't read Symlinks(%s), use cwd instead.\n", exe)
			name = filepath.Join("./", name)
		} else {
			name = filepath.Join(filepath.Dir(exe2), name)
		}
	}
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
	return c
}
