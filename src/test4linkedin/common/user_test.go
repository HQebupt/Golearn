package test4linkedin

import (
	"bytes"
	"io/ioutil"
	p_user "linkedin/proto/user"
	"linkedin/util"
	"net/http"
	p_common "test4linkedin/common"
	"testing"
	"github.com/golang/protobuf/proto"
)

//it only works for a non-registed phone
func TestRegist(t *testing.T) {
	userId := p_common.Register()
	userResp := p_common.LoginById(userId, p_common.Passwd)
	if p_common.Name != userResp.GetName() {
		t.Error("Get the returned Name error:" + userResp.GetName())
	}

	if p_common.Phone != userResp.GetPhone() {
		t.Error("Get the returned Phone error:" + userResp.GetPhone())
	}

	if p_common.Industry != userResp.GetIndustry() {
		t.Error(userResp.GetIndustry())
	}

	if p_common.Career != userResp.GetCareer() {
		t.Error(userResp.GetCareer())
	}

	if p_common.ImageUrl != userResp.GetImageURL() {
		t.Error("Get the returned Image error:" + userResp.GetImageURL())
	}
}

func testLogin(t *testing.T) {
	userResp := p_common.LoginByPhone(p_common.Phone, p_common.Passwd)
	t.Log(userResp)
	t.Log(userResp.GetToken())
	if 20 != len(userResp.GetToken()) {
		t.Error("the token lenght is invalid")
	}

	if p_common.Name != userResp.GetName() {
		t.Error("the returned user name is wrong")
	}
}

func testLoginByInvalidId(t *testing.T) {
	var testInvalidIDs = map[string]string{
		"11111111111":                   "110",
		"111":                           "11",
		"abcdea":                        "11",
		"asdfasdfadzxcvasdfwerasdfasdf": "asdfxvcasdfweafsdfzxcvasdfasdfasdf",
		"&&*....12312":                  "",
		"":                              "",
	}

	for key, value := range testInvalidIDs {
		userResp := p_common.LoginByPhone(key, value)
		if 3 != userResp.GetStatus() {
			t.Error("Login should be failed as the login phone/passwd are wrong!")
		}
	}
}

func testDupPhoneRegister(t *testing.T) {
	userInfo := p_user.RegisterRequest{
		Name:   proto.String("test1"),
		Phone:  proto.String(p_common.Phone),
		Passwd: proto.String(p_common.Passwd),
	}
	buf1 := util.MustMarshal(&userInfo)
	a, err1 := http.Post(p_common.Host+p_common.Regist, p_common.Metadata, bytes.NewReader(buf1))
	if err1 != nil {
		t.Log(err1)
	}

	phoneInfo := p_user.CellphoneCheckRequest{
		Phone: proto.String(p_common.Phone),
	}

	buf2 := util.MustMarshal(&phoneInfo)
	a, err2 := http.Post(p_common.Host+p_common.Check_phone_duplication, p_common.Metadata, bytes.NewReader(buf2))
	if err2 != nil {
		t.Log(err2)
	}

	resp, err3 := ioutil.ReadAll(a.Body)
	if err3 != nil {
		t.Log(err3)
	}

	checkPhoneResp := p_user.CellphoneCheckResponse{}
	proto.Unmarshal(resp, &checkPhoneResp)
	t.Log(checkPhoneResp.GetStatus())
	// 1 means duplication error code
	if 1 != checkPhoneResp.GetStatus() {
		t.Error("It should return duplication error code")

	}
}
