package controller

import (
	"gin_web/Bean"
	"gin_web/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 获取所有的社区列表
func GetCommunityList(c *gin.Context) {
	//查询所有社区,返回 id 名称 简介
	list, err := service.GetCommunityList()
	if err != nil {
		Bean.Error(c, Bean.SYSTEM_BUSY)
		return
	}
	Bean.Success(c, list)
}

// GetIntroductionById 根据id查找社区详情
func GetIntroductionById(c *gin.Context) {
	id := c.Param("id")
	communityid, err := strconv.Atoi(id)
	if err != nil {
		Bean.Error(c, Bean.SYSTEM_BUSY)
		return
	}
	community, err := service.GetIntroductionById(communityid)
	if err != nil {
		if err == Bean.COMMUNITY_ID_ERROR.MarkError {
			Bean.Success(c, nil)
			return
		}
		Bean.Error(c, Bean.SYSTEM_BUSY)
		return
	}
	Bean.Success(c, community)
}
