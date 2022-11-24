# USCS Liftover Server for 23andMe-formatted raw data

## Description
The project contains back-end and front-end source code for a web application providing positional liftover services
for 23andMe-formatted raw data files.

This web application is written in Golang and utilizes [pyliftover](https://github.com/konstantint/pyliftover) with the
code attached to a project. For more information on liftover dependencies see [documentation](pyliftover/README.md).

This web application uses authentication, cookies and sessions. Functionality is a subject of change and improvement.

## Input

Submitted files must comply with the following standard:
1. Non-compressed plain text
2. Tab-separated fields for body entries
3. `#` precedes header entries
4. rsID, CHROM, POS, GT fields where CHROM is either 1-22, X, Y or MT; POS is an integer

The example input is provided below:
```text
# rsid	chromosome	position	genotype
rs9701055	1	630053	CC
rs9651229	1	632287	CC
rs9701872	1	632828	TT
rs11497407	1	633147	GG
rs116587930	1	792461	GG
rs3131972	1	817341	GG
rs200599638	1	817538	GG
```

## Running

### Non-containerised

Run from inside the parent directory containing the project, set this directory as env `HOME` variable.

```bash
AUTH_USERNAME=alice AUTH_PASSWORD=alice HOME=/parent/dir \
CERT=/path/to/localhost.pem \
KEY=/path/to/localhost-key.pem \
STORAGE=./file-storage \
HOST=localhost PORT=8080 go run ./web-server-go-liftover/cmd/server/main.go
```

Accessing the main page is available at `https://<serverHost>:<serverPort>/index/`. Multiple client sessions can be run
simultaneously.

### Containerised

Build an image and run the container:

```bash
docker build .
docker container run --env AUTH_USERNAME=alice --env AUTH_PASSWORD=alice --env PORT=8080 --env HOST=0.0.0.0 \
--publish target=8080,published=<localIP>:8080,protocol=tcp <imageID>
```


## Notes
`CERT` and `KEY` files point to `.pem` files created locally via [`mkcert`](https://github.com/FiloSottile/mkcert)
or issued by any CA. `STORAGE` points to a folder which will handle uploads and downloads of files. The directory
will be created according to the provided path unless exists.

### File storage
Uploaded and processed files are stored in `uploaded-files` and `processed-files` folders in the `STORAGE`
directory according to configuration parameters. Client and server-provided files are immediately deleted upon processing
completion.
