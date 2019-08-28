package service

type DB interface {
	Ping() error
}

type Cache interface {
	Ping() error
}
