// Emit the first line read from stdin (terminated by newline),
// then run grep with the arguments passed in.
//
// Useful for printing the header line of ps/lsof/netstat/etc.
// output when searching for a process pattern.
//
// Example:
//    $ ps -efw | hgrep ssh
//    UID        PID  PPID  C STIME TTY          TIME CMD
//    root      1216     1  0 Jan03 ?        00:00:00 /usr/sbin/sshd -D
//    zackse    3474  3010  0 Jan03 pts/23   00:00:09 ssh dev
//    zackse    3499  3013  0 Jan03 pts/24   00:00:00 ssh augustus

package main

import (
	"os"
	"os/exec"
	"syscall"
)

func main() {
	grep, err := exec.LookPath("grep")
	if err != nil {
		panic(err)
	}

	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		os.Stdout.Write(b)
		if rune(b[0]) == '\n' {
			break
		}
	}

	err = syscall.Exec(grep, os.Args, os.Environ())
	if err != nil {
		panic(err)
	}
}
