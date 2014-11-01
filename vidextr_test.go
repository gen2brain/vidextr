package vidextr

import (
	"fmt"
	"strings"
	"testing"
)

func TestYouTube(t *testing.T) {
	url, err := YouTube("x28mzDuwQyk")
	if err != nil {
		t.Error(err.Error())
	}
	if !strings.Contains(url, "googlevideo.com") {
		t.Error("Failed")
	}
	fmt.Println(url)
}

func TestDailyMotion(t *testing.T) {
	url, err := DailyMotion("xftgry")
	if err != nil {
		t.Error(err.Error())
	}
	if !strings.Contains(url, "dailymotion.com") {
		t.Error("Failed")
	}
	fmt.Println(url)
}

func TestVimeo(t *testing.T) {
	url, err := Vimeo("79596219")
	if err != nil {
		t.Error(err.Error())
	}
	if !strings.Contains(url, "vimeocdn.com") {
		t.Error("Failed")
	}
	fmt.Println(url)
}

func TestVK(t *testing.T) {
	url, err := VK("http://vk.com/video_ext.php?oid=204781865&id=167613610&hash=341cbf13949e8357&api_hash=141475443129710f13be729a4bb6")
	if err != nil {
		t.Error(err.Error())
	}
	if !strings.Contains(url, ".vk.me") {
		t.Error("Failed")
	}
	fmt.Println(url)
}
