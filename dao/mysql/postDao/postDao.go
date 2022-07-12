package postDao

import (
	"fmt"
	"gin_web/Bean"
	"gin_web/dao/mysql"
	"gin_web/model"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
)

// AddPost 新建帖子
func ReleasePost(post *model.PostMsg) (int64, error) {
	sqlStr := "insert into post(post_id,title,content,author_id,community_id,status,create_time,update_time) values(?,?,?,?,?,?,?,?)"
	result, err := mysql.MysqlDB.Exec(sqlStr, post.PostID, post.Title, post.Content, post.AuthorID, post.CommunityID, post.Status, post.CreateTime, post.UpdateTime)
	if err != nil {
		zap.L().Warn("userDao AddUser 新建帖子执行 失败了", zap.Error(err))
		return -1, Bean.SYSTEM_BUSY.MarkError
	}
	//插入成功会返回自增的id
	id, err := result.LastInsertId()
	if err != nil {
		zap.L().Warn(fmt.Sprintf("userDao AddUser 新建帖子 失败了 ,帖子id为:%d", post.PostID), zap.Error(err))
		return -1, Bean.SYSTEM_BUSY.MarkError
	}
	zap.L().Info(fmt.Sprintf("userDao AddUser 新建帖子 成功 ,帖子id为:%d", post.PostID), zap.Error(err))
	return id, nil
}

// GetPostById 根据ID获取单条帖子的详情
func GetPostById(id int64) (post *model.PostMsg, err error) {
	fmt.Println(id)
	post = new(model.PostMsg)
	strSql := "select post_id,title,content,author_id,community_id,status,create_time,update_time from post where  post_id=?"
	err = mysql.MysqlDB.Get(post, strSql, id)
	if post.PostID == 0 {
		zap.L().Warn("GetPostById 根据id获取帖子的详情,ID出错了", zap.Error(err))
		return post, Bean.POST_ID_ERROR.MarkError
	}
	if err != nil {
		zap.L().Warn("GetPostById 根据id获取帖子的详情,数据库查询出错了", zap.Error(err))
		return post, Bean.SYSTEM_BUSY.MarkError
	}
	return post, nil
}

// GetPostList 获取帖子列表 分页查询
func GetPostList(page int64, size int64) (posts []*model.PostMsg, err error) {
	posts = make([]*model.PostMsg, 0, 2)
	strSql := "select post_id,title,content,author_id,community_id,status,create_time,update_time from post order by create_time desc limit ?,?"
	err = mysql.MysqlDB.Select(&posts, strSql, page, size)
	fmt.Println("------", posts)
	if err != nil {
		zap.L().Warn("GetPostList 获取帖子列表 分页查询出错了", zap.Error(err))
		return posts, Bean.SYSTEM_BUSY.MarkError
	}
	return posts, nil
}

// GetPostListById 根据id切片查询帖子
func GetPostListById(ids []string) (postList []*model.PostMsg, err error) {
	strSql := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id, ?)
	`
	query, args, err := sqlx.In(strSql, ids, strings.Join(ids, ","))
	if err != nil {
		zap.L().Warn("GetPostList 获取帖子列表 分页查询出错了", zap.Error(err))
		return
	}

	rebind := mysql.MysqlDB.Rebind(query)
	err = mysql.MysqlDB.Select(&postList, rebind, args...)
	if err != nil {
		return
	}
	return
}
