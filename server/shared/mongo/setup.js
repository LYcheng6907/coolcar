db.account.createIndex({
    open_id: 1,
},{
    unique: true,
})

db.trip.createIndex({
    "trip.accountid": 1, // 从小到大
    "trip.status": 1,
},{
	unique: true,
	partialFilterExpression:{
		"trip.status": 1, // 只有IN_PROGRESS时候才不可建立重复trip
	}
})
