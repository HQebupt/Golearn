package common

import (
	//"bufio"
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"github.com/surgemq/surgemq/message"
	"linkedin/model/proto/register"
	"linkedin/proto/chat"
	"linkedin/proto/group"
	"linkedin/socket"
	"linkedin/util"
	"net"
	"os"
	"reflect"
	"strconv"
	"time"
)

type ConnImplHandler struct {
	Conn                    net.Conn
	Index                   int32
	Uid                     int64
	UserToken               string
	Msgid                   string
	GroupInvitation         proto.Message
	RemoveNotification      proto.Message
	ApplyNotification       proto.Message
	DropGroupNotification   proto.Message
	RemoveAdminNotification proto.Message
	ExitGroupNotification   proto.Message
	ChatMsg                 *chat.Msg
	EndChan                 chan int
}

func NewConnHander() *ConnImplHandler {
	return &ConnImplHandler{
		EndChan: make(chan int, 1),
	}
}

func (c *ConnImplHandler) InitSocket(server, port string) {
	serverAddress := fmt.Sprintf("%s:%s", server, port)
	tcpAddr, err1 := net.ResolveTCPAddr("tcp", serverAddress)
	if err1 != nil {
		println("Get Tcp address fail!", err1.Error())
		os.Exit(1)
	}

	conn, err2 := net.DialTCP("tcp", nil, tcpAddr)
	if err2 != nil {
		println("Dial failed:", err2.Error())
		os.Exit(1)
	}

	fmt.Println(conn, "Hello Socket Server!")
	c.Conn = conn
}

func (c *ConnImplHandler) CloseSocket() {
	c.Conn.Close()
	<-c.EndChan
}

// Login socket server to regist information in nats
func (c *ConnImplHandler) LoginSocket() {
	pullmsg := message.NewPublishMessage()
	pullmsg.SetTopic([]byte("abc"))
	var tbuf []byte
	//tbuf := make([]byte, 0)
	tt := util.Int16ToBytes(300)
	fmt.Println(tt[0], tt[1])

	//tb := make([]byte, 1024)
	tbuf = append(tbuf, tt[:2]...)

	login := &chat.Login{
		Uid:   proto.Int64(c.Uid),
		Token: proto.String(c.UserToken),
	}

	mha, _ := proto.Marshal(login)
	//th := []byte("hello world!!!")
	na := len(mha)
	tbuf = append(tbuf, mha[:na]...)
	pullmsg.SetPayload(tbuf)
	pullmsg.SetQoS(0)

	fmt.Println(pullmsg.String())
	dstba := make([]byte, 1024)
	if nm, err := pullmsg.Encode(dstba); err == nil {
		c.Conn.Write(dstba[:nm])
	}
}

func (c *ConnImplHandler) SendRequest() {
	fmt.Println("---Send Request------", c.Conn, c.Uid, c.Msgid, c.Uid)
	request := &chat.Request{
		Index: proto.Int32(c.Index),
		Uid:   proto.Int64(c.Uid),
		Msgid: proto.String(c.Msgid),
	}

	bytes, _ := register.StreamEncode(request)
	msg := message.NewPublishMessage()
	msg.SetTopic([]byte("T"))
	msg.SetPayload(bytes)

	dst := make([]byte, len(bytes)+1024)
	if n, err := msg.Encode(dst); err == nil {
		n, err = c.Conn.Write(dst[:n])
		fmt.Printf("send size %d, err :%s\n", n, err)
	} else {
		fmt.Printf("msg to send %d ;err %s\n", n, err)
	}
}

func (c *ConnImplHandler) SendChatMsg(toUId int64, content string, location string) {
	fmt.Println("---Send Msg------", c.Conn, c.Uid, c.Msgid, c.Uid)
	timestamp, _ := strconv.Atoi(time.Now().Local().String())
	request := &chat.Msg{
		MsgId:     []byte{'1'},
		Type:      proto.Int32(0), // 0 chat, 1 pic, 2 audio, 3 for video ,4 location information, 5 notification, 6 name card
		From:      proto.Int64(c.Uid),
		To:        proto.Int64(toUId),
		Content:   proto.String(content),
		Timestamp: proto.Int64(int64(timestamp)),
		Lat:       proto.Float64(0.1),
		Lng:       proto.Float64(0.1),
		Location:  proto.String(location),
	}

	bytes, _ := register.StreamEncode(request)
	msg := message.NewPublishMessage()
	msg.SetTopic([]byte("T"))
	msg.SetPayload(bytes)

	dst := make([]byte, len(bytes)+1024)
	if n, err := msg.Encode(dst); err == nil {
		n, err = c.Conn.Write(dst[:n])
		fmt.Printf("send size %d, err :%s\n", n, err)
	} else {
		fmt.Printf("msg to send %d ;err %s\n", n, err)
	}
}

func (c *ConnImplHandler) HandleConnectMsg(msg *message.ConnectMessage) {
	fmt.Println("Get Connection message", msg)
}

func (c *ConnImplHandler) HandlePingPong(msg *message.PingreqMessage) {
	fmt.Println("Get Pingpone message", msg)
}

