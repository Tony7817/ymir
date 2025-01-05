package logic

import (
	"context"
	"database/sql"
	"strconv"

	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/model"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/cache"
	"ymir.com/pkg/vars"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
)

type ProductCommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCommentListLogic {
	return &ProductCommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductCommentListLogic) ProductCommentList(in *product.ProductCommentListRequest) (*product.ProductCommentListResponse, error) {
	var (
		count int64
		pcmts = make([]model.ProductComment, 0)
	)
	err := mr.Finish(func() error {
		var err error
		pcmts, err = l.svcCtx.ProductCommentModel.FindProductCommentList(l.ctx, in.ProductId, in.Page, in.PageSize)
		if err != nil {
			return err
		}
		return nil
	}, func() error {
		var err error
		count, err = l.svcCtx.ProductCommentModel.CountTotalProductComment(l.ctx, in.ProductId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var pcmtRes = make([]*product.ProductComment, 0)
	if len(pcmts) != 0 {
		var err error
		pcmtRes, err = l.buildProductCommentImages(pcmts)
		if err != nil {
			return nil, err
		}
	}

	return &product.ProductCommentListResponse{
		Comments: pcmtRes,
		Total:    count,
	}, nil
}

func (l *ProductCommentListLogic) buildProductCommentImages(pcmts []model.ProductComment) ([]*product.ProductComment, error) {
	return mr.MapReduce(func(source chan<- *model.ProductComment) {
		for _, pcmt := range pcmts {
			source <- &pcmt
		}
	}, func(pcmt *model.ProductComment, writer mr.Writer[*product.ProductComment], cancel func(error)) {
		ids, err := l.cacheImageList(pcmt.Id)
		if err != nil {
			logx.Errorf("[ProductCommentListLogic] failed to get cache product comment images, err: %v", err)
		}
		var images []model.ProductCommentImage
		if len(ids) != 0 {
			images, err = l.imagesByIds(ids)
			if err != nil {
				cancel(err)
				return
			}
		} else {
			images, err = l.findImagesAndSetCache(pcmt.Id)
			if err != nil {
				cancel(err)
				return
			}
		}
		var imagesUrl = make([]string, len(images))
		var imagesThumbUrl = make([]string, len(images))
		for i := 0; i < len(images); i++ {
			imagesUrl[i] = images[i].ImageUrl
			imagesThumbUrl[i] = images[i].ImageThumbUrl
		}
		writer.Write(&product.ProductComment{
			Id:          pcmt.Id,
			ProductId:   pcmt.ProductId,
			UserId:      pcmt.UserId,
			Comment:     pcmt.Comment,
			Rate:        pcmt.Rate,
			LikeNum:     pcmt.LikeNum,
			CreateAt:    pcmt.CreatedAt.Unix(),
			Images:      imagesUrl,
			ImagesThumb: imagesThumbUrl,
			Size:        pcmt.Size,
			Color:       pcmt.Color,
		})

	}, func(pipe <-chan *product.ProductComment, writer mr.Writer[[]*product.ProductComment], cancel func(error)) {
		var pcmts []*product.ProductComment
		for pcmt := range pipe {
			pcmts = append(pcmts, pcmt)
		}
		writer.Write(pcmts)
	})
}

func (l *ProductCommentListLogic) cacheImageList(commentId int64) ([]int64, error) {
	idsRaw, err := l.svcCtx.Redis.LrangeCtx(l.ctx, cache.CacheKeyProductCommentList(commentId), 0, -1)
	if err != nil {
		return nil, err
	}

	var ids []int64
	for i := 0; i < len(idsRaw); i++ {
		id, _ := strconv.ParseInt(idsRaw[i], 10, 64)
		ids = append(ids, id)
	}

	return ids, nil
}

func (l *ProductCommentListLogic) imagesByIds(ids []int64) ([]model.ProductCommentImage, error) {
	res, err := mr.MapReduce(func(source chan<- int64) {
		for _, id := range ids {
			source <- id
		}
	}, func(id int64, writer mr.Writer[*model.ProductCommentImage], cancel func(error)) {
		img, err := l.svcCtx.ProductCommentImageModel.FindOne(l.ctx, id)
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(img)
	}, func(pipe <-chan *model.ProductCommentImage, writer mr.Writer[[]model.ProductCommentImage], cancel func(error)) {
		var imgs []model.ProductCommentImage
		for img := range pipe {
			imgs = append(imgs, *img)
		}
		writer.Write(imgs)
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (l *ProductCommentListLogic) findImagesAndSetCache(commentId int64) ([]model.ProductCommentImage, error) {
	images, err := l.svcCtx.ProductCommentImageModel.FindImagesbyCommentId(l.ctx, commentId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return make([]model.ProductCommentImage, 0), nil
	}

	// set cache
	threading.GoSafe(func() {
		for i := 0; i < len(images); i++ {
			// don't use RpushCtx. The execution can't be done before context timeout.
			_, err = l.svcCtx.Redis.Rpush(cache.CacheKeyProductCommentList(commentId), images[i].Id)
			if err != nil {
				logx.Errorf("[ProductCommentListLogic] failed to set cache product comment images, err: %v", err)
			}
		}
		err = l.svcCtx.Redis.Expire(cache.CacheKeyProductCommentList(commentId), vars.CacheExpireIn1W)
		if err != nil {
			logx.Errorf("[ProductCommentListLogic] failed to set cache expire product comment images, err: %v", err)
		}
	})

	return images, nil
}
