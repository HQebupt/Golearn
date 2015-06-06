package common

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	p_user "linkedin/proto/user"
	"linkedin/util"
	"net/http"
	"strconv"
)

const (
	// User Related information
	Passwd     = "dddd"
	NewPwd     = "aaaa"
	Phone      = "15652264009"
	Name       = "test3"
	Industry   = 111111
	Career     = 222222
	WebiboId   = "weibo1"
	LinkedinID = "linkedinId1"
	ImageUrl   = "image1"

	// Http endpoints
	//Host = "http://master.inicn.com"
	Host = "http://stable.inicn.com"
	//Host = "http://localhost:8080"
	//Host                    = "http://192.168.200.76:8080"
	Metadata                = "application/x-protobuf"
	Regist                  = "/user/phone_register_step2"
	Login_phone             = "/user/login_phone"
	Login_id                = "/user/login_id"
	Check_phone_duplication = "/user/phone_register_step0"
	Forget_pwd_verify       = "/user/forget_pwd_verify"
	Update_pwd              = "/user/update_pwd"
)

type UserInfor struct {
	Phone     string
	Pwd       string
	Name      string
	UserId    int64
	UserToken string
}

func InitTestUserInformation(numOfUsers int32) []UserInfor {
	userList := make([]UserInfor, numOfUsers)
	for index := 0; index < len(userList); index++ {
		userList[index].Phone = fmt.Sprintf("120%08d", index)
		userList[index].Pwd = strconv.FormatInt(int64(120), 10)
		userList[index].Name = fmt.Sprintf("TestUser%d", index)

		// if the phone has been registerd, it should return fail but no matter for this test code,
		// as it will login successfully and return the user id.
		RegisterWithPhone(userList[index].Phone, userList[index].Pwd, userList[index].Name)
		resp := LoginByPhone(userList[index].Phone, userList[index].Pwd)
		userList[index].UserId = resp.GetUserID()
		userList[index].UserToken = resp.GetToken()
		fmt.Printf("User %s have been created,userToken is:%s. Password is: %s \n", userList[index].Phone, userList[index].UserToken, userList[index].Pwd)
	}

	return userList
}

func Register() int64 {
	userInfo := p_user.RegisterRequest{
		Name:       proto.String(Name),
		Phone:      proto.String(Phone),
		Passwd:     proto.String(Passwd),
		Industry:   proto.Int64(Industry),
		Career:     proto.Int64(Career),
		ImageURL:   proto.String(ImageUrl),
	}

	buf := util.MustMarshal(&userInfo)
	a, err1 := http.Post(Host+Regist, Metadata, bytes.NewReader(buf))
	if err1 != nil {

		fmt.Println(err1)
	}

	resp, err2 := ioutil.ReadAll(a.Body)
	if err2 != nil {

		fmt.Println(err2)
	}

	userRegisterInformation := p_user.RegisterResponse{}
	proto.Unmarshal(resp, &userRegisterInformation)
	userId := userRegisterInformation.GetUserID()
	return userId
}

func RegisterWithPhone(phone, passwd, name string) int64 {
	userInfo := p_user.RegisterRequest{
		Name:       proto.String(name),
		Phone:      proto.String(phone),
		Passwd:     proto.String(passwd),
		Industry:   proto.Int64(1),
		Career:     proto.Int64(1),
		ImageURL:   proto.String(""),
	}

	buf := util.MustMarshal(&userInfo)
	a, err1 := http.Post(Host+Regist, Metadata, bytes.NewReader(buf))
	if err1 != nil {

		fmt.Println(err1)
	}

	resp, err2 := ioutil.ReadAll(a.Body)
	if err2 != nil {
		fmt.Println(err2)
	}

	userRegisterInformation := p_user.RegisterResponse{}
	proto.Unmarshal(resp, &userRegisterInformation)
	userId := userRegisterInformation.GetUserID()
	return userId
}

func LoginByPhone(phone, passwd string) p_user.LoginResponse {
	userInfo := p_user.LoginByPhoneRequest{
		Phone:  proto.String(phone),
		Passwd: proto.String(passwd),
	}
	buf := util.MustMarshal(&userInfo)
	a, err := http.Post(Host+Login_phone, Metadata, bytes.NewReader(buf))
	if err != nil {
		fmt.Println(err)
	}
	resp, err2 := ioutil.ReadAll(a.Body)
	if err2 != nil {
		fmt.Println(err2)
	}

	userResp := p_user.LoginResponse{}
	proto.Unmarshal(resp, &userResp)

	return userResp
}

func LoginById(userId int64, passwd string) p_user.LoginResponse {
	userInfo := p_user.LoginByIDRequest{
		UserID: proto.Int64(userId),
		Passwd: proto.String(passwd),
	}

	buf := util.MustMarshal(&userInfo)
	a, err := http.Post(Host+Login_id, Metadata, bytes.NewReader(buf))
	if err != nil {
		fmt.Println(err)
	}
	resp, err2 := ioutil.ReadAll(a.Body)
	if err2 != nil {
		fmt.Println(err2)
	}

	userResp := p_user.LoginResponse{}
	proto.Unmarshal(resp, &userResp)

	return userResp
}
