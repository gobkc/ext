package jwt

import (
	"testing"
)

func TestNewJwt(t *testing.T) {
	jwt := NewJwt("AF9-C=AF,FJN+RVV(DDD(SFF")
	userClaims := UserClaims{
		Username:    "xiong",
		Nickname:    "aaa",
		CountryCode: "+86",
		Phone:       "13198520987",
		Email:       "10044375@qq.com",
		State:       1,
		Type:        1,
		LastIP:      "127.0.0.1",
	}
	token, _ := jwt.Token(&userClaims)
	testClaims := UserClaims{}
	jwt.Verify(token, &testClaims)
	oldToken := token
	jwt.Refresh(&token)
	if oldToken != token {
		t.Error("token error")
	}
}
