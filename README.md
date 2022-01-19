# USCS Liftover Server for 23andMe-formatted raw data

## Description
The project contains back-end and front-end source code for a web application providing positional liftover services
for 23andMe-formatted raw data files.

This web application is written in Golang and utilizes [pyliftover](https://github.com/konstantint/pyliftover) with the
code attached to a project. For more information on liftover dependencies see [documentation](./liftover/README.md).

This web application uses authentication, cookies and sessions. Functionality is a subject of change and improvement.

Currently, the herein described service is available at https://65.108.154.166:4000/index/ (authentication is required
upon file submission, credentials are currently being provided only upon request here or elsewhere).

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
To start the local server run the following command in a terminal specifying `AUTH_USERNAME`, `AUTH_PASSWORD` and a
path (either relative or absolute) to a configuration file `CONFIG`. Note that [example-config.yaml](./resources/example-config.yaml)
is a template and has to be amended.
```bash
AUTH_USERNAME=alice AUTH_PASSWORD=alice CONFIG=/path/to/config.yaml go run main.go
```

Accessing the main page is available at `https://<serverIP>:<serverPort>/index/`. Multiple client sessions can be run
simultaneously.

## Notes

### Config
The [config.yaml](./resources/example-config.yaml) file implements the following self-explanatory structure:
```yaml
constants:
  certFile: /path/to/localhost.pem
  keyFile: /path/to/localhost-key.pem
  serverIP: X.X.X.X
  serverPort: 4000
  fileStrorage: /path/to/filestoragedir
```
where `certFile` and `keyFile` point to `.pem` files created locally via [`mkcert`](https://github.com/FiloSottile/mkcert)
or issued by any CA. `fileStrorage` points to a folder which will handle uploads and downloads of files. The directory
will be created according to the provided path unless exists.

### File storage
Uploaded and processed files are stored in `uploaded-files` and `processed-files` folders in the `/path/to/filestoragedir`
directory according to configuration parameters. Client and server-provided files are immediately deleted upon processing
completion.
