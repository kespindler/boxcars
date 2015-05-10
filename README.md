Easy-to-configure Static Web/Reverse Proxy Server in Go.

![](http://i.cloudup.com/i5Tpn00lCc.png)

* [Install](#install)
* [Usage](#usage)
* [Configuration Examples](#configuration-examples)
* [Production](#production)
* [Logging](#logging)
* [Benchmarks](#benchmarks)
* [Troubleshooting](#troubleshooting)
* [TODO](#todo)

## Install

```bash
go get github.com/kespindler/boxcars/boxcars
go install github.com/kespindler/boxcars/boxcars
```

[Works with Go v1.1+](https://github.com/azer/boxcars/issues/14)

## Usage

Create a configuration file *(it'll be auto-loading changes once you start)*  like the below example;

```yaml
foo.com: /home/you/sites/foo.com
"*.bar.net": localhost:8080
qux.org:
  /static: /home/you/qux.org/static
  /favicon.ico: /home/you/qux.org/static/favicon.ico
  /: localhost:3000
```

And start the server:

```bash
$ boxcars config.yaml
```

To specify the port:

```bash
$ boxcars -port=8001 config.yaml
```

## Production

There is a `-secure` flag that allows you to drop privileges down after opening port 80.

```bash
$ sudo boxcars -port=80 -secure config.yaml
```

However, there is a bug that caused this to fail with an `operation not supported` error.
Instead, do the following

`sudo setcap 'cap_net_bind_service=+ep' $(which boxcars)`

And then run boxcars on port 80 as usual, without using sudo.

You can change the configuration anytime during boxcars running. 
It'll be watching your file and reloading it only if it parses with no error.

## Configuration Examples

I use below configuration for a static single-page app that connects to an HTTP API:

```yaml
singlepage.com:
  /api: localhost:1337
  "*": sites/singlepage.com
```

To catch any domain:

```yaml
foo.com: localhost:1234
"*": /home/you/404.html
```

To set a custom 404 page for a static server:

```yaml
"foo.com:
  /: /home/you/sites/foo.com
  "*": "/home/you/404.html
```

## Security

Once you enable `-secure`, boxcars switches from root user to a basic user after starting the server on specified port. 

```bash
$ sudo boxcars -port=80 -secure example.yaml
```

UID and GID is set to 1000 by default. Use `-uid` and `-gid` parameters to specify your own in case you need.

## Logging

Boxcars uses [debug](http://github.com/azer/debug) for logging. To enable logging for specific modules: 

```bash
$ DEBUG=server,sites boxcars config.yaml
```

To see how boxcars setup the HTTP handlers for your configuration;

```bash
$ DEBUG=handlers-of,sites boxcars config.yaml
```

To enable very verbose mode (not recommended):

```bash
$ DEBUG=* boxcars config.yaml
```

To silentize:

```bash
$ DEBUG=. boxcars config.yaml
```

It'll be outputting to stderr.

## Benchmarks

* [Nginx VS Boxcars](https://gist.github.com/azer/5955772)

## Troubleshooting

#### The "Too Many Open Files" Error

Boxcars creates a lot of files on `/proc/$pid/fd`. In case you see boxcars crashing, you can see how many files are open by;

```bash
$ sudo ls -l /proc/`pgrep boxcars`/fd | wc -l
```

To find out your personal limit:

```bash
$ ulimit -n
```

To change it:

```bash
$ ulimit -n 64000
```


You can change soft - hard limits by editing `/etc/security/limits.conf`.

## TODO

* Add -daemon option.

![](http://i.cloudup.com/rH_0UwNYg1.jpg)
