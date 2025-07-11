package copier

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Money struct {
	Count int
	Name  string
}

type AdminORM struct {
	AdminId   int
	AdminName string
}

type UserORM struct {
	AdminORM
	Id       int
	Username string
	Phone    []string
	Money    []Money
}

type UserVo struct {
	AdminId   int
	AdminName string
	Id        int
	Username  string
	Phone     []string
	Money     []Money
}

func initUser() UserORM {

	admin := AdminORM{AdminId: 9999, AdminName: "adminname1"}
	user := UserORM{
		Id:       12345,
		Username: "user1",
		Phone:    []string{"18826404531", "18565678343"},
		Money: []Money{
			{Name: "100元", Count: 3},
			{Name: "50元", Count: 5},
		},
		AdminORM: admin,
	}
	return user
}

func checkUser(t *testing.T, user UserORM, userVo UserVo) {

	t.Logf("user: %v", user)
	t.Logf("userVo: %v", userVo)
	if userVo.AdminId != user.AdminId {
		t.Errorf("dest AdminId %d != source AdminId %d", userVo.AdminId, user.AdminId)
	}
	if userVo.AdminName != user.AdminName {
		t.Errorf("dest AdminName %s != source AdminName %s", userVo.AdminName, user.AdminName)
	}
	if userVo.Id != user.Id {
		t.Errorf("dest Id %d != source Id %d", userVo.Id, user.Id)
	}
	if userVo.Username != user.Username {
		t.Errorf("dest Username %s != source Username %s", userVo.Username, user.Username)
	}
	if len(userVo.Phone) != len(user.Phone) {
		t.Errorf("Phone copy error: %v", userVo.Phone)
	}
	if userVo.Phone[0] != user.Phone[0] ||
		userVo.Phone[1] != user.Phone[1] {
		t.Errorf("dest Phone %v != source Phone %v", userVo.Phone, user.Phone)
	}
	if len(userVo.Money) != len(user.Money) {
		t.Errorf("Money copy error: %v", userVo.Phone)
	}
	if userVo.Money[0].Name != user.Money[0].Name ||
		userVo.Money[0].Count != user.Money[0].Count ||
		userVo.Money[1].Name != user.Money[1].Name ||
		userVo.Money[1].Count != user.Money[1].Count {
		t.Errorf("dest Money %v != source Money %v", userVo.Phone, user.Phone)
	}
}

func initUsers() []UserORM {
	users := []UserORM{
		{
			Id: 11111, Username: "user1", Phone: []string{"18826404531", "18565678343"},
			Money: []Money{
				{Name: "100元", Count: 3},
				{Name: "50元", Count: 5},
			},
			AdminORM: AdminORM{AdminId: 9998, AdminName: "adminname1"},
		},
		{
			Id: 66666, Username: "user2", Phone: []string{"18826404532", "18565678344"},
			Money: []Money{
				{Name: "100元", Count: 5},
				{Name: "20元", Count: 1},
			},
			AdminORM: AdminORM{AdminId: 9999, AdminName: "adminname2"},
		},
	}
	return users
}

func checkUsers(t *testing.T, users []UserORM, userVos []UserVo) {

	t.Logf("users: %v", users)
	t.Logf("userVos: %v", userVos)
	if len(userVos) != len(users) {
		t.Error("dst len error")
	}
	if userVos[0].AdminId != users[0].AdminId {
		t.Error("dst AdminId != source AdminId")
	}
	if userVos[0].AdminName != users[0].AdminName {
		t.Error("dst AdminName != source AdminName")
	}
	if userVos[1].Id != users[1].Id {
		t.Error("dst Id != source Id")
	}
	if userVos[1].Username != users[1].Username {
		t.Error("dst Username != source Username")
	}
	if len(userVos[1].Phone) != 2 ||
		userVos[1].Phone[0] != users[1].Phone[0] ||
		userVos[1].Phone[1] != users[1].Phone[1] {
		t.Error("dst Phone != source Phone")
	}
	if len(userVos[1].Money) != 2 ||
		userVos[1].Money[0].Name != users[1].Money[0].Name ||
		userVos[1].Money[0].Count != users[1].Money[0].Count ||
		userVos[1].Money[1].Name != users[1].Money[1].Name ||
		userVos[1].Money[1].Count != users[1].Money[1].Count {
		t.Error("dst Money != source Money")
	}
}

func TestCopyObject(t *testing.T) {

	user := initUser()
	var userVo UserVo
	err := Copy(user, &userVo)
	if err != nil {
		t.Error(err)
	}
	checkUser(t, user, userVo)
}

func TestCopyObjectNotPtr(t *testing.T) {

	user := initUser()
	var userVo UserVo
	err := Copy(user, userVo)
	if err == nil {
		t.Error("receive destination not a pointer, but no error")
	}
	t.Log(err)
}

func TestCopySlice(t *testing.T) {

	users := initUsers()
	var userVos []UserVo
	err := Copy(&users, &userVos)
	if err != nil {
		t.Error(err)
		return
	}
	checkUsers(t, users, userVos)
}

func TestCopySliceNotPtr(t *testing.T) {

	users := initUsers()
	var userVos []UserVo
	err := Copy(&users, userVos)
	if err == nil {
		t.Error("receive destination not a pointer, but no error")
	}
	t.Log(err)
}

func TestCopyWithIgnore(t *testing.T) {
	user := initUser()
	var userVo UserVo
	err := CopyWithIgnore(user, &userVo, []string{})
	if err != nil {
		t.Error(err)
	}
	checkUser(t, user, userVo)
}

func TestCopyWithIgnore_ignore(t *testing.T) {
	user := initUser()
	var userVo UserVo
	err := CopyWithIgnore(user, &userVo, []string{"Id"})
	if err != nil {
		t.Error(err)
	}

	assert.NotEqual(t, user.Id, userVo.Id)
}

func TestCopyWithIgnore_2(t *testing.T) {
	user := initUser()
	var userVo UserVo
	err := CopyWithIgnore(user, &userVo, []string{"Id", "Phone"})
	if err != nil {
		t.Error(err)
	}

	assert.NotEqual(t, user.Id, userVo.Id)
	assert.NotEqual(t, len(user.Phone), len(userVo.Phone))
}

func TestCopyWithIgnore_3(t *testing.T) {
	user := initUser()
	var userVo UserVo
	err := CopyWithIgnore(user, &userVo, []string{"Id", "Phone", "AdminId", "AdminName", "Username", "Money"})
	if err != nil {
		t.Error(err)
	}

	assert.NotEqual(t, user.Id, userVo.Id)
	assert.NotEqual(t, len(user.Phone), len(userVo.Phone))
	assert.NotEqual(t, user.AdminId, userVo.AdminId)
	assert.NotEqual(t, user.Username, userVo.Username)
	assert.NotEqual(t, len(user.Money), len(userVo.Money))
}

func TestCopyWithIgnore_notExistField(t *testing.T) {
	user := initUser()
	var userVo UserVo
	err := CopyWithIgnore(user, &userVo, []string{"notExistField"})
	if err != nil {
		t.Error(err)
	}

	checkUser(t, user, userVo)
}
