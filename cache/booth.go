package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.museum/proxy"
	"omo.msa.museum/proxy/nosql"
	"time"
)

type BoothInfo struct {
	baseInfo
	Remark  string
	Parent string
	Owner   string
	Exhibits []string
	Position proxy.VectorInfo
}

func (mine *cacheContext) CreateBooth(name, remark, owner, parent, operator string) (*BoothInfo, error) {
	db := new(nosql.Booth)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetBoothNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Name = name
	db.Remark = remark
	db.Exhibits = make([]string, 0, 1)
	db.Parent = parent
	db.Owner = owner
	if owner == "" {
		owner = DefaultOwner
	}
	err := nosql.CreateBooth(db)
	if err != nil {
		return nil, err
	}
	info := new(BoothInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetBooth(uid string) (*BoothInfo, error) {
	db, err := nosql.GetBooth(uid)
	if err != nil {
		return nil, err
	}
	info := new(BoothInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) RemoveBooth(uid, operator string) error {
	if uid == "" {
		return errors.New("the panorama uid is empty")
	}
	err := nosql.RemoveBooth(uid, operator)
	return err
}

func (mine *cacheContext) GetBoothsByOwner(owner string) ([]*BoothInfo, error) {
	if owner == "" {
		owner = DefaultOwner
	}
	array, err := nosql.GetAllBoothsByOwner(owner)
	if err != nil {
		return make([]*BoothInfo, 0, 1), err
	}
	list := make([]*BoothInfo, 0, len(array))
	for _, item := range array {
		info := new(BoothInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list, nil
}

func (mine *cacheContext) GetBoothsByParent(parent string) ([]*BoothInfo, error) {
	array, err := nosql.GetAllBoothsByParent(parent)
	if err != nil {
		return make([]*BoothInfo, 0, 1), err
	}
	list := make([]*BoothInfo, 0, len(array))
	for _, item := range array {
		info := new(BoothInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list, nil
}

func (mine *BoothInfo) initInfo(db *nosql.Booth) {
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
	mine.Exhibits = db.Exhibits
	mine.Position = db.Position
}

func (mine *BoothInfo) UpdateExhibit(operator string, arr []string) error {
	err := nosql.UpdateBoothExhibit(mine.UID, operator, arr)
	if err == nil {
		mine.Exhibits = arr
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *BoothInfo) UpdateBase(name, remark, operator string) error {
	err := nosql.UpdateBoothBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *BoothInfo) UpdatePosition(operator string, pos proxy.VectorInfo) error {
	err := nosql.UpdateBoothPosition(mine.UID, operator, pos)
	if err == nil {
		mine.Position = pos
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}
	return err
}
