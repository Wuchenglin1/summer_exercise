package corm

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"testing"
)

type User struct {
	Name string `corm:"PRIMARY KEY"`
	Age  int
}

var db *DB

func init() {
	database, err := Open("", "")
	if err != nil {
		log.Fatalf("open sql error : %v", err)
		return
	}
	db = database
}

func TestDB_Raw(t *testing.T) {
	defer db.Close()
	//exec
	exec, err := db.Raw("insert into User(Name,Age) values(?,?)", "zhangsan", 18).Exec()
	fmt.Println(exec, err)
	exec, err = db.Raw("delete from User where name = ?", "zhangsan").Exec()
	fmt.Println(exec, err)
	//queryRow
	row := db.Raw("select Age from User where name = ?", "Jack").QueryRow()
	var age int
	err = row.Scan(&age)
	//query
	var get []int
	rows, err := db.Raw("select Age from User").Query()
	defer rows.Close()
	for rows.Next() {
		var tmp int
		err = rows.Scan(&tmp)
		if err != nil {
			fmt.Println(err)
			return
		}
		get = append(get, tmp)
	}
	fmt.Println(err, get)

}

func TestHasTable(t *testing.T) {
	_ = db.Model(&User{}).dropTable()
	_ = db.Model(&User{}).createTable()
	get := db.Model(&User{}).hasTable()
	if !get {
		t.Fatal("failed to check the table")
	}

}

func TestDB_Create(t *testing.T) {
	var wantRow int64 = 1
	_ = db.Model(&User{}).dropTable()
	//_ = db.Model(&User{}).createTable()
	row, err := db.Create(&User{"testUser3", 28})
	if row != wantRow {
		t.Fatalf("got : %v but want : %v", row, wantRow)
	} else if err != nil {
		t.Fatalf("got error : %v", err)
	}
}

func TestDB_Find(t *testing.T) {
	_ = db.Model(&User{}).dropTable()
	u := []any{
		User{"Celia", 19},
		User{"Tom", 20},
		User{"Jack", 28},
	}
	rows, err := db.Create(u...)
	if err != nil {
		t.Fatalf("create table error : %v", err)
	}
	t.Logf("effected : %v", rows)
	//FIND
	var got []User
	err = db.Find(&got)
	if err != nil {
		t.Fatalf("find objects error : %v", err)
	} else if reflect.DeepEqual(u, got) {
		t.Fatalf("got : %v but want : %v", got, u)
	}
	fmt.Println(got)
}

func TestDB_Update(t *testing.T) {
	_ = db.Model(&User{}).dropTable()
	u1 := User{"testUser", 18}
	_, err := db.Create(&u1)
	if err != nil {
		t.Fatalf("create table and columns error : %v", err)
	}
	m := make(map[string]any)
	m["Name"] = "xiaohong"
	m["Age"] = 20
	_, err = db.Model(&User{}).Where("Age = ?", 18).Update("Name", "xiaoming")
	var uA []User
	_ = db.Find(&uA)
	_, err = db.Model(&User{}).Where("Age = ?", 18).Update(m)

	if err != nil {
		t.Fatalf("update column error : %v", err)
	}
	fmt.Println(uA)
}

func TestDB_Delete(t *testing.T) {
	_ = db.Model(&User{}).dropTable()
	u1 := User{"xiaoming", 18}
	u2 := User{"xiaohong", 19}
	u3 := User{"xiaoqiang", 20}

	_, err := db.Create(&u1, &u2, &u3)
	if err != nil {
		t.Fatalf("create data error : %v", err)
	}
	i, err := db.Model(&User{}).Where("age = ?", 19).Delete()
	fmt.Println(i, err)
}

func TestDB_Count(t *testing.T) {
	_ = db.Model(&User{}).dropTable()
	u1 := User{"xiaoming", 18}
	u2 := User{"xiaohong", 19}
	u3 := User{"xiaoqiang", 20}
	_, _ = db.Create(&u1, &u2, &u3)

	var want = 3
	got, err := db.Model(&User{}).Count()
	if err != nil {
		t.Fatalf("Count error : %v", err)
	} else if got != want {
		t.Fatalf("got %v but want %v", got, want)
	}
}

func TestDB_Transaction(t *testing.T) {
	_ = db.Model(&User{}).dropTable()
	db.Model(&User{}).createTable()
	err := db.Transaction(func(db *DB) error {
		_, err := db.Raw("insert into User(Name,Age) values(?,?)", "xiaoxiong", 20).Exec()
		if err != nil {
			return err
		}
		return errors.New("error")
	})
	db.Begin()
	if err != nil {
		t.Logf("there is something wrong : %v", err)
	}

	err = db.Transaction(func(db *DB) error {
		_, err = db.Raw("insert into User(Name,Age) values(?,?)", "xiaozhang", 20).Exec()
		if err != nil {
			return err
		}
		return err
	})

	if err != nil {
		t.Fatalf("there is something wrong : %v", err)
	}

}

func TestDB_First(t *testing.T) {
	db.Model(&User{}).dropTable()
	db.Create(&User{"xiaoming", 18}, &User{"xiaohai", 20}, &User{"xiaozheng", 28})
	u := User{}
	err := db.Model(&User{}).First(u)
	if err != nil {
		t.Fatalf("First data error : %v", err)
	} else {

	}
}

func TestMigrator_CurrentDatabase(t *testing.T) {
	//var want = "test"
	//dataBaseName := db.Migrator().CurrentDatabase()
	//if dataBaseName != want {
	//	t.Fatalf("out current database name : %v but want : %v", dataBaseName, want)
	//}
}
