package common

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	p_base "linkedin/proto/base"
	p_group "linkedin/proto/group"
	"linkedin/util"
	"strconv"
)

const (
	Create_Group                = "/api/groups/create"
	Create_Multichat            = "/api/groups/create_multichat"
	Get_Group_Detail            = "/api/groups/groupdetail/"
	Delete_Group                = "/api/group-multi-chat/%d/owner/drop"
	Get_Group_List              = "/api/groups/grouplist"
	Apply_Join_Group            = "/api/group/%d/apply"
	Invite_To_Join_Group        = "/api/group/%d/admin/invite"
	Approve_Join_Group          = "/api/group/%d/member/approve"
	Set_Group_To_Compere_Mode   = "/api/group/%d/admin/start-moderator-mode"
	Search_Group                = "/api/groups/search"
	Set_Group_Visible           = "/api/group/%d/owner/setvisible"
	Set_Group_Invisible         = "/api/group/%d/owner/setinvisible"
	Get_Joined_User_List        = "/api/groups/grouplist"
	Remove_User                 = "/api/group-multi-chat/%d/admin/remove"
	Add_Memeber_Admin_Access    = "/api/group/%d/admin/add-admin"
	Remove_Memmber_Admin_Access = "/api/group/%d/admin/remove-admin"
	Exit_Group                  = "/api/group-multi-chat/%d/member/exit"
	Start_Moderator_Mode        = "/api/group/%d/admin/start-moderator-mode"
	Update_Moderator_Mode       = "/api/group/%d/admin/update-moderator-mode"
	Stop_Moderator_Mode         = "/api/group/%d/admin/stop-moderator-mode"
)

const (
	GroupOwner  int32 = 1
	GroupAdmin  int32 = 2
	GroupMember int32 = 3
)

func CreateGroup(userId int64, userToken string, requestData proto.Message) p_group.GroupInfoResponse {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	groupResponse := p_group.GroupInfoResponse{}
	_, status := util.TestPost(Host+Create_Group, userIdAsString, userToken, requestData, &groupResponse)
	fmt.Println(status)
	return groupResponse
}

func CreateMultiChat(userId int64, userToken string, requestData proto.Message) p_group.GroupInfoResponse {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	groupResponse := p_group.GroupInfoResponse{}
	_, status := util.TestPost(Host+Create_Multichat, userIdAsString, userToken, requestData, &groupResponse)
	fmt.Println(status)
	return groupResponse
}

func GetGroupDetails(userId int64, userToken string, groupId int64) p_group.GroupDetailResponse {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	groupResponse := p_group.GroupDetailResponse{}
	err := util.TestGet(Host+Get_Group_Detail+strconv.FormatInt(int64(groupId), 10), userIdAsString, userToken, &groupResponse)
	fmt.Println(err)
	return groupResponse
}

func RemoveGroup(userId int64, userToken string, groudId int64) (p_base.OkResponse, int) {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	drop_group_url := fmt.Sprintf(Delete_Group, groudId)
	fmt.Println(drop_group_url)

	response := p_base.OkResponse{}
	statusCode, err := util.TestPostWithoutInput(Host+drop_group_url, userIdAsString, userToken, &response)
	fmt.Println("remove group result", err, statusCode)
	return response, statusCode
}

func GetGroupListByUser(userId int64, userToken string) p_group.GroupListResponse {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	response := p_group.GroupListResponse{}

	err := util.TestGet(Host+Get_Group_List, userIdAsString, userToken, &response)
	fmt.Println("Get Group list fail", err)
	return response
}

func SetGroupVisible(userId int64, userToken string, groupId int64) p_base.OkResponse {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	set_group_visible := fmt.Sprintf(Set_Group_Visible, groupId)
	fmt.Println(set_group_visible)

	response := p_base.OkResponse{}
	err := util.TestGet(Host+set_group_visible, userIdAsString, userToken, &response)
	fmt.Println(err)
	return response
}

func SetGroupInvisible(userId int64, userToken string, groupId int64) p_base.OkResponse {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	set_group_invisible := fmt.Sprintf(Set_Group_Invisible, groupId)
	fmt.Println(set_group_invisible)

	response := p_base.OkResponse{}
	err := util.TestGet(Host+set_group_invisible, userIdAsString, userToken, &response)
	fmt.Println(err)
	return response
}

func InviteUserToGroup(userId int64, userToken string, groupId int64, requestData proto.Message) p_group.GroupInfoResponse {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	invite_url := fmt.Sprintf(Invite_To_Join_Group, groupId)
	fmt.Println(invite_url)

	response := p_group.GroupInfoResponse{}
	_, err := util.TestPost(Host+invite_url, userIdAsString, userToken, requestData, &response)
	fmt.Println(err, response)
	return response
}

