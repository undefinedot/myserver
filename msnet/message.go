package msnet

type Message struct {
	ID uint32  // 消息ID
	DataLen uint32
	Data []byte
}

// NewMsgPackage 初始化一个Message
func NewMsgPackage(id, len uint32, data []byte) *Message  {
	return &Message{
		ID:      id,
		DataLen: len,
		Data:    data,
	}
}

func (m *Message) GetMsgID() uint32 {
	return m.ID
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgID(msgID uint32) {
	m.ID = msgID
}

func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

