package service

type ServiceGroupEnter struct {
	DbService
	CronService
}

var (
	ServiceGroup = new(ServiceGroupEnter)
)
