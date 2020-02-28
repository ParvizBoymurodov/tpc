package rpc

import (
"bufio"
)

const Down  = "download"
const Up  = "upload"
const List  = "list"
const ServerPath  = "rpc/server/"
const ClientPath  = "rpc/client/"
const Ok  = "result ok"
const Error  = "result error"
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
