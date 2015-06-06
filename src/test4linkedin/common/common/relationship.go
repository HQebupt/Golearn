package common

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	p_base "linkedin/proto/base"
	p_profile "linkedin/proto/profile"
	p_relationship "linkedin/proto/relationship"
	"linkedin/util"
	"strconv"
)

const (
	Friend_Request     = "/api/connect_request"
	Friend_Accept      = "/api/accept_connect"
	Get_Friend_List    = "/api/friend_list"
	Get_Friend_Profile = "/api/friend_profiles"
	Disconnect_Friends = "/api/disconnect"
	Check_Invtation    = "/api/check_invitation"
	Follow_People      = "/api/follow"
	Unfollow_people    = "/api/unfollow"
	Get_Follow_List    = "/api/follow_list"
	Get_Follower_List  = "/api/follower_list"
)

func sendFriendRequest(userId int64, userToken string, requestData proto.Message) (p_relationship.CommonResponseStatus, int) {
	userIdAsString := strconv.FormatInt(int64(userId), 10)

	response := p_relationship.CommonResponseStatus{}
	statusCode, err := util.TestPost(Host+Friend_Request, userIdAsString, userToken, requestData, &response)
	fmt.Println(err, response)
	return response, statusCode
}

func acceptFriendRequest(userId int64, userToken string, requestData proto.Message) (p_relationship.CommonResponseStatus, int) {
	userIdAsString := strconv.FormatInt(int64(userId), 10)

	response := p_relationship.CommonResponseStatus{}
	statusCode, err := util.TestPost(Host+Friend_Accept, userIdAsString, userToken, requestData, &response)
	fmt.Println(err, response)
	return response, statusCode
}

func disconnectFriendConnection(userId int64, userToken string, requestData proto.Message) (p_base.OkResponse, int) {
	userIdAsString := strconv.FormatInt(int64(userId), 10)

	response := p_base.OkResponse{}
	statusCode, err := util.TestPost(Host+Disconnect_Friends, userIdAsString, userToken, requestData, &response)
	fmt.Println(err, response)
	return response, statusCode
}

func followPeople(userId int64, userToken string, requestData proto.Message) (p_base.OkResponse, int) {
	userIdAsString := strconv.FormatInt(int64(userId), 10)

	response := p_base.OkResponse{}
	statusCode, err := util.TestPost(Host+Follow_People, userIdAsString, userToken, requestData, &response)
	fmt.Println(err, response)
	return response, statusCode
}

func unfollowPeople(userId int64, userToken string, requestData proto.Message) (p_base.OkResponse, int) {
	userIdAsString := strconv.FormatInt(int64(userId), 10)

	response := p_base.OkResponse{}
	statusCode, err := util.TestPost(Host+Unfollow_people, userIdAsString, userToken, requestData, &response)
	fmt.Println(err, response)
	return response, statusCode
}

func GetFollowList(userId int64, userToken string) ([]int64, error) {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	response := p_relationship.UserListResponse{}
	error := util.TestGet(Host+Get_Follow_List, userIdAsString, userToken, &response)
	fmt.Println("Get Follow list is,", response)
	return response.GetUserId(), error
}

func GetFollowerList(userId int64, userToken string) ([]*p_profile.Profile, error) {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	response := p_profile.GetProfileListResponse{}
	error := util.TestGet(Host+Get_Follower_List, userIdAsString, userToken, &response)
	fmt.Println("Get Follower list is,", response.GetProfiles())
	return response.GetProfiles(), error
}

func GetFriendList(userId int64, userToken string) ([]int64, error) {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	response := p_relationship.UserListResponse{}
	error := util.TestGet(Host+Get_Friend_List, userIdAsString, userToken, &response)
	fmt.Println("Get Friend list is,", response.GetUserId())
	return response.GetUserId(), error
}

func GetFriendProfile(userId int64, userToken string) ([]*p_profile.Profile, error) {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	response := p_profile.GetProfileListResponse{}
	error := util.TestGet(Host+Get_Friend_Profile, userIdAsString, userToken, &response)
	fmt.Println("Get Friend Profile is,", response)
	return response.GetProfiles(), error
}

func MakeFriendConnection(inviter, invitee UserInfor) bool {
	inviteRequest := p_relationship.FriendRequest{
		Uid:  proto.Int64(inviter.UserId),
		Tid:  proto.Int64(invitee.UserId),
		Name: proto.String(inviter.Name),
		Msg:  proto.String("Test For invite"),
		Type: proto.String("invite"),
	}

	inviteResp, statusCode1 := sendFriendRequest(inviter.UserId, inviter.UserToken, &inviteRequest)
	fmt.Println("Invite reslt:", inviteResp, statusCode1)
	if statusCode1 != 200 {
		return false
	}

	acceptRequest := p_relationship.FriendRequest{
		Uid:  proto.Int64(invitee.UserId),
		Tid:  proto.Int64(inviter.UserId),
		Name: proto.String(invitee.Name),
		Msg:  proto.String("Test For accept"),
		Type: proto.String("accept"),
	}

	acceptResp, statusCode2 := acceptFriendRequest(invitee.UserId, invitee.UserToken, &acceptRequest)
	fmt.Println("Accept reslt:", acceptResp, statusCode2)

	if statusCode2 != 200 {
		return false
	}

	friendList, statusCode3 := GetFriendList(invitee.UserId, invitee.UserToken)
	fmt.Println(friendList, statusCode3, friendList[0], friendList[1])
	for i := 0; i < len(friendList); i++ {
		if int64(friendList[i]) == inviter.UserId {
			return true
		}
	}

	return false
}

func FollowConnection(follower, beFollow UserInfor) bool {
	request := p_relationship.Request{
		Dst: proto.Int64(beFollow.UserId),
	}

	_, statusCode1 := followPeople(follower.UserId, follower.UserToken, &request)
	if 200 != statusCode1 {
		return false
	}

	return true
}

func UnfollowConnection(follower, beFollow UserInfor) bool {
	request := p_relationship.Request{
		Dst: proto.Int64(beFollow.UserId),
	}

	_, statusCode1 := unfollowPeople(follower.UserId, follower.UserToken, &request)
	if 200 != statusCode1 {
		return false
	}

	return true
}

func DisconnectFriendConnection(friendA, friendB UserInfor) bool {
	request := p_relationship.Request{
		Dst: proto.Int64(friendB.UserId),
	}

	resp, statusCode := disconnectFriendConnection(friendA.UserId, friendA.UserToken, &request)
	fmt.Println("Disconnect friends,", resp, statusCode)
	if 200 != statusCode {
		return false
	}

	friendList, _ := GetFriendList(friendA.UserId, friendA.UserToken)
	for i := 0; i < len(friendList); i++ {
		if int64(friendList[i]) == friendB.UserId {
			return false
		}
	}

	return true
}
