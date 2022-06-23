package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.museum/proxy"
	"time"
)

//沙盘
type Sandtable struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deletedAt" bson:"deletedAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Status uint8 `json:"status" bson:"status"`
	Name   string `json:"name" bson:"name"`
	Remark string `json:"remark" bson:"remark"`
	Owner  string `json:"owner" bson:"owner"`
	Background string `json:"background" bson:"background"`
	Mask string `json:"mask" bson:"mask"`
	Width uint32 `json:"width" bson:"width"`
	Height uint32 `json:"height" bson:"height"`
	Narrate string `json:"narrate" bson:"narrate"`
	BGM string `json:"bgm" bson:"bgm"`
	Path []*proxy.FrameKeyInfo `json:"path" bson:"path"`
	Tags []string `json:"tags" bson:"tags"`
}

func CreateSandtable(info *Sandtable) error {
	_, err := insertOne(TableSandtable, info)
	if err != nil {
		return err
	}
	return nil
}

func GetSandtableNextID() uint64 {
	num, _ := getSequenceNext(TableSandtable)
	return num
}

func GetSandtable(uid string) (*Sandtable, error) {
	result, err := findOne(TableSandtable, uid)
	if err != nil {
		return nil, err
	}
	model := new(Sandtable)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetSandtablesByOwner(owner string) ([]*Sandtable, error) {
	msg := bson.M{"owner": owner, "deletedAt": new(time.Time)}
	cursor, err1 := findMany(TableSandtable, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Sandtable, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Sandtable)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateSandtableBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableSandtable, uid, msg)
	return err
}

func UpdateSandtableBG(uid, cover, operator string, width, height uint32) error {
	msg := bson.M{"background": cover, "width":width, "height":height, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableSandtable, uid, msg)
	return err
}

func UpdateSandtableStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableSandtable, uid, msg)
	return err
}

func UpdateSandtablePath(uid, operator string, path []*proxy.FrameKeyInfo) error {
	msg := bson.M{"path": path, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableSandtable, uid, msg)
	return err
}

func RemoveSandtable(uid, operator string) error {
	_, err := removeOne(TableSandtable, uid, operator)
	return err
}

func UpdateSandtableBGM(uid, bgm, operator string) error {
	msg := bson.M{"bgm": bgm, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(TableSandtable, uid, msg)
	return err
}
