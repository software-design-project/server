package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	. "../log"
	. "../models"
	"../utils"
)

func OrderGetOne(w http.ResponseWriter, r *http.Request) {
	// GET, 从URL中读取参数, 直接使用mux.Vars(r)
	vars := mux.Vars(r)
	orderId := vars["orderId"]

	order := Order{}
	err := Db["orders"].FindId(bson.ObjectIdHex(orderId)).One(&order)
	if err != nil {
		Log.Errorf("Get order id: %s failed, %v", orderId, err)
		utils.FailureResponse(&w, "获取订单信息失败", "")
		return
	}

	Log.Noticef("Get order successfully: %s", order)
	utils.SuccessResponse(&w, "获取订单信息成功", order)
}

func OrderGetAll(w http.ResponseWriter, r *http.Request) {
	var orders []Order
	err := Db["orders"].Find(nil).All(&orders)
	if err != nil {
		Log.Errorf("get all orders failed, %v", err)
		utils.FailureResponse(&w, "获取订单列表失败", "")
		return
	}
	Log.Notice("get all order successfully")
	utils.SuccessResponse(&w, "获取订单列表成功", orders)
}

func OrderGetAllByUserId(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["userId"]

	var orders []Order
	err := Db["orders"].Find(bson.M{"userId": userId}).All(&orders)
	if err != nil {
		Log.Errorf("get all orders failed, %v", err)
		utils.FailureResponse(&w, "根据userid获取订单列表失败", "")
		return
	}
	Log.Notice("get all order successfully")
	utils.SuccessResponse(&w, "根据userid获取订单列表成功", orders)
}

func OrderAddOne(w http.ResponseWriter, r *http.Request) {
	// 1. load request's body
	newOrder := Order{}
	ok := utils.LoadRequestBody(r, "insert order", &newOrder)
	if !ok {
		utils.FailureResponse(&w, "新建订单失败", "")
		return
	}
	// 3. set a new id
	// 由于将Order.Id解释成_id(见order定义), 所以order.Id需要自己指定, 没有这一步会导致插入失败
	newOrder.Id = bson.NewObjectId()
	// 4. insert into db
	err := Db["orders"].Insert(&newOrder)
	if err != nil {
		Log.Error("insert order falied: insert into db failed, ", err)
		utils.FailureResponse(&w, "添加订单失败", "")
		return
	}
	// 5. success
	Log.Notice("add one order successfully")
	utils.SuccessResponse(&w, "添加订单成功", "")
}

func OrderUpdateOne(w http.ResponseWriter, r *http.Request) {
	// 1. 获得URL中的参数
	vars := mux.Vars(r)
	orderId := vars["orderId"]
	// 2. 从request中解析出body数据
	newOrder := Order{}
	ok := utils.LoadRequestBody(r, "update order", &newOrder)
	if !ok {
		utils.FailureResponse(&w, "修改订单信息失败", "")
	}
	newOrder.Id = bson.ObjectIdHex(orderId)
	// 4. 修改数据
	// convert structure to bson.M, used to update
	updateData, _ := bson.Marshal(&newOrder)
	updateOrder := bson.M{}
	_ = bson.Unmarshal(updateData, &updateOrder)
	// 此处更新时如果没有"$set",会将整行直接覆盖，而不是按需修改
	err := Db["orders"].Update(bson.M{"_id": newOrder.Id}, bson.M{"$set": updateOrder})
	if err != nil {
		Log.Error("update order falied: failed to update data into db, ", err)
		utils.FailureResponse(&w, "修改订单信息失败", "")
		return
	}
	// 5. 成功返回
	Log.Notice("update order successfully")
	utils.SuccessResponse(&w, "修改订单成功", "")
}

func OrderDeleteOne(w http.ResponseWriter, r *http.Request) {
	orderId := mux.Vars(r)["orderId"]

	err := Db["orders"].Remove(bson.M{"_id": bson.ObjectIdHex(orderId)})
	if err != nil {
		Log.Error("delete order from db failed: ", err)
		utils.FailureResponse(&w, "删除订单失败", "")
		return
	}

	Log.Notice("delete order successfully")
	utils.SuccessResponse(&w, "删除订单成功", "")
}

var OrderRoutes Routes = Routes{
	Route{"OrderGetOne", "GET", "/order/{orderId}", OrderGetOne},
	Route{"OrderGetAll", "GET", "/order/", OrderGetAll},
	Route{"OrderGetAllByUserId", "GET", "/order/user/{userId}", OrderGetAllByUserId},
	Route{"OrderAddOne", "POST", "/order/", OrderAddOne},
	Route{"OrderUpdateOne", "PUT", "/order/{orderId}", OrderUpdateOne},
	Route{"OrderDeleteOne", "DELETE", "/order/{orderId}", OrderDeleteOne},
}
