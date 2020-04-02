package model

type SQLType int

const (
	INSERT SQLType = iota
	DELETE
	UPDATE
	SELECT
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
	default:
		return "Unknown"
	}
}
