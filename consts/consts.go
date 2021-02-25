package consts

const (
	VerifyRoute = "https://lkdr.nalog.ru/api/v1/auth/challenge/sms/verify"
	StartRoute  = "https://lkdr.nalog.ru/api/v1/auth/challenge/sms/start"
	CodeFromSMS = `{
   "challengeToken":"%s",
   "phone":"%d",
   "code": "%d",
   "deviceInfo":{
      "sourceDeviceId":"deviceid",
      "sourceType":"WEB",
      "appVersion":"1.0.0",
      "metaDetails":{
         "userAgent":"useragent"
      }
   }
}`
)
