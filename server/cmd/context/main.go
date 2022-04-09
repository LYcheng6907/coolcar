package main

import (
	"context"
	"fmt"
	"time"
)

type paramKey struct {
}

// context 示例代码
func main() {
	c := context.WithValue(context.Background(), paramKey{}, "abc")
	c, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	// 总任务
	go mainTask(c)
	// 防止主线程退出
	// time.Sleep(time.Hour)

	// 从键盘输入，主动cancel(),并在mainTask函数加上go，并行执行
	var cmd string
	for {
		fmt.Scan(&cmd)
		if cmd == "c" {
			cancel()
		}
	}
}

func mainTask(c context.Context) {
	fmt.Printf("main task started with param %q\n", c.Value(paramKey{}))
	// 总任务的步骤
	// smallTask(c, "task1", 4*time.Second)
	// smallTask(c, "task2", 2*time.Second)

	// 限制步骤的执行时间(但还是收到context的时间限制) c1 为子任务，受到主任务的时间限制
	// c1, cancel1 := context.WithTimeout(c, 2*time.Second)
	// defer cancel1()
	// smallTask(c1, "task1", 4*time.Second)
	// smallTask(c, "task2", 2*time.Second)

	// 总任务的后台任务(独立的环境)，子任务是传递context
	go func() {
		c1, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		smallTask(c1, "task1", 9*time.Second)
	}()

	smallTask(c, "task2", 8*time.Second)

}

// main task started with param "abc"
// task1 started with param "abc"
// task2 started with param %!q(<nil>)
func smallTask(c context.Context, name string, d time.Duration) {
	fmt.Printf("%s started with param %q\n", name, c.Value(paramKey{}))
	select {
	case <-time.After(d):
		fmt.Printf("%s done\n", name)
	case <-c.Done():
		fmt.Printf("%s cancelled\n", name)
	}

}
