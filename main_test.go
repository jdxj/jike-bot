package main

import (
	"fmt"
	"os"
	"testing"
)

func TestBot_PostWallpaper(t *testing.T) {
	b := New()
	b.PostWallpaper()
}

func TestEnv(t *testing.T) {
	s := os.ExpandEnv("")
	fmt.Printf("%s\n", s)
}
