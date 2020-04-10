package datasource

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/XH-JMC/cta/common/sqlparser/model"
	"github.com/XH-JMC/cta/util"
	"strconv"
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

// 生成UndoImage后关闭rows
func generateUndoImage(rows *sql.Rows) *Image {
	defer rows.Close()
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
			// 避免使用[]byte类型
			// 因为json.Marshal遇到[]byte类型时会做base64编码变成字符串类型，
			// 使用json.Unmarshal复原，用interface{}类型承接时只能得到字符串类型
			if val, ok := tmpRow[i].([]byte); ok {
				tmpRow[i] = string(val)
			}
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
	beforeImageRes, err := s.beforeImageStmt.QueryContext(ctx, beforeImageArgs...)
	if err != nil {
		return nil, fmt.Errorf("query before image error: %s", err)
	}
	beforeImage := generateUndoImage(beforeImageRes)

	res, err := s.stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, err
	}

	afterImageArgs := getImageArgs(s.afterImageArgsFunc, args)
	afterImageRes, err := s.afterImageStmt.QueryContext(ctx, afterImageArgs...)
	if err != nil {
		return nil, fmt.Errorf("query after image error: %s", err)
	}
	afterImage := generateUndoImage(afterImageRes)

	undoItem := &UndoItem{
		SQLType:     s.sqlParser.GetSQLType(),
		TableName:   s.sqlParser.GetTableName(),
		BeforeImage: beforeImage,
		AfterImage:  afterImage,
	}
	s.ltx.addUndoItem(undoItem)
	s.addLockKeyFromImage(beforeImage)
	return res, nil
}

func (s *Stmt) execDeleteContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	beforeImageArgs := getImageArgs(s.beforeImageArgsFunc, args)
	beforeImageRes, err := s.beforeImageStmt.QueryContext(ctx, beforeImageArgs...)
	if err != nil {
		return nil, fmt.Errorf("query before image error: %s", err)
	}
	beforeImage := generateUndoImage(beforeImageRes)

	res, err := s.stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, err
	}

	undoItem := &UndoItem{
		SQLType:     s.sqlParser.GetSQLType(),
		TableName:   s.sqlParser.GetTableName(),
		BeforeImage: beforeImage,
		AfterImage:  nil,
	}
	s.ltx.addUndoItem(undoItem)
	s.addLockKeyFromImage(beforeImage)
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
	for i := int64(0); i < rowsAffected; i++ {
		image.Rows = append(image.Rows, ImageRow{
			BusinessTablePK: ImageField{
				Name:  BusinessTablePK,
				Value: lastId + i,
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
	s.addLockKeyFromImage(image)
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
					// select for update should flush lockKeys
					beforeImageArgs := getImageArgs(s.beforeImageArgsFunc, args)
					beforeImageRes, err := s.beforeImageStmt.QueryContext(ctx, beforeImageArgs...)
					if err != nil {
						return nil, fmt.Errorf("query before image error: %s", err)
					}
					beforeImage := generateUndoImage(beforeImageRes)
					s.addLockKeyFromImage(beforeImage)
					_ = s.ltx.flushGlobalLock()
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
					// select for update should flush lockKeys
					beforeImageArgs := getImageArgs(s.beforeImageArgsFunc, args)
					beforeImageRow := s.beforeImageStmt.QueryRowContext(ctx, beforeImageArgs...)
					var pkId string
					err := beforeImageRow.Scan(&pkId)
					if err == nil && len(pkId) > 0 {
						s.addLockKey(pkId)
						_ = s.ltx.flushGlobalLock()
					}
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

func (s *Stmt) addLockKey(pkIds ...string) {
	s.ltx.addLockKey(s.sqlParser.GetTableName(), pkIds...) // lock the table
}

func (s *Stmt) addLockKeyFromImage(image *Image) {
	pkIds := make([]string, 0, len(image.Rows))
	for _, row := range image.Rows {
		pkId, err := util.Interface2Int64(row[BusinessTablePK].Value)
		if err != nil {
			continue
		}
		pkIds = append(pkIds, strconv.FormatInt(pkId, 10))

	}
	s.addLockKey(pkIds...)
}
