package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"path"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

const menumarker = "<<<ipfs-gopher-menu>>>"
const maxline = 4096

var host string
var port string
var daemon string
var listen string

func fix(line string, dir string) string {
	fields := strings.Split(line, "\t")
	switch len(fields) {
	case 1:
		if fields[0] == "." {
			break
		}
		if fields[0] != "" {
			fields = append(fields, "fake", host, port)
		}
	case 2:
		fields = append(fields, host, port)
		if strings.HasPrefix(fields[1], "./") {
			fields[1] = dir + fields[1][1:]

		}
	case 3:
		fields = append(fields, host, "70")
	}
	return strings.Join(fields, "\t") + "\r\n"
}

func forward(dst io.Writer, src io.Reader, selector string) error {
	dir := path.Dir(selector)
	// Read the marker if any from the destination
	markerbuf := make([]byte, len(menumarker))
	n, err := src.Read(markerbuf)
	if err != nil {
		return err
	}

	// If we don't have a menu marker, copy what we buffered to the destination
	// and then copy the rest of the src
	if string(markerbuf) != menumarker {
		if _, err := dst.Write(markerbuf[:n]); err != nil {
			return err
		}
		if _, err := io.Copy(dst, src); err != nil {
			return err
		}
		return nil
	}

	// Otherwise we have a menu.
	rd := bufio.NewReaderSize(src, maxline)

	// Discard the remaining line
	_, err = rd.ReadBytes('\n')
	if err != nil {
		return err
	}

	// Fix the menu lines
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		fixed := fix(string(line), dir)
		if _, err := dst.Write([]byte(fixed)); err != nil {
			return err
		}
	}
	return nil
}

func fetch(sh *shell.Shell, selector string, out io.Writer) error {
	rc, err := sh.Cat(selector)
	if err != nil {
		return err
	}
	defer rc.Close()
	err = forward(out, rc, selector)
	if err != nil {
		return err
	}

	return nil
}

func handle(conn net.Conn, sh *shell.Shell) {
	defer conn.Close()
	rd := bufio.NewReader(conn)
	selector, _, err := rd.ReadLine()
	if err != nil {
		log.Printf("failed to read selector: %v", err)
		return
	}
	err = fetch(sh, string(selector), conn)
	if err != nil {
		log.Printf("failed to fetch: %v", err)
		return
	}
}

func main() {
	flag.StringVar(&host, "host", "localhost", "The host to use in IPFS selectors")
	flag.StringVar(&port, "port", "7070", "The port to use in IPFS selectors")
	flag.StringVar(&daemon, "daemon", "localhost:5001", "The address of the IPFS daemon")
	flag.StringVar(&listen, "listen", "localhost:7070", "The address of the proxy")
	flag.Parse()

	sh := shell.NewShell(daemon)

	listener, err := net.Listen("tcp", listen)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
		}
		go handle(conn, sh)
	}
}
