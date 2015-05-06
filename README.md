[![Build Status](https://travis-ci.org/sentinel-tools/pod-connector.svg?branch=master)](https://travis-ci.org/sentinel-tools/pod-connector)


# Pod Connector
Pod connector is a CLI Tool for extracting information for a specific pod from
a sentinel config file.

Additionally, it can be used to connect directly to the pod's master given the podname.


# Usage


```shell
pod-connector [-sentinelconfig "/path/to/file"] (info [-j]|cli) podname
```

By default the expected Sentinel config file is "/etc/redis/sentinel.conf"

## Info Dump

By default it will dump whatever information is available from the source. You
can choose to let it be formatted for humans, or pass `-json` or `-j` to get it
as a JSON encoded string.

## Connecting Directly
If you use the `cli` command, pod-connector will execute redis-cli on your
behalf with all needed information to connect to the master - including
authentication. Essentially this saves you from needing to execute the
redis-cli yourself. You will be placed directly into the redis-cli shell on a
successful execution. Of course, this means redis-cli must be in your PATH.

# Which Pod Tool?
Why pod-connector and pod-manager? Pod Connector's primary purpose is to
provide connectivity to the specified pod, with a bit of info available. Pod
Manager, however, is designed to make changes to sentinel and in some cases
directly to the instances in the pod.

You might give access to pod-connector people whom you are willing to allow
connectivity to the pod, but reserve pod-manager access for those who manage
sentinel and the pod. Different needs, different tools. This is much simpler
than yet another user system in a tool.

# TODO
* Ability to pass redis-cli arguments to pod-connector to have it execute for you


# Bash Completion

Add the pod-connector_completion script to your completions directory and/or
source it directly to add bash completion.
