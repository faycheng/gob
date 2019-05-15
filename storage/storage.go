package storage

type KV struct {
	Key   string
	Value interface{}
}

func WithKV(key string, value interface{}) KV {
	return KV{
		Key:   key,
		Value: value,
	}
}

type Storage interface {
	Open(dsn string) error
	Storage(kvs ...KV) error
	LoadAll() (kvs []KV, err error)
	Close(dsn string) error
}
