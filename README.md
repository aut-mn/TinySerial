# TinySerial

TinySerial is a simple and easy-to-use serial port listener written in Go
It allows you to listen to data coming from serial ports on your system.

## Installation

> [!IMPORTANT]
> TinySerial does not work on Windows
> It can be modified to work on MacOS with minimal effort

Clone the repository and navigate to the project directory, then run the following command to build the project:

```bash
go build main.go -o TinySerial
```

## Usage

You can run the program with the following command:

```bash
./TinySerial
```

The program accepts the following command line arguments:

| Flag | Description              | Default Value |
|:-----|:-------------------------|--------------:|
| `-b` | Set the baud rate        |          9600 |
| `-d` | Set the data bits        |             8 |
| `-s` | Set the stop bits        |             2 |
| `-h` | Display help information |           N/A |

After starting the program, it will list all available serial ports. Enter the number of the port you want to listen to