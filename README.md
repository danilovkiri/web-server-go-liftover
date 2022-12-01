# USCS Liftover Server for 23andMe-formatted raw data

## Description
The project contains back-end and front-end source code for a web application providing positional liftover services
for 23andMe-formatted raw data files.

This web application is written in Golang and utilizes [pyliftover](https://github.com/konstantint/pyliftover) with the
code attached to a project. For more information on liftover dependencies see [documentation](pyliftover/README.md).

This web application uses authentication, cookies and sessions. Functionality is a subject of change and improvement.

Currently running on:
1. https://65.21.156.138:8080/index/ — as a standalone Docker container
2. https://65.21.156.138:8081/index/ — as a minikube single-node cluster with 3 replicas and liveness probes

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

Accessing the main page is available at `https://<HOST>:<PORT>/index/`. 

### Containerised

#### Docker

Build an image and run the container:

```bash
docker build .
docker container run --env AUTH_USERNAME=alice --env AUTH_PASSWORD=alice --env PORT=<service_PORT> --env HOST=0.0.0.0 \
--publish target=<service_PORT>,published=<vm_external_IP>:<vm_external_PORT>,protocol=tcp <imageID>
```

The service will be available worldwide at `https://<vm_external_IP>:<vm_external_PORT>/index/`.

#### Minikube

Install `kubectl`, Minikube and Docker, clone the current repository and navigate to its root directory. Run

```bash
# this will change the local docker daemon to the minikube's one
eval $(minikube docker-env)
# build an image inside the minikube's docker daemon
docker build -t go-liftover-server .
# change docker daemon back to local
eval $(minikube docker-env -u)
```

Create a deployment using the provided config via running

```bash
kubectl create -f deployment.yaml
```

where `deployment.yaml` has the following contents:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: go-liftover-server
  name: go-liftover-server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-liftover-server
  template:
    metadata:
      labels:
        app: go-liftover-server
    spec:
      containers:
      - name: go-liftover-server
        image: go-liftover-server:latest
        imagePullPolicy: Never
        ports:
        - containerPort: 8080
        livenessProbe:
          initialDelaySeconds: 30
          periodSeconds: 30
          timeoutSeconds: 1
          successThreshold: 1
          failureThreshold: 1
          httpGet:
            scheme: HTTPS
            path: /health
            port: 8080
        env:
        - name: PORT
          value: "8080"
        - name: AUTH_USERNAME
          value: "alice"
        - name: AUTH_PASSWORD
          value: "alice"
```

Expose a deployment via a NodePort service:

```bash
kubectl expose deployment go-liftover-server --name=liftover-service --type=NodePort
```

Set up a port forwarding:

##### kubectl port forwarding

Run the following code in the background where `vm_external_PORT` is the port you intend to listen to on your VM for
outside requests and `target_PORT` is the TargetPort of the k8s type=NodePort `liftover-service` service.

```bash
firewall-cmd --add-port=<vm_external_PORT>/tcp
kubectl port-forward --address 0.0.0.0 services/<k8s-service-name> <vm_external_PORT>:<target_PORT>
```

The service will be available worldwide at `https://<vm_external_IP>:<vm_external_PORT>/index/` where `vm_external_IP` is
the IPv4 address of your VM/server.

##### nginx proxy

Get `liftover-service`s externally available `service_HOST` and `service_PORT` parameters via running:

```bash
minikube service liftover-service --url
```

Install nginx and edit the `/etc/nginx/nginx.conf` to contain

```text
stream {
    server {
        listen 0.0.0.0:<vm_external_PORT>;
        proxy_pass <service_HOST>:<service_PORT>;
    }
}
```

Restart the nginx via running:

```bash
service nginx restart
```

The service will be available worldwide at `https://<vm_external_IP>:<vm_external_PORT>/index/` where `vm_external_IP` is
the IPv4 address of your VM/server.

## Miscellaneous

`CERT` and `KEY` files point to `.pem` files created locally via [`mkcert`](https://github.com/FiloSottile/mkcert)
or issued by any CA. `STORAGE` points to a folder which will handle uploads and downloads of files. The directory
will be created according to the provided path unless exists.

Uploaded and processed files are stored in `uploaded-files` and `processed-files` folders in the `STORAGE`
directory according to configuration parameters. Client and server-provided files are immediately deleted upon processing
completion.
