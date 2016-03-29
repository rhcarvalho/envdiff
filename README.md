# envdiff

Generates smart environment diffs.


## Installing

You will need to have the `go` tool [installed](https://golang.org/doc/install),
then:

```
go get -u github.com/rhcarvalho/envdiff
```

Now `envdiff` should be in your PATH.

## Using

`envdiff` takes two arguments that should contain a list of newline or null-byte
separated environment variables.

Examples:

1. Detecting new variables:

  ```console
$ envdiff TERM=xterm $'TERM=xterm\nFOO=bar'
FOO=bar
```

2. Detecting removed variables:

  ```console
$ envdiff $'TERM=xterm\nFOO=bar' TERM=xterm
FOO=
```

3. Detecting changing of a list:

  ```console
$ envdiff PATH=/bin:/sbin PATH=/usr/bin/:/usr/sbin:/bin:/sbin
PATH=/usr/bin/:/usr/sbin:$PATH
```
