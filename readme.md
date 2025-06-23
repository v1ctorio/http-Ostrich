# HTTP Ostrich

The ultimate fast ready to use HTTP server for easily distributing files inside a network.

## How it works (pseudocode)
Ok this is basically a schema of what I need it to do. It will be a cli with the following flags:

```
--port <port> # Port to listen on
--zip # if it is a directory, zip it and only server the zip
--no-expose # Do not reply to request outside of localhost
--passphrase # Simple HTTP basic auth to access the files

[target] (first positional argument) # The target file or directory to serve
```

pseudocode:
```
{port, zip, noExpose, passphrase } = args

type file {
    size: int,
    name: string,
    content: binary,
}

type directory {
    files: [file],
    name: string,
}


const readTarget = readDirectory(target)

const serving: directory = getServing(readTarget);

const server = createServer(port, serving, noExpose)

server.get('/', (req,res) => {
    res.send(serving.files.forEach(f=> {
        return `<a href="/f/${f.name}">${f.name} (${f.size} bytes)</a><br>`
    }))
})


server.get('/f/:name', (req, res) => {
    const file = serving.files.find(f => f.name === req.params.name)
    if (file) {
        res.setHeader('Content-Type', 'application/octet-stream')
        res.setHeader('Content-Disposition', `attachment; filename="${file.name}"`)
        res.send(file.content)
    } else {
        res.status(404).send('File not found')
    }
})





fun getServing(target) {


if target.isDir {
    if zip {
        const zipFile = createZipFromDirectory(target)
        return {
            files: [
                {
                    size: zipFile.size,
                    name: target.name,
                    content: zipFile.content
                }
            ],
            name: `{target.name} (zipped)`
        }
    } else {
        const files = mapFiles(target)
        return {
            files: files,
            name: target.name
        }
    }

} else {
    return {
        files: [
            {
                size: target.size,
                name: target.name,
                content: target.content
            }
        ],
        name: target.name 
    }
}
}

