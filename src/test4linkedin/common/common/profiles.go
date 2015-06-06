package common

import (
	"fmt"
	p_profile "linkedin/proto/profile"
	"linkedin/util"
	"strconv"

	"github.com/golang/protobuf/proto"
)

const (
	// Http endpoints
	Get_User_Info       = "/api/profile/id/"
	Update_User_Profile = "/api/profile/update_profile"
	Update_User_Setting = "/api/profile/setting/update"

	Get_User_Setting     = "/api/profile/setting/"
	Change_Password      = "/api/profile/change_pwd"
	Update_Privacy       = "/api/profile/privacy_update"
	Get_Privacy          = "/api/profile/privacy"
	Get_Profile_By_Phone = "/api/profile/phone/%s"
)

func GetUserProfileWithLogin(phone, passwd string) p_profile.Profile {

	loginResp := LoginByPhone(phone, passwd)
	userId := strconv.FormatInt(int64(loginResp.GetUserID()), 10)

	fmt.Println("user id is: " + userId)
	rsp := p_profile.Profile{}
	err := util.TestGet(Host+Get_User_Info+userId, userId, loginResp.GetToken(), &rsp)
	if err != nil {
		fmt.Println("err information:", err)
	}

	return rsp
}

func GetUserProfileByUserID(userId string) p_profile.Profile {
	rsp := p_profile.Profile{}
	err := util.TestGet(Host+Get_User_Info+userId, userId, "fake_token", &rsp)
	if err != nil {
		fmt.Println("err information:", err)
	}

	return rsp
}

func GetUserProfileByPhone(userId, token, phone string) p_profile.Profile {
	get_profile_by_phone := fmt.Sprintf(Get_Profile_By_Phone, phone)
	fmt.Println(get_profile_by_phone)
	rsp := p_profile.Profile{}
	err := util.TestGet(Host+get_profile_by_phone, userId, token, &rsp)
	if err != nil {
		fmt.Println("err information:", err)
	}

	return rsp
}

func UpdateUserProfile(phone string, passwd string, requestData proto.Message) (int, string) {
	loginResp := LoginByPhone(phone, passwd)
	userId := strconv.FormatInt(int64(loginResp.GetUserID()), 10)

	status, body := util.TestPostNormal(Host+Update_User_Profile, userId, loginResp.GetToken(), requestData)
	return status, body
}

func UpdateUserSettings(userId int64, userToken string, requestData proto.Message) (int, string) {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	status, body := util.TestPostNormal(Host+Update_User_Setting, userIdAsString, userToken, requestData)
	return status, body
}

func GetUserSettings(userId int64, userToken string) (error, p_profile.Settings) {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	resp := p_profile.Settings{}
	err := util.TestGet(Host+Get_User_Setting, userIdAsString, userToken, &resp)
	return err, resp
}

func ChangePassword(userId int64, userToken string, requestData proto.Message) (int, string) {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	//Todo: It should call util.TestPost to get the http response, however it is ok for this api as well.
	status, body := util.TestPostNormal(Host+Change_Password, userIdAsString, userToken, requestData)
	return status, body
}

func UpdateUserPrivacy(userId int64, userToken string, requestData proto.Message) (int, string) {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	status, body := util.TestPostNormal(Host+Update_Privacy, userIdAsString, userToken, requestData)
	return status, body
}

func GetUserPrivacy(userId int64, userToken string, requestData proto.Message) error {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	err := util.TestGet(Host+Get_Privacy, userIdAsString, userToken, requestData)
	return err
}
