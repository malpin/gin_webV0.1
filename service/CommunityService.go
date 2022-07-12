package service

import (
	"gin_web/dao/mysql/communityDao"
	"gin_web/model"
)

// GetCommunityList 查询所有的社区列表
func GetCommunityList() (communityList []model.Community, err error) {
	communityList, err = communityDao.GetCommunityList()
	if err != nil {
		return nil, err
	}
	return communityList, nil
}

// GetIntroductionById 根据id查找社区详情
func GetIntroductionById(communityID int) (community *model.Community, err error) {
	community, err = communityDao.GetIntroductionById(communityID)
	if err != nil {
		return community, err
	}
	return community, nil
}
