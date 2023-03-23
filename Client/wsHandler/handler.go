package wsHandler

import (
	"client/common"
	"client/myaes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

const (
	Server_domain string = "127.0.0.1"
	Server_port   string = "8080"
)

var Recvdata = make(chan []byte)

var mySocket *websocket.Conn

type Request_data struct {
	Status  int    `json:"Status"`  // 状态码
	Payload string `json:"Payload"` // 数据负载

}

func Send(plain []byte) {
	//只负责流加密并发送
	//log.Println("发送websocket消息:", data)
	var rowFlow common.RowFlowToServer
	enc := myaes.EncryptecbMode_withPadding(plain, Sessionkey)
	dec := myaes.DecryptecbMode_withUnpadding(enc, Sessionkey)
	fmt.Println("statuedPlain", string(plain))
	fmt.Println("dec", string(dec))
	b64Enc := base64.StdEncoding.EncodeToString(enc)

	rowFlow.SessionID = UserName
	rowFlow.Encrypted = b64Enc
	JsonRowFlow, _ := json.Marshal(rowFlow)

	mySocket.WriteMessage(websocket.TextMessage, JsonRowFlow)

}

func recv() {
	for {
		_, AESb64JsonDataBytes, err := mySocket.ReadMessage()
		if err != nil {
			break
		}
		var AESb64JsonData common.Feedback

		_ = json.Unmarshal(AESb64JsonDataBytes, &AESb64JsonData)
		var plaintext []byte

		//fmt.Println("feedback",AESb64JsonData.Flag)
		if AESb64JsonData.Flag == true {
			//var b64RowMsgList []byte
			enc, _ := base64.StdEncoding.DecodeString(AESb64JsonData.MsgList)
			rowMsglist := myaes.DecryptecbMode_withUnpadding(enc, Sessionkey)
			b64RowMsgList := base64.StdEncoding.EncodeToString(rowMsglist)
			AESb64JsonData.MsgList = b64RowMsgList
		} else {
			enc, _ := base64.StdEncoding.DecodeString(AESb64JsonData.MsgList)
			rowMsglist := myaes.DecryptecbMode_withUnpadding(enc, Sessionkey)
			b64RowMsgList := base64.StdEncoding.EncodeToString(rowMsglist)
			AESb64JsonData.MsgList = b64RowMsgList
		}
		plaintext, _ = json.Marshal(AESb64JsonData)
		//fmt.Println("socket接收", string(b64enc))
		//b64dec, _ := base64.StdEncoding.DecodeString(string(b64enc))
		//plaintext, _ := ECBDecrypt(b64dec, Key)
		Recvdata <- plaintext
	}
}

func InitHandler() {

	dl := websocket.Dialer{}
	conn, _, _ := dl.Dial("ws://127.0.0.1:"+Server_port, nil)
	mySocket = conn
	if initkey(conn) == false {
		fmt.Println("秘钥协商失败")
	} else {
		fmt.Println("秘钥协商成功")
	}
	//os.Exit(0)
	//Sessionkey = []byte("QWERTYUIOOIUYTRE")

	go recv()
	//

}
