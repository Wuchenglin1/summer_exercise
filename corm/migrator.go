package corm

import (
	"errors"
	"fmt"
	"summer/summer_exercise/corm/clog"
	"summer/summer_exercise/corm/schema"
	"summer/summer_exercise/corm/tool"
)

type Migrator struct {
	*DB
}

func (db *DB) Migrator() *Migrator {
	return &Migrator{DB: db.Session(&Session{})}
}

// CurrentDatabase 查看当前使用的数据库名
func (m *Migrator) CurrentDatabase() (name string, err error) {
	err = m.Raw("SELECT DATABASE()").QueryRow().Scan(&name)
	if err != nil {
		clog.Error("not found !")
	}
	return
}

//AutoMigrate 自动迁移数据
func (m *Migrator) AutoMigrate(dst ...any) error {
	//先解析传入的结构体
	for _, structure := range dst {
		//先创建一个map存储旧的column
		var (
			oldColumns    []*schema.Field
			oldCoulumnMap = make(map[string]*schema.Field)
			constraintMap = make(map[string]string)
		)
		//解析新的结构体
		newSchema := schema.ParseType(structure, m.Statement.dialector)
		//查看新结构体是否存在表格了
		has := m.Model(newSchema.Model).hasTable()
		if !has {
			m.Statement.schema = newSchema
			//如果没有的话就创一个新表
			return m.createTable()
		}
		//查询表结构
		rows1, err := m.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", newSchema.Name)).Query()
		if err != nil {
			clog.Error("%v", err)
			return errors.New(fmt.Sprintf("%v", err))
		}
		//查询所有的键
		rows2, err := m.Raw(fmt.Sprintf(`select CONSTRAINT_NAME,COLUMN_NAME from INFORMATION_SCHEMA.KEY_COLUMN_USAGE where TABLE_NAME = "%s"`, newSchema.Name)).Query()
		if err != nil {
			clog.Error("get CONSTRAINT_NAME error : %v", err)
			return err
		}
		for rows2.Next() {
			var (
				constraintName string
				columnName     string
			)
			err = rows2.Scan(&constraintName, &columnName)
			if err != nil {
				clog.Error("scan constraintName error : %v", err)
				return err
			}
			//将键值信息存储起来
			constraintMap[columnName] = constraintName
			//将原表中的键都删了
			if constraintName == "PRIMARY KEY" {
				_, err = m.Raw(fmt.Sprintf(`ALTER TABLE %s DROP PRIMARY KEY`, newSchema.Name)).Exec()
				if err != nil {
					return err
				}
			} else {
				_, err = m.Raw(fmt.Sprintf(`ALTER TABLE %s DROP FOREIGN KEY %s`, newSchema.Name, constraintName)).Exec()
				if err != nil {
					return err
				}
			}
		}

		//获取字段名信息
		types, err := rows1.ColumnTypes()
		if err != nil {
			clog.Error("get column type error : %v", err)
			return err
		}
		for _, v := range types {
			//	遍历类型列表，转换成dialect中对应的类型、键
			field := &schema.Field{}
			//赋值消息
			field.Name = v.Name()
			field.Type = v.DatabaseTypeName()
			if constraint, ok := constraintMap[field.Name]; ok {
				field.Tag = constraint
			}
			//将字段加入到oldColumn里面去
			oldColumns = append(oldColumns, field)
			oldCoulumnMap[field.Name] = field
		}
		//现在只需要比较旧的和新的的区别就行了，这个函数是获取旧的键、类型与新键、值不同的键、值
		modifyMap := tool.GetColumnDiff(newSchema.FieldsArr, oldColumns)
		for columnName := range modifyMap {
			//if oldCoulumnMap[columnName].Tag == "PRIMARY KEY" || modifyMap[columnName].Tag == "PRIMARY KEY" {
			//	continue
			//}
			if _, ok := oldCoulumnMap[columnName]; !ok {
				//在旧的字段里面没有该字段，那么就添加该字段
				_, err = m.Raw(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s %s", newSchema.Name, modifyMap[columnName].Name, modifyMap[columnName].Type, modifyMap[columnName].Tag)).Exec()
				if err != nil {
					clog.Error("add column %v error : %v", modifyMap[columnName].Name, err)
					return err
				}
			} else {
				//如果在旧的该字段里面有这个字段，那么就修改该字段的信息
				_, err = m.Raw(fmt.Sprintf("ALTER TABLE %s MODIFY %s %s %s", newSchema.Name, modifyMap[columnName].Name, modifyMap[columnName].Type, modifyMap[columnName].Tag)).Exec()
				if err != nil {
					clog.Error("modify column %v error : %v", modifyMap[columnName].Name, err)
					return err
				}
			}
		}
		rows1.Close()
		rows2.Close()
	}
	return nil
}
