package obs

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"time"
)

// Borrowed from aliyun doc
//   https://help.aliyun.com/document_detail/31926.html

// var accessKeyId string = "6MKOqxGiGU4AUk44"
// var accessKeySecret string = "ufu7nS8kS59awNihtjSonMETLI0KLy"
// var host string = "http://post-test.oss-cn-hangzhou.aliyuncs.com"
// var expire_time int64 = 60
// var upload_dir string = "user-dir/"

func get_gmt_iso8601(expire_end int64) string {
	var tokenExpire = time.Unix(expire_end, 0).Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type PolicyToken struct {
	AccessKeyId string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
}

func getPolicyToken(accessKeyId string, accessKeySecret string, host string, expire_time int64, upload_dir string) (string, error) {
	now := time.Now().Unix()
	expire_end := now + expire_time
	var tokenExpire = get_gmt_iso8601(expire_end)

	//create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, upload_dir)
	config.Conditions = append(config.Conditions, condition)

	//calucate signature
	result, err := json.Marshal(config)
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(accessKeySecret))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var policyToken = PolicyToken{
		AccessKeyId: accessKeyId,
		Host:        host,
		Expire:      expire_end,
		Signature:   string(signedStr),
		Directory:   upload_dir,
		Policy:      string(debyte),
	}

	response, err := json.Marshal(policyToken)
	if err != nil {
		msg := fmt.Sprintf("json err:", err)
		return "", errors.New(msg)
	}
	return string(response), nil
}
