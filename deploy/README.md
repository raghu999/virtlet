# Deploying Virtlet

## Deploying Virtlet on a real cluster

See [this document](real-cluster.md) for instructions.

## Deploying Virtlet as a DaemonSet on kubeadm-dind-cluster

The steps described here are performed automatically by
[demo.sh](demo.sh) script.

1. Start [kubeadm-dind-cluster](https://github.com/Mirants/kubeadm-dind-cluster)
   with Kubernetes version 1.9 (you're not required to download it to your home directory):
```
$ wget -O ~/dind-cluster-v1.9.sh https://cdn.rawgit.com/Mirantis/kubeadm-dind-cluster/master/fixed/dind-cluster-v1.9.sh
$ chmod +x ~/dind-cluster-v1.9.sh
$ ~/dind-cluster-v1.9.sh up
$ export PATH="$HOME/.kubeadm-dind-cluster:$PATH"
```
   The cluster script stores appropriate kubectl version in `~/.kubeadm-dind-cluster`.

1. Label a node to accept Virtlet pod:
```
kubectl label node kube-node-1 extraRuntime=virtlet
```
1. Make several mounts shared inside the dind node that will run Virtlet:
```
for p in /dind /dev /boot /sys/fs/cgroup; do docker exec kube-node-1 mount --make-shared $p; done
```
1. Add virtlet image translation configmap:
```
kubectl create configmap -n kube-system virtlet-image-translations --from-file images.yaml
```
1. Install CRI proxy on the node:
```
CRIPROXY_DEB_URL="https://github.com/Mirantis/criproxy/releases/download/v0.10.0/criproxy-nodeps_0.10.0_amd64.deb"
docker exec kube-node-1 /bin/bash -c "curl -sSL '${CRIPROXY_DEB_URL}' >/criproxy.deb && dpkg -i /criproxy.deb && rm /criproxy.deb"
```
1. Deploy Virtlet DaemonSet (assuming that you have [virtlet-ds.yaml](virtlet-ds.yaml) in the current directory):
```
kubectl create -f virtlet-ds.yaml
```
1. Wait for Virtlet pod to activate:
```
kubectl get pods -w -n kube-system
```
1. Go to `examples/` directory and follow [the instructions](../examples/README.md) from there.

## Configuring Virtlet

Virtlet can be customized through the `virtlet-config` ConfigMap Kuberenetes object.
The following keys in the config map are honored by the `virtlet-ds.yaml`:

  * `disable_kvm` - disables KVM support and forces QEMU instead. Use "1" as a value.
  * `download_protocol` - default image download protocol - either `http` or `https`. The default is https.
  * `loglevel` - integer log level value for the virtlet written as a string (e.g. "3", "2", "1").
  * `calico-subnet` - netmask width for the Calico CNI. Default is "24".
  * `image_regexp_translation` - enables regexp syntax for the image name translation rules.
  * `disable_logging` - disables log streaming from VMs. Use "1" to disable.

## Removing Virtlet

In order to remove Virtlet, first you need to delete all the VM pods.

You can remove Virtlet DaemonSet with the following command:
```bash
kubectl delete daemonset -R -n kube-system virtlet
```

After that you can remove CRI Proxy if you're not going to use the
node for Virtlet again by undoing the steps you made to install it
(see CRI Proxy
[documentation](https://github.com/Mirantis/criproxy/)).
