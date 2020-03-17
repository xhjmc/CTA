package sqlparser

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xwb1989/sqlparser"
	"testing"
)

func TestParser(t *testing.T) {
	db, _ := sql.Open("mysql", "jmc:chenjinming@tcp(127.0.0.1:3306)/cta?charset=utf8")
	q := "select * from test where id = :v1;"
	res := db.QueryRow(q, 1)
	var a, b interface{}
	_ = res.Scan(&a, &b)
	fmt.Println(a, b)

	sqls := []string{
		"select * from tbl where id = ?",
		"SELECT * FROM tbl WHERE a = :v1 and b = 'bbb' and c = 1 and d = 0.5 and e = true for update;",
		"insert into tbl(a, b, c, d, e) values(?, 'bbb', 1, 0.5, true);",
		"update tbl set f = null where a = :v1 and b = 'bbb' and c = 1 and d = 0.5 and e = true;",
		"delete from tbl where a = ? and b = 'bbb' and c = 1 and d = 0.5 and e = true;",
	}
	for _, sql := range sqls {
		stmt, err := sqlparser.Parse(sql)
		if err != nil {
			// Do something with the err
		}
		fmt.Println(stmt, err)
		// Otherwise do something with stmt
		switch stmt := stmt.(type) {
		case *sqlparser.Select:
			fmt.Println("select", stmt)
		case *sqlparser.Insert:
			fmt.Println("insert", stmt)
		case *sqlparser.Update:
			fmt.Println("update", stmt)
			query := sqlparser.Select{
				SelectExprs: []sqlparser.SelectExpr{
					&sqlparser.StarExpr{},
				},
				From:  stmt.TableExprs,
				Where: stmt.Where,
			}
			buf := sqlparser.NewTrackedBuffer(nil)
			query.Format(buf)
			querySql := buf.String()
			fmt.Println(querySql)
		case *sqlparser.Delete:
			fmt.Println("delete", stmt)
		}

		buf := sqlparser.NewTrackedBuffer(nil)
		stmt.Format(buf)
		fmt.Println(buf.String())
	}
}

func TestSQLParser(t *testing.T) {
	sql := "select * from tbl where id = ? and a = ':v1' and b = ?"
	stmt, _ := sqlparser.Parse(sql)
	buf := sqlparser.NewTrackedBuffer(nil)
	stmt.Format(buf)
	newSql := buf.String()
	fmt.Println(newSql)
}

/*
有限状态机，用于匹配没有被单引号引用的":v[0-9]*"
from	to			event	action
s0		s1			'
s0		s2			:
s0		s0			else

s1		s0			'
s1		s1			else

s2		s3			v
s2		s1			'
s2		s0			else

s3		s3			0-9
s3		s0			else	将从s2开始匹配到的子串替换为?
*/
//type Status = byte
//
//const (
//	S0 = Status(iota)
//	S1
//	S2
//	S3
//	S4
//)
//
//func fixSQLPlaceholders(sql string) string {
//	bytes := make([]byte, len(sql), len(sql))
//	beginId := 0
//	newLen := 0
//	status := S0
//	for i, c := range sql {
//		switch status {
//		case S0:
//			switch c {
//			case '\'':
//				status = S1
//			case ':':
//				status = S2
//				beginId = i
//			}
//		case S1:
//			switch c {
//			case '\'':
//				status = S0
//			}
//		case S2:
//			switch c {
//			case 'v':
//				status = S3
//			}
//		}
//		bytes[id] = byte(c)
//		id++
//	}
//}

// (((?!'|:v[0-9]*).)*('[^']*'((?!'|:v[0-9]*).)*)*)(:v[0-9]*)(((?!'|:v[0-9]*).)*('[^']*'((?!'|:v[0-9]*).)*)*)
// $1?$6
