package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-museum/proto/museum"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.museum/cache"
)

type BoothService struct{}

func switchBooth(info *cache.BoothInfo) *pb.BoothInfo {
	tmp := new(pb.BoothInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = uint64(info.CreateTime.Unix())
	tmp.Updated = uint64(info.UpdateTime.Unix())
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Exhibit = info.Exhibit
	tmp.Owner = info.Owner
	tmp.Parent = info.Parent
	tmp.Position = &pb.Vector3{X: info.Position.X, Y: info.Position.Y}
	return tmp
}

func (mine *BoothService) AddOne(ctx context.Context, in *pb.ReqBoothAdd, out *pb.ReplyBoothInfo) error {
	path := "panorama.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreateBooth(in.Name, in.Remark, in.Owner, in.Parent, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchBooth(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *BoothService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyBoothInfo) error {
	path := "panorama.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetBooth(in.Uid)
	if er != nil {
		out.Status = outError(path, "the panorama not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchBooth(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *BoothService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "panorama.getStatistic"
	inLog(path, in)
	if len(in.Field) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *BoothService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "panorama.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	er := cache.Context().RemoveBooth(in.Uid, in.Operator)
	if er != nil {
		out.Status = outError(path, "the panorama not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *BoothService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyBoothList) error {
	path := "panorama.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *BoothService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyBoothList) error {
	path := "panorama.getListByFilter"
	inLog(path, in)
	var list []*cache.BoothInfo
	var err error
	if in.Field == "" {
		list, err = cache.Context().GetBoothsByOwner(in.Owner)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.BoothInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchBooth(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *BoothService) UpdateBase(ctx context.Context, in *pb.ReqBoothUpdate, out *pb.ReplyInfo) error {
	path := "panorama.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetBooth(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *BoothService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "panorama.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	_, er := cache.Context().GetBooth(in.Uid)
	if er != nil {
		out.Status = outError(path, "the panorama not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error

	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}
