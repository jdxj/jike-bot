package main

import (
	"context"
	"math/rand"
	"time"

	pexels "github.com/jdxj/pexels-sdk-go"
	unsplash_sdk_go "github.com/jdxj/unsplash-sdk-go"
	wh "github.com/jdxj/wallhaven-sdk-go"
)

var (
	myRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type WallpaperSource interface {
	GetWallpaper(context.Context) (string, error)
}

const (
	pageSize = 24
)

func newWallhavenSource() *wallhavenSource {
	return &wallhavenSource{
		whc:   wh.NewClient(),
		total: 1,
	}
}

type wallhavenSource struct {
	whc *wh.Client

	total int
}

func (ws *wallhavenSource) GetWallpaper(ctx context.Context) (string, error) {
	randTotal := myRand.Intn(ws.total)

	searchRsp, err := ws.whc.Search(ctx, &wh.SearchReq{
		Category: wh.People | wh.Anime | wh.General,
		AIArt:    true,
		Purity:   wh.SFW,
		Sorting:  wh.TopList,
		Order:    wh.Desc,
		TopRange: wh.M1,
		Page:     (randTotal - randTotal%pageSize + pageSize) / pageSize,
	})
	if err != nil {
		return "", err
	}
	if len(searchRsp.Wallpapers) == 0 {
		return "", ErrWallpaperNotFound
	}

	url := searchRsp.Wallpapers[randTotal%pageSize].Path
	ws.total = searchRsp.Meta.Total
	return url, nil
}

func newPexelsSource(apiKey string) *pexelsSource {
	return &pexelsSource{
		pc: pexels.NewClient(apiKey, pexels.WithDebug(true)),
	}
}

type pexelsSource struct {
	pc *pexels.Client
}

func (ps *pexelsSource) GetWallpaper(ctx context.Context) (string, error) {
	pl, err := ps.pc.CuratedPhotos(ctx, &pexels.CuratedReq{
		Pagination: pexels.Pagination{
			Page:    myRand.Intn(8000) + 1,
			PerPage: 1,
		},
	})
	if err != nil {
		return "", err
	}
	if len(pl.Photos) == 0 {
		return "", ErrWallpaperNotFound
	}
	return pl.Photos[0].Src.Original, nil
}

func newUnsplashSource(ak, sk string) *unsplashSource {
	return &unsplashSource{
		uc: unsplash_sdk_go.New(ak, sk, unsplash_sdk_go.WithDebug(true)),
	}
}

type unsplashSource struct {
	uc *unsplash_sdk_go.Client
}

func (us *unsplashSource) GetWallpaper(ctx context.Context) (string, error) {
	rsp, err := us.uc.GetRandomPhoto(ctx, nil)
	if err != nil {
		return "", err
	}
	return rsp.Urls.Full, nil
}
