#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail
set -o errtrace

testmode=
if [[ ${1:-} == -testmode ]]; then
  testmode=1
fi

if [[ -f /dind/vmwrapper ]]; then
  ln -fs /dind/vmwrapper /vmwrapper
fi


chown root:root /etc/libvirt/libvirtd.conf
chown root:root /etc/libvirt/qemu.conf
chmod 644 /etc/libvirt/libvirtd.conf
chmod 644 /etc/libvirt/qemu.conf

# Without this hack qemu dies trying to unlink
# '/var/lib/libvirt/qemu/capabilities.monitor.sock'
# while libvirt is querying capabilities.
# Removal of the socket below helps but not always.

if [[ -e /var/lib/libvirt/qemu ]]; then
  mv /var/lib/libvirt/qemu /var/lib/libvirt/qemu.ok
  mv /var/lib/libvirt/qemu.ok /var/lib/libvirt/qemu
fi

# export LIBVIRT_LOG_FILTERS="1:qemu.qemu_process 1:qemu.qemu_command 1:qemu.qemu_domain"
# export LIBVIRT_DEBUG=1

# only make vmwrapper suid in libvirt container
chown root.root /vmwrapper
chmod ug+s /vmwrapper

if [[ ${testmode} ]]; then
  # leftover socket prevents libvirt from initializing correctly
  rm -f /var/lib/libvirt/qemu/capabilities.monitor.sock
  /usr/local/sbin/libvirtd --listen --daemon
else
  # FIXME: try using exec liveness probe instead
  while true; do
    /usr/local/sbin/libvirtd --listen
    sleep 1
  done
fi
