# render

[![Version](https://img.shields.io/badge/version-v0.1.3-brightgreen.svg)](https://github.com/VirtusLab/render/releases/tag/v0.1.3)
[![Travis CI](https://img.shields.io/travis/VirtusLab/render.svg)](https://travis-ci.org/VirtusLab/render)
[![Github All Releases](https://img.shields.io/github/downloads/VirtusLab/render/total.svg)](https://github.com/VirtusLab/render/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/VirtusLab/render "Go Report Card")](https://goreportcard.com/report/github.com/VirtusLab/render)
[![GoDoc](https://godoc.org/github.com/VirtusLab/render?status.svg "GoDoc Documentation")](https://godoc.org/github.com/VirtusLab/render/renderer)

Universal data-driven templates for generating textual output. Can be used as a single static binary (no dependencies)
or as a golang library.

Just some of the things to `render`:
- configuration files
- Infrastructure as Code files (e.g. CloudFormation templates)
- Kubernetes manifests

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
* [Contribution](README.md#contribution)
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
   v0.1.3-86eb03f

AUTHOR:
   VirtusLab

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d                   run in debug mode
   --indir value                 the input directory, can't be used with --out
   --outdir value                the output directory, the same as --outdir if empty, can't be used with --in
   --in value                    the input template file, stdin if empty, can't be used with --outdir
   --out value                   the output file, stdout if empty, can't be used with --indir
   --config value                optional configuration YAML file, can be used multiple times
   --set value, --var value      additional parameters in key=value format, can be used multiple times
   --unsafe-ignore-missing-keys  do not fail on missing map key and print '<no value>' ('missingkey=invalid')
   --help, -h                    show help
   --version, -v                 print the version
```

**Notes:**
- `--in`, `--out` take only files (not directories), `--in` will consume any file as long as it can be parsed
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

Also see a [more advanced template](examples/example.yaml.tmpl) example.

#### As a library

```go
package example

import (
    "github.com/VirtusLab/render/renderer"
    "github.com/VirtusLab/render/renderer/parameters"
)

func CustomRender(template string, opts []string, params parameters.Parameters) (string, error) {
    return renderer.New(
    	renderer.WithOptions(opts...),
        renderer.WithParameters(params),
        renderer.WithSprigFunctions(),
        renderer.WithExtraFunctions(),
        renderer.WithCryptFunctions(),
    ).Render(template)
}
```

See also [`other functions`](https://godoc.org/github.com/VirtusLab/render/renderer).

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

- `render` - calls the `render` from inside of the template, making the renderer recursive (also accepts an optional template parameters override)
- `readFile` - reads a file from a path, relative paths are translated to absolute paths, based on `root` function
- `root` - the root path, used for relative to absolute path translation in any file based operations; by default `PWD` is used
- `toYaml` - provides a configuration data structure fragment as a YAML format
- `gzip`, `ungzip` - use `gzip` compression and extraction inside the templates, for best results use with `b64enc` and `b64dec`

See also [example](examples/example.yaml.tmpl) template 
and a more [detailed documentation](https://godoc.org/github.com/VirtusLab/render/renderer#Renderer.ExtraFunctions).

Cloud KMS (AWS, Amazon, Google) based cryptography functions form [`crypt`](https://github.com/VirtusLab/crypt):
- `encryptAWS` - encrypts data using AWS KMS, for best results use with `gzip` and `b64enc`
- `decryptAWS` - decrypts data using AWS KMS, for best results use with `ungzip` and `b64dec`
- `encryptGCP` - encrypts data using GCP KMS, for best results use with `gzip` and `b64enc`
- `decryptGCP` - decrypts data using GCP KMS, for best results use with `ungzip` and `b64dec`
- `encryptAzure` - encrypts data using Azure Key Vault, for best results use with `gzip` and `b64enc`
- `decryptAzure` - decrypts data using Azure Key Vault, for best results use with `ungzip` and `b64dec`

#### Helm compatibility

As of now, there is a limited Helm 2 Chart compatibility, simple Charts will render just fine.

To mimic Helm behaviour regarding to missing keys use `--unsafe-ignore-missing-keys` option.

There is no plan to implement full compatibility with Helm, because of unnecessary complexity that would bring.

If you need full Helm compatilble rendering see: [`helm-nomagic`](https://github.com/giantswarm/helm-nomagic).

## Limitations and future work

#### Planned new features

- `.renderignore` files [`#12`](https://github.com/VirtusLab/render/issues/12)

#### Operating system support

We provide cross-compiled binaries for most platforms, but is currently used mainly with `linux/amd64`.

## Contribution

Feel free to file [issues](https://github.com/VirtusLab/render/issues) 
or [pull requests](https://github.com/VirtusLab/render/pulls).

## Development

    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    
    mkdir -p $GOPATH/src/github.com/VirtusLab
    cd $GOPATH/src/github.com/VirtusLab/render
    git clone git@github.com:VirtusLab/render.git
    cd render
    
    go get -u github.com/golang/dep/cmd/dep
    make all

## The name

We believe in obvious names. It renders. It's a *verb*. It's `render`.
