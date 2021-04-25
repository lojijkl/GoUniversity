package week2

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"errors"

	_ "github.com/go-sql-driver/mysql"
)

const (
	USERNAME = "root"
	PASSWORD = "123"
	NETWORK  = "tcp"
	SERVER   = "localhost"
	PORT     = 3306
	DATABASE = "test"
)

var DB *sql.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Open mysql failed,err:%v\n", err)
		return
	}
	DB.SetConnMaxLifetime(100 * time.Second) //最大连接周期，超过时间的连接就close
	DB.SetMaxOpenConns(100)                  //设置最大连接数
	DB.SetMaxIdleConns(16)                   //设置闲置连接数
}

type Temp struct {
	ID   int64          `db:"id"`
	Name sql.NullString `db:"name"` //由于在mysql的users表中name没有设置为NOT NULL,所以name可能为null,在查询过程中会返回nil，如果是string类型则无法接收nil,但sql.NullString则可以接收nil值
}

func Test_sqlRows(t *testing.T) {
	fmt.Println(controller())
}

func controller() (*Temp, error) {
	data, err := dao()
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			// fmt.Println(err)
			// fmt.Printf("---%v\n", err)
			// fmt.Printf("===%+v\n", err)
			// fmt.Printf(">>>>%v%v\n", err, errors.Unwrap(err))
			// 底部返还error，至于需不需要用err，让controller判断消化，或者返还
			err = nil
		default:
			fmt.Println("other error:", err)
		}
		return data, err
	}
	return data, err
}

func dao() (*Temp, error) {
	temp := new(Temp)
	row := DB.QueryRow("select * from temp where id=?", 1)
	if err := row.Scan(&temp.ID, &temp.Name); err != nil {

		return nil, fmt.Errorf("nodata:%w", err)
	}
	return temp, nil
}

/*
CREATE TABLE `temp` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
*/
