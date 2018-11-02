# render

[![Version](https://img.shields.io/badge/version-v0.0.3-brightgreen.svg)](https://github.com/VirtusLab/render/releases/tag/v0.0.3)
[![Travis CI](https://img.shields.io/travis/VirtusLab/render.svg)](https://travis-ci.org/VirtusLab/render)
[![Github All Releases](https://img.shields.io/github/downloads/VirtusLab/render/total.svg)](https://github.com/VirtusLab/render/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/VirtusLab/render "Go Report Card")](https://goreportcard.com/report/github.com/VirtusLab/render)
[![GoDoc](https://godoc.org/github.com/VirtusLab/render?status.svg "GoDoc Documentation")](https://godoc.org/github.com/VirtusLab/render/renderer)

Universal data-driven templates for generating textual output. Can be used as a single static binary (no dependencies)
or as a golang library.

The renderer extends 
[go-template](https://golang.org/pkg/text/template/) and [Sprig](http://masterminds.github.io/sprig/) functions. 

If you are interested in one of the use cases, take a look at this [blog post](https://medium.com/virtuslab/helm-alternative-d6568aa9d40b) 
about Kubernetes resources rendering. Also see [Helm compatibility](README.md#helm-compatibility).

* [Installation](README.md#installation)
  * [Binaries](README.md#binaries)
  * [Via Go](README.md#via-go)
* [Usage](README.md#usage)
  * [Command line](README.md#command-line)
  * [Notable standard and sprig functions](README.md#notable-standard-and-sprig-functions)
  * [Custom functions](README.md#custom-functions)
  * [Helm compatibility](README.md#helm-compatibility)
  * [Limitations and future work](README.md#limitations-and-future-work)
* [Development](README.md#development)
* [The Name](README.md#the-name)

## Installation

#### Binaries

For binaries please visit the [Releases Page](https://github.com/VirtusLab/render/releases).

The binaries are statically compiled and does not require any dependencies. 

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
- `--set`, `--var` are the same (one is used in Helm, the other in Terraform), we provide both for convenience, any values set here **will override** values form configuration files

#### Command line

Example usage of `render` with `stdin`, `stdout` and `--var`:
```console
$ echo "something {{ .value }}" | render --var "value=new"
something new
```

Example usage of `render` with `--in`, `--out` and `--config`:
```console
$ echo "something {{ .value }}" > test.txt.tmpl
$ echo "value: new" > test.config.yaml
$ ./render --in test.txt.tmpl --out test.txt --config test.config.yaml
$ cat test.txt
something new
```

Also see a [more advanced tempalte](examples/example.yaml.tmpl) example.

#### As a library

```go
package example

import (
    "github.com/VirtusLab/render/renderer"
    "github.com/VirtusLab/render/renderer/configuration"
)

func CustomRender(template string) (string, error) {
    config := configuration.Configuration{}
    r := renderer.New(config, renderer.MissingKeyErrorOption)
    return r.Render("nameless", template)
}
```

See also [`RenderWith`](https://godoc.org/github.com/VirtusLab/render/renderer#Renderer.RenderWith) function that takes a custom functions map.

Also see [tests](https://github.com/VirtusLab/render/blob/master/renderer/render_test.go) for more usage examples.

#### Notable standard and sprig functions

- [`indent`](https://masterminds.github.io/sprig/strings.html#indent)
- [`default`](https://masterminds.github.io/sprig/defaults.html#default)
- [`ternary`](https://masterminds.github.io/sprig/defaults.html#ternary)
- [`toJson`](https://masterminds.github.io/sprig/defaults.html#tojson)
- [`b64enc`, `b64dec`](https://masterminds.github.io/sprig/encoding.html)

All syntax and functions:
- [Go template functions](https://golang.org/pkg/text/template)
- [Sprig functions](http://masterminds.github.io/sprig)

#### Custom functions

- `render` - calls the `render` from inside of the template, making the renderer recursive
- `readFile` - reads a file from a path, relative paths are translated to absolute paths, based on `root` function
- `root` - the root path, used for relative to absolute path translation in any file based operations; by default `PWD` is used
- `toYaml` - provides a configuration data structure fragment as a YAML format
- `gzip`, `ungzip` - use `gzip` compression and extraction inside the templates, for best results use with `b64enc` and `b64dec`

See also [example](examples/example.yaml.tmpl) template 
and a more [detailed documentation](https://godoc.org/github.com/VirtusLab/render/renderer#Renderer.ExtraFunctions).

#### Helm compatibility

As of now, there is a limited Helm 2 Chart compatibility, simple Charts will render just fine.

There is no plan to implement full compatibility with Helm, because of unnecessary complexity that would bring.

## Limitations and future work

#### Planned new functions

- `encrypt`, `decrypt` - cloud KMS (AWS, Amazon, Google) based encryption for any data

#### Planned new features

- directories as `--in` and `--out` arguments, currently only files are supported

#### Operating system support

We provide cross-compiled binaries for most platforms, but is currently used mainly with `linux/amd64`.

## Contribution

Feel free to file [issues](https://github.com/VirtusLab/render/issues) or [pull requests](https://github.com/VirtusLab/render/pulls).

## Development

    mkdir $GOPATH/src/github.com/VirtusLab/
    git clone 
    
    go get -u github.com/golang/dep/cmd/dep
    
    export PATH=$PATH:$GOPATH/bin
    cd $GOPATH/src/github.com/VirtusLab/render
    make all

## The name

We believe in obvious names. It renders. It's a *verb*. It's `render`.
