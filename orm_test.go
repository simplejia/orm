package orm

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"testing"
	"time"
)

type Demo struct {
	Id    int
	Name  string
	Value time.Time
}

func getDb() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:40001)/demo?charset=latin1&loc=Asia%2FShanghai&parseTime=true&columnsWithAlias=false")
	if err != nil {
		panic(err)
	}
	return db
}

func TestRows2Strus(t *testing.T) {
	db := getDb()

	rows, err := db.Query("select * from demo where id=? or id=?", 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var d []*Demo
	err = Rows2Strus(rows, &d)
	if err != nil {
		t.Fatal(err)
	}
	for p, v := range d {
		t.Log(p, v)
	}
}

func TestRows2Stru(t *testing.T) {
	db := getDb()

	rows, err := db.Query("select * from demo where id=? or id=?", 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var d *Demo
	err = Rows2Stru(rows, &d)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(d)
}

func TestRows2Cnts(t *testing.T) {
	db := getDb()

	rows, err := db.Query("select count(1) from demo where id=?", 1)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var d []int // 当返回值有可能为Null时(比如select max时)，应使用[]sql.NullInt64
	err = Rows2Cnts(rows, &d)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(d)
}

func TestRows2Cnt(t *testing.T) {
	db := getDb()

	rows, err := db.Query("select count(1) from demo where id=?", 1)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var d int
	err = Rows2Cnt(rows, &d)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(d)
}
