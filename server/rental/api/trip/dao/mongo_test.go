package dao

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	mongotesting "coolcar/shared/mongo/testing"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestCreateTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)

	}
	db := mc.Database("coolcar")
	err = mongotesting.SetupIndexes(c, db)
	if err != nil {
		t.Fatalf("cannot setip indexes: %v", err)
	}
	m := NewMongo(db)

	// 想要account返回正确的trip，处在In_PROGRESS状态下只能创建一个trip，不同用户可以分别创建
	cases := []struct {
		name       string
		tripID     string
		accountID  string
		tripStatus rentalpb.TripStatus
		wantErr    bool
	}{
		{
			name:       "finished",
			tripID:     "624e8c975e83600ca095d325",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_FINISHED,
		},
		{
			name:       "another_finished",
			tripID:     "624e8c975e83600ca095d326",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_FINISHED,
		},
		{
			name:       "in_progress",
			tripID:     "624e8c975e83600ca095d327",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
		},
		{
			name:       "another_in_progress",
			tripID:     "624e8c975e83600ca095d328",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
			wantErr:    true,
		},
		{
			name:       "in_progress_by_another_account",
			tripID:     "624e8c975e83600ca095d329",
			accountID:  "account2",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
		},
	}

	for _, cc := range cases {
		// mgutil.NewObjID = func() primitive.ObjectID {
		// 	return objid.MustFromID(id.TripID(cc.tripID))
		// }
		mgutil.NewObjIDWithValue(id.TripID(cc.tripID))

		tr, err := m.CreateTrip(c, &rentalpb.Trip{
			AccountId: cc.accountID,
			Status:    cc.tripStatus,
		})

		if cc.wantErr {
			if err == nil {
				t.Errorf("%s error expected; got none", cc.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s error creating trip: %v", cc.name, err)
			continue
		}

		if tr.ID.Hex() != cc.tripID {
			t.Errorf("%s incorrect trip id; want: %q; got: %q", cc.name, cc.tripID, tr.ID.Hex())
		}
	}
}

func TestGetTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}

	m := NewMongo(mc.Database("coolcar"))
	acct := id.AccountID("account1")
	mgutil.NewObjID = primitive.NewObjectID
	tr, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: acct.String(),
		CarId:     "car1",
		Start: &rentalpb.LocationStatus{
			PoiName: "startpoint",
			Location: &rentalpb.Location{
				Latitude:  30,
				Longitude: 120,
			},
		},
		End: &rentalpb.LocationStatus{
			PoiName:  "endpoint",
			FeeCent:  10000,
			KmDriven: 35,
			Location: &rentalpb.Location{
				Latitude:  35,
				Longitude: 115,
			},
		},
		Status: rentalpb.TripStatus_FINISHED,
	})
	if err != nil {
		t.Fatalf("cannot create trip: %v", err)
	}
	// t.Errorf("%+v", tr)

	got, err := m.GetTrip(c, objid.ToTripID(tr.ID), acct)
	if err != nil {
		t.Errorf("cannot get trip: %v", err)
	}

	if diff := cmp.Diff(tr, got, protocmp.Transform()); diff != "" {
		t.Errorf("result differs;-want +got: %s", diff)
	}
	// t.Errorf("got trip: %+v", got)
}

