# LGTVShutDowner

I want to shut down my Screen at the same time as the PC

## Description

> **This tool is only intended to work on Linux.**
> For Windows users, please go see [Maassoft/ColorControl](https://github.com/Maassoft/ColorControl)

If used correctly, this tool create a system tray icon with wich you can send different commands to your WebOS screen.

The tool can be runned as a service with a cmd argument at login to turn the screen on automatically.(`-cmd PowerOn`)

## Installation

Download latest [release](https://github.com/Prouk/LGTVShutDowner/releases) or build binary from source.

Classic installation should be like that :

- Follow the `.config` like shown in the folder tree of the archive.
- Create a LGTVShutDown folder in '/opt/' and put the file located in the binary folder of archive.
- Complete the config file with your LGTVWebOs IP (normal communication) and MAC adress (wake on LAN).
- Verify / Complete the LGTVShutDown.service file to match your files paths.
- `systemctl --user enable LGTVShutDowner --now` to register and launch the service.

## Usage

### !You need to launch the service while the TV is ON, and connect to it at least one time !

- Start the service and click on connect in the tray menu, then accept the connection on your TV.
- Automatic screen turn on, on user session login.
- Automatic screen turn on, on PC ShutDown.
- Tray icon with interactions (connect, turn on, turn off, ping, list api for development purpose)
