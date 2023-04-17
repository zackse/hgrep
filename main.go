// Emit the first N line read from stdin (terminated by newline),
// then run grep with the arguments passed in.
//
// Useful for printing the header line of ps/lsof/netstat/etc.
// output when searching for a process pattern.
//
// $ hgrep [-n LINES] ... grep args
//
// Example:
//    $ ps -efw | hgrep ssh
//    UID        PID  PPID  C STIME TTY          TIME CMD
//    root      1216     1  0 Jan03 ?        00:00:00 /usr/sbin/sshd -D
//    zackse    3474  3010  0 Jan03 pts/23   00:00:09 ssh dev
//    zackse    3499  3013  0 Jan03 pts/24   00:00:00 ssh augustus

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func main() {
	grep, err := exec.LookPath("grep")
	if err != nil {
		panic(err)
	}

	headerLines := 1
	grepArgs := os.Args

	// Sad hand-rolled command-line option parsing here vs. `flag` to avoid requiring
	// the user to include `--` as a separator between this program's args and the
	// ones intended for grep(1)
	if len(os.Args) > 2 && os.Args[1] == "-n" {
		if headerLines, err = strconv.Atoi(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't understand -n LINES [%s]\n", os.Args[2])
			os.Exit(1)
		}
		grepArgs = []string{os.Args[0]}
		grepArgs = append(grepArgs, os.Args[3:]...)
	}

	linesSeen := 0
	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		os.Stdout.Write(b)
		if rune(b[0]) == '\n' {
			linesSeen++
			if linesSeen == headerLines {
				break
			}
		}
	}

	err = syscall.Exec(grep, grepArgs, os.Environ())
	if err != nil {
		panic(err)
	}
}
