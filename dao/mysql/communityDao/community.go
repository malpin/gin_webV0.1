package communityDao

import (
	"database/sql"
	"fmt"
	"gin_web/Bean"
	"gin_web/dao/mysql"
	"gin_web/model"
	"go.uber.org/zap"
)

// GetCommunityList 查询所有的社区列表
func GetCommunityList() (communityList []model.Community, err error) {
	sqlStr := "select community_id,community_name from community"
	err = mysql.MysqlDB.Select(&communityList, sqlStr)
	fmt.Println(communityList)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("社区列表是空")
			err = nil
		} else {
			return nil, Bean.SYSTEM_BUSY.MarkError
		}
	}
	return communityList, nil
}

// GetIntroductionById 根据id查找社区详情
func GetIntroductionById(communityID int) (community *model.Community, err error) {
	community = new(model.Community)
	sqlStr := "select community_id,community_name,introduction,create_time from community where community_id=?"
	err = mysql.MysqlDB.Get(community, sqlStr, communityID)
	if community.CommunityId == 0 {
		zap.L().Warn("GetIntroductionById 根据id查找社区详情,ID出错了", zap.Error(err))
		return community, Bean.COMMUNITY_ID_ERROR.MarkError
	}
	if err != nil {
		zap.L().Warn("GetIntroductionById 根据id查找社区详情,查询数据库出错了", zap.Error(err))
		return community, Bean.SYSTEM_BUSY.MarkError
	}
	return community, nil
}
