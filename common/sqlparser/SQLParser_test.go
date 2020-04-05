package sqlparser_test

import (
	"database/sql"
	"fmt"
	"github.com/XH-JMC/cta/common/sqlparser"
	"github.com/XH-JMC/cta/common/sqlparser/model"
	_ "github.com/go-sql-driver/mysql"
	parser "github.com/xwb1989/sqlparser"
	"testing"
)

const MySQL = "mysql"

func TestParser(t *testing.T) {
	db, _ := sql.Open(MySQL, "jmc:chenjinming@tcp(127.0.0.1:3306)/cta?charset=utf8")
	q := "select * from test where id = :v1;"
	res := db.QueryRow(q, 1)
	var a, b interface{}
	_ = res.Scan(&a, &b)
	fmt.Println(a, b)

	sqls := []string{
		"SELECT * FROM tbl WHERE a = :v1 and b = 'bbb' and c = 1 and d = 0.5 and e = true for update;",
		"insert into tbl(a, b, c, d, e) values(?, 'bbb', 1, 0.5, true);",
		"update tbl set f = null where a = :v1 and b = 'bbb' and c = 1 and d = 0.5 and e = true;",
		"delete from tbl where a = ? and b = 'bbb' and c = 1 and d = 0.5 and e = true;",
	}
	for _, sql := range sqls {
		stmt, err := parser.Parse(sql)
		if err != nil {
			// Do something with the err
		}
		fmt.Println(stmt, err)
		// Otherwise do something with stmt
		switch stmt := stmt.(type) {
		case *parser.Select:
			fmt.Println("select", stmt)
		case *parser.Insert:
			fmt.Println("insert", stmt)
		case *parser.Update:
			fmt.Println("update", stmt)
			query := parser.Select{
				SelectExprs: []parser.SelectExpr{
					&parser.StarExpr{},
				},
				From:  stmt.TableExprs,
				Where: stmt.Where,
			}
			buf := parser.NewTrackedBuffer(nil)
			query.Format(buf)
			querySql := buf.String()
			fmt.Println(querySql)
		case *parser.Delete:
			fmt.Println("delete", stmt)
		}

		buf := parser.NewTrackedBuffer(nil)
		stmt.Format(buf)
		fmt.Println(buf.String())
	}
}

func TestSQLInsertParser(t *testing.T) {
	sql := "insert into tbl(a,b,c) values(?, 'bbb', 1), (0.5, true, nowtime());"
	stmt, _ := parser.Parse(sql)

	if stmt, ok := stmt.(*parser.Insert); ok {
		fmt.Println("table name:", stmt.Table.Name)
		for _, col := range stmt.Columns {
			fmt.Printf("\t%s", col.String())
		}
		fmt.Println()
		if rows, ok := stmt.Rows.(parser.Values); ok {
			buff := parser.NewTrackedBuffer(nil)
			for _, row := range rows {
				for _, expr := range row {
					expr.Format(buff)
					fmt.Printf("\t%s", buff.String())
					buff.Reset()
				}
				fmt.Println()
			}
		}
	}

	buf := parser.NewTrackedBuffer(func(buf *parser.TrackedBuffer, node parser.SQLNode) {
		if node, ok := node.(*parser.SQLVal); ok {
			switch node.Type {
			case parser.ValArg:
				buf.WriteArg("?")
				return
			}
		}
		node.Format(buf)
	})
	stmt.Format(buf)
	newSql := buf.String()
	fmt.Println(newSql)
}

func TestMySQLInsertParser(t *testing.T) {
	query := "insert into db.tbl(a,b,c) values(?, 'bbb', 1), (0.5, true, nowtime());"
	factory := sqlparser.GetSQLParserFactory(MySQL)
	parser, err := factory.NewSQLParser(query)
	if err != nil {
		panic(err)
	}
	insertParser := parser.(model.SQLInsertParser)
	fmt.Println(insertParser.GetSQL())
	fmt.Println(insertParser.GetSQLType().String())
	fmt.Println(insertParser.GetTableName())
	fmt.Println(insertParser.GetInsertColumns())
	fmt.Println(insertParser.GetRows())
}

func TestMySQLDeleteParser(t *testing.T) {
	query := "delete from db.tbl where a = ? and b = 'bbb' and c = 1 and d = 0.5 and e = true;"
	factory := sqlparser.GetSQLParserFactory(MySQL)
	parser, err := factory.NewSQLParser(query)
	if err != nil {
		panic(err)
	}
	deleteParser := parser.(model.SQLDeleteParser)
	fmt.Println(deleteParser.GetSQL())
	fmt.Println(deleteParser.GetSQLType().String())
	fmt.Println(deleteParser.GetTableName())
	fmt.Println(deleteParser.GetCondition())
}

func TestMySQLUpdateParser(t *testing.T) {
	query := "update t0, tbl1 t1, tbl2 as t2 set t1.a = t2.b, t2.a = t1.b, c = ? where t1.a = t2.a and c = ?;"
	factory := sqlparser.GetSQLParserFactory(MySQL)
	parser, err := factory.NewSQLParser(query)
	if err != nil {
		panic(err)
	}
	updateParser := parser.(model.SQLUpdateParser)
	fmt.Println(updateParser.GetSQL())
	fmt.Println(updateParser.GetSQLType().String())
	fmt.Println(updateParser.GetTableName())
	fmt.Println(updateParser.GetTableList())
	fmt.Println(updateParser.GetUpdateColumns())
	fmt.Println(updateParser.GetCondition())
	fmt.Println(updateParser.CountPlaceholderInCondition())
}

func TestMySQLSelectParser(t *testing.T) {
	query := "select t0.a a, t1.* from tbl0 t0, tbl1 as t1, tbl2 where t0.a = t1.a and t0.b = tbl2.b for UPDATE;"
	factory := sqlparser.GetSQLParserFactory(MySQL)
	parser, err := factory.NewSQLParser(query)
	if err != nil {
		panic(err)
	}
	selectParser := parser.(model.SQLSelectParser)
	fmt.Println(selectParser.GetSQL())
	fmt.Println(selectParser.GetSQLType().String())
	fmt.Println(selectParser.GetTableName())
	fmt.Println(selectParser.GetTableList())
	fmt.Println(selectParser.GetSelectColumns())
	fmt.Println(selectParser.GetCondition())
	fmt.Println(selectParser.IsSelectForUpdate())
}
