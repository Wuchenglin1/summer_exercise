# Corm

## 简介

这是参考[极客兔](https://github.com/geektutu)的[7天用Go从零实现ORM框架—GeeORM](https://geektutu.com/post/geeorm.html)以及[GORM](https://gorm.cn/)源码来实现的corm框架

注意：并没有添加筛字处理，所以理论上不会防sql注入！！！~~待添加~~

## :file_folder: 文件树

├─clause `负责构造 SQL 语句`
├─clog `打印日志文件的自定义封装包`
├─dialect `通过注册函数来屏蔽不同数据库的差别，实现对不同数据库的支持`
├─schema `将结构体的字段、值反射成不同数据库中的字段和值`
└─session `连接数据库，对数据库进行操作`

## :rocket: 实现功能

- [x] 自定义格式打印信息

- [x] Mysql原生语句实现连接、操作数据库
- [x] 对不同数据库的操作进行隔离，数据能够单独操作
- [x] 对象结构映射（通过反射获取任意Struct对象的名称和字段创建表格）
- [x] 数据表的创建和删除
- [x] 实现Create、Find、First、Update、Delete、Count函数功能
- [x] \`corm:"\<field\>"\`对单个字段的解析并应用于创建表格
- [x] \`corm:""\<field>;\<field>"`对多个字段解析
- [x] 链式调用—更删改查
- [ ] 钩子(Hook)【待实现】
- [x] 事务
- [x] 数据库的迁移(Migrate)（简陋的迁移）
- [ ] 连接池，不用手动关闭连接
- [ ] 预加载

## :computer:示例

## 实现步骤

总体思路是：

1.实现原生语句操作数据库

2.在1的基础上封装一系列函数来进行映射实现orm（Object Relational Mapping）(对象关系映射)

### 实现自定义打印日志封装

位于`clog/clog.go`下面

这里分为了三层打印，

**Debug等级** ~~写着写着给写忘了~~

I**nfo等级**

**Error等级**

可以主动设置只打印哪个等级的日志

还有一个就是SQL语句的打印

### Mysql原生语句实现连接、操作数据库

在根目录下的`statement.go`

实现了对`database/sql`包Exec、Query、QueryRow函数的封装，能够实现原生语句的加入，之后的全部函数都是建立在这三个函数之上的，只是实现方法不同。

### 对象结构映射

#### 类型映射

对于类型的映射，gorm的源码（截取了一小段，原文太长了）：

```go
	switch reflect.Indirect(fieldValue).Kind() {
	case reflect.Bool:
		field.DataType = Bool
		if field.HasDefaultValue && !skipParseDefaultValue {
			if field.DefaultValueInterface, err = strconv.ParseBool(field.DefaultValue); err != nil {
				schema.err = fmt.Errorf("failed to parse %s as default value for bool, got error: %v", field.DefaultValue, err)
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.DataType = Int
		if field.HasDefaultValue && !skipParseDefaultValue {
			if field.DefaultValueInterface, err = strconv.ParseInt(field.DefaultValue, 0, 64); err != nil {
				schema.err = fmt.Errorf("failed to parse %s as default value for int, got error: %v", field.DefaultValue, err)
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		field.DataType = Uint
		if field.HasDefaultValue && !skipParseDefaultValue {
			if field.DefaultValueInterface, err = strconv.ParseUint(field.DefaultValue, 0, 64); err != nil {
				schema.err = fmt.Errorf("failed to parse %s as default value for uint, got error: %v", field.DefaultValue, err)
			}
		}
	case reflect.Float32, reflect.Float64:
		field.DataType = Float
		if field.HasDefaultValue && !skipParseDefaultValue {
			if field.DefaultValueInterface, err = strconv.ParseFloat(field.DefaultValue, 64); err != nil {
				schema.err = fmt.Errorf("failed to parse %s as default value for float, got error: %v", field.DefaultValue, err)
			}
		}
	case reflect.String:
		field.DataType = String
		if field.HasDefaultValue && !skipParseDefaultValue {
			field.DefaultValue = strings.Trim(field.DefaultValue, "'")
			field.DefaultValue = strings.Trim(field.DefaultValue, `"`)
			field.DefaultValueInterface = field.DefaultValue
		}
	case reflect.Struct:
		if _, ok := fieldValue.Interface().(*time.Time); ok {
			field.DataType = Time
		} else if fieldValue.Type().ConvertibleTo(TimeReflectType) {
			field.DataType = Time
		} else if fieldValue.Type().ConvertibleTo(TimePtrReflectType) {
			field.DataType = Time
		}
		if field.HasDefaultValue && !skipParseDefaultValue && field.DataType == Time {
			if t, err := now.Parse(field.DefaultValue); err == nil {
				field.DefaultValueInterface = t
			}
		}
	case reflect.Array, reflect.Slice:
		if reflect.Indirect(fieldValue).Type().Elem() == ByteReflectType && field.DataType == "" {
			field.DataType = Bytes
		}
	}
