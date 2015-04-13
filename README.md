# Pod Connector
Pod connector is a CLI Tool for extracting information for a specific pod from
a sentinel config file or via a Redskull server.

Additionally, it can be used to connect directly to the pod's master given the podname.

# Configuration
Source information is configured via environment variables:

```shell
PODCONNECTOR_SENTINELCONFIGFILE = /path/to/config/file
PODCONNECTOR_REDSKULLADDRESS = IP:PORT
```

Currently, Redskull support is awaiting completion of the Red Skull RPC server.

# Usage


```shell
pod-connector -podname <podname> [-cli]
```

If neither of the source environment variables are available, pod-connector
will assume you want to use /etc/redis/sentinel.conf as the source.

To be useful you must use the flag `-podname <podname>` or `-podname=<podname>`
where '<podname>' is replaced by the name of the pod as it was registered in
Sentinel via the `MONITOR` command

## Info Dump

By default it will dump whatever information is available from the source.

## Connecting Directly
If you pass the flag `-cli`, pod-connector will execute redis-cli on your
behalf with all needed information to connect to the master - including
authentication. Essentially this saves you from needing to execute the
redis-cli yourself. You will be placed directly into the redis-cli shell on a
successful execution. Of course, this means redis-cli must be in your PATH.


# TODO
* Better formatting of pod information
* Option to set various sentinel variables such as down-in-milliseconds
* Update password option to update the password in sentinel (and optionally the instances?)
* Ability to pass redis-cli arguments to pod-connector to have it execute for you
* Implement Redskull RPC support

The options to update might go into a seaprate tool 'pod-configurator' to
better isolate capabilities.
