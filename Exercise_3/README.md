# 一个简单实现了闹钟功能的clock包

:clap:模仿cron包写的一个定时器包

:watermelon:可以添加函数来使用，目前只有添加函数，添加配置（目前只支持是否开启秒级计时器）

## :movie_camera: 使用示例

```go
//cron := clock.New(clock.WithSecond(true)) 可以使用
cron := clock.New()//可以加上clock.withSecond(true)或者clock.withSecond(false)来开启关闭是否开启解析秒

i, err := cron.AddFunction("1/10 */2 * * * *", func() {
		fmt.Println("This is a test functionn")
})
if err != nil {
	fmt.Println(i, err)
	return
}

cron.Run()//开启定时任务
for {
	fmt.Println(time.Now())
	time.Sleep(time.Second * 1)
}
```

目前已知bug：在开启定时器之后，距离下一次等待之间会反复运行函数（计时的函数太难了！官方也只是用的goto来达到迭代时间逼近目前的时间（8年没有更新了））

- [ ] 添加取消函数功能

- [ ] 改善Next函数的实现方法
- [ ] 更改sort.sort排序方法，~~使用优先排列方法~~