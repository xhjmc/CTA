package sqlparser

type SQLType int

const (
	INSERT SQLType = iota
	DELETE
	UPDATE
	SELECT_FOR_UPDATE
)
