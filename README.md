# [orm](http://github.com/simplejia/orm) (配合sql.Rows使用的超简单数据到对象映射功能函数)
## 实现初衷
* database/sql包，Db.Query返回的sql.Rows，通过Rows.Scan方式示例代码如下：
```
rows, err := db.Query("SELECT ...")
defer rows.Close()
for rows.Next() {
    var id int
    var name string
    err = rows.Scan(&id, &name)
}
err = rows.Err()
...
```
但实际项目场景里，我们更想这样：
```
rows, err := db.Query("SELECT ...")
defer rows.Close()
var d []*stru
err = Rows2Strus(rows, &d)
```
这就是一种简单的对象映射，通过转为对象的方式，我们的代码更方便处理了

## 功能
* 一共提供四种场景的使用方法：

> Rows2Strus, sql.Rows转为struct slice

> Rows2Stru, sql.Rows转为struct，等同db.QueryRow

> Rows2Cnts, sql.Rows转为int slice

> Rows2Cnt, sql.Rows转为int，用于select count(1)操作

* 支持tag: orm，如下：
```
type Demo struct {
    Id int
    DemoName string `orm:"demo_name"` // 映射成demo_name字段
}
```
* 支持匿名成员，如下：
```
type C struct {
    Id int
}
type P struct {
    C  // 映射成id字段
    Name string
}
```

## demo
[orm_test.go](http://github.com/simplejia/orm/tree/master/orm_test.go)