```

可以看到gorm的转换过程如下：

将所有的类型反射出来，转换成规定的const，比如将int、int16...都转换成规定的Int，然后将长度记录，再在另外一个函数当中去转换类型。

corm当中是直接对值分别进行映射，然后返回对应的sql类型语句

类型映射是`schema/schema.go`下的`ParseType`函数

#### 值映射

corm的值映射依赖与类型映射

```go
func (s *Schema) ParseValue(dest any) []any {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() == reflect.Ptr {
		destValue = destValue.Elem()
	}
	var fieldValues []any
	for _, field := range s.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())//这里可以看到对于
	}
	return fieldValues
}
```



#### 创建表

在实现了原生的sql语句输入之后就可以调用该接口来封装一些函数了，先是创建表

### 链式调用

链式调用是一种简化代码的编程方式，能够使代码更简洁、易读。链式调用的原理也非常简单，某个对象调用某个方法后，将该对象的引用/指针返回，即可以继续调用该对象的其他方法。

在所有方法调用之后，都会返回一个`*DB`类型以可以继续调用其他函数

### 插入值

#### Create

```mysql
思路：
DB.Create(&MODEL)
这里等效于
1.SELECT SELECT table_name FROM information_schema.TABLES WHERE table_name = "<MODEL.Name>"
如果表不存在的话，先创表
2.CREATE TABLE <MODEL.Name>(<column1> <type1>,<column2> <type2>);
如果表存在的话，先查询两表字段名是否相同，不相同的话更新表（我这里直接创建新表）
....
再插入值
3.INSERT INTO 
<Table_name> 
(<column1>,<column2>,...) 
VALUES 
(<values1>,<values2>,...),(<values5>,<values6>)
```

##### 使用方法

```go
type A struct{...}
var a = A{...}
db, err := Open("<driverName>", "<dsn>")//获取一个db
row,err := db.Create(&a,&A{...})//这里只能是相同的结构体类型否则会报错
//返回影响的行数以及error(如果不为nil的话)
```



### 查找值

这是官方mysql给出的参数顺序

```mysql
SELECT
    [ALL | DISTINCT | DISTINCTROW ]
    [HIGH_PRIORITY]
    [STRAIGHT_JOIN]
    [SQL_SMALL_RESULT] [SQL_BIG_RESULT] [SQL_BUFFER_RESULT]
    [SQL_NO_CACHE] [SQL_CALC_FOUND_ROWS]
    select_expr [, select_expr] ...
    [into_option]
    [FROM table_references
      [PARTITION partition_list]]
    [WHERE where_condition]
    [GROUP BY {col_name | expr | position}, ... [WITH ROLLUP]]
    [HAVING where_condition]
    [WINDOW window_name AS (window_spec)
        [, window_name AS (window_spec)] ...]
    [ORDER BY {col_name | expr | position}
      [ASC | DESC], ... [WITH ROLLUP]]
    [LIMIT {[offset,] row_count | row_count OFFSET offset}]
    [into_option]
    [FOR {UPDATE | SHARE}
        [OF tbl_name [, tbl_name] ...]
        [NOWAIT | SKIP LOCKED]
      | LOCK IN SHARE MODE]
    [into_option]

