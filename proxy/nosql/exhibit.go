package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.museum/proxy"
	"time"
)

type Exhibit struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deletedAt" bson:"deletedAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name   string   `json:"name" bson:"name"`
	Owner  string   `json:"owner" bson:"owner"`
	Entity string `json:"entity" bson:"entity"`
	SN string `json:"sn" bson:"sn"`
	Tags []string `json:"tags" bson:"tags"`
	Locals []*proxy.LocalInfo `json:"locals" bson:"locals"`
	Specials []*proxy.SpecialInfo `json:"specials" bson:"specials"`
}

func CreateExhibit(info *Exhibit) error {
	_, err := insertOne(TableExhibit, info)
	if err != nil {
		return err
	}
	return nil
}

func GetExhibitNextID() uint64 {
	num, _ := getSequenceNext(TableExhibit)
	return num
}

func GetExhibit(uid string) (*Exhibit, error) {
	result, err := findOne(TableExhibit, uid)
	if err != nil {
		return nil, err
	}
	model := new(Exhibit)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetExhibitByName(name string) (*Exhibit, error) {
	filter := bson.M{"name": name, "deletedAt": new(time.Time)}
	result, err := findOneBy(TableExhibit, filter)
	if err != nil {
		return nil, err
	}
	model := new(Exhibit)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetExhibitBySN(owner, sn string) (*Exhibit, error) {
	filter := bson.M{"owner":owner, "sn": sn, "deletedAt": new(time.Time)}
	result, err := findOneBy(TableExhibit, filter)
	if err != nil {
		return nil, err
	}
	model := new(Exhibit)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllExhibitsByOwner(owner string) ([]*Exhibit, error) {
	msg := bson.M{"owner": owner, "deletedAt": new(time.Time)}
	cursor, err1 := findMany(TableExhibit, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Exhibit, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Exhibit)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllExhibits() ([]*Exhibit, error) {
	cursor, err1 := findAll(TableExhibit, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Exhibit, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Exhibit)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateExhibitSN(uid, sn, operator string) error {
	msg := bson.M{"sn": sn, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableExhibit, uid, msg)
	return err
}

func UpdateExhibitLocals(uid, operator string, list []*proxy.LocalInfo) error {
	msg := bson.M{"locals": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableExhibit, uid, msg)
	return err
}

func UpdateExhibitSpecials(uid, operator string, list []*proxy.SpecialInfo) error {
	msg := bson.M{"specials": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableExhibit, uid, msg)
	return err
}

func RemoveExhibit(uid, operator string) error {
	_, err := removeOne(TableExhibit, uid, operator)
	return err
}
