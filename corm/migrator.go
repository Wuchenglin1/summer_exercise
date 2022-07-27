package corm

import "summer/summer_exercise/corm/clog"

type Migrator struct {
	*DB
}

func (db *DB) Migrator() Migrator {
	return Migrator{DB: db.Session(&Session{})}
}

// CurrentDatabase 查看当前使用的数据库名
func (m Migrator) CurrentDatabase() (name string, err error) {
	err = m.Raw("SELECT DATABASE()").QueryRow().Scan(&name)
	if err != nil {
		clog.Error("not found !")
	}
	return
}

//AutoMigrate 自动迁移数据
//func (m Migrator) AutoMigrate(dst ...any) error {
//
//}
