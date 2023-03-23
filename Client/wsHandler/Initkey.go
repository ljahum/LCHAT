package wsHandler

import (
	"client/myaes"
	"client/sign"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/gorilla/websocket"
	"math/big"
	"os"
)

var PubKey *rsa.PublicKey
var PriKey *rsa.PrivateKey

var Sessionkey []byte

const UserName = "user1"

type msg_to_server struct {
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

func initkey(conn *websocket.Conn) bool {

	var server_hello DHExchange
	//recv server hello
	_, plaintext, _ := conn.ReadMessage()
	_ = json.Unmarshal(plaintext, &server_hello)

	//send client hello
	P := server_hello.P
	G := server_hello.G
	bigA := server_hello.A
	b, _ := rand.Prime(rand.Reader, 10)
	bigB := myPow(G, b, P)
	server_hello.B = bigB
	server_hello.Username = UserName
	client_hello := server_hello

	plaintext, _ = json.Marshal(&client_hello)
	conn.WriteMessage(websocket.TextMessage, plaintext)

	//compelet key
	k2 := myPow(bigA, b, P)
	tmphash := md5.New()
	tmphash.Write(k2.Bytes())
	Sessionkey = tmphash.Sum(nil)

	//test key
	_, plaintext, _ = conn.ReadMessage()
	checkEnc := myaes.DecryptecbMode_withUnpadding(plaintext, Sessionkey)
	fmt.Println("decrypto test:", string(checkEnc))
	if string(checkEnc) == "hello" {

		fmt.Println("done")
	} else {
		fmt.Println("秘钥交换出错")
		return false
	}
	//发送公钥
	key, _ := os.ReadFile("./rsa_private.key")
	pkcs8keyStr, _ := pem.Decode(key)
	//解析成pkcs8格式私钥
	//privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)

	rowKey, _ := x509.ParsePKCS8PrivateKey(pkcs8keyStr.Bytes)
	PriKey = rowKey.(*rsa.PrivateKey)
	PubKey = &PriKey.PublicKey
	JsonPubKey, _ := json.Marshal(PubKey)
	conn.WriteMessage(websocket.TextMessage, JsonPubKey)
	//测试签名
	msg := "verifiable message"
	encBytes := myaes.EncryptecbMode_withPadding([]byte(msg), Sessionkey)
	S := sign.RsaSign(PriKey, []byte(msg))

	//反正b64一下
	b64encbytes := base64.StdEncoding.EncodeToString(encBytes)
	b64signature := base64.StdEncoding.EncodeToString(S)
	data := msg_to_server{UserName, b64encbytes, b64signature}

	json_data, _ := json.Marshal(data)
	conn.WriteMessage(websocket.TextMessage, json_data)

	return true
}
