package star

import (
	"context"

	"ymir.com/app/bffs/internal/svc"
	"ymir.com/app/bffs/internal/types"
	"ymir.com/app/product/admin/product"
	"ymir.com/pkg/id"
	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type CreateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateProductLogic) CreateProduct(req *types.CreateProductRequest) (*types.CreateProductResponse, error) {
	if req.Name == "" || req.Description == "" || !isSizesAdnColorInvalid(req.Color) {
		return nil, xerr.NewErrCode(xerr.ReuqestParamError)
	}

	sId, err := id.DecodeId(req.StarId)
	if err != nil {
		return nil, err
	}
	var pId = id.SF.GenerateID()
	respbProduct, err := l.svcCtx.ProductAdminRPC.CreateProduct(l.ctx, &product.CreateProductRequest{
		ProductId:   pId,
		StarId:      sId,
		Name:        req.Name,
		Description: req.Description,
		Detail:      req.Detail,
	})
	if err != nil {
		return nil, err
	}

	err = mr.MapReduceVoid(func(source chan<- *types.ProductColor) {
		for i := 0; i < len(req.Color); i++ {
			source <- &req.Color[i]
		}
	}, func(item *types.ProductColor, writer mr.Writer[int64], cancel func(error)) {
		cId := id.SF.GenerateID()
		respb, err := l.svcCtx.ProductAdminRPC.CreateProductColor(l.ctx, &product.CreateProductColorRequeset{
			ProductColorId:  cId,
			ProductId:       respbProduct.Id,
			ColorName:       item.ColorName,
			CoverUrl:        item.CoverUrl,
			ImageUrl:        item.ImagesUrl,
			DetailImagesUrl: item.DetailImagesUrl,
			Price:           item.Price,
			Unit:            item.Unit,
			Size:            composeSizeStr(item.Sizes),
			IsDefault:       item.IsDefault,
		})
		if err != nil {
			cancel(errors.Wrap(err, "failed to create product color"))
		}
		err = l.createProductColorStock(pId, respb.Id, item.Sizes)
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(respb.Id)
	}, func(pipe <-chan int64, cancel func(error)) {
		for range pipe {
		}
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateProductResponse{
		ProductId: id.EncodeId(pId),
	}, nil
}

func (l *CreateProductLogic) createProductColorStock(pId int64, cId int64, sizes []types.ProductColorSize) error {
	return mr.MapReduceVoid(func(source chan<- *types.ProductColorSize) {
		for i := 0; i < len(sizes); i++ {
			source <- &sizes[i]
		}
	}, func(s *types.ProductColorSize, writer mr.Writer[int64], cancel func(error)) {
		var err error
		sId := id.SF.GenerateID()
		_, err = l.svcCtx.ProductAdminRPC.CreateProductColorStock(l.ctx, &product.CreateProductColorStockRequest{
			ProductColorStockId: sId,
			ProductId:           pId,
			ColorId:             cId,
			InStock:             s.InStock,
			Size:                s.Size,
		})
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(0)
	}, func(pipe <-chan int64, cancel func(error)) {
		for range pipe {
		}
	})
}

func isSizesAdnColorInvalid(cs []types.ProductColor) bool {
	var flag = 0
	for i := 0; i < len(cs); i++ {
		if len(cs[i].Sizes) == 0 {
			return false
		}
		if cs[i].IsDefault {
			flag++
		}
	}
	if flag != 1 {
		return false
	}

	return true
}

func composeSizeStr(sizes []types.ProductColorSize) string {
	var res string
	for i := 0; i < len(sizes); i++ {
		res += sizes[i].Size
		res += ","
	}
	return res[:len(res)-1]
}
