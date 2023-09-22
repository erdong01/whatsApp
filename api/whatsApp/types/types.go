package types

type SendReq struct {
	// 发送者
	UserId string `json:"user_id"`

	// 接收者
	ReceiverPhone string `json:"receiver_phone"`

	// 消息内容
	Content string `json:"content"`

	// 消息类型
	Type string `json:"type"`
}
