package cache

import (
	"time"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/repositories"
)

var (
	articleRecommendCacheKey = "recommend_articles_cache"
	articleHotCacheKey       = "hot_articles_cache"
)

var ArticleCache = newArticleCache()

type articleCache struct {
	recommendCache cache.LoadingCache
	hotCache       cache.LoadingCache
}

func newArticleCache() *articleCache {
	return &articleCache{
		recommendCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.ArticleRepository.Find(simple.DB(),
					simple.NewSqlCnd().Where("status = ?", model.StatusOk).Desc("id").Limit(50))
				return
			},
			cache.WithMaximumSize(1),
			cache.WithRefreshAfterWrite(30*time.Minute),
		),
		hotCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, err error) {
				createTime := simple.Timestamp(time.Now().AddDate(0, 0, -3))
				value = repositories.ArticleRepository.Find(simple.DB(),
					simple.NewSqlCnd().Gt("create_time", createTime).Eq("status", model.StatusOk).Desc("view_count").Limit(5))
				return
			},
			cache.WithMaximumSize(1),
			cache.WithRefreshAfterWrite(10*time.Minute),
		),
	}
}

func (this *articleCache) GetRecommendArticles() []model.Article {
	val, err := this.recommendCache.Get(articleRecommendCacheKey)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]model.Article)
	}
	return nil
}

func (this *articleCache) InvalidateRecommend() {
	this.recommendCache.Invalidate(articleRecommendCacheKey)
}

func (this *articleCache) GetHotArticles() []model.Article {
	val, err := this.hotCache.Get(articleHotCacheKey)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]model.Article)
	}
	return nil
}
