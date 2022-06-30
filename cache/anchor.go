package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.museum/proxy"
	"omo.msa.museum/proxy/nosql"
	"time"
)

// 沙盘锚点
type AnchorInfo struct {
	baseInfo
	Remark   string
	Cover    string
	Parent string
	Panorama string
	Link string
	Owner string
	Position proxy.VectorInfo
	Tags    []string
}

func (mine *cacheContext) CreateAnchor(name, remark, owner, parent, operator string, tags []string) (*AnchorInfo, error) {
	db := new(nosql.Anchor)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetAnchorNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Name = name
	db.Remark = remark
	db.Owner = owner
	db.Parent = parent
	db.Cover = ""
	db.Tags = tags
	err := nosql.CreateAnchor(db)
	if err != nil {
		return nil, err
	}
	info := new(AnchorInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetAnchor(uid string) (*AnchorInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the museum uid is empty")
	}
	db, err := nosql.GetAnchor(uid)
	if err != nil {
		return nil, err
	}
	info := new(AnchorInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetAnchorsByParent(uid string) ([]*AnchorInfo,error) {
	array, err := nosql.GetAnchorsByParent(uid)
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

func (mine *AnchorInfo) initInfo(db *nosql.Anchor) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.UpdateTime = db.UpdatedTime
	mine.Name = db.Name
	mine.Remark = db.Remark
	mine.Cover = db.Cover
	mine.Link = db.Link
	mine.Owner = db.Owner
	mine.Parent = db.Parent
	mine.Panorama = db.Panorama
	mine.Position = db.Position
	mine.Tags = db.Tags
}

func (mine *AnchorInfo) UpdateBase(name, remark, operator string) error {
	err := nosql.UpdateAnchorBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *AnchorInfo) Remove(operator string) error {
	return nosql.RemoveAnchor(mine.UID, operator)
}

func (mine *AnchorInfo) UpdateCover(cover, operator string) error {
	err := nosql.UpdateAnchorCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
	}
	return err
}

func (mine *AnchorInfo) UpdatePosition(operator string, pos proxy.VectorInfo) error {
	err := nosql.UpdateAnchorPosition(mine.UID, operator, pos)
	if err == nil {
		mine.Position = pos
		mine.Operator = operator
	}
	return err
}
