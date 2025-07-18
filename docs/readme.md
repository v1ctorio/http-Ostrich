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
