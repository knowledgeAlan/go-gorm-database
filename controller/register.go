package controller

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"vscode/go-gorm-database/authorization"
	"vscode/go-gorm-database/dao"
	"vscode/go-gorm-database/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	dao.SetKey("golang", "hello golang")
	c.JSON(http.StatusOK, dao.GetKey("golang"))

}

func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json provided")
		return
	}

	token, err := authorization.CreateToken(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, token)
}

func Login01(c *gin.Context) {

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json provided")
		return
	}

	ts, err := authorization.CreateTokenAndRefresh(user.ID)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := dao.CreateAuth(user.ID, ts)

	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}

	tokens := map[string]interface{}{
		"accessToken":  ts.AccessToken,
		"refreshToken": ts.RefreshToken,
	}

	c.JSON(http.StatusOK, tokens)

}

func CreateTodo(c *gin.Context) {

	var td *models.Todo

	if err := c.ShouldBindJSON(&td); err != nil {

		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}

	tokenAuth, err := authorization.ExtractTokenMetaData(c.Request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "unauthorized")
		return
	}

	userId, err := dao.FetchAuth(tokenAuth)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "unauthorized")
		return
	}

	td.UserId = userId

	c.JSON(http.StatusCreated, td)

}

//退出登录
func LoginOut(c *gin.Context) {
	au, err := authorization.ExtractTokenMetaData(c.Request)

	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	deleted, delErr := dao.DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	c.JSON(http.StatusOK, "successfully logged out")
}

//刷新token
func RefreshToken(c *gin.Context) {

	mapToken := make(map[string]string)

	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	refreshToken := mapToken["refreshToken"]

	os.Setenv("REFRESH_SECRET", "rrrr")

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method :%v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, "refresh token expired")
		return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {

		refreshUuid, ok := claims["refreshUiid"].(string)
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["userId"]), 10, 64)

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "error occured")
			return
		}

		deleted, delErr := dao.DeleteAuth(refreshUuid)

		if delErr != nil || deleted == 0 {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}

		tk, _ := authorization.CreateToken(userId)

		// if createErr != nil {
		// 	c.JSON(http.StatusForbidden, createErr.Error())
		// 	return
		// }

		// saveErr := dao.CreateAuth(userId, ts)

		c.JSON(http.StatusOK, tk)

	}

}
