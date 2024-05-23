package main

import (
	"bufio"
	"flag"
	"fmt"
	"go.bug.st/serial"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	helpFlag := flag.Bool("h", false, "Show help")
	baudRate := flag.Int("b", 9600, "BaudRate")
	dataBits := flag.Int("d", 8, "DataBits")
	stopBits := flag.Int("s", 2, "StopBits")
	flag.Parse()
	if *helpFlag {
		fmt.Println("TinySerial is a simple serial port listener.")
		fmt.Println("Usage: TinySerial [-b baudRate] [-d dataBits] [-s stopBits] [-h]")
		flag.PrintDefaults()
		return

	}
	files, err := os.ReadDir("/dev")
	if err != nil {
		fmt.Println("Error reading /dev directory: ", err)
		return
	}
	var ports []string
	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, "tty") || strings.HasPrefix(name, "cu.") {
			ports = append(ports, "/dev/"+name)
		}
	}
	if len(ports) == 0 {
		fmt.Println("No serial ports found.")
		return
	}
	for i, port := range ports {
		fmt.Printf("%d#\t%s\n", i+1, port)
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("TinySerial> ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	num, err := strconv.Atoi(text)
	if err != nil {
		fmt.Println("Invalid selection.")
		return
	}
	selectedPort := ports[num-1]
	options := serial.Mode{
		BaudRate: *baudRate,
		DataBits: *dataBits,
		StopBits: serial.StopBits(*stopBits),
	}
	port, err := serial.Open(selectedPort, &options)
	if err != nil {
		fmt.Println("Error opening port: ", err)
		return
	}
	fmt.Println("Listening on port: ", selectedPort)
	defer func() {
		err = port.Close()
		if err != nil {
			fmt.Println("Error closing port: ", err)
		}
	}()
	quit := make(chan struct{})
	go func() {
		<-sigChan
		err := port.Close()
		if err != nil {
			fmt.Println("Error closing port: ", err)
		}
		close(quit)
		os.Exit(0)
	}()
	buf := make([]byte, 128)
	for {
		select {
		case <-quit:
			return
		default:
			input := bufio.NewReader(os.Stdin)
			readInput, _ := input.ReadString('\n')
			if readInput == "q!" {
				return
			}
			_, _ = port.Write([]byte(readInput))
			n, err := port.Read(buf)
			if err != nil {
				fmt.Println("Error reading from port: ", err)
				return
			}
			fmt.Println(string(buf[:n]))
		}
	}
}
