# http-Ostrich

The ultimate fast ready to use HTTP server for easily distributing files inside a network.

http-Ostrich is a simple cli that will spin up a HTTP server with the provided files allowing you to serve them easily and fast from the terminal to any device in your network only using a web browser. Requires no setup on the client side and minimum setup on the server side (the whole program is bundled in a single binary).

Have you ever wante to share files quickly and easily with no more than a terminal command? not any install on the client and no configuration but all the golang speed and simplicity.

![Usage example video](https://vhs.charm.sh/vhs-1rQ7gJbxRc3ZAyH8GiGZko.gif)

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


### With nix flakes 
```bash
nix profile install github:v1ctorio/http-ostrich
```

### With github releases (manual installation)

##### Linux bash or whatever x86_64
Download to `~/.bin` 
```bash
mkdir -p ~/.bin && cd ~/.bin
curl -L https://github.com/v1ctorio/http-ostrich/releases/latest/download/http-ostrich_Linux-x86_64 -o http-ostrich
chmod +x http-ostrich
```
and add it to your PATH
```bash
# For bash
echo 'export PATH="$HOME/.bin:$PATH"' >> ~/.bashrc
# For fish
fish_add_path ~/.bin
```

##### Windows Powershell (x86_64)
```powershell
$bindir = "$env:USERPROFILE\.bin"
if(!(Test-Path $bindir)){md $bindir} # create the .bin dir if it doesn't exist
cd $bindir
iwr https://github.com/v1ctorio/http-ostrich/releases/latest/download/http-ostrich_Windows-x86_64.exe -o http-ostrich.exe
$p = [Environment]::GetEnvironmentVariable("PATH","User")
if($p -notlike "*$bindir*"){[Environment]::SetEnvironmentVariable("PATH","$p;$bindir","User")} # Add the.bin dir to the PATH if it's not already there
```

### With Go (build it yourself, needs Go 1.22 installed)
```bash
go install github.com/v1ctorio/http-ostrich@latest
```


# Features
- Serve selected files
- Use the `recursive` flag to serve directories
- Use the `port` flag to specify the port to listen on (default is `8069`).
- Use the `zip` flag to compress the provided files or directory into a zip file before serving it.
- Use the `expose` flag to allow access from foreign IPs (default is local only).
- Use the `passphrase` flag to set basic authentication for http.
- Use the `verbose` flag to enable a verbose debug output log.