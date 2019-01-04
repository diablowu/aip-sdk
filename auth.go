package baidu_ai_sdk

import (
	"fmt"
	"github.com/diablowu/baidu-ai-sdk/internal/http"
	"strings"
)


// refs: https://cloud.baidu.com/doc/Reference/AuthenticationMechanism.html
type BceCredential struct {
	ApiKey    string
	SecretKey string
	AppID     string
}

// 生成access_token获取的query string
func (cred BceCredential) QueryParam() string {
	return fmt.Sprintf("?grant_type=%s&client_id=%s&client_secret=%s", aip_grant_type, ai.ak, ai.sk)
}

// unit access token 结构
type AccessTokenResponse struct {
	RefreshToken     string `json:"refresh_token" bson:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in" bson:"expires_in"`
	Scope            string `json:"scope" bson:"scope"`
	SessionKey       string `json:"session_key" bson:"session_key"`
	AccessToken      string `json:"access_token" bson:"access_token"`
	SessionSecret    string `json:"session_secret" bson:"session_secret"`
	Error            string `json:"error" bson:"error"`
	ErrorDescription string `json:"error_description" bson:"error_description"`
}

func (rsp AccessTokenResponse) Get() (bool, string) {
	if len(strings.TrimSpace(rsp.Error)) > 0 {
		return false, ""
	} else {
		return true, rsp.AccessToken
	}
}

func (rsp AccessTokenResponse) Ok() (bool) {
	return len(strings.TrimSpace(rsp.Error)) == 0
}

type AuthService struct {
}

func (auth AuthService) Token(url string) string {
	if bs, err := http.GetJson(url); err == nil {
		return string(bs)
	} else {
		return ""
	}
}
