package main

import (
	"context"
	"math/rand"
	"time"

	pexels "github.com/jdxj/pexels-sdk-go"
	wh "github.com/jdxj/wallhaven-sdk-go"
)

type WallpaperSource interface {
	GetWallpaper(context.Context) (string, error)
}

func newWallhavenSource() *wallhavenSource {
	return &wallhavenSource{
		whc: wh.NewClient(),
	}
}

type wallhavenSource struct {
	whc *wh.Client
}

func (ws *wallhavenSource) GetWallpaper(ctx context.Context) (string, error) {
	searchRsp, err := ws.whc.Search(ctx, &wh.SearchReq{
		Category: wh.People | wh.Anime | wh.General,
		Purity:   wh.SFW,
		Sorting:  wh.Random,
		Order:    wh.Desc,
	})
	if err != nil {
		return "", err
	}
	if len(searchRsp.Wallpapers) == 0 {
		return "", ErrWallpaperNotFound
	}
	return searchRsp.Wallpapers[0].Path, nil
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
	rand.Seed(time.Now().UnixNano())
	pl, err := ps.pc.CuratedPhotos(ctx, &pexels.CuratedReq{
		Pagination: pexels.Pagination{
			Page:    rand.Intn(8000) + 1,
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
