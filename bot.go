package main

import (
	"context"
	"errors"
	"path"
	"path/filepath"
	"time"

	"github.com/go-resty/resty/v2"
	jike "github.com/jdxj/jike-sdk-go"
	wh "github.com/jdxj/wallhaven-sdk-go"
	"github.com/robfig/cron/v3"
)

var (
	ErrWallpaperNotFound = errors.New("wallpaper not found")
)

func New() *Bot {
	bot := &Bot{
		rc:   resty.New(),
		whc:  wh.NewClient(),
		jkc:  jike.NewClient(jike.WithDebug(true)),
		cron: cron.New(),
	}
	return bot
}

type Bot struct {
	rc   *resty.Client
	whc  *wh.Client
	jkc  *jike.Client
	cron *cron.Cron
}

func (b *Bot) PostWallpaper() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	searchRsp, err := b.whc.Search(ctx, &wh.SearchReq{
		Category: wh.People | wh.Anime | wh.General,
		Purity:   wh.SFW,
		Sorting:  wh.TopList,
		Order:    wh.Desc,
		TopRange: wh.D1,
	})
	if err != nil {
		logger.Errorf("%s", err)
		return
	}
	if len(searchRsp.Wallpapers) == 0 {
		logger.Errorf("%s", ErrWallpaperNotFound)
		return
	}

	wallpaper := searchRsp.Wallpapers[0]
	urlPath := wallpaper.Path
	filename := path.Base(urlPath)
	fullPath := filepath.Join(conf.CachePath, filename)
	_, err = resty.New().R().
		SetOutput(fullPath).
		Get(urlPath)
	if err != nil {
		logger.Errorf("%s", err)
		return
	}

	_, err = b.jkc.LoginWithPhoneAndPassword(ctx, &jike.LoginWithPhoneAndPasswordReq{
		AreaCode:          conf.AreaCode,
		MobilePhoneNumber: conf.Phone,
		Password:          conf.Password,
	})
	if err != nil {
		logger.Errorf("%s\n", err)
		return
	}

	utReq, err := jike.NewUploadTokenReq(fullPath, jike.PIC)
	if err != nil {
		logger.Errorf("%s", err)
		return
	}
	utRsp, err := b.jkc.UploadToken(ctx, utReq)
	if err != nil {
		logger.Errorf("%s", err)
		return
	}

	uRsp, err := b.jkc.Upload(ctx, &jike.UploadReq{
		UploadToken: utRsp.UpToken,
		Filename:    fullPath,
	})
	if err != nil {
		logger.Errorf("%s", err)
		return
	}
	_, err = b.jkc.CreatePost(ctx, &jike.CreatePostReq{
		PictureKeys:          []string{uRsp.Key},
		SyncToPersonalUpdate: true,
		SubmitToTopic:        "59e58bea89ee3f0016b4d2c6",
	})
	if err != nil {
		logger.Errorf("%s", err)
		return
	}
	logger.Info("post %s ok", urlPath)
}

func (b *Bot) Run() {
	_, err := b.cron.AddFunc(conf.Spec, b.PostWallpaper)
	if err != nil {
		logger.Errorf("%s", err)
		return
	}
	b.cron.Run()
}
