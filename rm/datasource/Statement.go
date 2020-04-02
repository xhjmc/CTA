package datasource

import (
	"context"
	"cta/common/sqlparser/model"
	"database/sql"
	"fmt"
)

type Stmt struct {
	ltx                 *LocalTransaction
	stmt                *sql.Stmt
	beforeImageStmt     *sql.Stmt
	beforeImageArgsFunc func(args []interface{}) []interface{}
	afterImageStmt      *sql.Stmt
	afterImageArgsFunc  func(args []interface{}) []interface{}
	sqlParser           model.SQLParser
}

func generateUndoImage(rows *sql.Rows) *Image {
	image := &Image{Rows: make([]ImageRow, 0)}
	beforeImageCols, _ := rows.Columns()
	tmpLen := len(beforeImageCols)
	tmpDest := make([]interface{}, tmpLen, tmpLen)
	tmpRow := make([]interface{}, tmpLen, tmpLen)
	for i := 0; i < tmpLen; i++ {
		tmpDest[i] = &tmpRow[i]
	}
	for rows.Next() {
		err := rows.Scan(tmpDest...)
		if err != nil {
			continue
		}
		imageRow := ImageRow{}
		for i, colName := range beforeImageCols {
			imageRow[colName] = ImageField{
				Name:  colName,
				Value: tmpRow[i],
			}
		}
		image.Rows = append(image.Rows, imageRow)
	}
	return image
}

func getImageArgs(imageArgsFunc func(args []interface{}) []interface{}, args []interface{}) []interface{} {
	if imageArgsFunc != nil {
		return imageArgsFunc(args)
	}
	return args
}

func (s *Stmt) execUpdateContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	beforeImageArgs := getImageArgs(s.beforeImageArgsFunc, args)
	beforeImage, err := s.beforeImageStmt.QueryContext(ctx, beforeImageArgs...)
	if err != nil {
		return nil, fmt.Errorf("query before image error: %s", err)
	}
	res, err := s.stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	afterImageArgs := getImageArgs(s.afterImageArgsFunc, args)
	afterImage, err := s.afterImageStmt.QueryContext(ctx, afterImageArgs...)
	if err != nil {
		return nil, fmt.Errorf("query after image error: %s", err)
	}

	undoItem := &UndoItem{
		SQLType:     s.sqlParser.GetSQLType(),
		TableName:   s.sqlParser.GetTableName(),
		BeforeImage: generateUndoImage(beforeImage),
		AfterImage:  generateUndoImage(afterImage),
	}
	s.ltx.addUndoItem(undoItem)
	s.addLockKey()
	return res, nil
}

func (s *Stmt) execDeleteContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	beforeImageArgs := getImageArgs(s.beforeImageArgsFunc, args)
	beforeImage, err := s.beforeImageStmt.QueryContext(ctx, beforeImageArgs...)
	if err != nil {
		return nil, fmt.Errorf("query before image error: %s", err)
	}
	res, err := s.stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, err
	}

	undoItem := &UndoItem{
		SQLType:     s.sqlParser.GetSQLType(),
		TableName:   s.sqlParser.GetTableName(),
		BeforeImage: generateUndoImage(beforeImage),
		AfterImage:  nil,
	}
	s.ltx.addUndoItem(undoItem)
	s.addLockKey()
	return res, nil
}

// 约定业务表主键固定为名为pk_id，且会自动自增
func (s *Stmt) execInsertContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	res, err := s.stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	lastId, _ := res.LastInsertId()
	rowsAffected, _ := res.RowsAffected()

	image := &Image{Rows: make([]ImageRow, 0)}
	for i := lastId - rowsAffected + 1; i <= lastId; i++ {
		image.Rows = append(image.Rows, ImageRow{
			BusinessPK: ImageField{
				Name:  BusinessPK,
				Value: i,
			},
		})
	}

	undoItem := &UndoItem{
		SQLType:     s.sqlParser.GetSQLType(),
		TableName:   s.sqlParser.GetTableName(),
		BeforeImage: nil,
		AfterImage:  image,
	}
	s.ltx.addUndoItem(undoItem)
	s.addLockKey()
	return res, nil
}

func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	if s.sqlParser != nil {
		switch s.sqlParser.GetSQLType() {
		case model.INSERT:
			return s.execInsertContext(ctx, args...)
		case model.DELETE:
			return s.execDeleteContext(ctx, args...)
		case model.UPDATE:
			return s.execUpdateContext(ctx, args...)
		}
	}
	return s.stmt.ExecContext(ctx, args...)
}

func (s *Stmt) Exec(args ...interface{}) (sql.Result, error) {
	return s.ExecContext(context.Background(), args...)
}

func (s *Stmt) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	if s.sqlParser != nil {
		switch s.sqlParser.GetSQLType() {
		case model.SELECT:
			if parser, ok := s.sqlParser.(model.SQLSelectParser); ok {
				if parser.IsSelectForUpdate() {
					s.addLockKey()
				}
			}
		}
	}
	return s.stmt.QueryContext(ctx, args...)
}

func (s *Stmt) Query(args ...interface{}) (*sql.Rows, error) {
	return s.QueryContext(context.Background(), args...)
}

func (s *Stmt) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	if s.sqlParser != nil {
		switch s.sqlParser.GetSQLType() {
		case model.SELECT:
			if parser, ok := s.sqlParser.(model.SQLSelectParser); ok {
				if parser.IsSelectForUpdate() {
					s.addLockKey()
				}
			}
		}
	}
	return s.stmt.QueryRowContext(ctx, args...)
}

func (s *Stmt) QueryRow(args ...interface{}) *sql.Row {
	return s.QueryRowContext(context.Background(), args...)
}

func (s *Stmt) Close() error {
	if s.beforeImageStmt != nil {
		_ = s.beforeImageStmt.Close()
	}
	if s.afterImageStmt != nil {
		_ = s.afterImageStmt.Close()
	}
	return s.stmt.Close()
}

func (s *Stmt) addLockKey() {
	s.ltx.addLockKey(s.sqlParser.GetTableName()) // lock the table
}
