# Talos Hetzner Dedicated Control CLI - thdctrl


Initialise Hetzner servers by using their serverNumber


thdctrl init 123456

talosctl bootstrap

talosctl apply

talosctl dashboard

. ./init-env-sh
./gen-cilium.sh

kubectl apply -f gen/cilium.yaml

# Reboot servers after applying Cilium configuration
talosctl reboot

# Watch pods get healthy
kubectl get pods -A



## thdctrl

Use `thdctrl --help` to get commands and arguments


## New features
Get disk name during SSH sessions (e.g. if disk is not specified in the command line).
Or add command to list disks and sizes.

