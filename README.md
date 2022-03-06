#Tread on Me
A Go library and integrated web server for interacting with BLE equipped treadmills by Sole (at least the Sole F80)
from Linux.

I wanted to be able to get realtime telemetry from my treadmill and found the official app to be lacking flexibility.

## Warning -- Please Read
This software allows you to change the configuration of your treadmill while it is operating; **this could be dangerous**.

While developing the software I ran across configurations that could cause the treadmill to disable the physical
controls on the treadmill while the motor was engaged and running. The code and documentation presented here is based off
trial and error while testing with my own treadmill, yours may act differently. Please use this software at your own risk.

## Requirements
You'll need a Bluetooth LE adapter and a modern version of Linux. I'm using a Raspberry Pi 4 with its integrated BLE
module. You'll also need Go version 1.16 or higher installed.

## Installation
You can install this by running the following command:

    go install github.com/swedishborgie/treadonme/webserver@latest

You can then run the application like this (you can use `hcitool lescan` if you need to discover the mac address):

    webserver <mac address of treadmill>

Loading http://localhost:8080 (or http://your.ip.address:8080) should result in seeing a dashboard showing the current
state of the application.

The treadmill does not listen for connections while in low power mode, the display must be active in order to connect.

## Bluetooth LE Technical Details
The treadmill appears to use a fairly common integrated BLE to UART module. It advertises the following:

* Service ID: `49535343-FE7D-4AE5-8FA9-9FAFD205E455`

The service exposes two different characteristics that are required to communicate:

 * Read (Notify) - `49535343-1E4D-4BD9-BA61-23C647249616`
 * Write - `49535343-8841-43F4-A8D4-ECBE34729BB3`

To perform serial communication you'll subscribe to notifications for the read characteristic to receive messages from
the treadmill, and you'll write messages to the write characteristic when you want to communicate.

The messages sent back and forth between the treadmill use a fairly straight forward serial protocol:
 * All messages begin with `0x5b` followed immediately by a byte indicating the message length. All messages end with `0x5d`.
 * Most messages seem to require an acknowledgement and there are two types: an ACK message (`0x00`) or repeating the
   received command back to the treadmill.
 * The first command you send should be `Get Device Info (0xf0)` which will result in the treadmill establishing
   communication.
 * Some messages are sent from the treadmill without any prompting from the host. These messages need to be acknowledged
   for communication to continue.
 * The treadmill cannot handle messages too quickly, even if messages have been promptly acknowledged. A short sleep (300-500ms)
   is usually sufficient between writes to let the treadmill catch up.

## Messages
| Name             | Code   | Request                                | Response                 | ACK Type     | Direction         |
|------------------|--------|----------------------------------------|--------------------------|--------------|-------------------|
| Acknowledge      | `0x00` | N/A                                    | `5b0400094f4b5d`         | None         | Both              |
| Set Workout Mode | `0x02` | `5B0202025D`                           | `5B0202025D`             | Echo         | Host -> Treadmill |
| Workout Mode     | `0x03` | `5B0203015D`                           | `5B0203015D`             | Echo         | Treadmill -> Host |
| Workout Target   | `0x04` | `5B05040A0000005D`                     | `5b0400044f4b5d`         | ACK          | Host -> Treadmill |
| Workout Data     | `0x06` | `5b0f06093b0000000000050000000000015d` | `5b0400064f4b5d`         | ACK          | Treadmill -> Host |
| User Profile     | `0x07` | `5B06070123009B435D`                   | `5b0400074f4b5d`         | ACK          | Host -> Treadmill |
| Program Type     | `0x08` | `5b030810015d`                         | `5b0400084f4b5d`         | ACK          | Host -> Treadmill |
| Heart Rate Type  | `0x09` | `5b030901005d`                         | `5b0400094f4b5d`         | ACK          | Treadmill -> Host |
| Error Code       | `0x10` | `5b0210005d`                           | `5b0400104f4b5d`         | ACK          | Treadmill -> Host |
| End Workout      | `0x32` | `5b0a320013000000000800005d`           | `5b0400324f4b5d`         | ACK          | Treadmill -> Host |
| Get Device Info  | `0xF0` | `5B01F05D`                             | `5B08F092000178050F125D` | Echo (Kinda) | Host -> Treadmill |
