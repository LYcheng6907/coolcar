package dao

import (
	"context"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// openID类型
const openIDField = "open_id"

// 定义表
type Mongo struct {
	col *mongo.Collection
	// newObjID func() primitive.ObjectID // 固定_id
}

// 初始化表
func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col: db.Collection("account"),
		// newObjID: primitive.NewObjectID,
	}
}

// 将open_id转换为account openID，从数据库中拿到open_id
func (m *Mongo) ResolveAccountID(c context.Context, openID string) (id.AccountID, error) {

	insertedID := mgutil.NewObjID()
	res := m.col.FindOneAndUpdate(c, bson.M{
		openIDField: openID,
	}, mgutil.SetOnInsert(bson.M{
		mgutil.IDFieldName: insertedID,
		openIDField:        openID,
	}), options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After))

	if err := res.Err(); err != nil {
		return "", fmt.Errorf("cannot findOneAndUpdate : %v", err)
	}

	var row mgutil.IDField
	err := res.Decode(&row)
	if err != nil {
		return "", fmt.Errorf("cannot decode result: %v", err)
	}

	return objid.ToAccountID(row.ID), nil
}

// 逻辑保证：1、测试保证 2、联调保证
