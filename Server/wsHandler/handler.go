package wsHandler

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"server/common"
	"server/myaes"
	"server/userapi"
)

var UP = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// var userKey map[string]string = map[string]string{"admin": "202cb962ac59075b964b07152d234b70"}
var connectionsList map[string]*websocket.Conn
var SessionKeyList map[string][]byte
var PubKeyList map[string]*rsa.PublicKey

func handler(w http.ResponseWriter, r *http.Request) {

	conn, err := UP.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	SessionID, shearedkey, Pubkey := Initkey(conn)
	//inittest
	//os.Exit(0)
	connectionsList[SessionID] = conn
	//SessionKeyList["user1"] = []byte("QWERTYUIOOIUYTRE")
	//SessionKeyList["user2"] = []byte("QWERTYUIOOIUYTRE")
	SessionKeyList[SessionID] = shearedkey
	PubKeyList[SessionID] = Pubkey
	for {
		var recv_json common.RowFlowToServer

		_, data, _ := conn.ReadMessage()
		_ = json.Unmarshal(data, &recv_json)
		fmt.Println(recv_json)
		SessionKey := SessionKeyList[recv_json.SessionID]
		//b64
		encBytes, _ := base64.StdEncoding.DecodeString(recv_json.Encrypted)
		fmt.Println("发送者", recv_json.SessionID)
		fmt.Println("对应密钥", SessionKey)
		//aes
		rowStatuedPlain := myaes.DecryptecbMode_withUnpadding(encBytes, SessionKey)
		//fmt.Println("rowStatuedPlain", rowStatuedPlain)
		var StatuedPlain common.StatusFlow

		var feedback common.Feedback
		_ = json.Unmarshal(rowStatuedPlain, &StatuedPlain)
		if StatuedPlain.Status == common.StatusLogin { //login
			fmt.Println("申请登录")
			var LoginForm common.UserForm
			jsonLoginForm, _ := base64.StdEncoding.DecodeString(StatuedPlain.Payload)
			_ = json.Unmarshal(jsonLoginForm, &LoginForm)
			feedback = userapi.CheckLogin(LoginForm)

		}
		if StatuedPlain.Status == common.StatusRegedit {
			fmt.Println("申请注册")
			var LoginForm common.UserForm
			jsonLoginForm, _ := base64.StdEncoding.DecodeString(StatuedPlain.Payload)
			_ = json.Unmarshal(jsonLoginForm, &LoginForm)
			feedback = userapi.CheckRegedit(LoginForm)
		}

		if StatuedPlain.Status == common.StatusBroadcast {
			//todo
			fmt.Println("申请消息入库")
			var chatForm common.ChatMsg
			jsonLoginForm, _ := base64.StdEncoding.DecodeString(StatuedPlain.Payload)
			_ = json.Unmarshal(jsonLoginForm, &chatForm)
			pubkey := PubKeyList[recv_json.SessionID]
			feedback = userapi.CheckInsertMsg(pubkey, chatForm)

		}
		if StatuedPlain.Status == common.StatusChat {
			//不想做捏
		}
		var TestData []byte
		fmt.Println("feedback.Flag", feedback.Flag)
		if feedback.Flag == true {
			//TestData = []byte("testString")
			TestData = common.GetComment()
		} else {
			//bytes
			TestData = []byte("nil")

		}
		//aes
		//getMsgFromDB

		EncryptedTestData := myaes.EncryptecbMode_withPadding(TestData, SessionKey)
		//base64
		b64enc := base64.StdEncoding.EncodeToString(EncryptedTestData)
		//load to json
		feedback.MsgList = b64enc
		jsonFeedback, _ := json.Marshal(feedback)
		conn.WriteMessage(websocket.TextMessage, jsonFeedback)

	}
	log.Println("未收到消息 ,与该client链接关闭")
}

func InitHandler() {
	connectionsList = make(map[string]*websocket.Conn)
	SessionKeyList = make(map[string][]byte)
	PubKeyList = make(map[string]*rsa.PublicKey)
	http.HandleFunc("/", handler)
	fmt.Println("startListen")
	http.ListenAndServe(":8080", nil)
}
