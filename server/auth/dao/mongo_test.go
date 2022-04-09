package dao

import (
	"context"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	mongotesting "coolcar/shared/mongo/testing"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestResolveAccountID(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}

	m := NewMongo(mc.Database("coolcar"))
	m.col.InsertMany(c, []interface{}{
		bson.M{
			mgutil.IDFieldName: objid.MustFromID(id.AccountID("6248270ed529270184bd94a9")),
			openIDField:        "openid_1",
		},
		bson.M{
			mgutil.IDFieldName: objid.MustFromID(id.AccountID("6248270ed529270184bd9470")),
			openIDField:        "openid_2",
		},
	})
	// mgutil.NewObjID = func() primitive.ObjectID {
	// 	return objid.MustFromID(id.AccountID("6248270ed529270184bd9471"))
	// }
	mgutil.NewObjIDWithValue(id.AccountID("6248270ed529270184bd9471"))

	// 表格驱动测试
	cases := []struct {
		name   string
		openID string
		want   string
	}{
		{
			name:   "existing_user",
			openID: "openid_1",
			want:   "6248270ed529270184bd94a9",
		},
		{
			name:   "another_existing_user",
			openID: "openid_2",
			want:   "6248270ed529270184bd9470",
		},
		{
			name:   "new_user",
			openID: "openid_3",
			want:   "6248270ed529270184bd9471",
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			id, err := m.ResolveAccountID(context.Background(), cc.openID)
			if err != nil {
				t.Errorf("failed resolve account id for %q: %v", cc.openID, err)
			}
			if id.String() != cc.want {
				t.Errorf("resolve account id:want: %q,got: %q", cc.want, id)
			}
		})
	}

	// id, err := m.ResolveAccountID(c, "123")
	// if err != nil {
	// 	t.Errorf("failed resolve account id for 123: %v", err)
	// } else {
	// 	want := "6248270ed529270184bd94a9"
	// 	if id != want {
	// 		t.Errorf("resolve account id:want: %q,got: %q", want, id)
	// 	}
	// }
}

// 保证是新环境下的测试
func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoDocker(m))
}
