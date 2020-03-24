package model

type SQLType int

const (
	INSERT SQLType = iota
	DELETE
	UPDATE
	SELECT
	SELECT_FOR_UPDATE
)

func (t SQLType) String() string {
	switch t {
	case INSERT:
		return "INSERT"
	case DELETE:
		return "DELETE"
	case UPDATE:
		return "UPDATE"
	case SELECT:
		return "SELECT"
	case SELECT_FOR_UPDATE:
		return "SELECT_FOR_UPDATE"
	default:
		return "Unknown"
	}
}
