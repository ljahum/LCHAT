package common

// 第一层
type RowFlowToServer struct {
	//payload和sign默认b64传输
	SessionID string `json:"SessionID"`
	Encrypted string `json:"Encrypted"` //b64
}

// 第二层
type StatusFlow struct {
	//payload和sign默认b64传输
	Status  int    `json:"Status"`
	Payload string `json:"Payload"` // b64
}

// 第三层 登录注册
type UserForm struct {
	UserID   string `json:"userID"`   // Id
	Password string `json:"password"` // 密码
}

// 第三层 消息魔板
type ChatMsg struct {
	Msg       string `json:"Msg"`       // Id
	Signature string `json:"Signature"` // 密码
}

// 返回的feedback
type Feedback struct {
	Flag    bool   `json:"flag"`
	MsgList string `json:"msgList"` //b64
}

// 留言的结构体
type Comment struct {
	Id      int64
	Name    string
	Content string
	Mail    string
	Time    string
}
