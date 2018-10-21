# render

Universal file renderer based on [go-template](https://golang.org/pkg/text/template/) 
and [Sprig](http://masterminds.github.io/sprig/) functions.

## Usage

    NAME:
       render - Simple go-template files render

    USAGE:
       render [global options] command [command options] [arguments...]

    VERSION:
       0.0.1

    COMMANDS:
         help, h  Shows a list of commands or help for one command

    GLOBAL OPTIONS:
       --in value      the template file
       --out value     the output file, stdout if empty
       --config value  the config file
       --set value     an additional parameters in key=value format
       --help, -h      show help
       --version, -v   print the version

## Development

    mkdir $GOPATH/src/github.com/VirtusLab/
    git clone 
    
    go get -u github.com/golang/dep/cmd/dep
    
    export PATH=$PATH:$GOPATH/bin
    cd $GOPATH/src/github.com/VirtusLab/render
    make all

## The name

We believe in obvious names.