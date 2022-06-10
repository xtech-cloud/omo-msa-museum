package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.museum/proxy"
	"time"
)

type Booth struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deletedAt" bson:"deletedAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name    string `json:"name" bson:"name"`
	Remark  string `json:"remark" bson:"remark"`
	Exhibit   string `json:"exhibit" bson:"exhibit"`
	Owner string `json:"owner" bson:"owner"`
	Parent string `json:"parent" bson:"parent"`
	Position proxy.PositionInfo `json:"position" bson:"position"`
}

func CreateBooth(info *Booth) error {
	_, err := insertOne(TableBooth, info)
	if err != nil {
		return err
	}
	return nil
}

func GetBoothNextID() uint64 {
	num, _ := getSequenceNext(TableBooth)
	return num
}

func GetBooth(uid string) (*Booth, error) {
	result, err := findOne(TableBooth, uid)
	if err != nil {
		return nil, err
	}
	model := new(Booth)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllBoothsByOwner(owner string) ([]*Booth, error) {
	msg := bson.M{"owner": owner, "deletedAt": new(time.Time)}
	cursor, err1 := findMany(TableBooth, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Booth, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Booth)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllBoothsByParent(parent string) ([]*Booth, error) {
	msg := bson.M{"parent": parent, "deletedAt": new(time.Time)}
	cursor, err1 := findMany(TableBooth, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Booth, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Booth)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateBoothBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(TableBooth, uid, msg)
	return err
}

func UpdateBoothExhibit(uid, exhibit, operator string) error {
	msg := bson.M{"exhibit": exhibit, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(TableBooth, uid, msg)
	return err
}

func UpdateBoothPosition(uid, operator string, pos proxy.PositionInfo) error {
	msg := bson.M{"position": pos, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(TableBooth, uid, msg)
	return err
}

func RemoveBooth(uid, operator string) error {
	_, err := removeOne(TableBooth, uid, operator)
	return err
}
