# PiTona
<a title="NathanLee at English Wikipedia, Public domain, via Wikimedia Commons" href="https://commons.wikimedia.org/wiki/File:TriumphDaytona675Side.jpg"><img width="512" alt="TriumphDaytona675Side" src="https://upload.wikimedia.org/wikipedia/commons/thumb/0/04/TriumphDaytona675Side.jpg/512px-TriumphDaytona675Side.jpg"></a>

Welcome to the PiTona project! This project is an educational project to learn more about the
[OBD-II](https://en.wikipedia.org/wiki/On-board_diagnostics#OBD-II) protocol and the
[Raspberry Pi](https://www.raspberrypi.com/) microprocessor.

## Prerequisites
- Raspberry Pi with WiFi and USB connectivity.
- Pre 2013 Triumph Daytona 675 (these ECUs are not encrypted).

## Motivation
When I'm not writing code, I enjoy being out and about on my Triump Daytona 675.

Unfortunately, I recently saw the "check engine" light come on. Since I do not have access to any
OBD2 tools to read the ECU, I figured it'd be fun to write a little something myself. And thus, the
"PiTona" project was born!

I don't know how far I'm going to take this, but ultimately I'd like to at least be able to read my
bike's error code. Once I manage to do that, I'll see how far I can push it.

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
- [Go](https://go.dev/): host server programming language
- [React](https://reactjs.org/): front-end framework
- [Bulma](https://bulma.io/): super neat CSS framework to make everything look pretty

## Endpoints
The following endpoints can be used to interface with the Daytona 675's ECU:
- <span style="color:#27ae60">GET</span> `/api/v1/obdii/debug` - ⚠ **DANGER** ⚠ send raw data to the ECU (use with caution)
- <span style="color:#27ae60">GET</span> `/api/v1/obdii/01` - Send a supported PID from [mode 01](https://en.wikipedia.org/wiki/OBD-II_PIDs#Service_01_-_Show_current_data)
- <span style="color:#27ae60">GET</span> `/api/v1/obdii/03` - Request active DTC  from [mode 03](https://en.wikipedia.org/wiki/OBD-II_PIDs#Service_03_-_Show_stored_Diagnostic_Trouble_Codes_(DTCs))
- <span style="color:#27ae60">GET</span> `/api/v1/system/status` - Show Raspberry Pi system information
- <span style="color:#e74c3c">DELETE</span> `/server` - Gracefully stop the server

## Supported PIDs
The list below contains all PIDs that have been confirmed to work on my 2008 Triumph Daytona 675.

### Service 01 - show current data
| PID (HEX) | Description                        | Comments |
| --------- | ---------------------------------- | -------- |
| 00        | List supported PIDs [0x01 to 0x20] |          |

### Service 03 - show stored DTC
| PID (HEX) | Description         | Comments                                                                                                      |
| --------- | ------------------- | ------------------------------------------------------------------------------------------------------------- |
| N/A       | List all stored DTC | Not implemented yet - need to implement the [ISO 15765-2](https://en.wikipedia.org/wiki/ISO_15765-2) protocol |

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

# Development log
## September
### 25<sup>th</sup> of September 2022
- Add file serving capability to webserver
- Add client React project boilerplate code
- Add build script to automatically create a [tarball](https://en.wikipedia.org/wiki/Tar_(computing))

### 24<sup>th</sup> of September 2022
- Code clean-up
- Update README to include endpoint and PID documentation
- Add debug endpoint to send arbitrary data to the ECU

### 10<sup>th</sup> of September 2022
- Discovered that the `3033` "DTC" is not really a fault code. Instead, this response is most
  likely the start of a ISO-TP frame. Parsing this data is relatively difficult, which is why I
  will work on it once the application is a bit more mature.

### 5<sup>th</sup> of September 2022
- Major refactor of the codebase
- New structure makes it easier to add new functionality

### 4<sup>th</sup> of September 2022 🏆
- Implemented simple serial communication logic
- Implemented endpoints to communicate with the server
- Managed to send commands to the ECU
  - `0100` returns `30313030`, which indicates the PIDs supported by this ECU
  - `03` returns `3033`, which refers to a status code, but I have not been able to decode it yet

![first ECU response](media/first_time_reading_ecu.png)

### 2<sup>nd</sup> of September 2022
- Switched from Kotlin / Spring Boot to Go

## August
### 28<sup>th</sup> of August 2022
- Switched from C# / .NET to Kotlin / Spring Boot

### 21<sup>st</sup> of August 2022
- Add circular buffer implementation
- Add unit tests

### 20<sup>th</sup> of August 2022
- Add serial port reading logic

### 16<sup>th</sup> of August 2022
- Project set-up
- Simple .NET server
- Tried to turn the Raspberry Pi into a local access point
