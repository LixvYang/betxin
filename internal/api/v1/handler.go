package v1

import (
	"net/http"
	"time"

	"github.com/lixvyang/betxin/internal/utils/errmsg"

	"github.com/gin-gonic/gin"
)

const (
	REDISEXPIRE = time.Hour * 2

	CATEGORY_LIST  = "category_list"
	CATEGORY_TOTAL = "category_total"
	CATEGORY_GET   = "category_get_"

	// FI
	// topic_list_${limit}_${offset}
	// topic_list_from_cate_${cid}_${limit}_${offset}
	// topic_list_from_cate_total_${cid}
	TOPIC_LIST                = "topic_list"
	TOPIC_TOTAL               = "topic_total"
	TOPIC_GET                 = "topic_get_"
	TOPIC_LIST_FROMCATE       = "topic_list_from_cate_"
	TOPIC_LIST_FROMCATE_TOTAL = "topic_list_from_cate_total_"

	COLLECT_LIST           = "collect_list"
	COLLECT_TOTAL          = "collect_total"
	COLLECT_GET_USER_LIST  = "collect_get_user_list_"
	COLLECT_GET_USER_TOTAL = "collect_get_user_total_"

	// 存储用户购买的topic
	USERTOTOPIC_USER_TOTAL = "usertotopic_user_total_"
	USERTOTOPIC_USER_LIST  = "usertotopic_user_list_"
	// 存储哪些用户购买的topic
	USERTOTOPIC_TOPIC_TOTAL = "usertotopic_topic_total_"
	USERTOTOPIC_TOPIC_LIST  = "usertotopic_topic_list_"

	// user_to_topic
	// ===================
	//  usertotopic_list_:uid
	USERTOTOPIC_LIST = "usertotopic_list_"

	//
	USER_LIST  = "user_list"
	USER_TOTAL = "user_total"
	//
	USER_INFO = "user_info_"

	// currency
	CURRENCY_LIST  = "currency_list"
	CURRENCY_TOTAL = "currency_total"

	// comment list by tid
	// ===================
	// comment_list_:tid_
	COMMENT_LIST    = "comment_list_"
	COMMENT_MAXTIME = 9999999999999

	// check person parise comment
	// ===================
	// parise_comment_cid_:cid
	// key: uid
	PARISECOMMENT     = "parise_comment_cid_"
	PARISECOMMENT_KEY = "uid_"

	// USER INFO
	// ===================
	// key: userinfo_:uid
	// value: marshal useinfo
	USERINFO = "userinfo_"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendResponse(c *gin.Context, code int, data interface{}) {
	message := errmsg.GetErrMsg(code)
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
