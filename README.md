# Talos Hetzner Dedicated Control CLI - thdctrl

## Overview

`thdctrl` is a command-line tool to manage Hetzner dedicated servers with Talos. It provides various commands to initialize, configure, and manage your servers.

## Installation

To build and install `thdctrl`, use the provided Dockerfile:

```sh
docker build -t thdctrl .
docker run --rm -v $(pwd):/root thdctrl:latest /app/thdctrl --help
```

## Usage

Use `thdctrl --help` to get a list of available commands and arguments.  
Username and password for the Hetzner Robot API must be set using environment variables:
```
export HETZNER_USERNAME='myAPIuser'
export HETZNER_PASSWORD='password'
```


### Commands

#### `init`

Initialize Hetzner servers by using their server number.

```sh
thdctrl init <serverNumber>
```

Example:

```sh
thdctrl init 123456
```


### Flags

- `--help`: Show help information for `thdctrl` commands.
- `--version`: Show the version of `thdctrl`.

## Example Workflow

1. Initialize the server:

    ```sh
    thdctrl init 123456
    ```

2. Wait for "waiting for bootstrap" and then bootstrap Talos:

    ```sh
    talosctl bootstrap
    ```

3. Wait for the API server to be ready, then apply the configuration:

    ```sh
    . ./init-env-sh
    . ./generate-config.sh
    talosctl apply
    ```

4. Apply the Cilium configuration:

    ```sh
    ./gen-cilium.sh
    kubectl apply -f gen/cilium.yaml
    ```

5. Reboot the servers:

    ```sh
    talosctl reboot
    ```

6. Wait for the nodes to be ready and open the Talos dashboard:

    ```sh
    talosctl dashboard
    ```

7. Watch the pods get healthy:

    ```sh
    kubectl get pods -A
    ```

## New Features

- Get disk name during SSH sessions (e.g., if the disk is not specified in the command line).
- Add a command to list disks and sizes.

## TODO

- Add shutdown command.
- Change node-2 to a worker node.
- Re-initialize nodes.
- Add VIP address (in case of more control plane nodes).
- Install Hetzner Load Balancer operator.
- Test load balancer.
