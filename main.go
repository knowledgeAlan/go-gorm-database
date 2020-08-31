package main

import (
	"vscode/go-gorm-database/dao"
	"vscode/go-gorm-database/routers"
)

func main() {

	dao.InitMySql()
	// controller.SaveLike()
	// controller.SelectLikeList()
	// like := &dao.Like{
	// 	ID: 1,
	// }
	// controller.DeleteLike(like)
	dao.InitRedisClient()
	router := routers.SetRouter()
	router.Run(":8090")
}
