package server


import "sync/atomic"


var id int64


func GenerateID() int64 {


	return atomic.AddInt64(
		&id,
		1,//每次执行一次
	)

}