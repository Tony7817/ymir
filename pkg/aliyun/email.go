package aliyun

import (
	"os"

	"ymir.com/pkg/vars"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dm "github.com/alibabacloud-go/dm-20151123/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/pkg/errors"
)

type EmailClient struct {
	Email *dm.Client
}

func NewClientWrapper() (clientWrapper *EmailClient, err error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dm
	config.Endpoint = tea.String("dm.us-east-1.aliyuncs.com")
	email, err := dm.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &EmailClient{
		Email: email,
	}, err
}

func (c *EmailClient) SendNoReplyEmail(destination string, subject string, htmlBody string) (err error) {
	resp, err := c.Email.SingleSendMail(&dm.SingleSendMailRequest{
		AccountName:    tea.String(vars.EmailNoReplySenderName),
		AddressType:    tea.Int32(1),
		ReplyToAddress: tea.Bool(false),
		Subject:        tea.String(subject),
		HtmlBody:       tea.String(htmlBody),
		FromAlias:      tea.String(vars.EmailAlias),
		ToAddress:      tea.String(destination),
	})
	if err != nil {
		return errors.Wrapf(err, "[SendNoReplyEmail] send email failed.")
	}
	if resp != nil && (*resp.StatusCode < 200 || *resp.StatusCode >= 300) {
		return errors.Wrapf(nil, "[SendNoReplyEmail] send email failed, response: %+v", *resp)
	}

	return nil
}
