package db

type Connecter interface {
	Connect() interface{}
	Close() error
}
