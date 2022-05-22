package main

import "context"

type poller struct {
	index int
	wss   []WallpaperSource
}

func (p *poller) registerWallpaperSource(ws WallpaperSource) {
	p.wss = append(p.wss, ws)
}

func (p *poller) GetWallpaper(ctx context.Context) (string, error) {
	defer func() {
		p.index++
	}()

	i := p.index % len(p.wss)
	ws := p.wss[i]
	return ws.GetWallpaper(ctx)
}
