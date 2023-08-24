package tools

import "github.com/haobinfei/ginner/config"

// 密码解密
func NewParPassword(passwd string) string {
	pass, _ := RSADecrypt([]byte(passwd), config.Conf.System.RSAPrivateBytes)
	return string(pass)
}
