package database

type Comics struct {
	ID       int      `db:"id" json:"-"`
	Url      string   `db:"url" json:"url"`
	Keywords []string `db:"keywords" json:"keywords"`
}

type Database interface {
	AddComics([]*Comics) error
	GetComic(int) *Comics
	GetIds(string) []int
	BuildIndex() error
	CountComics() int
}
