package handler

import (
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func (h *Handler) GetWxEcho(c echo.Context) error {
	verifyMsgSign := c.QueryParam("msg_signature")
	verifyTimestamp := c.QueryParam("timestamp")
	verifyNonce := c.QueryParam("nonce")
	verifyEchoStr := c.QueryParam("echoStr")
	str := h.sv.WxVerify(c.Request().Context(), verifyMsgSign, verifyTimestamp, verifyNonce, verifyEchoStr)
	return c.String(http.StatusOK, str)
}

func (h *Handler) PostWxEcho(c echo.Context) error {
	reqMsgSign := c.QueryParam("msg_signature")
	reqTimestamp := c.QueryParam("timestamp")
	reqNonce := c.QueryParam("nonce")
	bodyReader := c.Request().Body
	reqData, err := io.ReadAll(bodyReader)
	if err != nil {
		return err
	}

	// 解密
	body, err := h.sv.WxDecryptMsg(c.Request().Context(), reqMsgSign, reqTimestamp, reqNonce, reqData)
	if err != nil {
		return err
	}

	// 处理
	resp, err := h.sv.Report(c.Request().Context(), body.Content)
	if err != nil {
		return err
	}

	// 加密响应
	respData := h.sv.WxEncryptMsg(c.Request().Context(), resp, reqTimestamp, reqNonce)
	return c.String(http.StatusOK, respData)
}
