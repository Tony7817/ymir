package product

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/product/rpc/product"
	"ymir.com/app/star/rpc/star"
	"ymir.com/pkg/id"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type ProductDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductDetailLogic {
	return &ProductDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProductDetailLogic) ProductDetail(req *types.ProductDetailRequest) (*types.ProductDetailResponse, error) {
	var pIdDecoded = l.svcCtx.Hash.DecodedId(req.Id)

	var (
		p  *product.ProductDetailResponse
		cl *product.ProductColorListResponse
	)

	err := mr.Finish(func() error {
		var err error
		p, err = l.svcCtx.ProductRPC.ProductDetail(l.ctx, &product.ProductDetailReqeust{
			Id: pIdDecoded,
		})
		if err != nil {
			return errors.Wrapf(err, "[ProductDetail] failed to get product detail")
		}
		return nil
	}, func() error {
		var err error
		cl, err = l.svcCtx.ProductRPC.ProductColorList(l.ctx, &product.ProductColorListRequest{
			ProductId: pIdDecoded,
		})
		if err != nil {
			return errors.Wrap(err, "[ProductDetail] failed to get product color list")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var (
		s  *star.StarDetailResponse
		pc *product.ProductColorResponse
	)

	err = mr.Finish(func() error {
		var err error
		s, err = l.svcCtx.StarRPC.StarDetail(l.ctx, &star.StarDetailRequest{
			Id: p.StarId,
		})
		if err != nil {
			return errors.Wrapf(err, "[ProductDetail] failed to get star detail")
		}
		return nil
	}, func() error {
		var err error
		pc, err = l.svcCtx.ProductRPC.ProductColor(l.ctx, &product.ProductColorRequest{
			ColorId: p.DefaultColorId,
		})
		if err != nil {
			return errors.Wrapf(err, "[ProductDetail] failed to get product color")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sizes, err := mr.MapReduce(func(source chan<- string) {
		for _, size := range pc.AvaliableSizes {
			source <- size
		}
	}, func(size string, writer mr.Writer[*types.ProductSize], cancel func(error)) {
		stock, err := l.svcCtx.ProductRPC.ProductStock(l.ctx, &product.ProductStockRequest{
			ProductId: pIdDecoded,
			ColorId:   pc.Id,
			Size:      size,
		})
		if err != nil {
			cancel(errors.Wrapf(err, "[ProductDetail] failed to get product stock"))
			return
		}
		productStock := types.ProductSize{
			SizeName: size,
			InStock:  stock.Stock,
		}
		writer.Write(&productStock)
	}, func(pipe <-chan *types.ProductSize, writer mr.Writer[[]types.ProductSize], cancel func(error)) {
		var sizes []types.ProductSize
		for size := range pipe {
			sizes = append(sizes, *size)
		}
		writer.Write(sizes)
	})
	if err != nil {
		return nil, err
	}

	pIdEncoded, err := l.svcCtx.Hash.EncodedId(p.Id)
	if err != nil {
		return nil, err
	}
	sIdEncoded, err := l.svcCtx.Hash.EncodedId(s.Id)
	if err != nil {
		return nil, err
	}

	var color = types.ProductColor{
		Id:            pIdEncoded,
		ColorName:     pc.Color,
		Images:        pc.Images,
		Detail_Images: pc.DetailImages,
		Price:         pc.Price,
		Unit:          pc.Unit,
		Size:          sizes,
	}

	var clres []types.ProductColorListItem
	for i := 0; i < len(cl.Colors); i++ {
		cIdEncoded, err := id.Hash.EncodedId(cl.Colors[i].ColorId)
		if err != nil {
			return nil, err
		}
		clres = append(clres, types.ProductColorListItem{
			ColorId:  cIdEncoded,
			CoverUrl: cl.Colors[i].CoverUrl,
		})
	}

	var res = &types.ProductDetailResponse{
		Id:          pIdEncoded,
		Description: p.Description,
		Rate:        p.Rate,
		RateCount:   p.ReteCount,
		Color:       color,
		ColorList:   clres,
		SoldNum:     p.SoldNum,
		Detail:      &p.Detail,
		StarAvatar:  s.AvatarUrl,
		StarName:    s.Name,
		StarId:      sIdEncoded,
	}

	return res, nil
}
