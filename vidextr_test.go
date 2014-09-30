package videoparse

import (
	"fmt"
	"strings"
	"testing"
)

func TestYouTube(t *testing.T) {
	url, err := YouTube("x28mzDuwQyk")
	if(err != nil) {
		t.Error(err.Error())
	}
	if !strings.Contains(url, "googlevideo.com") {
		t.Error("Failed")
	}
	fmt.Println(url)
}

func TestDailyMotion(t *testing.T) {
	url, err := DailyMotion("xftgry")
	if(err != nil) {
		t.Error(err.Error())
	}
	if !strings.Contains(url, "dailymotion.com") {
		t.Error("Failed")
	}
	fmt.Println(url)
}

func TestVimeo(t *testing.T) {
	url, err := Vimeo("79596219")
	if(err != nil) {
		t.Error(err.Error())
	}
	if !strings.Contains(url, "vimeocdn.com") {
		t.Error("Failed")
	}
	fmt.Println(url)
}
