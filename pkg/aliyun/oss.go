package aliyun

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts "github.com/alibabacloud-go/sts-20150401/v2/client"
	"github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/zeromicro/go-zero/core/logx"
)

type OssClient struct {
	Oss *sts.Client
}

func NewOssClient() (*OssClient, error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Sts
	config.Endpoint = tea.String("sts.us-east-1.aliyuncs.com")
	client, err := sts.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &OssClient{
		Oss: client,
	}, nil
}

func (a *OssClient) GetSTSToken(userName string) (*sts.AssumeRoleResponseBodyCredentials, error) {
	var assumeRoleRequest = &sts.AssumeRoleRequest{
		DurationSeconds: tea.Int64(60 * 60),
		RoleArn:         tea.String(os.Getenv("ALIBABA_CLOUD_ROLE_ARN")),
		RoleSessionName: tea.String(fmt.Sprintf("userId-%s", userName)),
	}
	res, err := a.Oss.AssumeRoleWithOptions(assumeRoleRequest, &service.RuntimeOptions{})
	if err != nil {
		var SDKError = &tea.SDKError{}
		if e, ok := err.(*tea.SDKError); ok {
			SDKError = e
		} else {
			SDKError.Message = tea.String(err.Error())
		}
		logx.Errorf("[GetSTSToken] assume role fail, userName: %s, err: %+v", userName, SDKError)
		var data any
		var d = json.NewDecoder(strings.NewReader(tea.StringValue(SDKError.Data)))
		if err := d.Decode(&data); err != nil {
			logx.Errorf("[GetSTSToken] decode data fail, userName: %s, err: %s", userName, err)
		}
		if m, ok := data.(map[string]any); ok {
			recommend := m["Recommend"]
			logx.Errorf("Recommend: %+v", recommend)
		}
		return nil, err
	}

	return res.Body.Credentials, nil
}
