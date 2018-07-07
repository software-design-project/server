package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	. "../log"
	. "../models"
	"../utils"
)

func SellerGetOne(w http.ResponseWriter, r *http.Request) {
	// GET, 从URL中读取参数, 直接使用mux.Vars(r)
	vars := mux.Vars(r)
	sellerId := vars["sellerId"]

	seller := Seller{}
	err := Db["sellers"].FindId(bson.ObjectIdHex(sellerId)).One(&seller)
	if err != nil {
		Log.Errorf("Get seller id: %s failed, %v", sellerId, err)
		utils.FailureResponse(&w, "获取商家信息失败", "")
		return
	}

	Log.Noticef("Get seller successfully: %s", seller)
	utils.SuccessResponse(&w, "获取商家信息成功", seller)
}

func SellerGetAll(w http.ResponseWriter, r *http.Request) {
	var sellers []Seller
	err := Db["sellers"].Find(nil).All(&sellers)
	if err != nil {
		Log.Errorf("get all sellers failed, %v", err)
		utils.FailureResponse(&w, "获取商家列表失败", "")
		return
	}
	Log.Notice("get all seller successfully")
	utils.SuccessResponse(&w, "获取商家列表成功", sellers)
}

func SellerAddOne(w http.ResponseWriter, r *http.Request) {
	// 1. load request's body
	newSeller := Seller{}
	ok := utils.LoadRequestBody(r, "insert seller", &newSeller)
	if !ok {
		utils.FailureResponse(&w, "新建商家失败", "")
		return
	}
	// 3. set a new id
	// 由于将Seller.Id解释成_id(见seller定义), 所以seller.Id需要自己指定, 没有这一步会导致插入失败
	newSeller.Id = bson.NewObjectId()
	// 4. insert into db
	err := Db["sellers"].Insert(&newSeller)
	if err != nil {
		Log.Error("insert seller falied: insert into db failed, ", err)
		utils.FailureResponse(&w, "添加商家失败", "")
		return
	}
	// 5. success
	Log.Notice("add one seller successfully")
	utils.SuccessResponse(&w, "添加商家成功", "")
}

func SellerUpdateOne(w http.ResponseWriter, r *http.Request) {
	// 1. 获得URL中的参数
	vars := mux.Vars(r)
	sellerId := vars["sellerId"]
	// 2. 从request中解析出body数据
	newSeller := Seller{}
	ok := utils.LoadRequestBody(r, "update seller", &newSeller)
	if !ok {
		utils.FailureResponse(&w, "修改商家信息失败", "")
	}
	newSeller.Id = bson.ObjectIdHex(sellerId)
	// 4. 修改数据
	// convert structure to bson.M, used to update
	updateData, _ := bson.Marshal(&newSeller)
	updateSeller := bson.M{}
	_ = bson.Unmarshal(updateData, &updateSeller)
	// 此处更新时如果没有"$set",会将整行直接覆盖，而不是按需修改
	err := Db["sellers"].Update(bson.M{"_id": newSeller.Id}, bson.M{"$set": updateSeller})
	if err != nil {
		Log.Error("update seller falied: failed to update data into db, ", err)
		utils.FailureResponse(&w, "修改商家信息失败", "")
		return
	}
	// 5. 成功返回
	Log.Notice("update seller successfully")
	utils.SuccessResponse(&w, "修改商家成功", "")
}

func SellerDeleteOne(w http.ResponseWriter, r *http.Request) {
	sellerId := mux.Vars(r)["sellerId"]

	err := Db["sellers"].Remove(bson.M{"_id": bson.ObjectIdHex(sellerId)})
	if err != nil {
		Log.Error("delete seller from db failed: ", err)
		utils.FailureResponse(&w, "删除商家失败", "")
		return
	}

	Log.Notice("delete seller successfully")
	utils.SuccessResponse(&w, "删除商家成功", "")
}

var SellerRoutes Routes = Routes{
	Route{"SellerGetOne", "GET", "/seller/{sellerId}", SellerGetOne},
	Route{"SellerGetAll", "GET", "/seller/", SellerGetAll},
	Route{"SellerAddOne", "POST", "/seller/", SellerAddOne},
	Route{"SellerUpdateOne", "PUT", "/seller/{sellerId}", SellerUpdateOne},
	Route{"SellerDeleteOne", "DELETE", "/seller/{sellerId}", SellerDeleteOne},
}
