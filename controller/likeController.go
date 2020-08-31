package controller

import (
	"fmt"
	"time"
	"vscode/go-gorm-database/dao"
)

/**
新增喜欢
**/
func SaveLike() {

	like := &dao.Like{
		ID:        1,
		Ip:        "127.0.0.2",
		Ua:        "测试数据2",
		Title:     "第一批数据2",
		Hash:      1231,
		CreatedAt: time.Now(),
	}
	dao.SaveLike(like)
}

func SelectLikeList() (likeList []dao.Like) {
	likeList, err := dao.SelectLike(&dao.Like{ID: 2})
	if err != nil {
		fmt.Println(" 查询报错 err", err)
	}
	fmt.Println("select likeList", likeList)
	return likeList
}

func DeleteLike(like *dao.Like) {
	dao.DeleteLike(like)
}
