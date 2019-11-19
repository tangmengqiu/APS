package v1

import (
	src "APS/src"
	vm "APS/src/api/vm"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetUsers get users
// @Summary  get all user
// @Tags 用户
// @Description get all user s
// @ID  get users
// @Accept  json
// @Produce  json
// @Success 200 {string} string "success"
// @Failure 404 {string} string "failure"
// @Router /user [get]
func GetUsers(c *gin.Context) {
	//auth
	users := src.GetUsers()
	vm.MakeSuccess(c, 200, users)
	return
}

// AddUser add user
// @Summary  add user
// @Tags 用户
// @Description add user s
// @ID  add user
// @Accept  json
// @Produce  json
// @Param body body vm.ReqUser true "user register"
// @Success 200 {string} string "success"
// @Failure 404 {string} string "failure"
// @Router /user/add [post]
func AddUser(c *gin.Context) {
	var u vm.ReqUser
	if err := c.ShouldBindJSON(&u); err != nil {
		logrus.Info(err.Error())
		vm.MakeFail(c, 400, err.Error())
		return
	}
	if err := src.AddUser(u); err != nil {
		logrus.Info(err.Error())
		vm.MakeFail(c, 500, err.Error())
		return
	}
	vm.MakeSuccess(c, 200, "添加用户成功")
	return
}

// DeleteUser delete users
// @Summary  delete user
// @Tags 用户
// @Description delete  user 
// @ID  delete user
// @Accept  json
// @Produce  json
// @Param user_name path string true "user_name"
// @Success 200 {string} string "success"
// @Failure 404 {string} string "failure"
// @Router /user/{user_name} [delete]
func DeleteUser(c *gin.Context) {
	//auth
	userName := c.Params.ByName("user_name")
	if err := src.DeleteUser(userName); err != nil {
		logrus.Info(err.Error())
		vm.MakeFail(c, 500, err.Error())
		return
	}
	vm.MakeSuccess(c, 200, "删除: "+userName+" 成功")
	return
}
