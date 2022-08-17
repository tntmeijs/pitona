# PiTona
<a title="NathanLee at English Wikipedia, Public domain, via Wikimedia Commons" href="https://commons.wikimedia.org/wiki/File:TriumphDaytona675Side.jpg"><img width="512" alt="TriumphDaytona675Side" src="https://upload.wikimedia.org/wikipedia/commons/thumb/0/04/TriumphDaytona675Side.jpg/512px-TriumphDaytona675Side.jpg"></a>

Welcome to the PiTona project! This project is an educational project to learn more about the OBD2
protocol and the [Raspberry Pi](https://www.raspberrypi.com/) microprocessor.

## Prerequisites
- Raspberry Pi with WiFi and USB connectivity.
- Pre 2013 Triumph Daytona 675 (these ECUs are not as locked-down as the newer ones).

## Installation
1. Download [Raspberry Pi Imager](https://www.raspberrypi.com/software/).
2. Install the Debian Lite 64-bit operating system.
   1. Before installing, open the "advanced settings" (settings icon in the bottom right).
   2. Make sure to enable SSH and configure a network connection.
3. Download a webserver [release](https://github.com/tntmeijs/pitona/releases).
4. Copy `/install_pitona` and all of its contents to the `/home/pi` folder on your Raspberry Pi.
   1. On Windows, you could use [WinSCP](https://winscp.net/) to transfer your files.
5. Copy the `server` binary to `/home/pi`.
6. Ensure you can execute the webserver and installation script.
   1. Run `sudo chmod +x install_pitona/install.sh`
   2. Run `sudo chmod +x server`
7. With everything in place, execute `sudo ./install.sh` as on your Raspberry Pi.
8. Press `[Enter]` to reboot the Raspberry Pi after the script finishes executing.
9.  Give the device a minute or so to boot.
10. You should now see a new network with `PiTona_675` as its SSID.
11. Congratulations, your Raspberry Pi is now configured to run PiTona!

## Motivation
When I'm not writing code, I enjoy being out and about on my beloved Triump Daytona 675.

Unfortunately, I recently saw the infamous "check engine" light come on. Since I do not have access
to any OBD2 tools to read the ECU, I figured it'd be fun to write a little something myself. And
thus, the "PiTona" project was born!

I don't know how far I'm going to take this, but ultimately I'd like to at least be able to read my
bike's error code. Once I manage to do that, I'll see how far I can push my skills.

## Goals
The main purpose of this project is to see if it is possible to read OBD2 data from the onboard
ECU. However, simply logging data to a console is rather boring and too easy. To make things a
little more interesting, I'd like to eventually build an Android application that can display the
real-time data of the motorcycle's ECU.

If all goes well, I'd like to turn the application into a very neat tool to analyse my riding data.
It'd be awesome to have access to statistics such as lean angle, top speed, averge speed, GPS data,
and fuel efficiency!

## Technology stack
- [Raspberry Pi](https://www.raspberrypi.com/): hardware on which this whole thing runs
- [.NET](https://dotnet.microsoft.com/): programming language in which the webserver is written
- [React Native](https://reactnative.dev/): framework used to write the Android application
- [Bulma](https://bulma.io/): super neat CSS framework to make everything look pretty

## Help
> Services will not boot after updating the configuration files, even though the files are correct.

Check if your line endings are correct. The files should use `LF` line endings. If you save the
files using a Windows machine, chances are they are using `CRLF` line endings.

> I do not see the Raspberry Pi's network after rebooting.

Connect a keyboard and monitor to your Raspberry Pi to troubleshoot. The `journalctl` command might
come in handy to determine what exactly is failing.

> How can I SSH into my Raspberry Pi after installing PiTona?

PiTona is built to run in an isolated local network. Simply connect to your Raspberry Pi and SSH
into it using your favourite SSH agent. If you have used the default settings, try to SSH into
`pi@gw.wlan`.

> How can I change my DNS, AP, or other settings?

Either modify the configuration files in `/install_pitona` before you run the installation
script, or SSH into your Raspberry Pi and manually update the relevant configuration file(s).

## Disclaimer
Use this project at your own risk. There is a very real possibility that sending OBD2 commands,
without an understanding of what they do, will result in a broken ECU. This project is not
malicious in any way, shape, or form... however, I will not be held responsible for any damage,
issues, or other problems that might arise from the use of this software.
