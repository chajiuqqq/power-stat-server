package service

import (
	"context"
	"encoding/xml"
	"mqtt-wx-forward/types"
)

// WxVerify 验证url
func (s *Service) WxVerify(ctx context.Context, verifyMsgSign, verifyTimestamp, verifyNonce, verifyEchoStr string) string {
	echoStr, cryptErr := s.wxcpt.VerifyURL(verifyMsgSign, verifyTimestamp, verifyNonce, verifyEchoStr)
	if nil != cryptErr {
		s.Logger.Println("verifyUrl fail", cryptErr)
		return ""
	}
	s.Logger.Println("verifyUrl success echoStr", string(echoStr))
	return string(echoStr)
}

// 解密用户的消息
func (s *Service) WxDecryptMsg(ctx context.Context, reqMsgSign, reqTimestamp, reqNonce string, reqData []byte) (*types.MsgContent, error) {
	msg, cryptErr := s.wxcpt.DecryptMsg(reqMsgSign, reqTimestamp, reqNonce, reqData)
	if nil != cryptErr {
		s.Logger.Println("DecryptMsg fail", cryptErr)
	}
	s.Logger.Println("after decrypt msg: ", string(msg))

	var msgContent types.MsgContent
	err := xml.Unmarshal(msg, &msgContent)
	if nil != err {
		s.Logger.Println("Unmarshal fail")
		return nil, err
	} else {
		s.Logger.Println("struct", msgContent)
	}
	return &msgContent, nil
}

// WxEncryptMsg 加密返回的消息
func (s *Service) WxEncryptMsg(ctx context.Context, respData, reqTimestamp, reqNonce string) string {
	encryptMsg, cryptErr := s.wxcpt.EncryptMsg(respData, reqTimestamp, reqNonce)
	if nil != cryptErr {
		s.Logger.Println("DecryptMsg fail", cryptErr)
	}

	sEncryptMsg := string(encryptMsg)
	s.Logger.Println("after encrypt sEncryptMsg: ", sEncryptMsg)
	return sEncryptMsg
}
