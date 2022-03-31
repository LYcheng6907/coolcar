package main

import (
	trippb "coolcar/proto/gen/go"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
)

func main() {
	trip := trippb.Trip{
		Start:       "abc",
		End:         "def",
		DurationSec: 3600,
		FeeCent:     10000,
		StartPos: &trippb.Location{
			Latitude:   30,
			Longtitude: 120,
		},
		EndPos: &trippb.Location{
			Latitude:   35,
			Longtitude: 115,
		},
		PathLocations: []*trippb.Location{
			{
				Latitude:   31,
				Longtitude: 119,
			},
			{
				Latitude:   32,
				Longtitude: 118,
			},
		},
		Status: trippb.TripStatus_IN_PROGRESS,
	}
	fmt.Println("初始化trip")
	fmt.Println(&trip)
	// 编码成二进制流
	b, err := proto.Marshal(&trip)

	if err != nil {
		panic(err)
	}

	fmt.Println("proto.Marshal 转换成二进制流：")
	fmt.Printf("%X\n", b)

	var trip2 trippb.Trip
	// 解码数据类型
	err = proto.Unmarshal(b, &trip2)
	if err != nil {
		panic(err)
	}

	fmt.Println("proto.Unmarshal 解析为对应数据类型：")
	fmt.Println(&trip2)

	b, err = json.Marshal(&trip2)

	if err != nil {
		panic(err)
	}

	fmt.Println("json.Marshal 转换为json：")
	fmt.Printf("%s\n", b)
	err = json.Unmarshal(b, &trip2)
	if err != nil {
		panic(err)
	}
	fmt.Println("json.Unmarshal 解析为对应数据类型：")
	fmt.Println(&trip2)

}
