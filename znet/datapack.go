package znet

import (
	"bytes"
	"errors"
	"encoding/binary"
	"workspace/src/zinx/ziface"
	"workspace/src/zinx/utils"
)

type DataPack struct {}

func NewDataPack() ziface.IDataPack {
	return &DataPack{}
}

// 获取包头部的长度
func (dp *DataPack) GetHaedLen() uint32 {
	// 头部包含信息：Datalen + id，这两个都是uint32(4字节)，那么需要8个字节
	return 8
}

// 将输入的消息封包
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes数据的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	/*
		写入的顺序要固定，不能交换，因为读取是按照顺序读取的
		使用二进制写入时需要指定大端或小端，读取的时候保持一致
	*/

	// 第一步：写入DataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 第二步：写入消息的Id
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 第三步：写入消息的内容
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包，先读取头部，根据头部中记录的消息长度来读取消息
func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的读取ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 创建一个初始消息类型，用于存储读取的消息
	msg := &Message{}

	// 第一步：从数据中读取DataLen，由于msg.DataLen是uint32，可以直接使用该属性读取，保证读取4个字节
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 第二步：从数据中再读取4个字节，即Id
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 第三步：判断读取出来的DataLen是否合法
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("Too large msg data received!")
	}

	// // 这里也以设置为不读取消息，在外面去写第二次读取
	// msg.Data = make([]byte, msg.DataLen)
	// if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Data); err != nil {
	// 	return nil, err
	// }

	return msg, nil
}	