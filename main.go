package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const GopherPort = 70

type ItemKind int

const (
	ItemKindFILE      ItemKind = '0'
	ItemKindDIR       ItemKind = '1'
	ItemKindPHONE     ItemKind = '2'
	ItemKindERROR     ItemKind = '3'
	ItemKindMACBINHEX ItemKind = '4'
	ItemKindDOSBIN    ItemKind = '5'
	ItemKindUNIXUU    ItemKind = '6'
	ItemKindIDXSEARCH ItemKind = '7'
	ItemKindTELNET    ItemKind = '8'
	ItemKindBIN       ItemKind = '9'
	ItemKindRSERVER   ItemKind = '+'
	ItemKindTN3270    ItemKind = 'T'
	ItemKindGIF       ItemKind = 'g'
	ItemKindIMG       ItemKind = 'I'
	ItemKindINFO      ItemKind = 'i'
)

type gopherLine struct {
	kind     ItemKind
	content  string
	selector string
	domain   string
	port     int
}

func processLine(line string) gopherLine {
	var l gopherLine
	l.kind = ItemKind(line[0])
	line = line[1:]

	ptr := 0
	col := 0
	for i, r := range line {
		if r == '\t' {
			switch col {
			case 0:
				l.content = line[:i]
				ptr = i + 1
			case 1:
				l.selector = line[ptr:i]
				ptr = i + 1
			case 2:
				l.domain = line[ptr:i]
				port, _ := strconv.Atoi(line[i+1:])
				l.port = port
			}

			col++
		}
	}
	
	return l
}

func openconn(addr string) (*bufio.ReadWriter, error) {
	log.Println(fmt.Sprintf("Dial %s\n", addr))

	hp := fmt.Sprintf("%s:%v", addr, GopherPort)
	log.Println(hp)
	conn, err := net.Dial("tcp", hp)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial %s\n")
	}

	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

func processpage(rw *bufio.ReadWriter) ([]gopherLine, error) {
	var page []gopherLine
	scanner := bufio.NewScanner(rw)
	for scanner.Scan() {
		page = append(page, processLine(scanner.Text()))
	}
	
	if scanner.Err() != nil {
		return nil, fmt.Errorf("processpage: %s", scanner.Err());
	}
	
	return page, nil
}

func main() {
	rw, err := openconn("gopher.floodgap.com")
	if err != nil {
		log.Fatalf("failed to dial floodgap\n")
	}

	_, err = rw.WriteString("\n")
	if err != nil {
		log.Fatalf("failed to write newline\n")
	}

	err = rw.Flush()
	if err != nil {
		log.Fatalf("failed to flush tcp readwriter\n")
	}

	page, err := processpage(rw)
	if err != nil {
		log.Fatalf("failed to process page\n")
	}
	
	log.Println(page[0])

	// Could read $PAGER rather than hardcoding the path.
	cmd := exec.Command("/usr/bin/less")

	// Feed it with the string you want to display.
	cmd.Stdin = strings.NewReader("ass")

	// This is crucial - otherwise it will write to a null device.
	cmd.Stdout = os.Stdout

	// Fork off a process and wait for it to terminate.
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("pager closed")
}
