# Veritas

Veritas is a cli for getting at Diego's truth.

## Downloading on a BOSH VM

For a linux build on a bosh vm (the Cells are best):

```bash
 pushd $HOME
 wget http://onsi-public.s3.amazonaws.com/veritas -O ./veritas
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

For an OS X build (mainly for chugging logs locally):

```bash
  mkdir -p $HOME/bin

  pushd $HOME/bin

  wget http://onsi-public.s3.amazonaws.com/veritas-osx
  mv veritas-osx veritas
  chmod +x ./veritas

  popd
```

## Connecting to the BBS

- As Vertias no longer detects the location of the Diego BBS, or has a default, you must tell it where the BBS server is with the environment variable `BBS_ENDPOINT` with each command. With TLS disabled, BBS can be reached at `http://bbs.service.cf.internal:8889`

 Example:

  `$ BBS_ENDPOINT=http://bbs.service.cf.internal:8889 veritas dump-store`

- When TLS is enabled, the `BBS_CERT_FILE` and `BBS_KEY_FILE` environment variables must also be provided.

 Example:

```bash
  BBS_ENDPOINT=http://bbs.service.cf.internal:8889 \
  BBS_CERT_FILE=path/to/cert \
  BBS_KEY_FILE=path/to/key \
  veritas dump-store
```

## Launch and update an LRP

Veritas can submit/remove DesiredLRPs and DesiredLRPUpdates with the `veritas desire-lrp`, `veritas update-lrp` and `veritas remove-lrp` subcommands.

### Desiring an LRP

`veritas desire-lrp <path to json file>` takes the path to a file.  This file should contain a JSON representation of the DesiredLRP.  For example:

```
{
    "process_guid":"92bcf571-630f-4ad3-bfa6-146afd40bded",
    "domain":"redis-example",
    "root_fs":"docker:///redis",
    "instances":1,
    "ports":[
        6379
    ],
    "action":{
        "run_action":{
            "path":"/entrypoint.sh",
            "args":[
                "redis-server"
            ],
            "dir":"/data",
            "user":"root"
        }
    },
    "routes":{
        "tcp-router":[
            {
                "external_port":50000,
                "container_port":6379
            }
        ]
    }
}
```

### Updating an LRP

`veritas update-lrp <process-guid> <path to json file>` take a process guid and a path to a file.  This file should contain a JSON representation of a `DesiredLRPUpdate`.  For example:

```
{
    "instances": 3,
    "routes":{
        "tcp-router":[
            {
                "external_port":50001,
                "container_port":6379
            }
        ]
    }
}
```

### Removing an LRP

`veritas remove-lrp <process-guid>` will remove the LRP with associated process guid.  This will shut down any associated containers.

### Fetching data

- `veritas get-desired-lrp <process-guid>` fetches and outputs the DesiredLRP with the associated process guid
- `veritas get-actual-lrp <process-guid>` fetches all ActualLRPs associated with the process guid
- `veritas get-actual-lrp <process-guid> <index>` fetches the ActualLRP with index `<index>` associated with the process guid
- `veritas dump-store` emits a formatted representation of the contents of the cluster
