package test4linkedin

import (
	"github.com/golang/protobuf/proto"
	p_group "linkedin/proto/group"
	"strconv"
	p_common "test4linkedin/common"
	"testing"
	"time"
)

func testCreateGroupAndValidateGroupInformation(t *testing.T) {
	// define test data
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

	// start to test
	loginResp := p_common.LoginByPhone(p_common.Phone, p_common.Passwd)
	userId := loginResp.GetUserID()
	userToken := loginResp.GetToken()

	groupResponse := p_common.CreateGroup(userId, userToken, &groupRequest)
	if 1 == groupResponse.GetStatus() {
		t.Error("create group fails", groupResponse.GetStatus())
	}

	groupDetails := p_common.GetGroupDetails(userId, userToken, groupResponse.GetGroupId())
	if "group1" != groupDetails.GetName() {
		t.Error("group name error", groupDetails.GetName())
	}

	if "desc1" != groupDetails.GetDescription() {
		t.Error("group description error.", groupDetails.GetDescription())
	}

	if "Beijing" != groupDetails.GetLocationName() {
		t.Error("group Location error.", groupDetails.GetLocationName())
	}

	if "pic1" != groupDetails.GetImageUrl() {
		t.Error("group image url error.", groupDetails.GetImageUrl())
	}

	if 1 != groupDetails.GetRole() {
		t.Error("user role for group error.", groupDetails.GetRole())
	}
}

func testCreateGroupWithoutGroupName(t *testing.T) {
	// define test data
	tags := make([]int64, 2)
	tags[0] = 0
	tags[1] = 1

	coordinates := make([]float64, 2)
	coordinates[0] = 0.1
	coordinates[1] = 0.2

	groupRequest := p_group.GroupInfoRequest{
		Name:         proto.String(""),
		Desc:         proto.String("desc1"),
		Industry:     proto.Int32(123),
		Picture:      proto.String("pic1"),
		Tags:         tags,
		Coordinate:   coordinates,
		LocationName: proto.String("Beijing"),
	}

	// start to test
	loginResp := p_common.LoginByPhone(p_common.Phone, p_common.Passwd)
	userId := loginResp.GetUserID()
	userToken := loginResp.GetToken()
	groupResponse := p_common.CreateGroup(userId, userToken, &groupRequest)
	t.Log(groupResponse.GetGroupId())
	if 1 != groupResponse.GetStatus() {
		t.Error("create group should be failed as no group name!", groupResponse.GetStatus())
	}
}