into_option: {
    INTO OUTFILE 'file_name'
        [CHARACTER SET charset_name]
        export_options
  | INTO DUMPFILE 'file_name'
  | INTO var_name [, var_name] ...
}
```

#### Find

```mysql
使用Find方法就等效于
SELECT <column1,column2,...> FROM <table_name> [WHERE <Option>] [ORDER BY ...] [LIMIT ...]
```

##### 使用方法

```go
type A struct
var arr []A
db, err := Open("<driverName>", "<dsn>")//获取一个db
err = db.Find(&arr)//这里一定要传地址，如果不传地址将会报错
```

#### First

只返回一条数据，可以加上条件OrderBy、Where

##### 使用方法

```go
type A struct
var a A
db, err := Open("<driverName>", "<dsn>")//获取一个db
db.Model(&A{}).Where(...).First(&a)
```

### 更新值

#### Update

不使用 WHERE 子句将数据表的全部数据进行更新，所以一定要加`Were`

```mysql
UPDATE <table_name> SET <column1>=<value1>[<column2> =<value2>,…] [WHERE ...]
```

mysql官网给出的语法：

单表：

```mysql
UPDATE [LOW_PRIORITY] [IGNORE] table_reference
    SET assignment_list
    [WHERE where_condition]
    [ORDER BY ...]
    [LIMIT row_count]

value:
    {expr | DEFAULT}

assignment:
    col_name = value

assignment_list:
    assignment [, assignment] ...
```

多表：

```mysql
UPDATE [LOW_PRIORITY] [IGNORE] table_references
    SET assignment_list
    [WHERE where_condition]
```

所以现在思路就很清晰了

```mysql
UPDATE <table_name>
```



我们只需要多建一个可以建立`Update`语法的函数，再进行封装就行了

##### 使用方法

更新必须携带Where方法，否则会报错

也可以添加OrderBy、Limit方法

```go
type A struct{...}
db, err := Open("<driverName>", "<dsn>")//获取一个db
db.Model(&A{}).Where("...?","..").Update("...","...")//这里可以使用map传值，也可以使用键值对传值
//map：
m := make(map[string]any)
m["column1"] = value1
m["column2"] = value2
...
db.Model(&A).Update(m)
//键值对：
db.Model(&A{}).Where("...?","..").Update("k1","v1","k2","v2",...)//必须为偶数个，因为kv键值对嘛
```

### 删除值

#### Delete

官方文档

```mysql
DELETE [LOW_PRIORITY] [QUICK] [IGNORE] FROM tbl_name [[AS] tbl_alias]
    [PARTITION (partition_name [, partition_name] ...)]
    [WHERE where_condition]
    [ORDER BY ...]
    [LIMIT row_count]
```

所以只需要实现

```mysql
DELETE FROM <table_name>
```

##### 使用方法

必须携带Where方法

可以携带OrderBy、Limit方法

```go
type A struct{}
db, err := Open("<driverName>", "<dsn>")//获取一个db
db.Model(&A{}).Where("...?","..").Delete() //一键删除 :)
```

### 数据条数

```mysql
SELECT COUNT(*) FROM <table_name>
```

#### Count

查询某张表有多少行数据

##### 使用方法

```go
type A struct{}
db, err := Open("<driverName>", "<dsn>")//获取一个db
row,err := db.Model(&A{}).Count()
```

### 事务

#### gorm

在gorm中，事务的开启是这样的

```go
//这里的fc就是需要执行的事务操作的函数，传入一个DB
func (db *DB) Transaction(fc func(tx *DB) error, opts ...*sql.TxOptions) (err error) {
    //开头在这里设置了一个panicked := true 在文末设置 panicked = false
    //最后有个延迟函数检测panicked的状态，如果panic的话就会rollback，否则就commit
   panicked := true
	//这里是连接
   if committer, ok := db.Statement.ConnPool.(TxCommitter); ok && committer != nil {
      // nested transaction
      if !db.DisableNestedTransaction {
         err = db.SavePoint(fmt.Sprintf("sp%p", fc)).Error
         if err != nil {
            return
         }

         defer func() {
            // Make sure to rollback when panic, Block error or Commit error
            if panicked || err != nil {
               db.RollbackTo(fmt.Sprintf("sp%p", fc))
            }
         }()
      }

      err = fc(db.Session(&Session{}))
   } else {
      tx := db.Begin(opts...)
      if tx.Error != nil {
         return tx.Error
      }

      defer func() {
         // Make sure to rollback when panic, Block error or Commit error
         if panicked || err != nil {
            tx.Rollback()
         }
      }()

      if err = fc(tx); err == nil {
         panicked = false
         return tx.Commit().Error
      }
   }

   panicked = false
   return
}
```

那么仿照官方写的话就是这样写了

`finisher_api.go`下的`Transaction`函数进行事务操作

就需要对之前的Exec、Query、QueryRow函数进行修改，先增加以下这个函数

```go
////这里是为了支持事务而写的一个接口，如果sql.tx不为空，那么就返回sql.tx否则返回sql.db
func (db *DB) dataBase() CommonDB {
	if db.Statement.tx != nil {
		return db.Statement.tx
	}
	return db.Statement.db
}
```

例如，需要执行Exec函数的话，原来是调用db.Statement.db来获取db对象

```go
//Exec mysql中的Exec方法，只返回一行参数
func (db *DB) Exec() (sql.Result, error) {
	//清除缓存，复用接口
	defer db.clear()
	//打印sql语句
	clog.Sql(db.Statement.sql.String(), db.Statement.values...)
	Result, err := db.Statement.db.Exec(db.Statement.sql.String(), db.Statement.values...)
	if err != nil {
		clog.Error("Exec error : %v", err)
		return nil, err
	}
	return Result, nil
}
```

现在只需要修改成da.Statement.dataBase.Exec来自动选取sql.DB或者是sql.Tx

那么修改后的代码就如下所示（Exec函数）

```go
func (db *DB) Exec() (sql.Result, error) {
	//清除缓存，复用接口
	defer db.clear()
	//打印sql语句
	clog.Sql(db.Statement.sql.String(), db.Statement.values...)
	Result, err := db.dataBase().Exec(db.Statement.sql.String(), db.Statement.values...)//修改重点
	if err != nil {
		clog.Error("Exec error : %v", err)
		return nil, err
	}
	return Result, nil
}
```

#### 使用方法

##### 壹

**通过封装的函数来进行自动commit或rollback**

gorm源码以及[database/sql Tx - detecting Commit or Rollback](https://stackoverflow.com/questions/16184238/database-sql-tx-detecting-commit-or-rollback)都有这个函数，我这里直接加以引用的这个函数

```go
db, err := Open("<driverName>", "<dsn>")//获取一个db

