package db

type Migrator interface {
	Up(to string) error
	Down(to string) error
}
