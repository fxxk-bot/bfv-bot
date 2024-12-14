package api

import (
	"bfv-bot/service"
	"regexp"
)

type ApiGroupEnter struct {
	EventApi
}

var (
	// ApiGroup public
	ApiGroup       = new(ApiGroupEnter)
	GroupAnswerReg = regexp.MustCompile(`答案：(.*)`)

	// dbService private
	dbService = service.ServiceGroup.DbService
)
