package corm

import (
	"errors"
	"fmt"
	"reflect"
	"summer/summer_exercise/corm/clause"
	"summer/summer_exercise/corm/clog"
)

//Create 接收一个或多个结构体的插入
func (db *DB) Create(values ...any) (int64, error) {
	//创建一个新的[]any来装value值
	v := make([]any, 0)
	for _, value := range values {
		//如果是多个结构体的话，遍历每个结构体，对每个结构体都进行操作
		//查询表结构
		table := db.Model(value).Statement.schema
		has := db.hasTable()
		if !has {
			//不存在该表，就创建表
			err := db.Model(value).createTable()
			if err != nil {
				return 0, err
			}
		}
		//存在表之后就继续操作
		//INSERT INTO <tableName>
		db.Statement.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		//如果只是一个结构体没有值的话就需要再处理了 todo
		v = append(v, table.ParseValue(value))
	}
	//Values(?,?,...),(?,?,...)
	db.Statement.clause.Set(clause.VALUES, v...)
	//	构建完整的sql语句
	sql, args := db.Statement.clause.Build(clause.INSERT, clause.VALUES)
	result, err := db.Raw(sql, args...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//Where 这里规定where中前面的描述性句子和args打包到一个[]any中发送过去，_where函数只需要解析[]any就行
func (db *DB) Where(desc string, args ...any) *DB {
	var vas []any
	db.Statement.clause.Set(clause.WHERE, append(append(vas, desc), args...))
	return db
}

// Find 找到该结构体所在表的所有行，传入的必须是一个[]struct，否则会报错
func (db *DB) Find(value any) error {
	//总体思路：通过values映射得到表名 -> 查询表里的所有字段(遍历) -> 赋值到中间变量切片 -> 将切片里的值全部赋值给values -> 返回
	destSlice := reflect.ValueOf(value)
	if destSlice.Kind() == reflect.Ptr {
		destSlice = destSlice.Elem()
	}
	//拿到表结构赋值给table，方便使用
	//映射得到结构体的类型，先查看是否存在该表，如果不存在的话就返回error
	table := db.Model(reflect.New(destSlice.Type().Elem()).Elem().Interface()).Statement.schema
	has := db.hasTable()
	if !has {
		return errors.New("not exists this table")
	}
	//设置sql语句
	db.Statement.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	//设置可添加的参数，比如说可以添加Where、OrderBy、Limit，后续还可以继续添加一些
	sql, vas := db.Statement.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := db.Raw(sql, vas...).Query()
	if err != nil {
		return err
	}
	for rows.Next() {
		//遍历所有记录，全部存储起来
		//先创建一个指向传过来切片，类型为切片里结构体类型的结构体(Elem取了reflect.New创建的对象所指向的值)
		dest := reflect.New(destSlice.Type().Elem()).Elem()
		//声明一个切片用来存储查到的单行的所有字段值
		var v []any
		//遍历所有的字段名，将所有字段名平铺开来
		for _, name := range table.FieldNames {
			//将dest对应name的字段的地址存储到v中去
			v = append(v, dest.FieldByName(name).Addr().Interface())
		}
		//这里赋值就相当于赋值的是dest的字段
		if err = rows.Scan(v...); err != nil {
			return err
		}
		//将dest(结构体)添加到destSlice中去，然后一起归destSlice所有，这里的destSlice指向的传入的value的值
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	//循环直到所有的记录都添加到 destSlice 中，也就是传入的value中
	return rows.Close()
}

// Update 更新字段
func (db *DB) Update(kv ...any) (int64, error) {
	//这里是两种传参方式，一种是直接入参，<key1>,<value1>,<key2>,<value2>,...
	//第二种是map[key]any	 ,存储方式为 map[k1:v1 k2:v2 ...]
	//先判断是否是map类型
	m, ok := kv[0].(map[string]any)
	if !ok {
		//如果不为偶数个就返回错误
		if len(kv)%2 != 0 {
			return 0, errors.New("args number can not be base")
		}
		//先实例化map
		m = make(map[string]any)
		for i := 0; i < len(kv); i += 2 {
			//存入键值对
			m[kv[i].(string)] = kv[i+1]
		}
	}
	//生成sql语句
	fmt.Println(db.Statement.schema.Name)
	db.Statement.clause.Set(clause.UPDATE, db.Statement.schema.Name, m)
	//建立sql的语句，得sql语句和values
	sql, args := db.Statement.clause.Build(clause.UPDATE, clause.WHERE, clause.ORDERBY, clause.LIMIT)

	result, err := db.Raw(sql, args...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete 不需要传参，但是需要带上Model一起使用
func (db *DB) Delete() (int64, error) {
	db.Statement.clause.Set(clause.DELETE, db.Statement.schema.Name)
	sql, args := db.Statement.clause.Build(clause.DELETE, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	result, err := db.Raw(sql, args...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//Count 不需要传入参数，返回数据表中的行数
func (db *DB) Count() (line int, err error) {
	//构建COUNT语句
	db.Statement.clause.Set(clause.COUNT, db.Statement.schema.Name)
	//构建sql语句
	sql, args := db.Statement.clause.Build(clause.COUNT)
	//执行语句
	row := db.Raw(sql, args...).QueryRow()
	err = row.Scan(&line)
	return
}

// Transaction 开启事务
func (db *DB) Transaction(fc func(db *DB) error) (err error) {
	if err = db.Begin(); err != nil {
		return err
	}
	//执行延时函数，在gorm源码中
	//以及https://stackoverflow.com/questions/16184238/database-sql-tx-detecting-commit-or-rollback
	//中也有对次函数的介绍，这里就是引用这个函数
	defer func() {
		//如果捕捉到panic不为空，就执行回滚
		if p := recover(); p != nil {
			_ = db.Rollback()
			//执行完recover恢复之后还是需要panic掉，只是保证能够回滚事务
			panic(p)
		} else if err != nil {
			//错误不为nil需要回滚
			_ = db.Rollback()
		} else {
			//不然就提交
			err = db.Commit()
		}
	}()
	return fc(db)
}

// Begin 事务开始
func (db *DB) Begin() (err error) {
	//开启事务，并将tx赋值给db.Statement.tx中，之后就可以直接使用这个东西了
	if db.Statement.tx, err = db.Statement.db.Begin(); err != nil {
		clog.Error("%v", err)
		return
	}
	clog.Info("transaction begin ")
	return
}

// Commit 提交事务
func (db *DB) Commit() (err error) {
	if err = db.Statement.tx.Commit(); err != nil {
		clog.Error("%v", err)
	}
	db.Statement.tx = nil
	clog.Info("transaction commit successful")
	return
}

// Rollback 事务回滚
func (db *DB) Rollback() (err error) {
	if err = db.Statement.tx.Rollback(); err != nil {
		clog.Error("%v", err)
		return
	}
	db.Statement.tx = nil
	clog.Info("transaction rollback successful")
	return
}

//Limit sql语句中的Limit，查询num条数据
func (db *DB) Limit(num int) *DB {
	db.Statement.clause.Set(clause.LIMIT, num)
	return db
}

//First 查询符合条件的第一条数据
func (db *DB) First(value any) (err error) {
	//先创建一个指向value的dest
	dest := reflect.ValueOf(value)
	//处理指针
	if dest.Kind() == reflect.Ptr {
		dest = dest.Elem()
	}
	//创建一个以dest为类型的切片类型，为了后续find操作找到所有的数据
	destSlice := reflect.New(reflect.SliceOf(dest.Type()))

	//复用limit函数然后将destSlice的指针放进去 因为这里的destSlice已经是一个指针了
	err = db.Limit(1).Find(destSlice.Interface())
	if err != nil {
		return errors.New("not found data")
	}
	//查询到之后，只需要将切片中的第一条数据返回给dest就行了
	dest.Set(destSlice.Elem().Index(0))
	return
}
