package proxy

type SpecialInfo struct {
	ID uint32 `json:"id" bson:"id"`
	Key string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

type LocalInfo struct {
	Language string `json:"language" bson:"language"`
	Name string `json:"name" bson:"name"`
	Remark string `json:"remark" bson:"remark"`
}

type VectorInfo struct {
	X float32 `json:"x" bson:"x"`
	Y float32 `json:"y" bson:"y"`
	Z float32 `json:"z" bson:"z"`
}

type FrameKeyInfo struct {
	Key string `json:"key" bson:"key"`
	Scale float32 `json:"scale" bson:"scale"`
	Position VectorInfo `json:"position" bson:"position"`
	Rotation VectorInfo `json:"rotation" bson:"rotation"`
}