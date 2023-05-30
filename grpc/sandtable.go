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

type SandtableService struct{}

func switchSandtable(info *cache.SandtableInfo) *pb.SandtableInfo {
	tmp := new(pb.SandtableInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = uint64(info.CreateTime.Unix())
	tmp.Updated = uint64(info.UpdateTime.Unix())
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Owner = info.Owner
	tmp.Background = info.Background
	tmp.Bgm = info.BGM
	tmp.Mask = info.Mask
	tmp.Narrate = info.Narrate
	tmp.Width = info.Width
	tmp.Height = info.Height
	tmp.Tags = info.Tags
	tmp.Paths = make([]*pb.PathInfo, 0, len(info.Paths))
	for _, item := range info.Paths {
		tmp.Paths = append(tmp.Paths, switchPath(item))
	}
	return tmp
}

func switchPath(path *proxy.PathInfo) *pb.PathInfo {
	tmp := new(pb.PathInfo)
	tmp.Uid = path.UID
	tmp.Name = path.Name
	tmp.Color = path.Color
	tmp.Points = make([]*pb.PathKeyInfo, 0, len(path.Points))
	for _, item := range path.Points {
		pos := cache.SwitchVector2(&item.Position)
		ro := cache.SwitchVector2(&item.Rotation)
		tmp.Points = append(tmp.Points, &pb.PathKeyInfo{Key: item.Key, Scale: item.Scale, Position: pos, Rotation: ro})
	}
	return tmp
}

func (mine *SandtableService) AddOne(ctx context.Context, in *pb.ReqSandtableAdd, out *pb.ReplySandtableInfo) error {
	path := "sandtable.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	if len(in.Background) < 1 {
		out.Status = outError(path, "the background asset is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreateSandtable(in.Name, in.Remark, in.Background, in.Owner, in.Operator, in.Width, in.Height)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchSandtable(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SandtableService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplySandtableInfo) error {
	path := "sandtable.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	var info *cache.SandtableInfo
	var er error
	if in.Flag == "" {
		info, er = cache.Context().GetSandtable(in.Uid)
	}

	if er != nil {
		out.Status = outError(path, "the sandtable not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchSandtable(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SandtableService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "sandtable.getStatistic"
	inLog(path, in)
	if len(in.Field) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *SandtableService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "sandtable.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetSandtable(in.Uid)
	if er != nil {
		out.Status = outError(path, "the sandtable not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	er = info.Remove(in.Operator)
	if er != nil {
		out.Status = outError(path, "the sandtable not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *SandtableService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplySandtableList) error {
	path := "sandtable.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *SandtableService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplySandtableList) error {
	path := "sandtable.getListByFilter"
	inLog(path, in)
	var list []*cache.SandtableInfo
	var err error
	if in.Field == "" {
		list, err = cache.Context().GetSandTablesByOwner(in.Owner)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.SandtableInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchSandtable(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *SandtableService) UpdateBase(ctx context.Context, in *pb.ReqSandtableBase, out *pb.ReplyInfo) error {
	path := "sandtable.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetSandtable(in.Uid)
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

func (mine *SandtableService) UpdateBackground(ctx context.Context, in *pb.ReqSandtableBG, out *pb.ReplyInfo) error {
	path := "sandtable.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetSandtable(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateBackground(in.Asset, in.Operator, in.Width, in.Height)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *SandtableService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "sandtable.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetSandtable(in.Uid)
	if er != nil {
		out.Status = outError(path, "the sandtable not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Field == "bgm" {
		err = info.UpdateBGM(in.Value, in.Operator)
	} else if in.Field == "narrate" {
		err = info.UpdateNarrate(in.Value, in.Operator)
	} else {
		err = errors.New("the field not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *SandtableService) UpdatePath(ctx context.Context, in *pb.ReqSandtablePath, out *pb.ReplyInfo) error {
	path := "sandtable.updatePath"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetSandtable(in.Uid)
	if er != nil {
		out.Status = outError(path, "the sandtable not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	path, er = info.UpdatePath(in.Operator, in.Path, in.Name, in.Color, in.Points)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = path
	out.Status = outLog(path, out)
	return nil
}
