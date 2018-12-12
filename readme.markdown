# Scylla Test

The following repository aims to provide a Command and Control (C&C) system
for managing "unlimited" number of servers to perform tasks and report back.

## The idea behind the implementation:

There are three levels of services

0. Dispatcher - CnC server
1. Orchestration - Sends tasks to nodes or sub orchestration servers
3. Nodes - Actual workers that does the code.

### Flow idea

The dispatcher has X Orchestration servers that it knows about, having them
registered and reporting an heartbeat every second to know that they are up.

The Orchestration server has registered X amount of either nodes or additional
orchestration servers.
The server sends heartbeat to the Dispather if it is connected directly to it,
and sends actions to nodes or other orchestration servers.

It gathers all information arrived from all nodes and other orchestration servers
and pass it up to the dispatcher.

The node server is a worker that sends a heartbeat to a known orchestration server
every second, and listen to an action to do.

When an action is given, it execute it, and reports back what was the result.

# How to use

There is the `cmd` directory. It contains `nodes` and `orchestration` sub directories.

The `nodes` directory are the workers nodes described under this readme.

The `orchestration` directory contains the orchestration server described under this readme.

# How I tested

Created the following file:

    l: 216
    t: 2

    {
      "check_etc_hosts_has_8888": {
        "path": "/etc/hosts",
        "type": "file_contains",
        "check": "8.8.8.8"
      },
      "check_log_file_exists": {
        "path": "/var/log/messages.log",
        "type": "file_exists"
      }
    }

Then used netcat like so:

    nc 127.0.0.2:8082 < test.request

