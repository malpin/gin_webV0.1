package tool

import (
	"crypto/md5"
	"errors"
	"fmt"
	"gin_web/settings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

//密码加密
func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(settings.Conf.MD5salt))
	return fmt.Sprintf("%x", hash.Sum([]byte(str)))
}

// 自定义结构体内嵌jwt.StandardClaims
// jwt.StandardClaims只包含了官方字段
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// jwt token 过期时间
const TokenExpireDuration = time.Hour * 2

// token签名密钥
var mySigningKey = []byte(settings.Conf.MySigningKey)

// GenerateToken 生成jwt Token
func GenerateToken(userID int64, username string) (string, error) {
	c := MyClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			//Issuer:    "",//发行人
			//Subject:   "",//主题
			//Audience:  nil,//观众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)), //到期时间
			//NotBefore: nil,//不是以前
			//IssuedAt:  nil,//IssuedAt
			//ID:        "",//ID
		},
	}
	//创建一个新的令牌对象，指定签名方法和您希望它包含的声明。
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	//使用秘密签名并以字符串形式获取完整的编码令牌
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*MyClaims, error) {
	mClaims := new(MyClaims)
	// 解析token
	token, err := jwt.ParseWithClaims(
		tokenString,
		mClaims, //解析到
		func(token *jwt.Token) (i interface{}, err error) {
			//加盐了-mySigningKey
			return mySigningKey, nil
		})

	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*MyClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		fmt.Println(mc.UserID)
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set("userID", mc.UserID)

		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
