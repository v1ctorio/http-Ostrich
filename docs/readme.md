# http-Ostrich

The ultimate fast ready to use HTTP server for easily distributing files inside a network.

Have you ever wante to share files quickly and easily with no more than a terminal command? not any install on the client and no configuration but all the golang speed and simplicity.

![Usage example video](https://vhs.charm.sh/vhs-3dMowiKbK1rMRP0Nokand3.gif)

## Usage

```bash
$ http-ostrich [flags] <file1> <file2> ...
```
```
NAME:
   http-ostrich - The http file-sharing ostrich

USAGE:
   http-ostrich [global options] [arguments...]

DESCRIPTION:
   The easy and fast http file sharing ostrich.

GLOBAL OPTIONS:
   --port int, -p int              Port to listen on (default: 0)
   --expose, -e                    Wether to expose the server to foreign IPs (default: false)
   --passphrase string, -a string  Passphrase for basic authentication
   --zip, -z                       Wether to compress the files into a zip file (default: false)
   --recursive, -r                 (default: false)
   --verbose, -v                   (default: false)
   --help, -h                      show help
```


## Installation


### with nix flakes 
```bash
nix profile install github:v1ctorio/http-ostrich
```

### Download with github releases (todo)


### With Go (build it yourself, needs Go 1.22 installed)
```bash
go install github.com/v1ctorio/http-ostrich@latest
```
