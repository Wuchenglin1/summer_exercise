package corm

import (
	"database/sql"
	"summer/summer_exercise/corm/clause"
	"summer/summer_exercise/corm/clog"
	"summer/summer_exercise/corm/dialect"
)

type DB struct {
	Statement *Statement
}

// Session 可以在这里加一些配置选项 todo
type Session struct {
}

// Open 建立连接
func Open(driverName string, dsn string) (database *DB, err error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		clog.Error("open sql error : %v", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		clog.Error("server ping field ", err)
		return nil, err
	}

	dial, ok := dialect.GetDialect(driverName)
	if !ok {
		clog.Error("not support %v driver", driverName)
	}
	database = &DB{}
	database.Statement = &Statement{
		db:        db,
		dialector: dial,
	}

	clog.Info("connect database successful")
	return
}

func (db *DB) Close() {
	if err := db.Statement.db.Close(); err != nil {
		clog.Error("close database error : %v", err)
		return
	}
	clog.Info("close connection successful")
}

// Session 创建一个新的db对象
func (db *DB) Session(config *Session) *DB {
	tx := &DB{
		Statement: db.Statement,
	}
	tx.Statement.tx = nil
	tx.Statement.clause = clause.Clause{}
	tx.Statement.sql.Reset()
	tx.Statement.schema = nil
	tx.Statement.values = nil
	//可以加一些变量设置 todo
	return tx
}
