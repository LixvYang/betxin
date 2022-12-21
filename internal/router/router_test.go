package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/lixvyang/betxin/api/v1"
	"github.com/lixvyang/betxin/api/v1/topic"

	"github.com/gin-gonic/gin"
	"github.com/zeebo/assert"
)

func SetupRouter() *gin.Engine {
	return gin.Default()
}

func TestListTopicsHandler(t *testing.T) {
	r := SetupRouter()
	r.POST("/api/v1/topic/list", topic.ListTopics)

	postBody := []byte(`{"offset": 0,"limit": 10,"title": "","content": ""}`)
	req, _ := http.NewRequest("POST", "http://localhost:3000/api/v1/topic/list", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var res v1.Response
	json.Unmarshal(body, &res)
	assert.Equal(t, http.StatusOK, w.Code)
}
