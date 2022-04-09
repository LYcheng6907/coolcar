package mgutil

import (
	"coolcar/shared/mongo/objid"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mongo document id
const (
	IDFieldName        = "_id"
	UpdatedAtFieldName = "updatedat"
)

type IDField struct {
	ID primitive.ObjectID `bson:"_id"`
}

type UpdatedAtField struct {
	UpdatedAt int64 `bson:"updatedat"` // 插入数据时看到不用下划线原因
}

// 固定下来，方便测试进行替换,go语言中常见做法,例如jwt.TimeFunc
var NewObjID = primitive.NewObjectID

//
func NewObjIDWithValue(id fmt.Stringer) {
	NewObjID = func() primitive.ObjectID {
		return objid.MustFromID(id)
	}
}

var UpdatedAt = func() int64 {
	return time.Now().UnixNano()
}

// 给其他服务共享MongoDB的Set方法
func Set(v interface{}) bson.M {
	return bson.M{
		"$set": v,
	}
}

func SetOnInsert(v interface{}) bson.M {
	return bson.M{
		"$setOnInsert": v,
	}
}
