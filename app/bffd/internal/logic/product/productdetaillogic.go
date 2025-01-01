package product

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/product/rpc/product"
	"ymir.com/app/star/rpc/star"
	"ymir.com/app/user/rpc/user"
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
		p        *product.ProductDetailResponse
		cl       *product.ProductColorListResponse
		pcmt     *product.ProductCommentListResponse
		pcmtSize int64 = 10
	)

	// find product basic info by pId
	// find product color list by pId
	// find product comments by pId
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
	}, func() error {
		var err error
		pcmt, err = l.svcCtx.ProductRPC.ProductCommentList(l.ctx, &product.ProductCommentListRequest{
			ProductId: pIdDecoded,
			Page:      1,
			PageSize:  pcmtSize,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var (
		s         *star.StarDetailResponse
		pcolor    *product.ProductColorResponse
		pcmtFinal = make([]types.ProductComment, 0)
	)

	// find star detail by starId
	// find default color by defaultColorId
	// find user info in product comment by user id
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
		pcolor, err = l.svcCtx.ProductRPC.ProductColor(l.ctx, &product.ProductColorRequest{
			ColorId: p.DefaultColorId,
		})
		if err != nil {
			return errors.Wrapf(err, "[ProductDetail] failed to get product color")
		}
		return nil
	}, func() error {
		var err error
		if len(pcmt.Comments) > 0 {
			pcmtFinal, err = l.buildProductCommentUserInfo(pcmt.Comments)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// find sizes by productId and colorId
	sizes, err := mr.MapReduce(func(source chan<- string) {
		for _, size := range pcolor.AvaliableSizes {
			source <- size
		}
	}, func(size string, writer mr.Writer[*types.ProductSize], cancel func(error)) {
		stock, err := l.svcCtx.ProductRPC.ProductStock(l.ctx, &product.ProductStockRequest{
			ProductId: pIdDecoded,
			ColorId:   pcolor.Id,
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
		ColorName:     pcolor.Color,
		Images:        pcolor.Images,
		Detail_Images: pcolor.DetailImages,
		Price:         pcolor.Price,
		Unit:          pcolor.Unit,
		Size:          sizes,
	}

	var clres = make([]types.ProductColorListItem, len(cl.Colors))
	for i := 0; i < len(cl.Colors); i++ {
		cIdEncoded, err := id.Hash.EncodedId(cl.Colors[i].ColorId)
		if err != nil {
			return nil, err
		}
		clres[i] = types.ProductColorListItem{
			ColorId:  cIdEncoded,
			CoverUrl: cl.Colors[i].CoverUrl,
		}
	}

	var pcmtRes = types.ProductCommentResponse{
		Comments: pcmtFinal,
		Total:    pcmt.Total,
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
		StarRate:    s.Rate,
		Comments:    pcmtRes,
	}

	return res, nil
}

func (l *ProductDetailLogic) buildProductCommentUserInfo(pcmts []*product.ProductComment) ([]types.ProductComment, error) {
	res, err := mr.MapReduce(func(source chan<- *product.ProductComment) {
		for _, pcmt := range pcmts {
			source <- pcmt
		}
	}, func(pcmt *product.ProductComment, writer mr.Writer[types.ProductComment], cancel func(error)) {
		user, err := l.svcCtx.UserRPC.GetUserInfo(l.ctx, &user.GetUserInfoRequest{
			UserId: &pcmt.UserId,
		})
		if err != nil {
			cancel(errors.Wrapf(err, "[ProductDetail] failed to get user info"))
			return
		}
		pcmtIdEncoded, err := id.Hash.EncodedId(pcmt.Id)
		if err != nil {
			cancel(errors.Wrapf(err, "[ProductDetail] failed to encode product comment id"))
			return
		}
		writer.Write(types.ProductComment{
			Id:          pcmtIdEncoded,
			UserName:    user.User.Username,
			UserAvatar:  user.User.AvatarUrl,
			Rate:        pcmt.Rate,
			Comment:     pcmt.Comment,
			LikeNum:     pcmt.LikeNum,
			Images:      pcmt.Images,
			ImagesThumb: pcmt.ImagesThumb,
			CreatedAt:   pcmt.CreateAt,
			Size:        pcmt.Size,
			Color:       pcmt.Color,
		})
	}, func(pipe <-chan types.ProductComment, writer mr.Writer[[]types.ProductComment], cancel func(error)) {
		var pcmts []types.ProductComment
		for pcmt := range pipe {
			pcmts = append(pcmts, pcmt)
		}
		writer.Write(pcmts)
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
