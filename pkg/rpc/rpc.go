package rpc

import (
	"bufio"
)

const (
	Tcp        = "tcp"
	Down       = "download"
	Up         = "upload"
	List       = "list"
	ServerPath = "rpc/server/"
	ClientPath = "rpc/client/"
	Ok         = "result ok"
	Error      = "result error"
)

func ReadLine(reader *bufio.Reader) (line string, err error) {
	return reader.ReadString('\n')
}

func WriteLine(line string, writer *bufio.Writer) (err error) {
	_, err = writer.WriteString(line + "\n")
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
	return
}
