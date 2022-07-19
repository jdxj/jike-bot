package main

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
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

func TestPathBase(t *testing.T) {
	u1 := "https://images.unsplash.com/photo-1570913149827-d2ac84ab3f9a?ixid=MnwzMzM4NTR8MHwxfGFsbHx8fHx8fHx8fDE2NTQxMzgwMDI&ixlib=rb-1.2.1"
	u1b := path.Base(u1)
	fmt.Printf("u1b: %s\n", u1b)
	u2 := "https://images.unsplash.com/photo-1570913149827-d2ac84ab3f9a?crop=entropy&cs=tinysrgb&fm=jpg&ixid=MnwzMzM4NTR8MHwxfGFsbHx8fHx8fHx8fDE2NTQxMzgwMDI&ixlib=rb-1.2.1&q=80"
	u2b := path.Base(u2)
	fmt.Printf("u2b: %s\n", u2b)

	url1, err := url.Parse(u1)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("url1: %s\n", url1.Path)
	fmt.Printf("join: %s\n", filepath.Join("/tmp", url1.Path))

}
