package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"server/URWWPacketProtocol"
	"server/logutil"
	"strconv"

	"github.com/golang/protobuf/proto"
)

// protocol header
const (
	ConstHeader       = "com.ur.URPackageHeader"
	ConstHeaderLength = 22
	ConstdataLength   = 4
	constTagLength    = 1
)

func main() {

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":7777")
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	logutil.Writelog("socket listen")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		// daytime := time.Now().Format("2006-01-02 15:04:05")
		// conn.Write([]byte(daytime))
		// logutil.Writelog("finish send file")
		// conn.Close()

		logutil.Writelog("handle request")
		go handleClient(conn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleClient(conn net.Conn) {
	request := make([]byte, 1024)

	defer conn.Close()

	for {
		readlen, err := conn.Read(request)

		if err != nil {
			fmt.Println(err)
			logutil.Writelog("send content break:")
			break
		}

		if readlen == 0 {
			logutil.Writelog("send content break:")
			break
		} else {
			// content := string(request)

			Unpack(request, readlen)

			// messageData := BytesToInt(request[0:4])

			// logutil.Writelog("send content to int : " + strconv.Itoa(messageData))

			// if 4+messageData > read_len {
			// messageData = read_len - 4
			// }

			// realContent := string(request[4 : 4+messageData])

			// logutil.Writelog("origin content:" + realContent)

			// logutil.Writelog("send content:" + realContent + ":" + strconv.Itoa(messageData))

			// for i := 0; i < 5; i++ {
			// 	words := strconv.Itoa(i) + "This is a test for long conn"
			// 	conn.Write([]byte(words))
			// 	time.Sleep(2 * time.Second)
			// }

			// conn.Write([]byte(realContent))

		}

		logutil.Writelog("clear send data")
		request = make([]byte, 1024)
	}

	logutil.Writelog("handle request over")
}

// Unpack the package
func Unpack(buffer []byte, length int) {

	logutil.Writelog("recv package length to int : " + strconv.Itoa(length))

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+ConstHeaderLength {
			break
		}
		// headerTag := string(buffer[i : i+ConstHeaderLength])
		// logutil.Writelog("send rever header string : " + headerTag)
		// if headerTag == ConstHeader {

		tagLength := TagBytesToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+1])

		logutil.WriteLogInt(tagLength)

		// logutil.Writelog("send tagString : " + tagString)
		// logutil.Writelog("send tagLength to int : " + strconv.Itoa(tagLength))
		// log.Println("%d", tagLength)

		if length < i+ConstHeaderLength+constTagLength+ConstdataLength {
			break
		}

		dataLength := BytesToInt(buffer[i+ConstHeaderLength+constTagLength : i+ConstHeaderLength+constTagLength+ConstdataLength])
		logutil.Writelog("send dataLength to int : " + strconv.Itoa(dataLength))

		if length < i+ConstHeaderLength+constTagLength+ConstdataLength+dataLength {
			break
		}

		messageData := buffer[i+ConstHeaderLength+constTagLength+ConstdataLength : i+ConstHeaderLength+constTagLength+ConstdataLength+dataLength]

		message := &URWWPacketProtocol.URProtocol{}
		err := proto.Unmarshal(messageData, message)
		if err != nil {
			logutil.Writelog("Unmarshal error :" + err.Error())
		}

		uri := message.GetUri()

		switch uri {
		case URWWPacketProtocol.URPacketType_kUriPLoginReq:
			logutil.Writelog("recv  kUriPLoginReq")
			logutil.Writelog(message.GetLoginReq().GetPassport())
			logutil.Writelog(message.GetLoginReq().GetPassword())

			// result := pack()
			// conn.Write(result)

		case URWWPacketProtocol.URPacketType_kUriPLogoutReq:
			logutil.Writelog("recv kUriPLogoutReq")
		}
	}
}

// IntToBytes int 转化
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// BytesToInt 转化 int
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	// log.Println(bytesBuffer)
	// log.Println(x)

	return int(x)
}

// TagBytesToInt 将byte变为int
func TagBytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	// log.Println(b)
	//log.Println(bytesBuffer)
	var x int8
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	// log.Println(bytesBuffer)
	// log.Println(x)

	return int(x)
}

func pack() []byte {
	message := &URWWPacketProtocol.URProtocol{
		Uri: URWWPacketProtocol.URPacketType_kUriPLoginRes,
		Header: &URWWPacketProtocol.PHeader{
			Result: &URWWPacketProtocol.Result{
				Code:   URWWPacketProtocol.ResultType_ResultTypeOK,
				ResMsg: "success",
			},
		},
	}
	data, err := proto.Marshal(message)

	checkError(err)

	return data
}
