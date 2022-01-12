# USCS Liftover Server for 23andMe-formatted raw data

## Description
The project contains back-end and front-end source code for a web application providing positional liftover services
for 23andMe-formatted raw data files.

This web application is written in Golang and utilizes [pyliftover](https://github.com/konstantint/pyliftover) with the
code attached to a project. For more information on liftover dependencies see [documentation](./liftover/README.md).

This web application uses authentication, cookies and JS. Functionality is a subject of change and continuous
improvement.

## Running
To start the local server run the following command in a terminal specifying `AUTH_USERNAME`, `AUTH_PASSWORD` and a
path (either relative or absolute) to a configuration file `CONFIG`.
```bash
AUTH_USERNAME=alice AUTH_PASSWORD=alice CONFIG=./resources/config.yaml go run main.go
```

Accessing the main page is available at `https://<serverIP>:<serverPort>/index/`

## Notes

### Config
The [config.yaml](./resources/config.yaml) file implements the following self-explanatory structure:
```yaml
constants:
  certFile: ./resources/localhost.pem
  keyFile: ./resources/localhost-key.pem
  serverIP: 192.168.1.68
  serverPort: 4000
```
where `certFile` and `keyFile` point to `.pem` files created locally via [`mkcert`](https://github.com/FiloSottile/mkcert).

### File storage
Uploaded and processed files are stored in `temp-files` and `output-files` folders in the project directory
(to be created manually). No autoamtic file removal is implemented. Further updates will enhance file handling logic.
