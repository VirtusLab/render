# render

[![Version](https://img.shields.io/badge/version-v0.0.3-brightgreen.svg)](https://github.com/VirtusLab/render/releases/tag/v0.0.3)
[![Travis CI](https://img.shields.io/travis/VirtusLab/render.svg)](https://travis-ci.org/VirtusLab/render)
[![Github All Releases](https://img.shields.io/github/downloads/VirtusLab/render/total.svg)](https://github.com/VirtusLab/render/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/VirtusLab/render "Go Report Card")](https://goreportcard.com/report/github.com/VirtusLab/render)

Universal file renderer based on [go-template](https://golang.org/pkg/text/template/) 
and [Sprig](http://masterminds.github.io/sprig/) functions. 

If you want to read more about the underlying reason for creating this library and our initial use of it, take a look at this [blog post](https://medium.com/virtuslab/helm-alternative-d6568aa9d40b) (hint: Kubernetes config files rendering). Keep in mind however that it is actually only one use case, `render` can be used for any text file rendering.

* [Installation](README.md#installation)
  * [Binaries](README.md#binaries)
  * [Via Go](README.md#via-go)
* [Usage](README.md#usage)
  * [Custom functions](README.md#custom-functions)
  * [Helm compatibility](README.md#helm-compatibility)
  * [Limitations and future work](README.md#limitations-and-future-work)
* [Development](README.md#development)
* [The Name](README.md#the-name)

## Installation

#### Binaries

For binaries please visit the [Releases Page](https://github.com/VirtusLab/render/releases).

#### Via Go

```console
$ go get github.com/VirtusLab/render
```

## Usage

```console
$ render --help
NAME:
   render - Universal file renderer

USAGE:
   render [global options] command [command options] [arguments...]

VERSION:
   v0.0.3-9545028

AUTHOR:
   VirtusLab

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d               run in debug mode
   --in value                the input template file, stdin if empty
   --out value               the output file, stdout if empty
   --config value            optional configuration YAML file, can be used multiple times
   --set value, --var value  additional parameters in key=value format, can be used multiple times
   --help, -h                show help
   --version, -v             print the version
```

**Notes:**
- `--in`, `--out` take only files (not directories) at the moment, `--in` will consume any file as long as it can be parsed
- `stdin` and `stdout` can be used instead of `--in` and `--out`
- `--config` accepts any YAML file, can be used multiple times, the values of the configs will be merged
- `--set`, `--var` are the same (one is used in Helm, the other in Terraform), we provide both for convinience, any values set here **will override** values form configuration files

#### Custom functions

- `render` - invokes the `render` from inside of the template, making the renderer recursive, [see example](examples/example.yaml.tmpl#L10), can be combined with other functions, e.g. `b64dec`, `b64enc`, `readFile`
- `readFile` - reads a file from a given path, relative paths are translated to absolute paths, based on `root` function
- `root` - the root path for rendering, used relative to absolute path translation in any file based operations; by default `PWD` is used, can be overridden with a `--config` or `--set`

#### Helm compatibility

As of now, there is a limited Helm 2 Chart compatibility, simple Charts will render just fine.

There is no plan to implement full compatibility with Helm, because of unnecessary complexity that would bring.

## Limitations and future work

Planned new functions:

- `encrypt`, `decrypt` - cloud KMS (AWS, Amazon, Google) based encryption for any data
- `gzip`, `ungzip` - use `gzip` compression and extraction inside the templates

Planned new features:

- directories as `--in` and `--out` arguments, currently only files are supported

## Development

    mkdir $GOPATH/src/github.com/VirtusLab/
    git clone 
    
    go get -u github.com/golang/dep/cmd/dep
    
    export PATH=$PATH:$GOPATH/bin
    cd $GOPATH/src/github.com/VirtusLab/render
    make all

## The name

We believe in obvious names. It renders. It's a *verb*, `render`.
