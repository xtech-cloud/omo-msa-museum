package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.museum/proxy"
	"time"
)

type Anchor struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deletedAt" bson:"deletedAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name     string           `json:"name" bson:"name"`
	Remark   string           `json:"remark" bson:"remark"`
	Owner    string           `json:"owner" bson:"owner"`
	Parent   string           `json:"parent" bson:"parent"`
	Cover    string           `json:"cover" bson:"cover"`
	Panorama string           `json:"panorama" bson:"panorama"`
	Link     string           `json:"link" bson:"link"`
	Position proxy.VectorInfo `json:"position" bson:"position"`
	Tags     []string         `json:"tags" bson:"tags"`
	Assets   []string         `json:"assets" bson:"assets"`
}

func CreateAnchor(info *Anchor) error {
	_, err := insertOne(TableAnchor, info)
	if err != nil {
		return err
	}
	return nil
}

func GetAnchorNextID() uint64 {
	num, _ := getSequenceNext(TableAnchor)
	return num
}

func GetAnchor(uid string) (*Anchor, error) {
	result, err := findOne(TableAnchor, uid)
	if err != nil {
		return nil, err
	}
	model := new(Anchor)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAnchorsByParent(parent string) ([]*Anchor, error) {
	msg := bson.M{"parent": parent, "deletedAt": new(time.Time)}
	cursor, err1 := findMany(TableAnchor, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Anchor, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Anchor)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAnchorsByOwner(owner string) ([]*Anchor, error) {
	msg := bson.M{"owner": owner, "deletedAt": new(time.Time)}
	cursor, err1 := findMany(TableAnchor, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Anchor, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Anchor)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateAnchorBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAnchor, uid, msg)
	return err
}

func UpdateAnchorCover(uid, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAnchor, uid, msg)
	return err
}

func UpdateAnchorLink(uid, link, operator string) error {
	msg := bson.M{"link": link, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAnchor, uid, msg)
	return err
}

func UpdateAnchorPanorama(uid, cover, operator string) error {
	msg := bson.M{"panorama": cover, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAnchor, uid, msg)
	return err
}

func UpdateAnchorPosition(uid, operator string, pos proxy.VectorInfo) error {
	msg := bson.M{"position": pos, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAnchor, uid, msg)
	return err
}

func UpdateAnchorAssets(uid, operator string, assets []string) error {
	msg := bson.M{"assets": assets, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAnchor, uid, msg)
	return err
}

func RemoveAnchor(uid, operator string) error {
	_, err := removeOne(TableAnchor, uid, operator)
	return err
}

func UpdateAnchorTags(uid, operator string, tags []string) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"tags": tags, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAnchor, uid, msg)
	return err
}
