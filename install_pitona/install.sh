#!/bin/bash
if [[ $EUID > 0 ]]
  then echo "Please run this script as root"
  exit -1
fi

echo "Install PiTona on your system - your system will reboot at the end of the installation? [yes/no]"
read proceed

if [ $proceed != "yes" ]
then
echo "Type \"yes\" to install - since you did not, the installation will be aborted."
exit 0
fi

# Reference for this script:
# https://www.raspberrypi.com/documentation/computers/configuration.html#setting-up-a-routed-wireless-access-point

# HostAPD is used to turn this device into an access point
echo "====================="
echo "Installing HostAPD..."
echo "=====================\n"
apt install hostapd -y
systemctl unmask hostapd
systemctl enable hostapd

# Dnsmasq is used for network management services (DNS / DHCP)
echo "====================="
echo "Installing Dnsmasq..."
echo "=====================\n"
apt install dnsmasq -y

# Static IP address configuration
echo "============================================="
echo "Configuring device to use a static IP address"
echo "=============================================\n"
cp dhcpcd.conf /etc/dhcpcd.conf

# Dnsmasq configuration
mv /etc/dnsmasq.conf /etc/dnsmasq.conf.BACKUP
cp dnsmasq.conf /etc/dnsmasq.conf

# Ensure WiFi radio is not blocked
rfkill unblock wlan

# Access point configuration
echo "============================================="
echo "Configuring device as a wireless access point"
echo "=============================================\n"
cp hostapd.conf /etc/hostapd/hostapd.conf

# Add the webserver as a service
echo "===================================================================="
echo "Adding the webserver as a Linux service to ensure it runs on startup"
echo "====================================================================\n"
cp pitona.service /etc/systemd/system/pitona.service
systemctl enable pitona

# Reboot to ensure everything works as intended
echo "========="
echo "All done!"
echo "=========\n"
echo "The device will now reboot - please look for a new WiFi network once it comes back online"

read -p "Press enter to continue"
reboot
