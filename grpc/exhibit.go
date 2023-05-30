package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-museum/proto/museum"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.museum/cache"
	"omo.msa.museum/proxy"
)

type ExhibitService struct{}

func switchExhibit(info *cache.ExhibitInfo) *pb.ExhibitInfo {
	tmp := new(pb.ExhibitInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = uint64(info.CreateTime.Unix())
	tmp.Updated = uint64(info.UpdateTime.Unix())
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Sn = info.SN
	tmp.Owner = info.Owner
	tmp.Entity = info.Entity
	tmp.Tags = info.Tags
	tmp.Size = cache.SwitchVector2(&info.Size)
	tmp.Locals = switchLocals(info.Locals)
	tmp.Specials = switchSpecials(info.Specials)
	return tmp
}

func switchLocals(array []*proxy.LocalInfo) []*pb.LocalInfo {
	list := make([]*pb.LocalInfo, 0, len(array))
	for _, info := range array {
		list = append(list, &pb.LocalInfo{Language: info.Language, Name: info.Name, Remark: info.Remark})
	}
	return list
}

func switchSpecials(array []*proxy.SpecialInfo) []*pb.SpecialInfo {
	list := make([]*pb.SpecialInfo, 0, len(array))
	for _, info := range array {
		list = append(list, &pb.SpecialInfo{Id: info.ID, Key: info.Key, Value: info.Value})
	}
	return list
}

func (mine *ExhibitService) AddOne(ctx context.Context, in *pb.ReqExhibitAdd, out *pb.ReplyExhibitInfo) error {
	path := "exhibit.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	if cache.Context().HadExhibitByName(in.Name) {
		out.Status = outError(path, "the name is repeated", pbstatus.ResultStatus_Repeated)
		return nil
	}

	info, err := cache.Context().CreateExhibit(in.Name, in.Remark, in.Sn, in.Entity, in.Owner, in.Operator, cache.SwitchVector(in.Size))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchExhibit(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ExhibitService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyExhibitInfo) error {
	path := "exhibit.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetExhibit(in.Uid)
	if er != nil {
		out.Status = outError(path, "the exhibit not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchExhibit(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ExhibitService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "exhibit.getStatistic"
	inLog(path, in)
	if len(in.Field) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ExhibitService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "exhibit.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	er := cache.Context().RemoveExhibit(in.Uid, in.Operator)
	if er != nil {
		out.Status = outError(path, "the exhibit not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *ExhibitService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyExhibitList) error {
	path := "exhibit.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ExhibitService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyExhibitList) error {
	path := "exhibit.getListByFilter"
	inLog(path, in)
	var list []*cache.ExhibitInfo
	var err error
	if in.Field == "" {
		list, err = cache.Context().GetExhibitsByOwner(in.Owner)
	} else if in.Field == "array" {
		list = cache.Context().GetExhibits(in.Values)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.ExhibitInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchExhibit(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ExhibitService) UpdateLocals(ctx context.Context, in *pb.ReqExhibitLocals, out *pb.ReplyInfo) error {
	path := "exhibit.updateLocals"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetExhibit(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	var err error
	err = info.UpdateLocals(in.Operator, in.Locals)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ExhibitService) UpdateSpecials(ctx context.Context, in *pb.ReqExhibitSpecials, out *pb.ReplyInfo) error {
	path := "exhibit.updateSpecials"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetExhibit(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	var err error
	err = info.UpdateSpecials(in.Operator, in.Specials)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ExhibitService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "exhibit.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetExhibit(in.Uid)
	if er != nil {
		out.Status = outError(path, "the exhibit not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Field == "sn" {
		err = info.UpdateSN(in.Value, in.Operator)
	} else if in.Field == "size" {
		size := cache.ParseSize(in.Value)
		err = info.UpdateSize(in.Operator, size)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}