// 是否只能查找出accountID，是否能根据status筛选出想要的数据
func TestGetTrips(t *testing.T) {
	// 建立客户端
	// c := context.Background()
	// mc, err := mongotesting.NewClient(c)
	// if err != nil {
	// 	t.Fatalf("cannot connect mongodb: %v", err)
	// }

	// m := NewMongo(mc.Database("coolcar"))
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}

	m := NewMongo(mc.Database("coolcar"))
	rows := []struct {
		id        string
		accountid string
		status    rentalpb.TripStatus
	}{
		{
			id:        "624e8c975e83610ca095d325",
			accountid: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "624e8c975e83610ca095d326",
			accountid: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "624e8c975e83610ca095d327",
			accountid: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "624e8c975e83610ca095d328",
			accountid: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_IN_PROGRESS,
		},
		{
			id:        "624e8c975e83610ca095d329",
			accountid: "account_id_for_get_trips_1",
			status:    rentalpb.TripStatus_IN_PROGRESS,
		},
	}

	for _, r := range rows {
		mgutil.NewObjIDWithValue(id.TripID(r.id))
		_, err := m.CreateTrip(c, &rentalpb.Trip{
			AccountId: r.accountid,
			Status:    r.status,
		})
		if err != nil {
			t.Fatalf("connot create rows: %v", err)
		}
	}

	// 设计一些请求看是否返回想要的结果
	cases := []struct {
		name       string
		accountID  string
		status     rentalpb.TripStatus
		wantCount  int
		wantOnlyID string // 根据status和accountID查找
	}{
		{
			name:      "get_all",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_TS_NOT_SPECIFIED,
			wantCount: 4,
		},
		{
			name:       "get_in_progress",
			accountID:  "account_id_for_get_trips",
			status:     rentalpb.TripStatus_IN_PROGRESS,
			wantCount:  1,
			wantOnlyID: "624e8c975e83610ca095d328",
		},
	}

	for _, cc := range cases {
		res, err := m.GetTrips(context.Background(), id.AccountID(cc.accountID), cc.status)
		if err != nil {
			t.Errorf("cannot get trips: %v", err)
		}
		if cc.wantCount != len(res) {
			t.Errorf("incorrect result count; want: %d, got: %d", cc.wantCount, len(res))
		}

		if cc.wantOnlyID != "" && len(res) > 0 {

			if cc.wantOnlyID != res[0].ID.Hex() {
				t.Errorf("only_id incorrect; want: %q, got: %q", cc.wantOnlyID, res[0].ID.Hex())
			}
		}

	}

}

// 测试同时更新的情况
func TestUpdateTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}

	m := NewMongo(mc.Database("coolcar"))

	tid := id.TripID("625e8c975e83610ca095d325")
	aid := id.AccountID("account_for_update")

	var now int64 = 10000
	mgutil.NewObjIDWithValue(tid)
	mgutil.UpdatedAt = func() int64 {
		return now
	}

	tr, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: aid.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStatus{
			PoiName: "start_poi",
		},
	})

	if err != nil {
		t.Fatalf("cannot create trip: %v", err)
	}

	if tr.UpdatedAt != 10000 {
		t.Fatalf("wrong updatedat; want: 10000,got: %d", tr.UpdatedAt)
	}

	update := &rentalpb.Trip{
		AccountId: aid.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStatus{
			PoiName: "start_poi_updated",
		},
	}

	cases := []struct {
		name          string
		now           int64
		withUpdatedAt int64
		wantErr       bool
	}{
		{
			name:          "normal_update",
			now:           20000,
			withUpdatedAt: 10000, // 建好了就拿出来
		},
		{
			name:          "update_with_old_timestamp", // 带着老的时间戳
			now:           30000,
			withUpdatedAt: 10000, // 此时应携带20000才能update
			wantErr:       true,
		},
		{
			name:          "update_with_refetch",
			now:           40000,
			withUpdatedAt: 20000,
		},
	}

	for _, cc := range cases {
		// 保证测试顺序，故不用t.run()
		now = cc.now
		err := m.UpdateTrip(c, tid, aid, cc.withUpdatedAt, update)
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s: want error;got none", cc.name)
			} else {
				continue
			}
		} else {
			if err != nil {
				t.Errorf("%s: cannot update: %v", cc.name, err)
			}
		}

		updatedTrip, err := m.GetTrip(c, tid, aid)
		if err != nil {
			t.Errorf("%s: cannot get trip after update: %v", cc.name, err)
		}

		if cc.now != updatedTrip.UpdatedAt {
			t.Errorf("%s: incorrect updatedat: want %d,got %d", cc.name, cc.now, updatedTrip.UpdatedAt)
		}

	}

}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoDocker(m))
}
