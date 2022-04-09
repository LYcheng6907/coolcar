package id

// 无论是关系型数据库还是非关系型数据库都可以使用

type AccountID string // 设置强类型

func (a AccountID) String() string {
	return string(a)
}

type TripID string

func (t TripID) String() string {
	return string(t)
}
