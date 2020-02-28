package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"tpc/pkg/rpc"
)

func main() {
	file, err := os.OpenFile("serverLog.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("Can't close file; %v", err)
		}
	}()
	log.SetOutput(file)
	const address = "0.0.0.0:9999"
	log.Print("server starting")
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("can't listen on %s: %v", address, err)
	}
	defer func() {
		err := listener.Close()
		log.Fatalf("Can't close %v", err)
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("can't accept conection %v", conn)
			return
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatalf("can't close connection:%v", err)
		}
	}()
	reader := bufio.NewReader(conn)
	line, err := rpc.ReadLine(reader)
	if err != nil {
		log.Printf("error while reading: %v", err)
		return
	}
	index := strings.IndexByte(line, ':')
	writer := bufio.NewWriter(conn)
	if index == -1 {
		log.Printf("invalid line received %s", line)
		err := rpc.WriteLine("error: invalid line", writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return
		}
		return
	}
	cmd, options := line[:index], line[index+1:]
	log.Printf("command received: %s", cmd)
	log.Printf("options received: %s", options)
	switch cmd {
	case rpc.Down:
		options = strings.TrimSuffix(options, "\n")
		file, err := os.Open(rpc.ServerPath + options)
		if err != nil {
			log.Print("file does not exist")
			err = rpc.WriteLine(rpc.Error, writer)
			return
		}
		err = rpc.WriteLine(rpc.Ok, writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return
		}
		_, err = io.Copy(writer, file)
		err = writer.Flush()
		if err != nil {
			log.Printf("Can't flush: %v", err)
			return
		}
	case rpc.Up:
		options = strings.TrimSuffix(options, "\n")
		line, err := rpc.ReadLine(reader)
		if err != nil {
			log.Printf("can't read: %v", err)
			return
		}
		if line == rpc.Error+"\n" {
			log.Printf("file not such: %v", err)
			return
		}
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("can't read data: %v", err)
			}
		}
		err = ioutil.WriteFile(rpc.ServerPath+options, bytes, 0666)
		if err != nil {
			log.Printf("can't write filename: %v", err)
		}
	case rpc.List:
		write := bufio.NewWriter(conn)
		dir, err := ioutil.ReadDir("rpc/server")
		for _, info := range dir {
			if !info.IsDir() && !strings.HasSuffix(info.Name(), ".go") {
				_, err = write.Write([]byte(info.Name() + "\n"))
				if err != nil {
					log.Printf("can't write to client %e", err)
					return
				}
				err = write.Flush()
				if err != nil {
					log.Printf("can't write to client %e", err)
					return
				}
			}
		}
		return
	default:
		err := rpc.WriteLine(rpc.Error, writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return
		}
	}
}
