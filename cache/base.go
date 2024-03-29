package cache

import (
	pb "github.com/xtech-cloud/omo-msp-museum/proto/museum"
	"omo.msa.museum/config"
	"omo.msa.museum/proxy"
	"omo.msa.museum/proxy/nosql"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const DefaultOwner = "system"

type baseInfo struct {
	ID         uint64 `json:"-"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	Creator    string
	Operator   string
	CreateTime time.Time
	UpdateTime time.Time
}

type cacheContext struct {
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if err == nil {

	}
	return err
}

func Context() *cacheContext {
	return cacheCtx
}

func checkPage(page, number uint32, all interface{}) (uint32, uint32, interface{}) {
	if number < 1 {
		number = 10
	}
	array := reflect.ValueOf(all)
	total := uint32(array.Len())
	maxPage := total / number
	if total%number != 0 {
		maxPage = total/number + 1
	}
	if page < 1 {
		return total, maxPage, all
	}
	if page > maxPage {
		page = maxPage
	}

	var start = (page - 1) * number
	var end = start + number
	if end > total {
		end = total
	}

	list := array.Slice(int(start), int(end))
	return total, maxPage, list.Interface()
}

func SwitchVector(point *pb.Vector3) proxy.VectorInfo {
	return proxy.VectorInfo{X: point.X, Y: point.Y, Z: point.Z}
}

func SwitchVector2(point *proxy.VectorInfo) *pb.Vector3 {
	return &pb.Vector3{X: point.X, Y: point.Y, Z: point.Z}
}

func ParseSize(str string) proxy.VectorInfo {
	arr := strings.Split(str, ";")
	vec := proxy.VectorInfo{
		X: 0, Y: 0, Z: 0,
	}
	if len(arr) != 3 {
		return vec
	}
	x, _ := strconv.ParseFloat(arr[0], 32)
	y, _ := strconv.ParseFloat(arr[1], 32)
	z, _ := strconv.ParseFloat(arr[2], 32)
	vec.X = float32(x)
	vec.Y = float32(y)
	vec.Z = float32(z)
	return vec
}