func ApplyToJoinAGroup(userId int64, userToken string, groupId int64, requestData proto.Message) int {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	apply_group_url := fmt.Sprintf(Apply_Join_Group, groupId)
	fmt.Println(apply_group_url)

	response := p_base.OkResponse{}
	statusCode, err := util.TestPost(Host+apply_group_url, userIdAsString, userToken, requestData, &response)
	fmt.Println(err, response)
	return statusCode
}

func ApprovalJoinGroupRequest(userId int64, userToken string, groupId int64, requestData proto.Message) p_group.ApproveUserResponse {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	approve_group_url := fmt.Sprintf(Approve_Join_Group, groupId)
	fmt.Println(approve_group_url)

	response := p_group.ApproveUserResponse{}
	_, err := util.TestPost(Host+approve_group_url, userIdAsString, userToken, requestData, &response)
	fmt.Println(err, response)
	return response
}

func RemoveUser(adminId int64, userToken string, groupId int64, removedUserId int64) (p_base.OkResponse, int) {
	adminIdAsString := strconv.FormatInt(int64(adminId), 10)
	remove_url := fmt.Sprintf(Remove_User, groupId)
	fmt.Println(remove_url)
	removeUserRequest := p_group.RemoveUserRequest{
		UserId: proto.Int64(removedUserId),
	}
	response := p_base.OkResponse{}
	statusCode, err := util.TestPost(Host+remove_url, adminIdAsString, userToken, &removeUserRequest, &response)
	fmt.Println(err, response, statusCode)
	return response, statusCode
}

func AddMemberAdminAccess(userId int64, userToken string, groupId int64, requestData proto.Message) int {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	add_member_admin_access := fmt.Sprintf(Add_Memeber_Admin_Access, groupId)
	fmt.Println(add_member_admin_access)

	response := p_base.OkResponse{}
	statusCode, err := util.TestPost(Host+add_member_admin_access, userIdAsString, userToken, requestData, &response)
	fmt.Println(err, response)
	return statusCode
}

func RemoveMemberAdminAccess(userId int64, userToken string, groupId int64, requestData proto.Message) int {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	remove_member_admin_access := fmt.Sprintf(Remove_Memmber_Admin_Access, groupId)
	fmt.Println(remove_member_admin_access)

	response := p_base.OkResponse{}
	statusCode, err := util.TestPost(Host+remove_member_admin_access, userIdAsString, userToken, requestData, &response)
	fmt.Println("Remove Member response is: ", err, response)
	return statusCode
}

func ExitGroup(userId int64, userToken string, groupId int64) int {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	exit_group_url := fmt.Sprintf(Exit_Group, groupId)
	fmt.Println(exit_group_url)

	response := p_base.OkResponse{}
	statusCode, err := util.TestPostWithoutInput(Host+exit_group_url, userIdAsString, userToken, &response)
	fmt.Println("Exit Group response is: ", err, response)
	return statusCode
}

func StartModeratorMode(userId int64, userToken string, groupId int64, requestData proto.Message) int {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	start_moderator_mode := fmt.Sprintf(Start_Moderator_Mode, groupId)
	fmt.Println(start_moderator_mode)

	response := p_base.OkResponse{}
	statusCode, err := util.TestPost(Host+start_moderator_mode, userIdAsString, userToken, requestData, &response)
	fmt.Println("Start moderator mode response is: ", err, response)
	return statusCode
}

func UpdateModeratorMode(userId int64, userToken string, groupId int64, requestData proto.Message) int {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	update_moderator_mode := fmt.Sprintf(Update_Moderator_Mode, groupId)
	fmt.Println(update_moderator_mode)

	response := p_base.OkResponse{}
	statusCode, err := util.TestPost(Host+update_moderator_mode, userIdAsString, userToken, requestData, &response)
	fmt.Println("Update moderator mode response is: ", err, response)
	return statusCode
}

func StopModeratorMode(userId int64, userToken string, groupId int64) int {
	userIdAsString := strconv.FormatInt(int64(userId), 10)
	stop_moderator_mode := fmt.Sprintf(Stop_Moderator_Mode, groupId)
	fmt.Println(stop_moderator_mode)

	response := p_base.OkResponse{}
	statusCode, err := util.TestPostWithoutInput(Host+stop_moderator_mode, userIdAsString, userToken, &response)
	fmt.Println("Stop moderator mode response is: ", err, response)
	return statusCode
}
