# hashi-report

Generate a usage report for a given hashicorp tool cluster.

## Building

```shell
# substitute the version number as desired
go build -ldflags "-X main.Version=0.1.0
```

## Usage

```
Usage: hashi-report [--version] [--help] <command> [<args>]

Available commands are:
    nomad      Generate a nomad report
    version    Return the version of the binary
```