func (c *ConnImplHandler) HandlePublishMsg(msg *message.PublishMessage) {
	fmt.Println("----HandlePublishMsg------", c.Conn, c.Index, c.Msgid, c.Uid)
	payload := msg.Payload()
	finalMsg, _ := register.StreamDecode(payload)

	switch finalMsg.(type) {
	case *chat.RadarLocation:
		fmt.Println("Get Radar location msg")
	case *chat.Notifier:
		fmt.Println("Get chat notifier msg")

		c.SendRequest()
	case *chat.Response:
		c.Index = finalMsg.(*chat.Response).GetIndex()
		c.Msgid = finalMsg.(*chat.Response).GetMsgid()
		remain := finalMsg.(*chat.Response).GetRemain()
		for _, data := range finalMsg.(*chat.Response).GetMsglist() {
			innerMsg, _ := register.StreamDecode(data)
			fmt.Println("------innerMsg------", innerMsg)
			fmt.Printf("type %s\n", reflect.TypeOf(innerMsg))
			switch innerMsg.(type) {
			case *chat.Msg:
				c.ChatMsg = innerMsg.(*chat.Msg)
				fmt.Println("----------Get Chat talking Msg----------", c.ChatMsg)
			case *chat.GroupMsg:
				fmt.Println("----------Get chat group----------")
			case *chat.RadarResponse:
				fmt.Println("----------Get rador response----------")
			case *group.InvitationNotification:
				c.GroupInvitation = innerMsg
				fmt.Println("------------Get Group Invitation----------", c.GroupInvitation)
			case *group.RemoveUserNotification:
				c.RemoveNotification = innerMsg
				fmt.Println("------------Remove User Notification Received----------", c.RemoveNotification)
			case *group.ApplicationNotification:
				c.ApplyNotification = innerMsg
				fmt.Println("------------Apply to join group Notification Received----------", c.ApplyNotification)
			case *group.DropGroupNotification:
				c.DropGroupNotification = innerMsg
				fmt.Println("------------Drop group Notification Received----------", c.DropGroupNotification)
			case *group.RemoveAdminNotification:
				c.RemoveAdminNotification = innerMsg
				fmt.Println("------------Remove Admin Notification Received----------", c.RemoveNotification)
			case *group.ExitGroupNotification:
				c.ExitGroupNotification = innerMsg
				fmt.Println("------------Exit Group Notification Received----------", c.ExitGroupNotification)
			default:
				fmt.Println("----------default----------")
			}
		}

		if remain != 0 {
			c.Index = finalMsg.(*chat.Response).GetIndex()
			c.Msgid = finalMsg.(*chat.Response).GetMsgid()
			c.SendRequest()

		}

	default:
		fmt.Println("No type detect")
	}
}

func (c *ConnImplHandler) GetSocketMessage() {
	defer func() {
		fmt.Println("All sockets information received!")
		if r := recover(); r != nil {
			c.EndChan <- 1
			return
		}
	}()

	// init object
	fmt.Println("----C Object Infor-----", c.Conn, c.Index, c.Msgid, c.Uid)
	for {
		if nil == c.Conn {
			break
		}
		mqttMsg := socket.ReadMqttMessage(c.Conn)
		socket.HandleRecvMsg(mqttMsg, c)

	}
	//	reader := readMqtMessage(conn)
	//	fmt.Println(n)
	//	register.StreamDecode(reader)
	// receiveData := message.NewPublishMessage()
	// receiveData.Decode(socketData)
	// fmt.Println("the data after decode is:", receiveData)
	// receiveBytes := receiveData.Payload()
	// //the data header has 2 bytes for package header
	// payload := receiveBytes[2:]
	// notifier := chat.Notifier{}
	// err := proto.Unmarshal(payload, &notifier)
	// fmt.Println(err)
	// fmt.Println(notifier.GetRemain())
	// remain := notifier.GetRemain()

	//	if remain != 0 {
	//		n, _ := conn.Read(socketData)
	//		fmt.Println(n)
	//		receiveData := message.NewPublishMessage()
	//		receiveData.S(socketData)
	//		fmt.Println("the data after decode is:", receiveData)
	//		receiveBytes := receiveData.Payload()
	//		//the data header has 2 bytes for package header
	//		payload := receiveBytes[2:]
	//		fmt.Println("+++++++++++", payload[0], payload[1])
	//		fmt.Println(payload)
	//	}
	// resq := chat.Request{}
	// resp := chat.Request{
	// 	Index :proto.Int32(notifier.)
	// }

}

/*
func GetUnmarshalData(receiveData []byte, dataType proto.Message) {
	receiveData := message.NewPublishMessage()
	receiveData.Decode(receiveData)
	receiveBytes := receiveData.Payload()

	//the data header has 2 bytes for package header
	payload := receiveBytes[2:]
	fmt.Println(receiveBytes[1])
	notifier := chat.Notifier{}
	err := proto.Unmarshal(payload, &notifier)

	fmt.Println(notifier.GetRemain())
}
*/
// func GetSocketResponse(conn net.Conn, resp proto.Message) {

// 	data, err3 := bufio.NewReader(conn).ReadString('\n')
// 	if err3 != nil {
// 		panic(err3)
// 	}

// 	receiveData := message.NewPublishMessage()
// 	receiveData.Decode([]byte(data))
// 	fmt.Println("the data after decode is:", receiveData)
// 	receiveBytes := receiveData.Payload()

// 	//the data header has 2 bytes for package header
// 	payload := receiveBytes[2:]
// 	err4 := proto.Unmarshal(payload, resp)
// 	fmt.Println(err4)
// }
