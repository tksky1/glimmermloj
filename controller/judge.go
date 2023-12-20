package controller

import (
	"bytes"
	"fmt"
	"glimmermloj/judge"
	"glimmermloj/repository"
	"gorm.io/gorm"
	"strings"
	"time"
)

func Judge(userID int64, username string, userTarBuffer *bytes.Buffer) {
	ret := "未知错误200"

	log, err := judge.EvaluateUserCode(userTarBuffer, "glimmermloj")
	userInfo := &repository.UserInfo{
		UserID:        userID,
		NickName:      username,
		LastAccuracy:  0,
		LastSpeed:     0,
		LastTimestamp: time.Now().String(),
		LastMsg:       ret,
	}
	if err != nil && log != "timeout" {
		userInfo.LastMsg = "评测机错误201"
		println("【log】evaluator err for user", username, ":", err.Error())
	} else {
		lines := strings.Split(log, "\n")
		lastLine := ""
		if len(lines) > 2 {
			lastLine = lines[len(lines)-2]
		} else if len(lines) <= 2 {
			lastLine = lines[0]
		}
		accuracy := 0.0
		speed := 0.0
		if log == "timeout" {
			userInfo.LastMsg = "代码运行超时101"
		} else {
			_, err := fmt.Sscanf(lastLine, "%f %f", &accuracy, &speed)
			if err != nil {
				userInfo.LastMsg = "用户代码有误100"
				println("【log】user code output error for", username, ":", log)
			}
		}

		oldInfo, err := GetHighestRankBodyByUserID(repository.MLOJDB, userID)

		if err != nil {
			userInfo.LastMsg = "评测机错误202"
		} else {
			if userInfo.LastMsg == "未知错误200" {
				userInfo.LastMsg = "评测成功"
			}
			userInfo = &repository.UserInfo{
				UserID:        userID,
				NickName:      username,
				LastAccuracy:  accuracy,
				LastSpeed:     speed,
				LastTimestamp: time.Now().Format("2006-01-02 15:04:05"),
				LastMsg:       userInfo.LastMsg,
			}
			if speed > 0.0001 && (oldInfo == nil || accuracy > oldInfo.Accuracy) {
				err := UpdateOrCreateRanking(repository.MLOJDB, userID, username, accuracy, speed, userInfo.LastTimestamp)
				if err != nil {
					userInfo.LastMsg = "评测机错误203"
				}
			}
		}
	}
	err = CreateOrUpdateUserInfo(repository.MLOJDB, userInfo)
	if err != nil {
		println("database error: ", err.Error())
	}

}

func CreateOrUpdateUserInfo(db *gorm.DB, userInfo *repository.UserInfo) error {
	result := db.Create(userInfo)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return result.Error
		}
		return db.Save(userInfo).Error
	}
	return nil
}

func UpdateOrCreateRanking(db *gorm.DB, userID int64, nickName string, accuracy, speed float64, timestamp string) error {
	var ranking repository.Ranking
	err := db.First(&ranking, userID).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}

		ranking.UserID = userID
		ranking.NickName = nickName
		ranking.Accuracy = accuracy
		ranking.Speed = speed
		ranking.Timestamp = timestamp

		err = db.Create(&ranking).Error
		if err != nil {
			return err
		}
	} else {
		ranking.NickName = nickName
		ranking.Accuracy = accuracy
		ranking.Speed = speed
		ranking.Timestamp = timestamp

		err = db.Save(&ranking).Error
		if err != nil {
			return err
		}
	}

	return nil
}
