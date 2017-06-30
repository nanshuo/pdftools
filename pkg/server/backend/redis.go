package backend

import (
	"github.com/jiusanzhou/pdf2html/pkg/server/job"
	"github.com/mediocregopher/radix.v2/redis"
)

type RedisBackend struct {
	redis      *redis.Client
	reCntCount int
	address    string
	prefix     string
}

func (rb *RedisBackend) Connect() error {
	var err error
	rb.redis, err = redis.Dial("tcp", rb.address)
	rb.reCntCount++
	return err
}

func (rb *RedisBackend) DisConnect() error {
	return rb.redis.Close()
}

func (rb *RedisBackend) Init() error {

	if rb.redis == nil {
		rb.Connect()
	}

	return nil
}

func (rb *RedisBackend) GetJob(id string, hang bool) (*job.ConvertJob, error) {

	return
}

func (rb *RedisBackend) PutJob(job *job.ConvertJob) error {

	return nil
}

func (rb *RedisBackend) DelJob(id string) error {

	return nil
}

func NewRedisBackend(addr, prefix string) (Backend, error) {

	rb := &RedisBackend{
		address: addr,
		prefix:  prefix,
	}
	rb.Connect()

	return rb, nil
}
