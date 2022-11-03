package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.museum/proxy/nosql"
	"time"
)

type AreaInfo struct {
	baseInfo
	Remark   string
	Owner    string
	Parent   string
	Template string //产品配置模板
	Width    int32
	Height   int32
}

func (mine *cacheContext) CreateArea(name, remark, owner, parent, operator string) (*AreaInfo, error) {
	db := new(nosql.Area)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetAreaNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Name = name
	db.Remark = remark
	db.Owner = owner
	db.Parent = parent
	db.Width = 0
	db.Height = 0
	db.Template = ""

	err := nosql.CreateArea(db)
	if err != nil {
		return nil, err
	}
	info := new(AreaInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetArea(uid string) (*AreaInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the collective museum uid is empty")
	}
	db, err := nosql.GetArea(uid)
	if err != nil {
		return nil, err
	}
	info := new(AreaInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetAreas(parent string) []*AreaInfo {
	list := make([]*AreaInfo, 0, 20)
	array, err := nosql.GetAreasByParent(parent)
	if err != nil {
		return list
	}
	for _, item := range array {
		info := new(AreaInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetAreasByOwner(uid string) []*AreaInfo {
	array, err := nosql.GetAreasByOwner(uid)
	if err != nil {
		return make([]*AreaInfo, 0, 0)
	}
	list := make([]*AreaInfo, 0, len(array))
	for _, item := range array {
		info := new(AreaInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetAreaList(array []string) []*AreaInfo {
	if array == nil || len(array) < 1 {
		return make([]*AreaInfo, 0, 0)
	}
	list := make([]*AreaInfo, 0, len(array))
	for i := 0; i < len(array); i += 1 {
		db, err := nosql.GetArea(array[i])
		if err == nil {
			info := new(AreaInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *AreaInfo) initInfo(db *nosql.Area) {
	mine.Name = db.Name
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Remark = db.Remark
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Parent = db.Parent
	mine.Owner = db.Owner
	mine.Template = db.Template
	mine.Width = db.Width
	mine.Height = db.Height
}

func (mine *AreaInfo) UpdateBase(name, remark, operator string) error {
	err := nosql.UpdateAreaBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *AreaInfo) UpdateTemplate(template, operator string) error {
	err := nosql.UpdateAreaTemplate(mine.UID, template, operator)
	if err == nil {
		mine.Template = template
		mine.Operator = operator
	}
	return err
}

func (mine *AreaInfo) Remove(operator string) error {
	return nosql.RemoveArea(mine.UID, operator)
}
