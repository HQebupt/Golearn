package common

import (
	"fmt"
	p_config "linkedin/proto/config"
	"linkedin/util"
	"strconv"
)

const (
	Get_Socket_Information = "/config/get_host"
)

func GetSocketInformation(userId int64, userToken string) p_config.GetHostResponse {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	socketInfor := p_config.GetHostResponse{}
	err := util.TestGet(Host+Get_Socket_Information, userIdAsString, userToken, &socketInfor)
	fmt.Println(err)
	return socketInfor
}
