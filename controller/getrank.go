package controller

import (
	"github.com/gin-gonic/gin"
	"glimmermloj/repository"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func GetRank(c *gin.Context) {
	token, _ := c.GetQuery("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "token not found"})
		return
	}
	ch := make(chan interface{})
	go getUserinfo(ch, token)
	username1 := <-ch
	username := username1.(string)
	userID1 := <-ch
	userID, ok := userID1.(int64)
	if username == "" || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "token invalid"})
		return
	}

	highestRanking, err := getHighestRecordByUserID(repository.MLOJDB, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "get ranking db error"})
		return
	}
	userRank, err := getUserAccuracyRanking(repository.MLOJDB, userID)
	var top5 []repository.Ranking
	err = repository.MLOJDB.Order("accuracy desc").Limit(5).Find(&top5).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "get ranking db error"})
		return
	}
	retList := make([][]string, 5)
	for i, rank := range top5 {
		retList[i] = make([]string, 4)
		retList[i][0] = rank.NickName
		retList[i][1] = strconv.FormatFloat(rank.Accuracy, 'f', 5, 64)
		retList[i][2] = strconv.FormatFloat(rank.Speed, 'f', 5, 64)
		retList[i][3] = rank.Timestamp
	}
	retList = append(retList, make([]string, 4))
	if highestRanking != nil {
		retList[len(retList)-1][0] = strconv.Itoa(userRank)
		retList[len(retList)-1][1] = strconv.FormatFloat(highestRanking.Accuracy, 'f', 5, 64)
		retList[len(retList)-1][2] = strconv.FormatFloat(highestRanking.Speed, 'f', 5, 64)
		retList[len(retList)-1][3] = highestRanking.Timestamp
	}
	userInfo, err := GetUserInfoByUserID(repository.MLOJDB, userID)
	retList = append(retList, make([]string, 5))
	if userInfo != nil {
		retList[len(retList)-1][0] = "-"
		retList[len(retList)-1][1] = strconv.FormatFloat(userInfo.LastAccuracy, 'f', 5, 64)
		retList[len(retList)-1][2] = strconv.FormatFloat(userInfo.LastSpeed, 'f', 5, 64)
		retList[len(retList)-1][3] = userInfo.LastTimestamp
		if userInfo.LastMsg != "评测成功" {
			retList[len(retList)-1][4] = userInfo.LastMsg
		}
	}
	c.JSON(200, gin.H{"rank": retList})
}

func getHighestRecordByUserID(db *gorm.DB, userID int64) (*repository.Ranking, error) {
	var ranking repository.Ranking
	err := db.Where("user_id = ?", userID).First(&ranking).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 找不到记录时返回nil，而不是错误
		}
		return nil, err
	}
	return &ranking, nil
}

func getUserAccuracyRanking(db *gorm.DB, userID int64) (int, error) {
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM rankings WHERE accuracy > (SELECT accuracy FROM rankings WHERE user_id = ?)", userID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count) + 1, nil
}

func GetUserInfoByUserID(db *gorm.DB, userID int64) (*repository.UserInfo, error) {
	var userInfo repository.UserInfo
	err := db.Where("user_id = ?", userID).First(&userInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 找不到记录时返回nil，而不是错误
		}
		return nil, err
	}
	return &userInfo, nil
}

func GetHighestRankBodyByUserID(db *gorm.DB, userID int64) (*repository.Ranking, error) {
	var userInfo repository.Ranking
	err := db.Where("user_id = ?", userID).First(&userInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 找不到记录时返回nil，而不是错误
		}
		return nil, err
	}
	return &userInfo, nil
}
