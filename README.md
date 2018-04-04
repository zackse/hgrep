hgrep
=====

Preserve headers when `grep`-ping: emit the first line of input, then execute
`grep(1)`.

Installation
------------

    go get github.com/zackse/hgrep

Description
-----------

This program reads from standard input and emits the first line, then executes
`grep` with any supplied arguments. It assumes newline is the line delimiter.

This is useful for printing the header line of ps/lsof/netstat/etc.
output when searching for a process pattern.

Example
-------

```bash
$ ps -efw | hgrep ssh
UID        PID  PPID  C STIME TTY          TIME CMD
root      1216     1  0 Jan03 ?        00:00:00 /usr/sbin/sshd -D
zackse    3474  3010  0 Jan03 pts/23   00:00:09 ssh peaeye
zackse    3499  3013  0 Jan03 pts/24   00:00:00 ssh augustus
```

Before:

```bash
# lsof -nP -iTCP | grep smb
smbd      2947            root   37u  IPv4   27529      0t0  TCP *:445 (LISTEN)
smbd      2947            root   38u  IPv4   27530      0t0  TCP *:139 (LISTEN)
```

After:

```bash
# lsof -nP -iTCP | hgrep smb
COMMAND    PID            USER   FD   TYPE  DEVICE SIZE/OFF NODE NAME
smbd      2947            root   37u  IPv4   27529      0t0  TCP *:445 (LISTEN)
smbd      2947            root   38u  IPv4   27530      0t0  TCP *:139 (LISTEN)
```

Alternatives
------------

You could replace this code with a shell function wrapper around Perl, for example:

```bash
hgrep () {
    perl -e '
        my $c;
        while (sysread(STDIN, $c, 1)) {
            print $c;
            last if $c eq "\n";
        }
        exec("grep", @ARGV) or die "exec(grep): $!"
    ' -- "$@"
}
```

Or a Perl script alone:

```perl
#!/usr/bin/perl

$|++;
use strict;

my $c;
while (sysread(STDIN, $c, 1)) {
    print $c;
    last if $c eq "\n";
}

exec("grep", @ARGV) or die "exec(grep): $!"
```

License
-------

MIT
