package w2net

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/lnnlmario/w2-inx/w2iface"
	"github.com/lnnlmario/w2-inx/w2utils"
)

// 封包/拆包实例
type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d DataPack) GetHeadLen() uint32 {
	// id(unit32) + DataLen(uint32)
	return 8
}

// 封包方法 压缩数据
func (d DataPack) Pack(msg w2iface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 写msgId
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法 解压数据
func (d DataPack) Unpack(binaryData []byte) (w2iface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	msg := &Message{}

	// 读取dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读取msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen长度是否唱过最大包长度
	if w2utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > w2utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg data received")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}
