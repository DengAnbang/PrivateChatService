package socket

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gitee.com/DengAnbang/PrivateChatService/src/bean"
	"gitee.com/DengAnbang/PrivateChatService/src/socket/push"
	"gitee.com/DengAnbang/PrivateChatService/src/socket/socketConst"
	"net"

	"gitee.com/DengAnbang/goutils/loge"
)

const (
	HeaderLength = 4 //头的长度
	bytesSize    = 1024 * 4
)

type TcpConn struct {
	Id           string
	headerLength int //头的长度
	Conn         net.Conn
}

func (conn *TcpConn) SetId(id string) {
	conn.Id = id
}
func (conn *TcpConn) GetId() string {
	return conn.Id
}
func (conn *TcpConn) SendMessageToConn(msg interface{}) (err error) {
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	loge.W("tcp:发送的内容:", string(jsonBytes))
	_, err = conn.WriteMessage(jsonBytes)
	return err
}
func (conn *TcpConn) Response(err error, messageType string) {
	if resultData, ok := err.(*bean.ResultData); ok {
		resultData.Type = messageType
		sendErr := conn.SendMessageToConn(resultData)
		if sendErr != nil {
			loge.W(sendErr)
		}
	} else if err, ok := err.(error); ok {
		data := bean.NewErrorMessage("服务器内部错误")
		data.DebugMessage = fmt.Sprintf("%v", err)
		sendErr := conn.SendMessageToConn(data)
		if sendErr != nil {
			loge.W(sendErr)
		}
	}
}
func (conn *TcpConn) WriteMessage(message []byte) (int, error) {
	//messageLen := len(message)
	return conn.Conn.Write(conn.Packet(message))
}

//解包
func (conn *TcpConn) Unpack(data []byte) []byte {
	return data[HeaderLength:]
}

//打包
func (conn *TcpConn) Packet(data []byte) []byte {
	buffer := bytes.NewBuffer([]byte{})
	err := binary.Write(buffer, binary.BigEndian, IntToBytes(len(data)))
	fmt.Println(err)
	buffer.Write(data)
	return buffer.Bytes()
}

//读取包
func (conn *TcpConn) ReadConn(readerChannel chan<- []byte) error {
	readBytes := make([]byte, bytesSize)
	defer conn.Conn.Close()
	defer close(readerChannel)
	buffer := bytes.NewBuffer(make([]byte, 0, bytesSize))
	for {
		n, err := conn.Conn.Read(readBytes)
		if err != nil {
			return err
		}
		buffer.Write(readBytes[:n])
		//是否包含了头
		if buffer.Len() >= conn.headerLength {
			dataLen := BytesToInt(buffer.Bytes()[:conn.headerLength])
			//包含应该完整的包了
			if buffer.Len() >= conn.headerLength+dataLen {
				//conn.data = buffer.Bytes()[:conn.headerLength+dataLen]
				unpack := conn.Unpack(buffer.Bytes()[:conn.headerLength+dataLen])
				readerChannel <- unpack[:dataLen]
				last := buffer.Bytes()[conn.headerLength+dataLen:]
				buffer = bytes.NewBuffer(make([]byte, 0, bytesSize))
				buffer.Write(last)
			}
		}
	}

}
func TcpRun(port string) {
	netListen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		loge.W(err)
		return
	}
	defer netListen.Close()
	for {
		conn, err := netListen.Accept()
		if err != nil {
			loge.W(err)
			continue
		}
		loge.WD("连接成功请求地址:" + conn.RemoteAddr().String())
		go tcpSocketHandler(conn)
	}
}
func tcpSocketHandler(conn net.Conn) {
	tcpConn := TcpConn{headerLength: HeaderLength, Conn: conn}
	message := bean.NewSucceedMessage("连接成功!")
	message.Type = socketConst.TypeConnect
	err := tcpConn.SendMessageToConn(message)
	if err != nil {
		loge.W(err)
	}
	data := make(chan []byte)
	go tcpConn.ReadConn(data)
	for {
		select {
		case v, ok := <-data:
			if ok {
				var sm bean.SocketData
				err := json.Unmarshal(v, &sm)
				loge.WD("读取到的消息:" + string(v))
				if err == nil {
					go Dispense(&sm, &tcpConn)
				} else {
					tcpConn.Response(bean.NewErrorMessage(fmt.Sprintf("编码错误,%s", err)), "0")
				}
			} else {
				push.UnRegister(&tcpConn)
				loge.WD("关闭:" + tcpConn.Id)
				return
			}
		}
	}

}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}
