package tool

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"gin_web/Bean"
	"gin_web/dao/redis"
	"gin_web/settings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
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
const AccessTokenExpireDuration = time.Hour * 2    //短时间token
const RefreshTokenExpireDuration = time.Hour * 168 //长时间token

// token签名密钥
var mySigningKey = []byte(settings.Conf.MySigningKey)

// GenerateAccessToken 生成jwt AccessToken
func GenerateAccessToken(userID int64, username string) (string, error) {
	c := MyClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			//Issuer:    "",//发行人
			//Subject:   "",//主题
			//Audience:  nil,//观众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpireDuration)), //到期时间
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

// GenerateRefreshToken 生成jwt RefreshToken
func GenerateRefreshToken(userID int64, username string) (string, error) {
	c := MyClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			//Issuer:    "",//发行人
			//Subject:   "",//主题
			//Audience:  nil,//观众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpireDuration)), //到期时间
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
//func JWTAuthMiddleware() func(c *gin.Context) {
//	return func(c *gin.Context) {
//		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
//		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
//		// 这里的具体实现方式要依据你的实际业务情况决定
//		authHeader := c.Request.Header.Get("Authorization")
//		if authHeader == "" {
//			Bean.Error(c, Bean.USER_SESSION_EXPIRED)
//			c.Abort()
//			return
//		}
//		fmt.Println(authHeader)
//		//// 按空格分割
//		//parts := strings.SplitN(authHeader, " ", 2)
//		//if !(len(parts) == 2 && parts[0] == "Bearer") {
//		//	Bean.Error(c, Bean.USER_SESSION_EXPIRED)
//		//	c.Abort()
//		//	return
//		//}
//
//		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
//		//mc, err := ParseToken(parts[1])
//		mc, err := ParseToken(authHeader)
//
//		if err != nil {
//			Bean.Error(c, Bean.SESSION_EXPIRED)
//			c.Abort()
//			return
//		}
//		fmt.Println(mc.UserID)
//		// 将当前请求的userID信息保存到请求的上下文c上
//		c.Set(settings.Conf.ContextUserID, mc.UserID)
//
//		c.Next() // 后续的处理函数可以用过c.Get("userID")来获取当前请求的用户信息
//	}
//}

// GetContextUser 获取上下文之中的登录用户的userid
func GetContextUser(c *gin.Context) (uid int64, err error) {
	value, ok := c.Get(settings.Conf.ContextUserID)
	if !ok {
		err = Bean.SESSION_EXPIRED.MarkError
		return
	}
	uid, ok = value.(int64)
	if !ok {
		err = Bean.SESSION_EXPIRED.MarkError
		return
	}
	return
}

// OutputToken 新的生成token
func OutputToken(userID int64, username string) (token string, err error) {
	//生成access token(短时间的token)
	accesstoken, err := GenerateAccessToken(userID, username)
	if err != nil {
		return "", Bean.SYSTEM_BUSY.MarkError
	}
	//生成refresh token(长时间的token)
	refreshtoken, err := GenerateRefreshToken(userID, username)
	if err != nil {
		return "", Bean.SYSTEM_BUSY.MarkError
	}
	//将token(长时间的token)存进redis 用来验证过期时间
	formatInt := strconv.FormatInt(userID, 10)
	at := formatInt + ":accessToken"
	rt := formatInt + ":refreshToken"

	err = redis.SetUserToken(context.Background(), at, accesstoken, AccessTokenExpireDuration)
	if err != nil {
		return "", Bean.SYSTEM_BUSY.MarkError
	}
	err = redis.SetUserToken(context.Background(), rt, refreshtoken, RefreshTokenExpireDuration)
	if err != nil {
		return "", Bean.SYSTEM_BUSY.MarkError
	}
	//将access token(短时间的token)返回用作用户访问时携带
	return accesstoken, nil
}

// VerifyToken 新的中间件
func VerifyToken() func(c *gin.Context) {
	return func(c *gin.Context) {
		//解析获取请求头里的access token(短时间的token)
		authHeader := c.Request.Header.Get("AccessToken")
		if authHeader == "" {
			Bean.Error(c, Bean.USER_SESSION_EXPIRED)
			c.Abort()
			return
		}
		// ParseToken 解析Token
		mc, err := ParseToken(authHeader)
		if err != nil {
			Bean.Error(c, Bean.SESSION_EXPIRED)
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set(settings.Conf.ContextUserID, mc.UserID)
		//根据userid获取redis里的refresh token(长时间的token)
		tokensMap, err := redis.GetUserToken(context.Background(), mc.UserID)
		if err != nil {
			if err == Bean.SESSION_EXPIRED.MarkError {
				Bean.Error(c, Bean.SESSION_EXPIRED)
				c.Abort()
				return
			}
			Bean.Error(c, Bean.SYSTEM_BUSY)
			c.Abort()
			return
		}

		//  refresh token不存在就是过期了重新登录
		formatInt := strconv.FormatInt(mc.UserID, 10)
		rt := formatInt + ":refreshToken"
		at := formatInt + ":accessToken"
		_, ok := tokensMap[rt]
		if !ok {
			Bean.Error(c, Bean.SESSION_EXPIRED)
			c.Abort()
			return
		}
		//  refresh token存在
		//  请求头里的access token不存在,返回重新登录
		accessToken, ok := tokensMap[at]
		if !ok {
			Bean.Error(c, Bean.SESSION_EXPIRED)
			c.Abort()
			return
		}
		//  不一样返回重新登录
		if authHeader != accessToken {
			Bean.Error(c, Bean.SESSION_EXPIRED)
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set(settings.Conf.ContextUserID, mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get("userID")来获取当前请求的用户信息
	}

}
