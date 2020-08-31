package authorization

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"vscode/go-gorm-database/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

//通过用户信息获取token
func CreateToken(userId int64) (string, error) {

	var err error
	os.Setenv("ACCESS_SECRET", "KKKK")

	atCliams := jwt.MapClaims{}

	atCliams["authorized"] = true
	atCliams["userId"] = userId
	atCliams["exp"] = time.Now().Add(time.Minute * 15).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atCliams)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
		return "", err
	}

	return token, err
}

//创建token和刷新token
func CreateTokenAndRefresh(userId int64) (*models.TokenDetails, error) {

	td := &models.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error

	os.Setenv("ACCESS_SECRET_rt", "fhdsjkfhdsjkfhdjs")
	atCliams := jwt.MapClaims{}
	atCliams["authorized"] = true
	atCliams["accessUuid"] = td.AccessUuid
	atCliams["userId"] = userId
	atCliams["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atCliams)
	as := os.Getenv("ACCESS_SECRET_rt")
	td.AccessToken, err = at.SignedString([]byte(as))

	if err != nil {
	}

	os.Setenv("ACCESS_SECRET_rt1", "rrrr")
	rtClaims := jwt.MapClaims{}

	rtClaims["refreshUiid"] = td.RefreshUuid
	rtClaims["userId"] = userId
	rtClaims["exp"] = td.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("ACCESS_SECRET_rt1")))

	if err != nil {
		return nil, err
	}
	return td, err
}

func ExtractToken(r *http.Request) string {

	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")

	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

//校验token
func VerifyToken(r *http.Request) (*jwt.Token, error) {

	tokenStr := ExtractToken(r)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method :%v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET_rt")), nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

//获取jwt原信息
func ExtractTokenMetaData(r *http.Request) (*models.AccessDetails, error) {

	token, err := VerifyToken(r)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		accessUuid, ok := claims["accessUuid"].(string)
		if !ok {
			return nil, err
		}

		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["userId"]), 10, 64)

		if err != nil {
			return nil, err
		}

		return &models.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}
