package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-museum/proto/museum"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.museum/proxy"
	"omo.msa.museum/proxy/nosql"
	"time"
)

type ExhibitInfo struct {
	Status uint8
	Type   uint8
	baseInfo
	SN string
	Entity  string
	Owner  string
	Tags   []string
	Locals []*proxy.LocalInfo
	Specials []*proxy.SpecialInfo
}

func (mine *cacheContext) CreateExhibit(name, remark, sn, entity, owner, operator string) (*ExhibitInfo, error) {
	db := new(nosql.Exhibit)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetExhibitNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Name = name
	db.Entity = entity
	db.Owner = owner
	db.SN = sn
	db.Locals = make([]*proxy.LocalInfo, 0, 1)
	db.Locals = append(db.Locals, &proxy.LocalInfo{Language: "zh", Name: name, Remark: remark})
	db.Tags = make([]string, 0, 1)
	db.Specials = make([]*proxy.SpecialInfo, 0, 1)
	err := nosql.CreateExhibit(db)
	if err != nil {
		return nil, err
	}
	info := new(ExhibitInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) HadExhibitByName(name string) bool {
	db, err := nosql.GetExhibitByName(name)
	if err != nil {
		return false
	}
	if db == nil {
		return false
	} else {
		return true
	}
}

func (mine *cacheContext) GetExhibit(uid string) (*ExhibitInfo, error) {
	db, err := nosql.GetExhibit(uid)
	if err != nil {
		return nil, err
	}
	info := new(ExhibitInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetExhibitBySN(owner, sn string) (*ExhibitInfo, error) {
	db, err := nosql.GetExhibitBySN(owner, sn)
	if err != nil {
		return nil, err
	}
	info := new(ExhibitInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) RemoveExhibit(uid, operator string) error {
	if uid == "" {
		return errors.New("the Exhibit uid is empty")
	}
	err := nosql.RemoveExhibit(uid, operator)
	return err
}

func (mine *cacheContext) GetExhibits(array []string) []*ExhibitInfo {
	if array == nil {
		return make([]*ExhibitInfo, 0, 1)
	}
	list := make([]*ExhibitInfo, 0, len(array))
	for _, item := range array {
		info, _ := mine.GetExhibit(item)
		if info != nil {
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetExhibitsByOwner(owner string) ([]*ExhibitInfo, error) {
	if owner == "" {
		owner = DefaultOwner
	}
	array, err := nosql.GetAllExhibitsByOwner(owner)
	if err != nil {
		return nil, err
	}
	list := make([]*ExhibitInfo, 0, len(array))
	for _, item := range array {
		info := new(ExhibitInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list, nil
}

func (mine *ExhibitInfo) initInfo(db *nosql.Exhibit) {
	mine.UID = db.UID.Hex()
	mine.Name = db.Name
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.SN = db.SN
	mine.Owner = db.Owner
	mine.Entity = db.Entity
	mine.Specials = db.Specials
	mine.Locals = db.Locals
	mine.Tags = db.Tags
}

func (mine *ExhibitInfo) UpdateSN(sn, operator string) error {
	err := nosql.UpdateExhibitSN(mine.UID, sn, operator)
	if err == nil {
		mine.SN = sn
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *ExhibitInfo) UpdateLocals(operator string, list []*pb.LocalInfo) error {
	arr := make([]*proxy.LocalInfo, 0, len(list))
	for _, info := range arr {
		arr = append(arr, &proxy.LocalInfo{Language: info.Language, Name: info.Name, Remark: info.Remark})
	}
	err := nosql.UpdateExhibitLocals(mine.UID, operator, arr)
	if err == nil {
		mine.Locals = arr
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *ExhibitInfo) UpdateSpecials(operator string, list []*pb.SpecialInfo) error {
	arr := make([]*proxy.SpecialInfo, 0, len(list))
	for _, info := range arr {
		arr = append(arr, &proxy.SpecialInfo{ID: info.ID, Key: info.Key, Value: info.Value})
	}
	err := nosql.UpdateExhibitSpecials(mine.UID, operator, arr)
	if err == nil {
		mine.Specials = arr
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}
	return err
}

