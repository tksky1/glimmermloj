package repository

type UserInfo struct {
	UserID        int64   `gorm:"primaryKey"`
	NickName      string  `gorm:"size:255"`
	LastAccuracy  float64 `gorm:"type:double"`
	LastSpeed     float64 `gorm:"type:double"`
	LastTimestamp string  `gorm:"size:255"`
	LastMsg       string  `gorm:"size:255"`
}

type Ranking struct {
	UserID    int64   `gorm:"primaryKey"`
	NickName  string  `gorm:"size:255"`
	Accuracy  float64 `gorm:"type:double"`
	Speed     float64 `gorm:"type:double"`
	Timestamp string  `gorm:"size:255"`
}
