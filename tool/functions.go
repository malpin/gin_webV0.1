package tool

import (
	"crypto/md5"
	"fmt"
	"gin_web/settings"
)

func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(settings.Conf.MD5salt))
	return fmt.Sprintf("%x", hash.Sum([]byte(str)))
}
