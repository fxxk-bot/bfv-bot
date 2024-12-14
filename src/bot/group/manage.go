package group

import (
	"bfv-bot/common/des"
	"bfv-bot/common/global"
	"bfv-bot/common/http"
	"bfv-bot/model/dto"
	"errors"
	"go.uber.org/zap"
)

func SetGroupKick(groupId int64, userId int64) {
	data := map[string]interface{}{
		"group_id":           groupId,
		"user_id":            userId,
		"reject_add_request": false,
	}
	_, _ = http.Post(global.GConfig.QQBot.Address+"/set_group_kick_async", data)
}

func GetGroupMemberInfo(groupId int64, userId int64) (error, dto.GetGroupMemberInfoData) {
	data := map[string]interface{}{
		"group_id": groupId,
		"user_id":  userId,
		"no_cache": true,
	}
	post, err := http.Post(global.GConfig.QQBot.Address+"/get_group_member_info", data)
	if err != nil {
		global.GLog.Error("get_group_member_info", zap.Error(err))
		return err, dto.GetGroupMemberInfoData{}
	}
	var result dto.GetGroupMemberInfoResp
	err = des.StringToStruct(post, &result)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return err, dto.GetGroupMemberInfoData{}
	}
	if result.Status != "ok" || result.Retcode != 0 {
		global.GLog.Error("get_group_member_info", zap.String("message", result.Message),
			zap.String("wording", result.Wording))
		return errors.New("bot接口异常"), dto.GetGroupMemberInfoData{}
	}
	return nil, result.Data
}

func SetGroupWholeBan(groupId int64, enable bool) {
	data := map[string]interface{}{
		"group_id": groupId,
		"enable":   enable,
	}
	_, _ = http.Post(global.GConfig.QQBot.Address+"/set_group_whole_ban_async", data)
}

func GetGroupMemberList(groupId int64) (error, []dto.GetGroupMemberListData) {
	data := map[string]interface{}{
		"group_id": groupId,
		"no_cache": true,
	}
	post, err := http.Post(global.GConfig.QQBot.Address+"/get_group_member_list", data)
	if err != nil {
		global.GLog.Error("get_group_member_list", zap.Error(err))
		return err, []dto.GetGroupMemberListData{}
	}
	var result dto.GetGroupMemberListResp
	err = des.StringToStruct(post, &result)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return err, []dto.GetGroupMemberListData{}
	}
	if result.Status != "ok" || result.Retcode != 0 {
		global.GLog.Error("get_group_member_list", zap.String("message", result.Message),
			zap.String("wording", result.Wording))

		return errors.New("bot接口异常"), []dto.GetGroupMemberListData{}
	}
	return nil, result.Data
}

// GetActiveGroupMemberCardMap 获取所有启用机器人服务群的群成员名片
func GetActiveGroupMemberCardMap() map[string]bool {
	memberMap := make(map[string]bool)
	for _, item := range global.GConfig.QQBot.ActiveGroup {
		err, memberList := GetGroupMemberList(item)
		if err != nil {
			continue
		}
		for _, member := range memberList {
			if member.Card == "" {
				continue
			}
			memberMap[member.Card] = true
		}
	}
	return memberMap
}
