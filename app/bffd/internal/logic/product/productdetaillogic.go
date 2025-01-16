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
	var (
		p         *product.ProductDetailResponse
		respbCs   *product.ProductColorListResponse
		pcmt      *product.ProductCommentListResponse
		s         *star.StarDetailResponse
		pcmtFinal []types.ProductComment
		pColors   []types.ProductColor
		pcmtSize  int64 = 10
	)
	pId, err := id.DecodeId(req.Id)
	if err != nil {
		return nil, err
	}

	// find product basic info by pId
	// find product color list by pId
	// find product comments by pId
	err = mr.Finish(func() error {
		var err error
		p, err = l.svcCtx.ProductRPC.ProductDetail(l.ctx, &product.ProductDetailReqeust{
			Id: pId,
		})
		if err != nil {
			return errors.Wrapf(err, "[ProductDetail] failed to get product detail")
		}
		s, err = l.svcCtx.StarRPC.StarDetail(l.ctx, &star.StarDetailRequest{
			Id: p.StarId,
		})
		if err != nil {
			return errors.Wrapf(err, "[ProductDetail] failed to get star detail")
		}
		return nil
	}, func() error {
		var err error
		respbCs, err = l.svcCtx.ProductRPC.ProductColorList(l.ctx, &product.ProductColorListRequest{
			ProductId: pId,
		})
		if err != nil {
			return errors.Wrap(err, "[ProductDetail] failed to get product color list")
		}
		pColors, err = l.findProductColorStock(respbCs.Colors)
		if err != nil {
			return errors.Wrap(err, "[ProductDetail] failed to get product color stock")
		}
		return nil
	}, func() error {
		var err error
		pcmt, err = l.svcCtx.ProductRPC.ProductCommentList(l.ctx, &product.ProductCommentListRequest{
			ProductId: pId,
			Page:      1,
			PageSize:  pcmtSize,
		})
		if err != nil {
			return err
		}
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

	var pcmtRes = types.ProductCommentResponse{
		Comments: pcmtFinal,
		Total:    pcmt.Total,
	}

	var res = &types.ProductDetailResponse{
		Id:          id.EncodeId(p.Id),
		Description: p.Description,
		Rate:        p.Rate,
		RateCount:   p.ReteCount,
		Colors:      pColors,
		SoldNum:     p.SoldNum,
		Detail:      &p.Detail,
		StarAvatar:  s.AvatarUrl,
		StarName:    s.Name,
		StarId:      id.EncodeId(s.Id),
		StarRate:    s.Rate,
		Comments:    pcmtRes,
	}

	return res, nil
}

func (l *ProductDetailLogic) findProductColorStock(cs []*product.ProductColor) ([]types.ProductColor, error) {
	resColor, err := mr.MapReduce(func(source chan<- *product.ProductColor) {
		for i := 0; i < len(cs); i++ {
			source <- cs[i]
		}
	}, func(c *product.ProductColor, writer mr.Writer[*types.ProductColor], cancel func(error)) {
		resStock, err := l.findProductColorStockByColor(c.AvaliableSizes, c.ProductId, c.Id)
		if err != nil {
			cancel(errors.Wrapf(err, "[ProductDetail] failed to get product stock"))
			return
		}
		writer.Write(&types.ProductColor{
			Id:            id.EncodeId(c.Id),
			ColorName:     c.Name,
			Images:        c.Images,
			Detail_Images: c.DetailImages,
			Price:         c.Price,
			CoverUrl:      c.CoverUrl,
			Unit:          c.Unit,
			Size:          resStock,
			IsDefault:     c.IsDefault,
		})
	}, func(pipe <-chan *types.ProductColor, writer mr.Writer[[]types.ProductColor], cancel func(error)) {
		var cs []types.ProductColor
		for c := range pipe {
			cs = append(cs, *c)
		}
		writer.Write(cs)
	})
	if err != nil {
		return nil, err
	}

	return resColor, nil
}

func (l *ProductDetailLogic) findProductColorStockByColor(sizes []string, pId int64, cId int64) ([]types.ProductSize, error) {
	res, err := mr.MapReduce(func(source chan<- string) {
		for i := 0; i < len(sizes); i++ {
			source <- sizes[i]
		}
	}, func(size string, writer mr.Writer[*types.ProductSize], cancel func(error)) {
		ps, err := l.svcCtx.ProductRPC.ProductStock(l.ctx, &product.ProductStockRequest{
			ProductId: pId,
			ColorId:   cId,
			Size:      size,
		})
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(&types.ProductSize{
			SizeName: size,
			InStock:  ps.Stock,
		})
	}, func(pipe <-chan *types.ProductSize, writer mr.Writer[[]types.ProductSize], cancel func(error)) {
		var res []types.ProductSize
		for p := range pipe {
			res = append(res, *p)
		}
		writer.Write(res)
	})
	if err != nil {
		return nil, err
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
		writer.Write(types.ProductComment{
			Id:          id.EncodeId(pcmt.Id),
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
