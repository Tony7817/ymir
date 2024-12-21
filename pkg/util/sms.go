package util

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"ymir.com/pkg/vars"

	"github.com/zeromicro/go-zero/core/logx"
)

func SendCaptchaToPhonenumber(phonenumber string, captcha string) error {
	var body = make(map[string]string)
	var variable = fmt.Sprintf(`{"code":"%s"}`, captcha)

	body["appid"] = os.Getenv("SUBMAIL_SMS_APPID")
	body["signature"] = os.Getenv("SUBMAIL_SMS_APPKEY")
	body["project"] = os.Getenv("SUBMAIL_SMS_PROJECT")
	body["to"] = phonenumber
	body["vars"] = variable

	var bodyBuffer = bytes.NewBuffer([]byte{})
	var writer = multipart.NewWriter(bodyBuffer)
	defer writer.Close()

	for k, v := range body {
		if err := writer.WriteField(k, v); err != nil {
			return err
		}
	}
	var contentType = writer.FormDataContentType()

	resp, err := http.Post(vars.SmsUrl, contentType, bodyBuffer)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	resRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	logx.Debug(resRaw)
	return nil
}
