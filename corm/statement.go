package corm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"summer/summer_exercise/corm/clause"
	"summer/summer_exercise/corm/clog"
	"summer/summer_exercise/corm/dialect"
	"summer/summer_exercise/corm/schema"
)

type Statement struct {
	db        *sql.DB
	tx        *sql.Tx
	clause    clause.Clause
	dialector dialect.Dialector
	sql       strings.Builder //往里面塞sql语句
	values    []any           //用来存储sql语句后面的args
	schema    *schema.Schema
}

// CommonDB 定义一个接口，sql.tx和sql.db都实现了这个接口，所以可以直接在tx和db之间切换啦~
type CommonDB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

//清除里面的缓存
func (db *DB) clear() {
	db.Statement.sql.Reset()
	db.Statement.values = nil
	db.Statement.clause = clause.Clause{}
	db.Statement.schema = nil
}

// Raw 手写sql语句
func (db *DB) Raw(sql string, args ...any) *DB {
	db.Statement.sql.WriteString(sql)
	db.Statement.sql.WriteString(" ")
	db.Statement.values = append(db.Statement.values, args...)
	return db
}

//写好Raw之后就可以直接通过Exec、Query、QueryRow来调用了

//Exec mysql中的Exec方法，只返回一行参数
func (db *DB) Exec() (sql.Result, error) {
	//清除缓存，复用接口
	defer db.clear()
	//打印sql语句
	clog.Sql(db.Statement.sql.String(), db.Statement.values...)
	Result, err := db.dataBase().Exec(db.Statement.sql.String(), db.Statement.values...)
	if err != nil {
		clog.Error("Exec error : %v", err)
		return nil, err
	}
	return Result, nil
}

//Query mysql中的Query方法，返回多行信息
func (db *DB) Query() (*sql.Rows, error) {
	defer db.clear()
	clog.Sql(db.Statement.sql.String(), db.Statement.values...)
	rows, err := db.dataBase().Query(db.Statement.sql.String(), db.Statement.values...)
	if err != nil {
		clog.Error("query error : %v", err)
		return nil, err
	}
	return rows, nil
}

//QueryRow mysql中的QueryRow方法，只返回一行信息
func (db *DB) QueryRow() *sql.Row {
	defer db.clear()
	clog.Sql(db.Statement.sql.String(), db.Statement.values...)
	return db.dataBase().QueryRow(db.Statement.sql.String(), db.Statement.values...)
}

//Model 解析传入过来的value(结构体)一定不能是指针！
func (db *DB) Model(value any) *DB {
	//先对传入的model进行判断，如果这次的结构体与上一次的结构体相同的话，就不需要解析了（少浪费一次时间），直接就返回db
	if db.Statement.schema == nil || reflect.ValueOf(value).Type() != reflect.Indirect(reflect.ValueOf(db.Statement.schema.Model)).Type() {
		db.Statement.schema = schema.ParseType(value, db.Statement.dialector)
	}
	return db
}

//createTable 给finisher_api提供一个建表服务的，将db里面已经解析好了的schema用来创表
//注意！这里只能在没有表且结构体已解析好之后才能调用这个函数，否则会报错
func (db *DB) createTable() error {
	table := db.Statement.schema
	var columns []string
	for _, field := range table.Fields {
		//将要创建的值、类型、外键都加入进去
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	//添加逗号
	desc := strings.Join(columns, ",")
	_, err := db.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

//删除表
func (db *DB) dropTable() error {
	_, err := db.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s;", db.Statement.schema.Name)).Exec()
	return err
}

//查看是否有表
func (db *DB) hasTable() bool {
	//先拿到每个方言数据库的`独特的`是否有表的 sql语句 和 args ，然后传入sql中执行查看是否有表
	Sql, args := db.Statement.dialector.TableExistSql(db.Statement.schema.Name)
	row := db.Raw(Sql, args...).QueryRow()
	var tmp string
	err := row.Scan(&tmp)
	if err != nil {
		return false
	} else {
		return true
	}
}

//这里是为了支持事务而写的一个接口，如果sql.tx不为空，那么就返回sql.tx否则返回sql.db
func (db *DB) dataBase() CommonDB {
	if db.Statement.tx != nil {
		return db.Statement.tx
	}
	return db.Statement.db
}
