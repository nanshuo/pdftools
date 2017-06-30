package backend

import "github.com/jiusanzhou/pdf2html/pkg/server/job"

type Backend interface {

	Connect() error

	DisConnect() error

	Init() error

	GetJob(id string, hang bool) (*job.ConvertJob, error)

	PutJob(job *job.ConvertJob) error

	DelJob(id string) error
}