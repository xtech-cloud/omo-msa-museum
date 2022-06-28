package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-museum/proto/museum"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.museum/proxy"
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
	Paths []*proxy.PathInfo
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
	db.Paths = make([]*proxy.PathInfo, 0, 1)
	db.Tags = make([]string, 0, 1)
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
	mine.Operator = db.Operator
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
	mine.Paths = db.Paths
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

func (mine *SandtableInfo) UpdateBGM(operator, asset string) error {
	err := nosql.UpdateSandtableBGM(mine.UID, asset, operator)
	if err == nil {
		mine.BGM = asset
		mine.Operator = operator
	}
	return err
}

func (mine *SandtableInfo) UpdateNarrate(operator, asset string) error {
	err := nosql.UpdateSandtableNarrate(mine.UID, asset, operator)
	if err == nil {
		mine.Narrate = asset
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

func (mine *SandtableInfo) UpdatePath(operator, path, name, color string, points []*pb.PathKeyInfo) error {
	if len(path) > 0 {
		er := mine.RemovePath(operator, path)
		if er != nil {
			return er
		}
	}
	if len(name) < 1 && len(points) == 0{
		return nil
	}
	return mine.CreatePath(operator, path, name, color, points)
}

func (mine *SandtableInfo) CreatePath(operator, path, name, color string, points []*pb.PathKeyInfo) error {
	info := new(proxy.PathInfo)
	if len(path) < 1 {
		info.UID = primitive.NewObjectID().Hex()
	}else{
		info.UID = path
	}
	info.Name = name
	info.Color = color
	info.Points = make([]proxy.FrameKeyInfo, 0, len(points))
	for _, item := range points {
		pos := SwitchVector(item.Position)
		ro := SwitchVector(item.Rotation)
		info.Points = append(info.Points, proxy.FrameKeyInfo{Key: item.Key, Scale: item.Scale, Position: pos, Rotation: ro})
	}
	err := nosql.AppendSandtablePath(mine.UID, info)
	if err == nil {
		mine.Paths = append(mine.Paths, info)
		mine.Operator = operator
	}
	return err
}

func (mine *SandtableInfo) RemovePath(operator, path string) error {
	err := nosql.SubtractSandtablePath(mine.UID, path)
	if err == nil {
		mine.Operator = operator
		for i, info := range mine.Paths {
			if info.UID == path {
				if i == len(mine.Paths) - 1 {
					mine.Paths = append(mine.Paths[:i])
				}else{
					mine.Paths = append(mine.Paths[:i], mine.Paths[i+1:]...)
				}
				break
			}
		}
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
