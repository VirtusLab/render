# render

Universal file renderer based on [go-template](https://golang.org/pkg/text/template/) 
and [Sprig](http://masterminds.github.io/sprig/) functions.

Custom functions:

- `render` - invokes the `render` from inside of the template, making the renderer recursive, [see example](examples/example.yaml.tmpl#L10), can be combined with other functions, e.g. `b64dec`, `b64enc`, `readFile`
- `readFile` - reads a file from a given path, relative paths are translated to abolute paths, based on `root` function
- `root` - the root path for rendering, used relative to absolute path tranlsation in any file based operations; by default `PWD` is used, can be overriden with a `--config` or `--set`

## Usage

    NAME:
       render - Universal file renderer
    
    USAGE:
       main [global options] command [command options] [arguments...]
    
    VERSION:
       v0.0.2-e932f66
    
    AUTHOR:
       VirtusLab
    
    COMMANDS:
         help, h  Shows a list of commands or help for one command
    
    GLOBAL OPTIONS:
       --debug, -d               run in debug mode
       --in value                the template file
       --out value               the output file, stdout if empty
       --config value            the config file
       --set value, --var value  additional parameters in key=value format, can be used multiple times
       --help, -h                show help
       --version, -v             print the version

## Helm compatibility

As of now, there is a limited Helm 2 Chart compatibility, simple Charts will render just fine.

There is no plan to implement full compatibility with Helm, because of unnecesary complexity that would bring.

## Limitations and future work

Planned new functions:

- `encrypt`, `decrypt` - cloud KMS (AWS, Amazon, Google) based encryption for any data
- `gzip`, `ungzip` - use `gzip` compression and extraction inside the templates

Planned new features:

- diretories as `--in` and `--out` arguments, currently only files are supported

## Development

    mkdir $GOPATH/src/github.com/VirtusLab/
    git clone 
    
    go get -u github.com/golang/dep/cmd/dep
    
    export PATH=$PATH:$GOPATH/bin
    cd $GOPATH/src/github.com/VirtusLab/render
    make all

## The name

We believe in obvious names.
