package proc

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

func Setuid(b bool) {
	if b {
		runtime.LockOSThread() //lock thread and try setuid root if can.
		syscall.Syscall(syscall.SYS_SETUID, uintptr(0), 0, 0)
	}
}

func SetSignalExit(exitFunc func()) {
	var s = make(chan os.Signal)
	signal.Notify(s, syscall.SIGCONT)
	go signalProcess(s, exitFunc)
}

func signalProcess(s <-chan os.Signal, exitFunc func()) {
	for {
		_ = <-s
		exitFunc()
	}
}

func SendExitSignal() {
	pids := seachProcess("/proc")
	if len(pids) == 0 {
		fmt.Fprintf(os.Stderr, "Can't find sshd in /proc/!\n")
		return
	}
	for _, pid := range pids {
		syscall.Kill(pid, syscall.SIGCONT)
	}
	time.Sleep(time.Second / 10)
	if !checkPids(pids) {
		fmt.Fprintf(os.Stderr, "no sshd exit!\n")
		return
	}
}

func checkPids(pids []int) (res bool) {
	for _, pid := range pids {
		if _, err := os.Stat(fmt.Sprintf("/proc/%d", pid)); err != nil {
			fmt.Fprintf(os.Stderr, "%d exited.\n", pid)
			res = true
		}
	}
	return
}

func seachProcess(srcPath string) (pids []int) {
	pids = make([]int, 0, 2)
	dir, err := os.Open(srcPath)
	if err != nil {
		return
	}
	defer dir.Close()
	fis, err := dir.Readdir(0)
	if err != nil {
		return
	}
	for _, fi := range fis {
		if fi.IsDir() {
			pid := isNum(fi.Name())
			if pid == 0 {
				continue
			}
			exe, err := os.Readlink(path.Join(srcPath, fi.Name(), "exe"))
			if err != nil {
				continue
			}
			if path.Base(exe) == path.Base(os.Args[0]) {
				pids = append(pids, pid)
			}
		}
	}
	return
}

func isNum(pid string) int {
	p, _ := strconv.Atoi(pid)
	return p
}
