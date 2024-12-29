package cache

import "fmt"

func CacheKeyProductCommentList(commentId int64) string {
	return fmt.Sprintf("cache:ymir:list:product:comment:%d", commentId)
}