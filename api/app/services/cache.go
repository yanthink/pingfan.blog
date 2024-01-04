package services

import (
	"blog/app/cache"
	"blog/app/helpers"
	"context"
	"fmt"
	"sync"
)

type rememberCache struct {
	currentVersion int64
	data           any
}

var (
	rememberMap   sync.Map
	rememberMutex sync.Mutex
)

func remember[T any](key string, fn func() T) T {
	versionStore := cache.New[int64]()

	ctx := context.Background()
	versionKey := fmt.Sprintf("%s_version", key)

	if value, ok := rememberMap.Load(key); ok {
		rCache, _ := value.(*rememberCache)

		cacheVersion, _ := versionStore.Get(ctx, versionKey)
		// 判断缓存是否有效，如果有效直接返回缓存
		if rCache.data != nil && rCache.currentVersion == cacheVersion {
			return rCache.data.(T)
		}
	}

	rCache := &rememberCache{}

	rememberMutex.Lock()
	defer rememberMutex.Unlock()

	cacheVersion, _ := versionStore.Get(ctx, versionKey)

	// 再次判断缓存是否已经被其他 goroutine 更新
	if value, ok := rememberMap.Load(key); ok {
		rCache, _ = value.(*rememberCache)
		// 判断缓存是否有效，如果有效直接返回缓存
		if rCache.data != nil && rCache.currentVersion == cacheVersion {
			return rCache.data.(T)
		}
	}

	rCache.data = fn()

	// 更新缓存版本号
	version := helpers.Max(rCache.currentVersion+1, cacheVersion)
	if err := versionStore.Add(ctx, versionKey, version, 0); err != nil {
		version = cacheVersion
	}
	rCache.currentVersion = version

	rememberMap.Store(key, rCache)

	return rCache.data.(T)
}

func refresh(key string) {
	versionStore := cache.New[int64]()

	versionKey := fmt.Sprintf("%s_version", key)
	_, _ = versionStore.Increment(context.Background(), versionKey, 1)
}
