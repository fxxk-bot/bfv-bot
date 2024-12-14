package service

import (
	"bfv-bot/common/global"
	"bfv-bot/common/utils"
	"bfv-bot/model/po"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

type DbService struct{}

func (b *DbService) QueryAllBlackList() map[string]po.Blacklist {

	listMap := make(map[string]po.Blacklist)

	var blacklist []po.Blacklist
	err := global.GDb.Raw("SELECT `id`, `name`, `reason` FROM `blacklist`").Scan(&blacklist).Error
	if err != nil {
		global.GLog.Error("QueryAllBlackList", zap.Error(err))
		return listMap
	}

	for _, item := range blacklist {
		listMap[item.Id] = item
	}

	return listMap
}

func (b *DbService) QueryAllIgnoreList() map[string]bool {

	listMap := make(map[string]bool)

	var ignorelist []po.Ignorelist
	err := global.GDb.Raw("SELECT `id` FROM `ignorelist`").Scan(&ignorelist).Error
	if err != nil {
		global.GLog.Error("QueryAllIgnoreList", zap.Error(err))
		return listMap
	}

	for _, list := range ignorelist {
		listMap[strings.ToLower(list.Id)] = true
	}

	return listMap
}

func (b *DbService) AddBlack(name string, reason string) (string, error) {

	err, data := utils.CheckPlayer(name)
	if err != nil {
		return "", err
	}
	user := po.Blacklist{Id: data.PID, Name: name, Reason: reason}
	err = global.GDb.Save(&user).Error
	if err == nil {
		global.GBlackListMap[data.PID] = po.Blacklist{Id: data.PID, Name: data.Name, Reason: reason}
	}
	return data.PID, err
}

func (b *DbService) RemoveBlack(name string) error {
	err, data := utils.CheckPlayer(name)
	if err != nil {
		return err
	}
	err = global.GDb.Delete(&po.Blacklist{}, data.PID).Error
	if err == nil {
		delete(global.GBlackListMap, data.PID)
	}
	return err
}

func (b *DbService) SelectAllSensitive() []string {
	var sensitives []string
	err := global.GDb.Raw("SELECT `id` FROM `sensitive`").Scan(&sensitives).Error
	if err != nil {
		global.GLog.Error("SelectAllSensitive", zap.Error(err))
	}
	return sensitives
}

func (b *DbService) AddSensitive(word string) error {
	words := po.Sensitive{Id: word}
	return global.GDb.Save(&words).Error
}

func (b *DbService) RemoveSensitive(word string) error {
	return global.GDb.Delete(&po.Sensitive{Id: word}).Error
}

func (b *DbService) AddIgnore(id string) error {
	user := po.Ignorelist{Id: id}
	err := global.GDb.Save(&user).Error
	if err == nil {
		global.GIgnoreListMap[id] = true
	}
	return err
}

func (b *DbService) RemoveIgnore(id string) error {
	err := global.GDb.Delete(&po.Ignorelist{Id: id}).Error
	if err == nil {
		delete(global.GIgnoreListMap, id)
	}
	return err
}

func (b *DbService) DeleteAllIgnore() error {
	err := global.GDb.Exec("DELETE FROM ignorelist").Error
	if err == nil {
		global.GIgnoreListMap = make(map[string]bool)
	}
	return err
}

func (b *DbService) AddBind(qq int64, name string, pid string) error {

	bind := po.Bind{Qq: qq, Name: name, Pid: pid}

	err := global.GDb.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "Qq"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "pid"}),
	}).Create(&bind).Error

	global.GBindMap[qq] = name
	return err
}

func (b *DbService) getBindById(id int64) (error, po.Bind) {
	var user po.Bind
	result := global.GDb.First(&user, id)
	if result.RowsAffected == 0 {
		return errors.New("未绑定账号, 请使用bind=<id>绑定后再查询"), user
	}
	return result.Error, user
}

func (b *DbService) QueryAllBind() map[int64]string {

	listMap := make(map[int64]string)

	var bindArr []po.Bind
	err := global.GDb.Raw("SELECT `qq`, `name` FROM `bind`").Scan(&bindArr).Error
	if err != nil {
		global.GLog.Error("QueryAllBind", zap.Error(err))
		return listMap
	}

	for _, bind := range bindArr {
		listMap[bind.Qq] = bind.Name
	}

	return listMap
}

func (b *DbService) GetBindName(qq int64) (error, string) {
	name, ok := global.GBindMap[qq]
	if ok {
		return nil, name
	} else {
		err, bind := b.getBindById(qq)
		if err != nil {
			global.GLog.Error("b.getBindById(qq)", zap.Error(err))
			return err, ""
		} else {
			if bind.Qq == 0 {
				return errors.New("未绑定账号, 请使用bind=<id>绑定后再查询"), ""
			} else {
				return nil, bind.Name
			}
		}
	}
}

func (b *DbService) AddCardCheck(qq int64, groupId int64) error {
	now := time.Now()
	cardCheck := po.CardCheck{Qq: qq, GroupId: groupId, FailCnt: 1, NextCheckTime: now.Add(6 * time.Hour).UnixMilli()}
	return global.GDb.Save(&cardCheck).Error
}

func (b *DbService) QueryCardCheckByTime(queryTime int64) (error, []po.CardCheck) {
	var cardCheck []po.CardCheck
	err := global.GDb.Raw("SELECT * FROM `card_check` WHERE next_check_time < ?", queryTime).Scan(&cardCheck).Error
	return err, cardCheck
}

func (b *DbService) UpdateCardCheck(qq int64, cnt int, nextCheckTime int64) error {
	return global.GDb.Model(po.CardCheck{}).Where("qq = ?", qq).
		Updates(po.CardCheck{FailCnt: cnt, NextCheckTime: nextCheckTime}).Error
}

func (b *DbService) DeleteCardCheck(qq int64) error {
	return global.GDb.Delete(po.CardCheck{}, qq).Error
}
