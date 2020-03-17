package tx

type Transaction interface {
	Commit() error
	RollBack() error
}
