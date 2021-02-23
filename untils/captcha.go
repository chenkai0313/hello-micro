package untils

import (
	"errors"
	"fmt"
	"github.com/mojocn/base64Captcha"
	"hello-micro/app"
	"time"
)

type configJsonBody struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

var store = base64Captcha.DefaultMemStore

func GenerateCaptchaHandler(imgHeight, imgWidth int) (string, string) {
	var param configJsonBody
	param.DriverDigit = &base64Captcha.DriverDigit{
		Height:   imgHeight,
		Width:    imgWidth,
		Length:   4,
		MaxSkew:  1,
		DotCount: 1,
	}
	var driver base64Captcha.Driver
	driver = param.DriverDigit
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _ := c.Generate()
	CaptchaVal := c.Store.Get(id, true)
	fmt.Println("CaptchaVal", CaptchaVal)
	//将信息存入redis
	if _, err := app.RedisDB.Set(id, CaptchaVal, 600*time.Second).Result(); err != nil {
		app.ZapLog.Error("captcha-redis-error", " 验证码保存redis失败 "+err.Error())
		return "", ""
	}
	return id, b64s
}

func CaptchaVerifyHandle(id, VerifyValue string) error {
	val, _ := app.RedisDB.Get(id).Result()
	if val == "" {
		return errors.New("验证码过期")
	}
	//删除验证码
	app.RedisDB.Del(id)
	if val != VerifyValue {
		return errors.New("验证码不正确")
	}
	return nil
}
