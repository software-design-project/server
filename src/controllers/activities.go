package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	. "../log"
	. "../models"
	"../utils"
)

func ActivityGetOne(w http.ResponseWriter, r *http.Request) {
	// GET, 从URL中读取参数, 直接使用mux.Vars(r)
	vars := mux.Vars(r)
	activityId := vars["activityId"]

	activity := Activity{}
	err := Db["activitys"].FindId(bson.ObjectIdHex(activityId)).One(&activity)
	if err != nil {
		Log.Errorf("Get activity id: %s failed, %v", activityId, err)
		utils.FailureResponse(&w, "获取活动信息失败", "")
		return
	}

	Log.Noticef("Get activity successfully: %s", activity)
	utils.SuccessResponse(&w, "获取活动信息成功", activity)
}

func ActivityGetAll(w http.ResponseWriter, r *http.Request) {
	var activitys []Activity
	err := Db["activitys"].Find(nil).All(&activitys)
	if err != nil {
		Log.Errorf("get all activitys failed, %v", err)
		utils.FailureResponse(&w, "获取活动列表失败", "")
		return
	}
	Log.Notice("get all activity successfully")
	utils.SuccessResponse(&w, "获取活动列表成功", activitys)
}

func ActivityAddOne(w http.ResponseWriter, r *http.Request) {
	// 1. load request's body
	newActivity := Activity{}
	ok := utils.LoadRequestBody(r, "insert activity", &newActivity)
	if !ok {
		utils.FailureResponse(&w, "新建活动失败", "")
		return
	}
	// 3. set a new id
	// 由于将Activity.Id解释成_id(见activity定义), 所以activity.Id需要自己指定, 没有这一步会导致插入失败
	newActivity.Id = bson.NewObjectId()
	// 4. insert into db
	err := Db["activitys"].Insert(&newActivity)
	if err != nil {
		Log.Error("insert activity falied: insert into db failed, ", err)
		utils.FailureResponse(&w, "添加活动失败", "")
		return
	}
	// 5. success
	Log.Notice("add one activity successfully")
	utils.SuccessResponse(&w, "添加活动成功", "")
}

func ActivityUpdateOne(w http.ResponseWriter, r *http.Request) {
	// 1. 获得URL中的参数
	vars := mux.Vars(r)
	activityId := vars["activityId"]
	// 2. 从request中解析出body数据
	newActivity := Activity{}
	ok := utils.LoadRequestBody(r, "update activity", &newActivity)
	if !ok {
		utils.FailureResponse(&w, "修改活动信息失败", "")
	}
	newActivity.Id = bson.ObjectIdHex(activityId)
	// 4. 修改数据
	// convert structure to bson.M, used to update
	updateData, _ := bson.Marshal(&newActivity)
	updateActivity := bson.M{}
	_ = bson.Unmarshal(updateData, &updateActivity)
	// 此处更新时如果没有"$set",会将整行直接覆盖，而不是按需修改
	err := Db["activitys"].Update(bson.M{"_id": newActivity.Id}, bson.M{"$set": updateActivity})
	if err != nil {
		Log.Error("update activity falied: failed to update data into db, ", err)
		utils.FailureResponse(&w, "修改活动信息失败", "")
		return
	}
	// 5. 成功返回
	Log.Notice("update activity successfully")
	utils.SuccessResponse(&w, "修改活动成功", "")
}

func ActivityDeleteOne(w http.ResponseWriter, r *http.Request) {
	activityId := mux.Vars(r)["activityId"]

	err := Db["activitys"].Remove(bson.M{"_id": bson.ObjectIdHex(activityId)})
	if err != nil {
		Log.Error("delete activity from db failed: ", err)
		utils.FailureResponse(&w, "删除活动失败", "")
		return
	}

	Log.Notice("delete activity successfully")
	utils.SuccessResponse(&w, "删除活动成功", "")
}

var ActivityRoutes Routes = Routes{
	Route{"ActivityGetOne", "GET", "/activity/{activityId}", BasicAuth(ActivityGetOne)},
	Route{"ActivityGetAll", "GET", "/activity/", ActivityGetAll},
	Route{"ActivityAddOne", "POST", "/activity/", ActivityAddOne},
	Route{"ActivityUpdateOne", "PUT", "/activity/{activityId}", BasicAuth(ActivityUpdateOne)},
	Route{"ActivityDeleteOne", "DELETE", "/activity/{activityId}", ActivityDeleteOne},
}
