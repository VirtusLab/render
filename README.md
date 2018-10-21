# render

Universal file renderer based on [go-template](https://golang.org/pkg/text/template/) 
and [Sprig](http://masterminds.github.io/sprig/) functions.

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

## Development

    mkdir $GOPATH/src/github.com/VirtusLab/
    git clone 
    
    go get -u github.com/golang/dep/cmd/dep
    
    export PATH=$PATH:$GOPATH/bin
    cd $GOPATH/src/github.com/VirtusLab/render
    make all

## The name

We believe in obvious names.