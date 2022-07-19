package main

import (
	"reflect"
)

//push 入一个
func (t *TSlice) push(v interface{}) {

	//通过反射检查两种元素是否相同

	if t.size != 0 {
		v1 := reflect.ValueOf(v)
		v2 := reflect.ValueOf(t.top())
		if v1.Kind() != v2.Kind() {
			panic("value kind is not same")
		}
	}

	t.lock.Lock()

	//将元素放在数组的最后面
	t.i = append(t.i, v)

	//队中元素+1
	t.size++

	t.lock.Unlock()
}

//pop 出一个
func (t *TSlice) pop() interface{} {
	t.lock.Lock()
	defer t.lock.Unlock()

	//如果队列中元素已空
	if t.size == 0 {
		panic("queue is empty")
	}
	//将队列最前面的元素放出去
	v := t.i[0]
	newQueue := make([]interface{}, t.size-1, t.size-1)
	for i := 1; i < t.size; i++ {
		newQueue[i-1] = t.i[i]
	}
	t.i = newQueue
	//更新队列长度
	t.size--
	return v
}

//top 返回队头元素
func (t *TSlice) top() interface{} {
	//加锁 并发时需要
	t.lock.Lock()
	defer t.lock.Unlock()

	//队列元素已空
	if t.size == 0 {
		panic("queue is empty")
	}

	return t.i[0]
}

//返回长度
func (t *TSlice) length() int {
	return t.size
}
