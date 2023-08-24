package middleware

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/haobinfei/ginner/config"
	"github.com/haobinfei/ginner/model"
	"github.com/haobinfei/ginner/model/request"
	"github.com/haobinfei/ginner/public/tools"
	"github.com/haobinfei/ginner/service/isql"
)

func InitJwtMiddleware() (*jwt.GinJWTMiddleware, error) {
	var err error
	var authMiddleware *jwt.GinJWTMiddleware

	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:           config.Conf.Jwt.Realm,                                 // jwt标识
		Key:             []byte(config.Conf.Jwt.Key),                           // 服务端密钥
		Timeout:         time.Hour * time.Duration(config.Conf.Jwt.Timeout),    // token过期时间
		MaxRefresh:      time.Hour * time.Duration(config.Conf.Jwt.MaxRefresh), // token最大刷新时间(RefreshToken过期时间=Timeout+MaxRefresh)
		PayloadFunc:     payloadFunc,                                           // 有效载荷处理
		IdentityHandler: identityHandler,                                       // 解析Claims
		Authenticator:   login,                                                 // 校验token的正确性, 处理登录逻辑
		Authorizator:    authorizator,                                          // 用户登录校验成功处理
		Unauthorized:    unauthorized,                                          // 用户登录校验失败处理
		LoginResponse:   loginResponse,                                         // 登录成功后的响应
		LogoutResponse:  logoutResponse,                                        // 登出后的响应
		RefreshResponse: refreshResponse,                                       // 刷新token后的响应
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",    // 自动在这几个地方寻找请求中的token
		TokenHeadName:   "Bearer",                                              // header名称
		TimeFunc:        time.Now,
	})

	return authMiddleware, err
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		var user model.User
		tools.JsonI2Struct(v["user"], &user)
		return jwt.MapClaims{
			jwt.IdentityKey: user.ID,
			"user":          v["user"],
		}
	}
	return jwt.MapClaims{}
}

// 解析Claims
func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	// 此处返回值类型map[string]interface{}与payloadFunc和authorizator的data类型必须一致, 否则会导致授权失败还不容易找到原因
	return tools.H{
		"IdentityKey": claims[jwt.IdentityKey],
		"user":        claims["user"],
	}
}

// 解析token的正确性，登录处理逻辑
func login(c *gin.Context) (interface{}, error) {
	var req request.RegisterAndLoginReq
	if err := c.ShouldBind(&req); err != nil {
		return "", err
	}

	decodeData, err := tools.RSADecrypt([]byte(req.Password), config.Conf.System.RSAPrivateBytes)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		UserName: req.Username,
		Password: string(decodeData),
	}

	user, err := isql.User.Login(u)
	if err != nil {
		return nil, err
	}

	return tools.H{
		"user": tools.Struct2Json(user),
	}, nil
}
