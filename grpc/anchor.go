package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-museum/proto/museum"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.museum/cache"
)

type AnchorService struct{}

func switchAnchor(info *cache.AnchorInfo) *pb.AnchorInfo {
	tmp := new(pb.AnchorInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Cover = info.Cover
    tmp.Parent = info.Parent
    tmp.Panorama = info.Panorama
    tmp.Link = info.Link
    //tmp.Owner = info.Owner
    tmp.Position = &pb.Vector3{X: info.Position.X, Y: info.Position.Y}
	tmp.Tags = info.Tags
	return tmp
}

func (mine *AnchorService) AddOne(ctx context.Context, in *pb.ReqAnchorAdd, out *pb.ReplyAnchorOne) error {
	path := "anchor.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	if len(in.Parent) < 1 {
		out.Status = outError(path, "the parent is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreateAnchor(in.Name, in.Remark, in.Owner, in.Parent, in.Operator, in.Tags)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchAnchor(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *AnchorService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAnchorOne) error {
	path := "anchor.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetAnchor(in.Uid)
	if er != nil {
		out.Status = outError(path, "the anchor not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchAnchor(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *AnchorService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "anchor.getStatistic"
	inLog(path, in)
	if len(in.Field) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *AnchorService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "anchor.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetAnchor(in.Uid)
	if er != nil {
		out.Status = outError(path, "the anchor not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.Remove(in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *AnchorService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAnchorList) error {
	path := "anchor.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *AnchorService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyAnchorList) error {
	path := "anchor.getListByFilter"
	inLog(path, in)
	var list []*cache.AnchorInfo
	var err error
	if in.Field == "" {

	} else if in.Field == "status" {

	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.AnchorInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchAnchor(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *AnchorService) UpdateBase(ctx context.Context, in *pb.ReqAnchorUpdate, out *pb.ReplyInfo) error {
	path := "anchor.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetAnchor(in.Uid)
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

func (mine *AnchorService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "anchor.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetAnchor(in.Uid)
	if er != nil {
		out.Status = outError(path, "the anchor not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Field == "cover" {
		err = info.UpdateCover(in.Value, in.Operator)
	} else if in.Field == "targets" {
		err = errors.New("the field not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}
