package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.museum/proxy/nosql"
	"time"
)

const (
	SandtableIdle uint8 = 0
	SandtableFroze uint8 = 1
)

// 沙盘
type SandtableInfo struct {
	Status uint8
	baseInfo
	Remark   string
	Owner    string
	Background string
	Mask string
	Width uint32
	Height uint32
	Narrate string
	BGM string
	Tags     []string
}

func (mine *cacheContext) CreateSandtable(name, remark, bg, owner, operator string, w, h uint32) (*SandtableInfo, error) {
	db := new(nosql.Sandtable)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetSandtableNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Name = name
	db.Remark = remark
	db.Width = w
	db.Height = h
	db.Background = bg
	db.Owner = owner
	db.Status = SandtableIdle
	err := nosql.CreateSandtable(db)
	if err != nil {
		return nil, err
	}
	info := new(SandtableInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetSandtable(uid string) (*SandtableInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the photocopy uid is empty")
	}
	db, err := nosql.GetSandtable(uid)
	if err != nil {
		return nil, err
	}
	info := new(SandtableInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetSandTablesByOwner(owner string) ([]*SandtableInfo, error) {
	list := make([]*SandtableInfo, 0, 20)
	array, err := nosql.GetSandtablesByOwner(owner)
	if err != nil {
		return list, err
	}
	for _, item := range array {
		info := new(SandtableInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list, nil
}

func (mine *SandtableInfo) initInfo(db *nosql.Sandtable) {
	mine.Name = db.Name
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Remark = db.Remark
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Status = db.Status
	mine.Background = db.Background
	mine.Width = db.Width
	mine.Tags = db.Tags
	mine.Height = db.Height
	mine.Narrate = db.Narrate
	mine.BGM = db.BGM
	mine.Mask = db.Mask
	mine.Owner = db.Owner
	mine.Tags = db.Tags
}

func (mine *SandtableInfo) UpdateBase(name, remark, operator string) error {
	err := nosql.UpdateSandtableBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *SandtableInfo) UpdateBackground(asset, operator string, width, height uint32) error {
	err := nosql.UpdateSandtableBG(mine.UID, asset, operator, width, height)
	if err == nil {
		mine.Background = asset
		mine.Width = width
		mine.Height = height
		mine.Operator = operator
	}
	return err
}

func (mine *SandtableInfo) Remove(operator string) error {
	return nosql.RemoveSandtable(mine.UID, operator)
}

func (mine *SandtableInfo) GetAnchors() ([]*AnchorInfo,error) {
	array, err := nosql.GetAnchorsByParent(mine.UID)
	if err != nil {
		return nil,err
	}
	list := make([]*AnchorInfo, 0, 20)
	for _, item := range array {
		info := new(AnchorInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list,nil
}
