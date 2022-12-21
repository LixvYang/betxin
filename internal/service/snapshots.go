package service

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/lixvyang/betxin/internal/utils"

	"github.com/fox-one/pkg/uuid"
	"github.com/golang-jwt/jwt"
)

type Snapshots struct {
	Data []Datum `json:"data"`
}

type Datum struct {
	Amount     string `json:"amount"`
	Asset      Asset  `json:"asset"`
	CreatedAt  string `json:"created_at"`
	Data       string `json:"data"`
	SnapshotID string `json:"snapshot_id"`
	Source     string `json:"source"`
	Type       string `json:"type"`
	UserID     string `json:"user_id"`
	TraceID    string `json:"trace_id"`
	OpponentID string `json:"opponent_id"`
}

type Asset struct {
	AssetID string `json:"asset_id"`
	ChainID string `json:"chain_id"`
	IconURL string `json:"icon_url"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Type    string `json:"type"`
}

func GetSnapshots() {
	uri := "https://mixin-api.zeromesh.net"
	path := "/me"
	token, err := SignAuthenticationToken(utils.ClientId, utils.SessionId, utils.PrivateKey, "GET", path, "")
	if err != nil {
		println(err)
		return
	}

	req, err := http.NewRequest("GET", uri+path, bytes.NewReader(nil))
	if err != nil {
		println(err)
		return
	}

	httpClient := &http.Client{}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Request-Id", uuid.New())
	resp, err := httpClient.Do(req)
	if err != nil {
		println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		println(err)
		return
	}
	println(ioutil.ReadAll(resp.Body))
}

/*
* appID: 机器人的 id
* authorizationID: 用户授权完成后返回的 authorization_id
* privateKey: 本地生成的 private key
* method: HTTP 请求方法，GET, POST
* url: 例如 /me
* body：GET 是 ""
* scp: 用户授权时的 scope "PROFILE:READ PHONE:READ"
* requestID: 随机生成的 uuid
 */
func SignOauthAccessToken(appID, authorizationID, privateKey, method, uri, body, scp string, requestID string) (string, error) {
	expire := time.Now().UTC().Add(time.Hour * 24 * 30 * 3)
	sum := sha256.Sum256([]byte(method + uri + body))
	claims := jwt.MapClaims{
		"iss": appID,
		"aid": authorizationID,
		"iat": time.Now().UTC().Unix(),
		"exp": expire.Unix(),
		"sig": hex.EncodeToString(sum[:]),
		"scp": scp,
		"jti": requestID,
	}

	kb, err := base64.RawURLEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	priv := ed25519.PrivateKey(kb)
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(priv)
}

/*
* uid: 用户或机器人的 uuid
* sid: Session Id
* privateKey: 机器人私钥
* method: HTTP 请求方法 GET, POST
* url: HTTP 请求 URL 例如: /transfers
* body: HTTP 请求内容, 例如: {"pin": "encrypted pin token"}
 */
func SignAuthenticationToken(uid, sid, privateKey, method, uri, body string) (string, error) {
	expire := time.Now().UTC().Add(time.Hour * 24 * 30 * 3)
	sum := sha256.Sum256([]byte(method + uri + body))
	claims := jwt.MapClaims{
		"uid": uid,
		"sid": sid,
		"iat": time.Now().UTC().Unix(),
		"exp": expire.Unix(),
		"jti": uuid.New(),
		"sig": hex.EncodeToString(sum[:]),
		"scp": "FULL",
	}
	priv, err := base64.RawURLEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	// more validate the private key
	if len(priv) != 64 {
		return "", fmt.Errorf("bad ed25519 private key %s", priv)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(ed25519.PrivateKey(priv))
}
