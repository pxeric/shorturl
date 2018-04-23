package models

type Statistics struct {
	Url          string `xorm:"varchar(20) pk notnull unique 's_shorturl'"`
	LongTotal    int    `xorm:"default 0 notnull 's_longtotal'"`
	IOSTotal     int    `xorm:"default 0 notnull 's_iostotal'"`
	AndroidTotal int    `xorm:"default 0 notnull 's_androidtotal'"`
}

func InsertStatistics(model *Statistics) error {
	_, err := db.Insert(model)
	return err
}

func AddTotal(url, platform string, total int) error {
	col := "s_longtotal"
	if platform == "ios" {
		col = "s_iostotal"
	} else if platform == "android" {
		col = "s_androidtotal"
	}
	var statics Statistics
	_, err := db.ID(url).Incr(col, total).Update(&statics)
	return err
}

func GetStatisticsById(url string) *Statistics {
	var model Statistics
	has, err := db.ID(url).Get(&model)
	if err != nil || !has {
		return nil
	}
	return &model
}

func GetStatisticsList(page, size int, url string) []Statistics {
	list := make([]Statistics, 0)
	err := db.Where("s_shorturl = ?", url).Limit(size, page*size).Find(&list)
	if err != nil {
		panic(err)
	}
	return list
}
