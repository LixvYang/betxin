package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UserInfo struct {
	Data Data `json:"data"`
}

type Data struct {
	UserID         string `json:"user_id"`
	IdentityNumber string `json:"identity_number"`
	FullName       string `json:"full_name"`
	Biography      string `json:"biography"`
	AvatarURL      string `json:"avatar_url"`
	SessionID      string `json:"session_id"`
	PinToken       string `json:"pin_token"`
	PinTokenBase64 string `json:"pin_token_base64"`
	Phone          string `json:"phone"`
}

func GetUserInfo(access_token string) (Data, error) {
	// 形成请求
	var userInfoUrl = "https://api.mixin.one/me" // mixin
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, userInfoUrl, nil); err != nil {
		return Data{}, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))

	// 发送请求并获取响应
	var client = http.Client{}
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return Data{}, err
	}
	defer res.Body.Close()
	// 将响应的数据写入 userInfo 中，并返回
	var userInfo UserInfo
	if err = json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return Data{}, err
	}
	return userInfo.Data, nil
}
