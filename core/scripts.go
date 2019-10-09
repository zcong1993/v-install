package core

const BbrScript = `
#!/bin/bash
set -eu

SYSCTL_FILE=/etc/sysctl.conf

# check root
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 
   exit 1
fi

# check OS version
source /etc/lsb-release
KERNEL_VERSION=$(uname -r)

if [ "$DISTRIB_ID" != "Ubuntu" ]; then
   echo "This script must be run under Ubuntu" 
   exit 1
fi

# install newest kernel
if [ "$DISTRIB_RELEASE" == "16.04" ]; then
    apt-get update -y
    apt-get install -y --install-recommends linux-generic-hwe-16.04
    apt-get autoremove -y
elif [ "$DISTRIB_RELEASE" == "18.04" ]; then
    echo "Kernel version enough, no need to install anything"
else
    # check kernel version
    if dpkg --compare-versions "$KERNEL_VERSION" "ge" "4.9"; then
        echo "WARNING: Non-LTS versions are not supported. Continuing since you have a compatible kernel."
    else
        echo "ERROR: Kernel auto install on Non-LTS versions is not supported. Please manually install kernel >= 4.9."
        exit 1
    fi
fi

if grep -q "tcp_bbr" "/etc/modules-load.d/modules.conf"; then
    echo "tcp_bbr" >> /etc/modules-load.d/modules.conf
fi

echo "Current configuration: "
sysctl net.ipv4.tcp_available_congestion_control
sysctl net.ipv4.tcp_congestion_control

# apply new config
if ! grep -q "net.core.default_qdisc=fq" "$SYSCTL_FILE"; then
    echo "net.core.default_qdisc=fq" >> $SYSCTL_FILE
fi
if ! grep -q "net.ipv4.tcp_congestion_control=bbr" "$SYSCTL_FILE"; then
    echo "net.ipv4.tcp_congestion_control=bbr" >> $SYSCTL_FILE
fi

# check if we can apply the config now
if lsmod | grep -q "tcp_bbr"; then
    sysctl -p $SYSCTL_FILE
    echo "BBR is available now."
elif modprobe tcp_bbr; then
    sysctl -p $SYSCTL_FILE
    echo "BBR is available now."
else
    echo "Please reboot to enable BBR."
fi
`
