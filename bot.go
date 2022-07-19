package main

import (
	"context"
	"errors"
	"net/url"
	"path/filepath"
	"time"

	"github.com/go-resty/resty/v2"
	jike "github.com/jdxj/jike-sdk-go"
	"github.com/robfig/cron/v3"
)

var (
	ErrWallpaperNotFound = errors.New("wallpaper not found")
)

func New() *Bot {
	bot := &Bot{
		rc:     resty.New(),
		jkc:    jike.NewClient(),
		cron:   cron.New(),
		poller: &poller{},
	}

	bot.poller.registerWallpaperSource(newWallhavenSource())
	bot.poller.registerWallpaperSource(newPexelsSource(conf.PeAPIKey))
	bot.poller.registerWallpaperSource(newUnsplashSource(conf.UnsplashAK, conf.UnsplashSK))
	return bot
}

type Bot struct {
	rc   *resty.Client
	jkc  *jike.Client
	cron *cron.Cron

	poller *poller
}

func (b *Bot) PostWallpaper() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	urlPath, err := b.poller.GetWallpaper(ctx)
	if err != nil {
		logger.Errorf("%s\n", err)
		return
	}

	httpURL, err := url.Parse(urlPath)
	if err != nil {
		logger.Errorf("%s\n", err)
		return
	}
	absolutePath := filepath.Join(conf.CachePath, httpURL.Path)

	_, err = b.rc.R().
		SetOutput(absolutePath).
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

	utReq, err := jike.NewUploadTokenReq(absolutePath, jike.PIC)
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
		Filename:    absolutePath,
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
	logger.Infof("post %s ok", urlPath)
}

func (b *Bot) Run() {
	_, err := b.cron.AddFunc(conf.Spec, b.PostWallpaper)
	if err != nil {
		logger.Errorf("%s", err)
		return
	}
	b.cron.Run()
}
