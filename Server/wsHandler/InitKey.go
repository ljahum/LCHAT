package wsHandler

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"math/big"
	"server/myaes"
	"server/sign"
)

type msg_to_server struct {
	//payload和sign默认b64传输

	ID      string `json:"ID"`
	Payload string `json:"Payload"`
	Sign    string `json:"Sign"`
}

type DHExchange struct {
	P        *big.Int
	G        *big.Int
	A        *big.Int
	B        *big.Int
	Username string
}

func myPow(M *big.Int, E *big.Int, N *big.Int) *big.Int {
	var c big.Int
	c.Exp(M, E, N)

	return &c
}

func Initkey(conn *websocket.Conn) (string, []byte, *rsa.PublicKey) {

	G, _ := rand.Prime(rand.Reader, 100)
	P, _ := rand.Prime(rand.Reader, 128)
	a, _ := rand.Prime(rand.Reader, 10)
	bigA := myPow(G, a, P)
	server_hello := &DHExchange{
		P,
		G,
		bigA,
		nil,
		"",
	}

	//send server hello
	plaintext, _ := json.Marshal(&server_hello)
	conn.WriteMessage(websocket.TextMessage, plaintext)

	// recv client hello
	var client_hello DHExchange
	_, plaintext, _ = conn.ReadMessage()
	_ = json.Unmarshal(plaintext, &client_hello)

	//compelet key
	bigB := client_hello.B
	k2 := myPow(bigB, a, P)
	tmphash := md5.New()
	tmphash.Write(k2.Bytes())

	//存储秘钥
	var shearedkey []byte
	shearedkey = tmphash.Sum(nil)

	//SessionKeys[client_hello.Username] = shearedkey
	//test aes
	checkEnc := myaes.EncryptecbMode_withPadding([]byte("hello"), shearedkey)
	conn.WriteMessage(websocket.TextMessage, checkEnc)
	fmt.Println("session exchange done")
	//接受公钥
	_, plaintext, _ = conn.ReadMessage()
	var pubKey *rsa.PublicKey
	_ = json.Unmarshal(plaintext, &pubKey)
	//验证公钥
	var recv_json msg_to_server
	_, data, _ := conn.ReadMessage()
	_ = json.Unmarshal(data, &recv_json)
	encMsg, _ := base64.StdEncoding.DecodeString(recv_json.Payload)
	signature, _ := base64.StdEncoding.DecodeString(recv_json.Sign)
	plaintext = myaes.DecryptecbMode_withUnpadding(encMsg, shearedkey)
	if sign.RsaVerify(pubKey, signature, plaintext) == true {
		fmt.Println("身份认证成功")
	} else {
		fmt.Println("身份认证失败")
	}

	return client_hello.Username, shearedkey, pubKey
}
