package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"testing"
	"time"
	"tpc/pkg/rpc"
)

func Test_uploadFromServer(t *testing.T) {
	const address = "0.0.0.0:9998"
	go func() {
		listener, err := net.Listen(rpc.Tcp, address)
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
		}()
	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("can't connect to %s: %v", address, err)
	}
	writer := bufio.NewWriter(conn)
	fileName := "par"
	cmd := rpc.Up +":"+ fileName
	err = rpc.WriteLine(cmd, writer)
	if err != nil {
		t.Fatal("can't write", err)
	}
	file, err := ioutil.ReadFile("testdata/" +fileName)
	if err != nil {
		t.Fatalf("can't read file: %v\n", err)
	}
	_, err = writer.Write(file)
	if err != nil {
		t.Fatalf("can't copy the file %s\n", file)
	}
	err = writer.Flush()
	if err != nil {
		t.Fatalf("can't flush %v\n", err)
	}
	err = conn.Close()
	if err != nil {
		t.Fatalf("can't conn close %v\n", err)
	}
	upload, err := ioutil.ReadFile("testdata/" +fileName)
	if err != nil {
		t.Fatalf("can't read file upload error:  %v\n",err)
	}
	if  !bytes.Equal(file, upload){
		t.Fatalf("files are not equal: %v", err)
	}
}



func Test_downloadFromServer(t *testing.T){
	const address = "0.0.0.0:9999"
	go func() {
		listener, err := net.Listen(rpc.Tcp, address)
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
	}()
	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("can't connect to %s: %v", address, err)
	}
	writer := bufio.NewWriter(conn)
	fileName := "par"
	cmd := rpc.Down +":"+ fileName
	err = rpc.WriteLine(cmd, writer)
	if err != nil {
		t.Fatalf("can't write command %v\n", err)
	}
	reader := bufio.NewReader(conn)
	download, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("can't reader file error: %v\n", err)
	}
	err = conn.Close()
	if err != nil {
		t.Fatalf("can't conn close %v\n", err)
	}
	err = ioutil.WriteFile("testdata/"+fileName, download, 0666)
	if err != nil {
		t.Fatalf("can't write file: %v\n", err)
	}
	downloadFile, err := ioutil.ReadFile("testdata/"+ fileName)
	if err != nil {
		t.Fatalf("can't read file upload error:  %v\n",err)
	}
	if !bytes.Equal(download,downloadFile) {
		t.Fatalf("files are not equal: %v", err)
	}
}