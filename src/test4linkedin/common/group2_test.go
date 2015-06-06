package test4linkedin

import (
	"github.com/golang/protobuf/proto"
	p_group "linkedin/proto/group"
	"strconv"
	p_common "test4linkedin/common"
	"testing"
	"time"
)

func testAddAndRemoveAdmin(t *testing.T) {
	users := p_common.InitTestUserInformation(2)
	socketInfor := p_common.GetSocketInformation(users[0].UserId, users[0].UserToken)
	socketServer := socketInfor.GetSockerServerHost()
	socketPort := strconv.FormatInt(int64(socketInfor.GetSockerServerPort()), 10)

	socketMessages := &p_common.ConnImplHandler{
		Index:     -1,
		Uid:       users[0].UserId,
		UserToken: users[0].UserToken,
		Msgid:     "",
		EndChan:   make(chan int, 1),
	}

	socketMessages.InitSocket(socketServer, socketPort)
	socketMessages.LoginSocket()
	go socketMessages.GetSocketMessage()

	// define group test data
	tags := []int64{0, 1}
	coordinates := []float64{0.1, 0.2}
	groupRequest := p_group.GroupInfoRequest{
		Name:         proto.String("group1"),
		Desc:         proto.String("desc1"),
		Industry:     proto.Int32(123),
		Picture:      proto.String("pic1"),
		Tags:         tags,
		Coordinate:   coordinates,
		LocationName: proto.String("Beijing"),
	}

	// create a group
	groupResponse := p_common.CreateGroup(users[0].UserId, users[0].UserToken, &groupRequest)
	if 1 == groupResponse.GetStatus() {
		t.Error("create group fails", groupResponse.GetStatus())
	}

	groupApply := p_group.ApplyGroupRequest{
		Message: proto.String("Use for test 1"),
	}

	p_common.ApplyToJoinAGroup(users[1].UserId, users[1].UserToken, groupResponse.GetGroupId(), &groupApply)

	userIds := []int64{users[1].UserId}
	approve := p_group.ApproveUserRequest{
		UserId: userIds,
	}

	approveResp := p_common.ApprovalJoinGroupRequest(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId(), &approve)
	if 0 != approveResp.GetStatus() {
		t.Error("It should approval the join invitation")
	}

	grantAdminRequest := p_group.AddAdminRequest{
		UserId: proto.Int64(users[1].UserId),
	}
	grantAdminResp := p_common.AddMemberAdminAccess(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId(), &grantAdminRequest)
	t.Log(grantAdminResp)

	groupDetails1 := p_common.GetGroupDetails(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	if p_common.GroupOwner != groupDetails1.GetRole() {
		t.Error("It should have admin access!", groupDetails1.GetRole())
	}

	groupDetails2 := p_common.GetGroupDetails(users[1].UserId, users[1].UserToken, groupResponse.GetGroupId())
	if p_common.GroupAdmin != groupDetails2.GetRole() {
		t.Error("It should have admin access!", groupDetails2.GetRole())
	}

	removeAdminRequet1 := p_group.RemoveUserRequest{
		UserId: proto.Int64(users[0].UserId),
	}

	removeAdminResp1 := p_common.RemoveMemberAdminAccess(users[1].UserId, users[1].UserToken, groupResponse.GetGroupId(), &removeAdminRequet1)
	t.Log(removeAdminResp1)

	removeAdminRequet2 := p_group.RemoveUserRequest{
		UserId: proto.Int64(users[1].UserId),
	}

	removeAdminResp2 := p_common.RemoveMemberAdminAccess(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId(), &removeAdminRequet2)
	t.Log(removeAdminResp2)

	groupDetails3 := p_common.GetGroupDetails(users[1].UserId, users[1].UserToken, groupResponse.GetGroupId())
	if p_common.GroupMember != groupDetails3.GetRole() {
		t.Error("It should have member access!", groupDetails3.GetRole())
	}

	response, _ := p_common.RemoveGroup(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	t.Log(response)

	time.Sleep(5 * time.Second)
	socketMessages.CloseSocket()

}

func testExitGroup(t *testing.T) {
	users := p_common.InitTestUserInformation(2)
	socketInfor := p_common.GetSocketInformation(users[0].UserId, users[0].UserToken)
	socketServer := socketInfor.GetSockerServerHost()
	socketPort := strconv.FormatInt(int64(socketInfor.GetSockerServerPort()), 10)

	socketMessages := &p_common.ConnImplHandler{
		Index:     -1,
		Uid:       users[0].UserId,
		UserToken: users[0].UserToken,
		Msgid:     "",
		EndChan:   make(chan int, 1),
	}

	socketMessages.InitSocket(socketServer, socketPort)
	socketMessages.LoginSocket()
	go socketMessages.GetSocketMessage()

	// define group test data
	tags := []int64{0, 1}
	coordinates := []float64{0.1, 0.2}
	groupRequest := p_group.GroupInfoRequest{
		Name:         proto.String("group1"),
		Desc:         proto.String("desc1"),
		Industry:     proto.Int32(123),
		Picture:      proto.String("pic1"),
		Tags:         tags,
		Coordinate:   coordinates,
		LocationName: proto.String("Beijing"),
	}

	// create a group
	groupResponse := p_common.CreateGroup(users[0].UserId, users[0].UserToken, &groupRequest)
	if 1 == groupResponse.GetStatus() {
		t.Error("create group fails", groupResponse.GetStatus())
	}

	groupApply := p_group.ApplyGroupRequest{
		Message: proto.String("Use for test 1"),
	}

	p_common.ApplyToJoinAGroup(users[1].UserId, users[1].UserToken, groupResponse.GetGroupId(), &groupApply)

	userIds := []int64{users[1].UserId}
	approve := p_group.ApproveUserRequest{
		UserId: userIds,
	}

	approveResp := p_common.ApprovalJoinGroupRequest(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId(), &approve)
	if 0 != approveResp.GetStatus() {
		t.Error("It should approval the join invitation")
	}

	stausCode := p_common.ExitGroup(users[1].UserId, users[1].UserToken, groupResponse.GetGroupId())
	if 200 != stausCode {
		t.Error("it should return 200 code to identify the action successfully!", stausCode)
	}

	groupDetails3 := p_common.GetGroupDetails(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	if 1 != len(groupDetails3.GetUserList()) {
		t.Error("It should only have 1 memeber", len(groupDetails3.GetUserList()), groupDetails3.GetRole())
	}

	response, _ := p_common.RemoveGroup(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	t.Log(response)

	time.Sleep(5 * time.Second)
	socketMessages.CloseSocket()

}

func testModeratorMode(t *testing.T) {
	users := p_common.InitTestUserInformation(3)
	socketInfor := p_common.GetSocketInformation(users[0].UserId, users[0].UserToken)
	socketServer := socketInfor.GetSockerServerHost()
	socketPort := strconv.FormatInt(int64(socketInfor.GetSockerServerPort()), 10)

	socketMessages := &p_common.ConnImplHandler{
		Index:     -1,
		Uid:       users[0].UserId,
		UserToken: users[0].UserToken,
		Msgid:     "",
		EndChan:   make(chan int, 1),
	}

	socketMessages.InitSocket(socketServer, socketPort)
	socketMessages.LoginSocket()
	go socketMessages.GetSocketMessage()

	// define group test data
	tags := []int64{0, 1}
	coordinates := []float64{0.1, 0.2}
	groupRequest := p_group.GroupInfoRequest{
		Name:         proto.String("group1"),
		Desc:         proto.String("desc1"),
		Industry:     proto.Int32(123),
		Picture:      proto.String("pic1"),
		Tags:         tags,
		Coordinate:   coordinates,
		LocationName: proto.String("Beijing"),
	}

	// create a group
	groupResponse := p_common.CreateGroup(users[0].UserId, users[0].UserToken, &groupRequest)
	if 1 == groupResponse.GetStatus() {
		t.Error("create group fails", groupResponse.GetStatus())
	}

	groupApply := p_group.ApplyGroupRequest{
		Message: proto.String("Use for test 1"),
	}

	// ask to join group by users[1]
	p_common.ApplyToJoinAGroup(users[1].UserId, users[1].UserToken, groupResponse.GetGroupId(), &groupApply)

	// ask to join group by users[2]
	p_common.ApplyToJoinAGroup(users[2].UserId, users[2].UserToken, groupResponse.GetGroupId(), &groupApply)

	userIds := []int64{users[1].UserId, users[2].UserId}
	approve := p_group.ApproveUserRequest{
		UserId: userIds,
	}

	approveResp := p_common.ApprovalJoinGroupRequest(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId(), &approve)
	if 0 != approveResp.GetStatus() {
		t.Error("It should approval the join invitation")
	}

	groupDetails1 := p_common.GetGroupDetails(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	if 3 != len(groupDetails1.GetUserList()) {
		t.Error("It should only have 1 memeber", len(groupDetails1.GetUserList()))
	}

	// moderator test 1 - create a new moderator list
	moderatorList := []int64{users[0].UserId, users[1].UserId}
	groupModeratorRequest := p_group.StartModeratorModeRequest{
		AddUseridList: moderatorList,
	}

	statusCode1 := p_common.StartModeratorMode(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId(), &groupModeratorRequest)
	if 200 != statusCode1 {
		t.Error("Set the moderator error", statusCode1)
	}

	groupDetails2 := p_common.GetGroupDetails(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	moderatorResp1 := groupDetails2.GetModeratorIdList()
	bContain1 := 0
	for value := range moderatorList {
		for actual := range moderatorResp1 {
			if actual == value {
				bContain1++
				break
			}
		}
	}

	if 2 != bContain1 {
		t.Error("the moderatorResp has error user id")
	}

	// moderator test 2 - update moderator list
	newM := []int64{users[2].UserId}
	removeM := []int64{users[0].UserId}
	updateMRequest := p_group.UpdateModeratorModeRequest{
		AddUseridList:    newM,
		RemoveUseridList: removeM,
	}

	statusCode2 := p_common.UpdateModeratorMode(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId(), &updateMRequest)
	if 200 != statusCode2 {
		t.Error("Set the moderator error", statusCode2)
	}

	groupDetails3 := p_common.GetGroupDetails(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	moderatorResp2 := groupDetails3.GetModeratorIdList()
	bContain2 := 0
	for data := range moderatorResp2 {
		if newM[0] == int64(data) {
			bContain2++
		}
	}

	for data := range moderatorResp2 {
		if users[1].UserId == int64(data) {
			bContain2++
		}
	}

	if 2 != bContain1 {
		t.Error("the moderatorResp has error user id")
	}

	// moderator test 3 - invalid user id
	invalidM := []int64{1234512311}
	invalidMRequest := p_group.StartModeratorModeRequest{
		AddUseridList: invalidM,
	}

	statusCode3 := p_common.StartModeratorMode(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId(), &invalidMRequest)
	t.Log(statusCode3)
	if 200 == statusCode3 {
		t.Error("It should return error for a invalid user id out", statusCode3)
	}

	groupDetails4 := p_common.GetGroupDetails(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	moderatorResp3 := groupDetails4.GetModeratorIdList()
	bContain3 := 0
	for data := range moderatorResp3 {
		if newM[0] == int64(data) {
			bContain3++
		}
	}

	for data := range moderatorResp3 {
		if users[1].UserId == int64(data) {
			bContain3++
		}
	}

	if 2 != bContain3 {
		t.Error("the moderatorResp has error user id", bContain3)
	}

	// moderator test 4
	statusCode4 := p_common.StopModeratorMode(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	t.Log(statusCode4)
	if 200 == statusCode4 {
		t.Error("Set the moderator error", statusCode4)
	}

	groupDetails5 := p_common.GetGroupDetails(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	if 0 != len(groupDetails5.GetModeratorIdList()) {
		t.Error("It should have a empty moderator list!", len(groupDetails5.GetModeratorIdList()))
	}

	response, _ := p_common.RemoveGroup(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	t.Log(response)

	time.Sleep(5 * time.Second)
	socketMessages.CloseSocket()
}

func TestMultiChat(t *testing.T) {
	users := p_common.InitTestUserInformation(2)
	socketInfor := p_common.GetSocketInformation(users[0].UserId, users[0].UserToken)
	socketServer := socketInfor.GetSockerServerHost()
	socketPort := strconv.FormatInt(int64(socketInfor.GetSockerServerPort()), 10)

	socketMessages := &p_common.ConnImplHandler{
		Index:     -1,
		Uid:       users[0].UserId,
		UserToken: users[0].UserToken,
		Msgid:     "",
		EndChan:   make(chan int, 1),
	}

	socketMessages.InitSocket(socketServer, socketPort)
	socketMessages.LoginSocket()
	go socketMessages.GetSocketMessage()

	userList := []int64{users[0].UserId, users[1].UserId}
	multiChatRequest := p_group.MultiChatInfoRequest{
		UserIdList: userList,
	}

	// create a multichat
	multiChatResp := p_common.CreateMultiChat(users[0].UserId, users[0].UserToken, &multiChatRequest)
	if 1 == multiChatResp.GetStatus() {
		t.Error("create group fails", multiChatResp.GetStatus())
	}

	multiChatDetails1 := p_common.GetGroupDetails(users[0].UserId, users[0].UserToken, multiChatResp.GetGroupId())
	if p_common.GroupOwner != multiChatDetails1.GetRole() {
		t.Error("It should have admin access!", multiChatDetails1.GetRole())
	}

	if 1 != len(multiChatDetails1.GetUserList()) {
		t.Error("The number of users in the multichat is worng!", len(multiChatDetails1.GetUserList()))
	}

	multiChatDetails2 := p_common.GetGroupDetails(users[1].UserId, users[1].UserToken, multiChatResp.GetGroupId())
	if p_common.GroupMember != multiChatDetails2.GetRole() {
		t.Error("It should have no  mutlichat access!", multiChatDetails2.GetRole())
	}

	response, _ := p_common.RemoveGroup(users[0].UserId, users[0].UserToken, multiChatResp.GetGroupId())
	t.Log(response)

	time.Sleep(5 * time.Second)
	socketMessages.CloseSocket()
}
