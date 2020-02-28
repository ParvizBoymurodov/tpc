package main

import (
	"bufio"
	"flag"

	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"tpc/pkg/rpc"
)

var down = flag.String("down", "ddd", "Down")
var up = flag.String("up", "uuu", "Up")
var lists = flag.Bool("list", false, "List")

func main() {
	//const authorizedOperations = `Список доступных операций:
	//download
	//upload
	//list`
	file, err := os.OpenFile("clientLog.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
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
flag.Parse()
	var cmd, filename string
	if *down != "ddd" {
		filename = *down
		cmd = rpc.Down
	} else if *up != "uuu" {
		cmd = rpc.Up
		filename = *up
	} else if *lists != false {
		cmd = rpc.List
		filename = ""
	} else{
		return}
	//fmt.Println(authorizedOperations)
	address := "0.0.0.0:9999"
	log.Print("client connecting")
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("can't connect to %s: %v", address, err)
	}
	defer func() {
		err:= conn.Close()
		if err != nil {
			log.Printf("Can't close file; %v", err)
		}
	}()
	log.Print("client connected")
	writer := bufio.NewWriter(conn)
	//var fileName, cmd string
	//fmt.Scan(&cmd)
	log.Print("command sent")
	switch cmd {
	case rpc.Down:
		//fmt.Scan(&filename)
		line := cmd + ":" + filename
		err = rpc.WriteLine(line, writer)
		if err != nil {
			log.Fatalf("can't send command %s to server: %v", line, err)
		}
		log.Print("command sending")
		download(conn, filename)
	case rpc.Up:
		//fmt.Scan(&filename)
		line := cmd + ":" + filename
		err = rpc.WriteLine(line, writer)
		if err != nil {
			log.Fatalf("can't send command %s to server: %v", line, err)
		}
		log.Print("command sending")
		upload(conn, filename)
	case rpc.List:
		err = rpc.WriteLine(cmd+":", writer)
		if err != nil {
			log.Fatalf("can't send command %s to server: %v", cmd, err)
		}
		list(conn)
	default:
		fmt.Println("Не правильная операция!")
	}
}

func download(conn net.Conn, filename string) {
	reader := bufio.NewReader(conn)
	line, err := rpc.ReadLine(reader)
	if err != nil {
		log.Printf("can't read: %v", err)
		return
	}
	if line == rpc.Error+"\n" {
		log.Printf("file not such: %v", err)
		fmt.Printf("Файл с название %s не существует", filename)
		return
	}
	log.Print(line)
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		if err != io.EOF {
			log.Printf("can't read data: %v", err)
		}
	}
	log.Print(len(bytes))
	err = ioutil.WriteFile(rpc.ClientPath+filename, bytes, 0666)
	if err != nil {
		log.Printf("can't write filename: %v", err)
	}
	fmt.Printf("Файл %s успешно скаченно", filename)
}

func upload(conn net.Conn, filename string) {
	options := strings.TrimSuffix(filename, "\n")
	file, err := os.Open(rpc.ClientPath + options)
	writer := bufio.NewWriter(conn)
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
	log.Print(filename)
	written, err := io.Copy(writer, file)
	log.Print(written)
	err = writer.Flush()
	if err != nil {
		log.Printf("Can't flush: %v", err)
		return
	}
	fmt.Printf("Файл %s загруженно на сервер", filename)
}

func list(conn net.Conn) {
	reader := bufio.NewReader(conn)
	counter := 0
	fmt.Println("files:")
	for {
		readString, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Printf("Total: %d files", counter)
			return
		}
		if err != nil {
			log.Printf(
				"can't read frocm loalhost %v", err)
			fmt.Println("Не удалось получить ответ")
			return
		}
		counter++
		fmt.Println(readString)

	}
}
