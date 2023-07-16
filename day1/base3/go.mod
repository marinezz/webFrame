module example

go 1.19


// 下面是导入本地包的写法
require (
	gee v0.0.0
)

replace (
	gee => ./gee
)