package models

import (
	"shorturl/backend/utils"
	"time"
)

type Urlmap struct {
	ShortUrl   string    `xorm:"varchar(20) pk notnull unique 'u_shorturl'"`
	LongUrl    string    `xorm:"varchar(500) notnull 'u_longurl'"`
	IOSUrl     string    `xorm:"varchar(500) 'u_iosurl'"`
	AndroidUrl string    `xorm:"varchar(500) 'u_androidurl'"`
	CustomUrl  string    `xorm:"varchar(20) notnull 'u_customurl'"`
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
}

type UrlmapList struct {
	Urlmap     `xorm:"extends"`
	Statistics `xorm:"extends"`
}

func (UrlmapList) TableName() string {
	return "urlmap"
}

func InsertUrlmap(model *Urlmap) error {
	_, err := db.Insert(model)
	return err
}

func UpdateUrlmap(model *Urlmap) error {
	_, err := db.Where("u_shorturl = ?", model.ShortUrl).Update(model)
	return err
}

func DeleteUrlmap(shorturl string) error {
	_, err := db.Delete(&Urlmap{ShortUrl: shorturl})
	return err
}

func GetUrlmapByShortUrl(shorturl string) *Urlmap {
	var model Urlmap
	has, err := db.Where("u_shorturl = ?", shorturl).Get(&model)
	if err != nil || !has {
		return nil
	}
	return &model
}

func GetUrlmapByLongUrl(longurl string) *Urlmap {
	var model Urlmap
	has, err := db.Where("u_longurl = ?", longurl).Get(&model)
	if err != nil || !has {
		return nil
	}
	return &model
}

func GetUrlmapByCustomUrl(customurl string) *Urlmap {
	var model Urlmap
	has, err := db.Where("u_shorturl = ? or u_customurl = ?", customurl, customurl).Get(&model)
	if err != nil || !has {
		return nil
	}
	return &model
}

func GetUrlmapList(page, size int, shorturl string) utils.Page {
	list := make([]UrlmapList, 0)
	var umap UrlmapList
	qs := db.Join("INNER", "statistics", "urlmap.u_shorturl=statistics.s_shorturl").Limit(size, (page-1)*size)
	if len(shorturl) > 0 {
		qs = qs.Where("urlmap.u_shorturl = ?", shorturl)
		umap.Urlmap.ShortUrl = shorturl
	}
	count, _ := db.Count(&umap)
	err := qs.OrderBy("urlmap.created desc").Find(&list)
	if err != nil {
		panic(err)
	}
	return utils.PageUtil(int(count), page, size, list)
}
