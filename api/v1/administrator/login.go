package administrator

import (
	"time"

	"github.com/lixvyang/betxin/model"

	v1 "github.com/lixvyang/betxin/api/v1"

	"github.com/lixvyang/betxin/internal/utils/errmsg"
	myjwt "github.com/lixvyang/betxin/internal/utils/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type LoginResponse struct {
	Id       int    `json:"id"`
	Token    string `json:"token"`
	UserName string `json:"username"`
}

// Login 后台登陆
func Login(c *gin.Context) {
	var formData model.Administrator
	_ = c.ShouldBindJSON(&formData)
	var code int

	formData, code = model.CheckLogin(formData.Username, formData.Password)
	if code == errmsg.SUCCSE {
		setToken(c, formData)
	} else {
		v1.SendResponse(c, errmsg.ERROR, &LoginResponse{
			Id:       0,
			Token:    "",
			UserName: "",
		})
		return
	}
}

// token生成函数
func setToken(c *gin.Context, user model.Administrator) {
	j := myjwt.NewJWT()
	claims := myjwt.MyClaims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 100,
			ExpiresAt: time.Now().Unix() + 604800,
			Issuer:    "Lixv",
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		v1.SendResponse(c, errmsg.ERROR, &LoginResponse{
			Id:       0,
			Token:    "",
			UserName: "",
		})
		return
	}

	v1.SendResponse(c, errmsg.SUCCSE, &LoginResponse{
		Id:       user.Id,
		Token:    token,
		UserName: user.Username,
	})
}
