package service

import (
	"fmt"
	"gin_web/dao/mysql/communityDao"
	"gin_web/dao/mysql/postDao"
	"gin_web/dao/mysql/userDao"
	"gin_web/dao/redis"
	"gin_web/model"
	"strconv"
)

// releasePost 新建帖子
func ReleasePost(post *model.PostMsg) (int64, error) {
	id, err := postDao.ReleasePost(post)
	if err != nil {
		return id, err
	}
	fmt.Println("运行了")
	//将创建时间添加到redis
	err = redis.AddPortCreatTime(strconv.FormatInt(post.PostID, 10), post.CommunityID)
	if err != nil {
		return id, err
	}
	return id, nil
}

// ReceivePostById 根据id获取帖子相关的信息
func ReceivePostById(id int64) (post model.PostDetail, err error) {

	p, err := postDao.GetPostById(id)
	if err != nil {
		return post, err
	}

	c, err := communityDao.GetIntroductionById(p.CommunityID)
	if err != nil {
		return post, err
	}

	u, err := userDao.FindUserById(p.AuthorID)
	if err != nil {
		return post, err
	}
	post.PostMsg = p
	post.Community = c
	post.User = u
	return
}

// GetPostList 获取帖子列表 分页查询
func GetPostList(page int64, size int64) (posts []*model.PostDetail, err error) {
	postMeg, err := postDao.GetPostList(page, size)
	if err != nil {
		return posts, err
	}
	posts = make([]*model.PostDetail, 0, len(postMeg))

	for _, p := range postMeg {
		c, err := communityDao.GetIntroductionById(p.CommunityID)
		if err != nil {
			continue
		}
		u, err := userDao.FindUserById(p.AuthorID)
		if err != nil {
			continue
		}
		m := new(model.PostDetail)
		m.PostMsg = p
		m.Community = c
		m.User = u
		posts = append(posts, m)
	}
	return
}

func PostListOrder(p *model.ParamPostList) (posts []*model.PostDetail, err error) {
	if p.CommunityID == 0 {
		return GetPostListOrder(p)
	} else {
		return GetPostListOrderByCommunity(p)
	}

}

func GetPostListOrder(p *model.ParamPostList) (posts []*model.PostDetail, err error) {
	ids, err := redis.GetPostListOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		return
	}
	//根据redis得到的id从数据库查询帖子
	postMeg, err := postDao.GetPostListById(ids)
	if err != nil {
		return
	}

	//查询帖子的点赞数量
	voteDataList, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//查询其他社区作者详细信息
	posts = make([]*model.PostDetail, 0, len(postMeg))
	for i, p := range postMeg {
		c, err := communityDao.GetIntroductionById(p.CommunityID)
		if err != nil {
			continue
		}
		u, err := userDao.FindUserById(p.AuthorID)
		if err != nil {
			continue
		}
		m := new(model.PostDetail)
		m.PostMsg = p
		m.Community = c
		m.User = u
		m.VoteData = voteDataList[i]
		posts = append(posts, m)
	}
	return
}

func GetPostListOrderByCommunity(p *model.ParamPostList) (posts []*model.PostDetail, err error) {
	ids, err := redis.GetPostListOrderByCommunity(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		return
	}
	//根据redis得到的id从数据库查询帖子
	postMeg, err := postDao.GetPostListById(ids)
	if err != nil {
		return
	}

	//查询帖子的点赞数量
	voteDataList, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//查询其他社区作者详细信息
	posts = make([]*model.PostDetail, 0, len(postMeg))
	for i, p := range postMeg {
		c, err := communityDao.GetIntroductionById(p.CommunityID)
		if err != nil {
			continue
		}
		u, err := userDao.FindUserById(p.AuthorID)
		if err != nil {
			continue
		}
		m := new(model.PostDetail)
		m.PostMsg = p
		m.Community = c
		m.User = u
		m.VoteData = voteDataList[i]
		posts = append(posts, m)
	}

	return
}
