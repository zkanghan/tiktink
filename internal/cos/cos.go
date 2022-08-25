package cos

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/spf13/viper"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// PublishFileToServer  url := https://tiktink-1311932295.cos.ap-nanjing.myqcloud.com
//  上传视频，调用腾讯云的第三方库
func PublishFileToServer(r io.Reader, key string) (string, error) {
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。https://console.cloud.tencent.com/cos5/bucket
	// COS_REGION 可以在控制台查看，https://console.cloud.tencent.com/cos5/bucket, 关于地域的详情见 https://cloud.tencent.com/document/product/436/6224
	u, _ := url.Parse("https://tiktink-1311932295.cos.ap-nanjing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  viper.GetString("cos.SECRETID"),  // 替换为用户的 SecretId，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
			SecretKey: viper.GetString("cos.SECRETKEY"), // 替换为用户的 SecretKey，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
		},
	})
	// 对象键（Key）是对象在存储桶中的唯一标识，例如，在对象的访问域名 `examplebucket-1250000000.cos.COS_REGION.myqcloud.com/test/objectPut.go` 中，对象键为 test/objectPut.go

	// 3.通过文件流上传对象
	_, err := c.Object.Put(context.Background(), key, r, nil)
	if err != nil {
		return "", err
	}
	//   返回可供访问的url
	return c.Object.GetObjectURL(key).String(), nil

}
