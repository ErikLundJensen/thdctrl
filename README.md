# Talos Hetzner Dedicated Control CLI - thdctrl


Initialise Hetzner servers by using their serverNumber


thdctrl --command=init --serverNumber=123456

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