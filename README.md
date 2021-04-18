# elastically

A command line interface (CLI) to interact with OpenSearch (or Elasticsearch)

## Supported Features

- List Indexes
- Delete Index

## Requirements
- [Go](https://golang.org/) >= 1.13

## Build
```bash
$ go build -o elastically ./cmd
$ ./elastically -h
```

## Usage

```bash
Usage: elastically [--url URL] [--sniff] [--loglevel LOGLEVEL] <command> [<args>]

Options:
  --url URL, -u URL
  --sniff, -s            Sniff nodes, then connect to them
  --loglevel LOGLEVEL [default: info]
  --help, -h             display this help and exit

Commands:
  index
```

### Index Subcommand

```bash
Usage: elastically index <command> [<args>]

Options:
  --help, -h             display this help and exit

Commands:
  list
  delete
```

#### List Indices

```bash
$ elastically index list   
.opensearch_dashboards_1
fluentd
logstash-2021.04.18
(...)
```

#### Delete index
```bash
$ elastically index delete fluentd
fluentd delete successfully!
```