func testDeleteGroup(t *testing.T) {
	// define test data
	/*
		tags := make([]int64, 2)
		tags[0] = 0
		tags[1] = 1

		coordinates := make([]float64, 2)
		coordinates[0] = 0.1
		coordinates[1] = 0.2
	*/

	users := p_common.InitTestUserInformation(1)
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

	// start to test
	groupResponse := p_common.CreateGroup(users[0].UserId, users[0].UserToken, &groupRequest)

	response, _ := p_common.RemoveGroup(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	t.Log(response)
	//if delete group fail, it will popup a error msg. so if everything works well , it should be OK.
}

func testGetGroupList(t *testing.T) {
	users := p_common.InitTestUserInformation(4)

	response1 := p_common.GetGroupListByUser(users[0].UserId, users[0].UserToken)
	if 0 != len(response1.GetList()) {
		t.Error("it should has no group as no one has been created!")
	}

	tags := []int64{0, 1}
	coordinates := []float64{0.1, 0.2}

	groupRequest1 := p_group.GroupInfoRequest{
		Name:         proto.String("group1"),
		Desc:         proto.String("desc1"),
		Industry:     proto.Int32(123),
		Picture:      proto.String("pic1"),
		Tags:         tags,
		Coordinate:   coordinates,
		LocationName: proto.String("Beijing"),
	}
	groupInfor1 := p_common.CreateGroup(users[0].UserId, users[0].UserToken, &groupRequest1)

	response2 := p_common.GetGroupListByUser(users[0].UserId, users[0].UserToken)
	if 1 != len(response2.GetList()) {
		t.Error("it should has 1 group as 1 group has been created!")
	}

	groupRequest2 := p_group.GroupInfoRequest{
		Name:         proto.String("group2"),
		Desc:         proto.String("desc2"),
		Industry:     proto.Int32(123),
		Picture:      proto.String("pic2"),
		Tags:         tags,
		Coordinate:   coordinates,
		LocationName: proto.String("Beijing"),
	}
	groupInfor2 := p_common.CreateGroup(users[0].UserId, users[0].UserToken, &groupRequest2)

	response3 := p_common.GetGroupListByUser(users[0].UserId, users[0].UserToken)
	if 2 != len(response3.GetList()) {
		t.Error("it should has 2 groups as 2 groups has been created!")
	}

	p_common.RemoveGroup(users[0].UserId, users[0].UserToken, groupInfor1.GetGroupId())

	response4 := p_common.GetGroupListByUser(users[0].UserId, users[0].UserToken)
	if 1 != len(response4.GetList()) {
		t.Error("it should has 1 group as that group has been deleted!")
	}

	p_common.RemoveGroup(users[0].UserId, users[0].UserToken, groupInfor2.GetGroupId())

	response5 := p_common.GetGroupListByUser(users[0].UserId, users[0].UserToken)
	if 0 != len(response5.GetList()) {
		t.Error("it should has 0 group as that group has been deleted!")
	}
}

func testGroupVisibility(t *testing.T) {
	users := p_common.InitTestUserInformation(1)

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

	groupInfor := p_common.CreateGroup(users[0].UserId, users[0].UserToken, &groupRequest)

	resp1 := p_common.SetGroupVisible(users[0].UserId, users[0].UserToken, groupInfor.GetGroupId())
	t.Log(resp1)

	groupDetails1 := p_common.GetGroupDetails(users[0].UserId, users[0].UserToken, groupInfor.GetGroupId())
	t.Log(groupDetails1.GetVisible())
	if true != groupDetails1.GetVisible() {
		t.Error("the group should be visible!")
	}

	resp2 := p_common.SetGroupInvisible(users[0].UserId, users[0].UserToken, groupInfor.GetGroupId())
	t.Log(resp2)

	groupDetails2 := p_common.GetGroupDetails(users[0].UserId, users[0].UserToken, groupInfor.GetGroupId())
	t.Log(groupDetails2.GetVisible())
	if false != groupDetails2.GetVisible() {
		t.Error("the group should be visible!")
	}

	// clean the created db
	response, _ := p_common.RemoveGroup(users[0].UserId, users[0].UserToken, groupInfor.GetGroupId())
	t.Log(response)
}

func tstGroupAdminInviteUser(t *testing.T) {
	users := p_common.InitTestUserInformation(2)
	socketInfor := p_common.GetSocketInformation(users[1].UserId, users[1].UserToken)
	socketServer := socketInfor.GetSockerServerHost()
	socketPort := strconv.FormatInt(int64(socketInfor.GetSockerServerPort()), 10)

	socketMessages := &p_common.ConnImplHandler{
		Index:     -1,
		Uid:       users[1].UserId,
		UserToken: users[0].UserToken,
		Msgid:     "",
		EndChan:   make(chan int, 1),
	}

	t.Log(socketServer, socketPort, users[1].UserId)
	socketMessages.InitSocket(socketServer, socketPort)
	socketMessages.LoginSocket()

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

	invitedUser := p_group.InviteUserRequest{
		UserId: []int64{users[1].UserId},
	}

	inviteResp := p_common.InviteUserToGroup(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId(), &invitedUser)
	t.Log(inviteResp.GetStatus())
	t.Log(inviteResp.GetGroupId())

	// get socket response by asyn method
	go socketMessages.GetSocketMessage()
	time.Sleep(5 * time.Second)

	// start to validate resp data
	groupAdmin := p_common.GetGroupListByUser(users[0].UserId, users[0].UserToken)
	t.Log("Current group uesrs are: ", groupAdmin.GetList())
	if 2 != groupAdmin.GetList()[0].GetGroupMemberCount() {
		t.Error("It should contain 2 members in this group")
	}

	if 1 != groupAdmin.GetList()[0].GetRole() {
		t.Error("It should has admin role for this user!", groupAdmin.GetList()[0].GetRole())
	}

	groupUsers := p_common.GetGroupListByUser(users[1].UserId, users[1].UserToken)
	t.Log("Current group uesrs are: ", groupUsers.GetList())
	if 2 != groupUsers.GetList()[0].GetGroupMemberCount() {
		t.Error("It should contain 2 members in this group")
	}

	if 3 != groupUsers.GetList()[0].GetRole() {
		t.Error("It should has admin role for this user!", groupUsers.GetList()[0].GetRole())
	}

	_, statusCode1 := p_common.RemoveGroup(users[1].UserId, users[1].UserToken, groupResponse.GetGroupId())
	if 401 != statusCode1 {
		t.Error("Remove group should be fail as the user is a member")
	}

	_, statusCode2 := p_common.RemoveGroup(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	if 200 != statusCode2 {
		t.Error("Remove group should be fail as the user is a member")
	}

	group := p_common.GetGroupDetails(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	t.Log(group)

	socketMessages.CloseSocket()
}

func testRemoveUser(t *testing.T) {
	users := p_common.InitTestUserInformation(2)
	socketInfor := p_common.GetSocketInformation(users[1].UserId, users[1].UserToken)
	socketServer := socketInfor.GetSockerServerHost()
	socketPort := strconv.FormatInt(int64(socketInfor.GetSockerServerPort()), 10)

	socketMessages := &p_common.ConnImplHandler{
		Index:     -1,
		Uid:       users[1].UserId,
		UserToken: users[0].UserToken,
		Msgid:     "",
		EndChan:   make(chan int, 1),
	}

	t.Log(socketServer, socketPort, users[1].UserId)
	socketMessages.InitSocket(socketServer, socketPort)
	socketMessages.LoginSocket()

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

	invitedUser := p_group.InviteUserRequest{
		UserId: []int64{users[1].UserId},
	}

	inviteResp := p_common.InviteUserToGroup(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId(), &invitedUser)
	t.Log(inviteResp.GetStatus())
	t.Log(inviteResp.GetGroupId())

	go socketMessages.GetSocketMessage()
	time.Sleep(5 * time.Second)

	removeUesrResp, statusCode := p_common.RemoveUser(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId(), users[1].UserId)
	t.Log("Remove User resp:", removeUesrResp, statusCode)

	time.Sleep(5 * time.Second)
	groupDetails := p_common.GetGroupDetails(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	if 1 != len(groupDetails.GetUserList()) {
		t.Error("The group should only have 1 member now!", len(groupDetails.GetUserList()))
	}

	response, _ := p_common.RemoveGroup(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	t.Log(response)

	socketMessages.CloseSocket()
}

func testApplyAGroup(t *testing.T) {
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
	time.Sleep(5 * time.Second)

	userIds := []int64{users[1].UserId}
	approve := p_group.ApproveUserRequest{
		UserId: userIds,
	}

	approveResp := p_common.ApprovalJoinGroupRequest(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId(), &approve)
	if 0 != approveResp.GetStatus() {
		t.Error("It should approval the join invitation")
	}

	groupDetails := p_common.GetGroupDetails(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	if 2 != len(groupDetails.GetUserList()) {
		t.Error("The group should only have 1 member now!", len(groupDetails.GetUserList()))
	}

	response, _ := p_common.RemoveGroup(users[0].UserId, users[0].UserToken, groupResponse.GetGroupId())
	t.Log(response)

	time.Sleep(5 * time.Second)
	socketMessages.CloseSocket()
}
