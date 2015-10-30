# Veritas

Veritas is a CLI for getting at Diego's truth.

## Downloading on a BOSH VM

For a Linux build on a BOSH VM (the Cells are best):

```bash
  pushd $HOME
  wget https://github.com/pivotal-cf-experimental/veritas/releases/download/latest/veritas -O ./veritas
  chmod +x ./veritas

  echo "export PATH=$PATH:$PWD" > veritas.bash
  echo "export DROPSONDE_ORIGIN=veritas" >> veritas.bash
  echo "export DROPSONDE_DESTINATION=localhost:3457" >> veritas.bash
  ./veritas autodetect >> veritas.bash
  ./veritas completions >> veritas.bash

  source ./veritas.bash
  popd
```

Once this is done, you simply need to `source ~/veritas.bash` when you log in again.


## Downloading on an OS X Workstation

For an OS X build (mainly useful for chugging logs locally, or connecting to the BBS on a local BOSH-Lite):

```bash
mkdir -p $HOME/bin
pushd $HOME/bin
  wget https://github.com/pivotal-cf-experimental/veritas/releases/download/latest/veritas-osx
  mv veritas-osx veritas
  chmod +x ./veritas
popd
```

## Connecting to the BBS

`veritas` commands must include the location of the Diego BBS server. The BBS is not typically publically routable, so run `veritas` from a VM in the same subnet. On the private network the BBS can be found at `https://bbs.service.cf.internal:8889`. When testing locally against Bosh Lite, you can run `veritas` locally and use the IP address of the `database_z1` job, `https://10.244.16.130:8889`.

The URL for BBS is specified with the environment variable `BBS_ENDPOINT`. When SSL is disabled on the BBS, specify the URL scheme as `http`:

```bash
BBS_ENDPOINT=http://bbs.service.cf.internal:8889 veritas dump-store
```

BBS support for SSL uses mutual authentication, meaning the client must also provide a certificate. Create files containing the certificate and private key and reference them with each `veritas` using the `BBS_CERT_FILE` and `BBS_KEY_FILE` environment variables. Also, use `https` in the scheme for `BBS_ENDPOINT`. For the purposes of testing with BOSH Lite, you can use the certificate and key for any BBS client in the BOSH manifest.

```bash
BBS_ENDPOINT=https://bbs.service.cf.internal:8889 BBS_CERT_FILE=path/to/client/cert BBS_KEY_FILE=path/to/client/key veritas dump-store
```

Instead of environment variables, BBS configuration parameters may be supplied with the flags `--bbsEndpoint`, `--bbsCertFile`, and `--bbsKeyFile`. For commands with positional arguments, such as `desire-lrp` or `remove-lrp`, the flags must be given **after** the command but **before** the positional arguments. For example:

```bash
veritas remove-lrp \
  --bbsEndpoint=https://bbs.service.cf.internal:8889 \
  --bbsCertFile=path/to/client/cert \
  --bbsKeyFile=path/to/client/key \
  my-process-guid

```

### Common Errors

```bash
$ veritas distribution
Failed to print distribution
Post https://bbs.service.cf.internal:8889/v1/desired_lrps/list: tls: oversized record received with length 20527
```

- This error means that `veritas` was configured with an `https` URL for its BBS endpoint when the BBS expected plain HTTP.


```bash
$ veritas distribution
Failed to print distribution
Post http://bbs.service.cf.internal:8889/v1/desired_lrps/list: http: ContentLength=2 with Body length 0
```

- This error means that `veritas` was configured with an `http` URL for its BBS endpoint when the BBS expected HTTPS.



## Creating, Updating, and Removing LRPs

Veritas can submit, update, and remove DesiredLRPs with the `veritas desire-lrp`, `veritas update-lrp` and `veritas remove-lrp` commands.

### Desiring an LRP

`veritas desire-lrp <path-to-json-file>` takes the path to a file.  This file should contain a JSON representation of the DesiredLRP.  

Two examples:

```
{
    "process_guid":"grace-1",
    "domain":"test",
    "rootfs":"docker:///onsi/grace-busybox",
    "instances":1,
    "ports":[
        8080
    ],
    "action":{
        "run":{
            "path":"/grace",
            "args":[
                "-chatty"
            ],
            "dir":"/tmp",
            "user":"root"
        }
    },
    "routes":{
        "cf-router":[
            {
                "hostnames": [
                  "grace.app-domain.com"
                ],
                "port": 8080
            }
        ]
    }
}
```

```
{
    "process_guid": "92bcf571-630f-4ad3-bfa6-146afd40bded",
    "domain": "redis-example",
    "rootfs": "docker:///redis",
    "instances": 1,
    "ports": [
        6379
    ],
    "action": {
        "run": {
            "path": "/entrypoint.sh",
            "args": [
                "redis-server"
            ],
            "dir": "/data",
            "user": "root"
        }
    },
    "routes": {
        "tcp-router": [
            {
                "external_port": 50000,
                "container_port": 6379
            }
        ]
    }
}
```

### Updating an LRP

`veritas update-lrp <process-guid> <path-to-json-file>` take a process guid and a path to a file.  This file should contain a JSON representation of a `DesiredLRPUpdate`.  For example:

```
{
    "instances": 3,
    "routes": {
        "tcp-router": [
            {
                "external_port": 50001,
                "container_port": 6379
            }
        ]
    },
    "annotation": "Hey, don't forget to delete me when you're done!"
}
```

### Removing an LRP

`veritas remove-lrp <process-guid>` will remove the LRP with the specified process guid. As long as the domain on the LRP is fresh, this action will also stop and destroy any containers associated to the LRP.

### Updating a domain

`veritas set-domain <domain-name> <duration>` will mark the specified domain as fresh for the given duration. This allows Diego to take destructive actions on instances of LRPs in this domain.

Example:

```bash
$ veritas set-domain redis-example 120s
```

Setting the duration to `0` will keep the domain fresh indefinitely. To delete a domain, set it to have a very short duration (say, 1s), and wait.

### Fetching data

- `veritas get-desired-lrp <process-guid>` fetches and outputs the DesiredLRP with the associated process guid.
- `veritas get-actual-lrp <process-guid>` fetches all ActualLRPs associated with the process guid.
- `veritas get-actual-lrp <process-guid> <index>` fetches the ActualLRP with index `<index>` associated with the process guid.
- `veritas dump-store` emits a formatted representation of the contents of the cluster.
