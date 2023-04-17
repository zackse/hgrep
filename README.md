hgrep
=====

Preserve headers when `grep`-ping: emit the first N lines of input, then execute
`grep(1)`.

Installation
------------

    go install github.com/zackse/hgrep

Description
-----------

This program reads from standard input and emits the first N (default 1) lines, then executes
`grep` with any supplied arguments. It assumes newline is the line delimiter.

This is useful for preserving the header line(s) of ps/lsof/netstat/etc.
output when searching for a pattern in the body of output.

Usage
-----

```bash
$ hgrep [-n LINES] grep_args ... < input
```

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

You could implement this with Perl, for example:

```perl
#!/usr/bin/perl

BEGIN {
  # allow "-n 2", for example
  $HEADER_LINES = ($ARGV[0] =~ /^-+n/ && shift @ARGV && shift @ARGV) || 1;
}
my $c;
while (sysread(STDIN, $c, 1)) {
  print $c;
  $lines_seen++ if $c eq "\n";
  last if $lines_seen == $HEADER_LINES;
}
exec("grep", @ARGV) or die "exec(grep): $!"
```

Or bash:

```bash
hgrep() {
    local lines=1
    if [ "$1" = "-n" ]; then
        lines=${2:-1}
        shift
        shift
    fi
    exec 3>&1
    tee >(head -n $lines >&3) | grep "$@"
    exec 3>&-
}
```

License
-------

MIT
