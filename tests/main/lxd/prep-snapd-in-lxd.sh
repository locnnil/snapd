#!/bin/sh

set -ex

# In some lxd images snapd could not be installed
if ! command -v snap; then
    apt install -y snapd
fi

# TODO: Remove the refresh once the issue https://github.com/lxc/lxd/issues/10079 is release to 4.0/candidate
# Make sure the lxd snap is updated before removing it
for _ in $(seq 30); do
    if snap changes | grep -qE "Done.*Initialize device"; then
        break
    fi
    sleep 1
done
snap wait system seed.loaded
if snap list lxd; then
    snap refresh lxd --channel=latest/stable
    snap remove lxd
fi

# XXX: remove once the "umount /snap" change in postrm has propagated all
#      the way to the image
if [ -e /var/lib/dpkg/info/snapd.postrm ]; then
    # ensure we can umount /snap
    sed -i 's#echo "Final directory cleanup"#umount /snap || true#' /var/lib/dpkg/info/snapd.postrm
fi

apt autoremove --purge -y snapd ubuntu-core-launcher
apt update

# requires the snapd deb to already have been "lxd file push"d into the 
# container
apt install -y /root/snapd_*.deb

# reload to take effect of the proxy that may have been set before this script
# XXX: systemctl daemon-reload times out in 16.04:my-nesting-lxd but every
#      appears to be working normally
systemctl daemon-reload || true
systemctl restart snapd.service

# wait for snapd to finish seeding
snap wait system seed.loaded

# for debugging
cat /etc/environment