db.Transaction(func(db *DB) error {
    doSomeThing()
	...
})
//返回的参数要求有一个error，如果返回的err不为nil或者在函数内引发panic，都会自动回滚，如果error为nil，那么会自动提交，这里实现的函数使用了
```



##### 贰

通过手动begin和commit或者rollback来执行事务

```go
db, err := Open("<driverName>", "<dsn>")//获取一个db
err = db.Begin() //这里事务已经开始了
...
err = db.Commit() //提交 或者 db.Rollback 来回滚
```



### 迁移（migrate）

#### Gorm的Migrator 接口

GORM 提供了 Migrator 接口，该接口为每个数据库提供了统一的 API 接口，可用来为您的数据库构建独立迁移，例如：

SQLite 不支持 `ALTER COLUMN`、`DROP COLUMN`，当你试图修改表结构，GORM 将创建一个新表、复制所有数据、删除旧表、重命名新表。

Gorm的迁移实现是基于接口实现的，gorm中的DB有一个`Dialector`接口，其中有一个`Migrator(db *DB) Migrator`方法（不同的数据库都需要实现这个方法），在调用db.Migrator之后，gorm就会调用这个方法返回一个`Migrator`对象，能够对后面的函数进行调用。

在gorm中`Dialector`接口就是对各种数据库差异的屏蔽，实现能够加载不同的驱动来实现对不同数据库的操作。

#### corm的Migrator接口

依据于Gorm的Migrator的描述，Corm实现的Migrator接口实现的功能为：如果发现表结构有改变，则删除旧表，创建新的表结构，并将原数据全部迁移到新表上

而由于gorm是基于接口实现的，对Migrator包的调用是允许的，而我没有单独分一个包出去写数据库差异的隔离~~后面看源码才看到的~~，就不能实现单独分包写Migrate，就直接写在根目录下的`migrator.go`下了

![image-20220727120112614](http://110.42.184.72:8092/1658894472.png)

corn先解析传入的结构体，先对结构体的类型进行解析，查看是否存在该表，如果不存在，直接创建表格后就返回

如果存在表格的话，就会查询表格结构以及外键、主键情况，将该表的全部键都删掉。然后对比新结构体和旧结构体是否存在差异，如果存在差异的话就会将这些差异找出来（提到一个新的map中来，然后一并来处理），

如果有差异：

如果是原表不存在该字段，那么就新添加字段到上面去，如果原表存在该字段的话就修改原字段的信息为新字段的信息

#### 使用方法

```go
type A struct{}
db, err := Open("<driverName>", "<dsn>")//获取一个db
db.Migrator().AutoMigrate(&User{})
```

