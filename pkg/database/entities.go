package database

type Comics struct {
	ID       int      `db:"id"`
	Url      string   `db:"url"`
	Keywords []string `db:"keywords"`
}
