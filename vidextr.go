package videoparse

import (
	"time"
	"regexp"
	"strings"
	"net/url"
	"encoding/json"

	goreq "github.com/franela/goreq"
)


func YouTube(id string) (string, error) {
	res, err := goreq.Request{
		Uri:          "http://www.youtube.com/get_video_info?video_id=" + id + "&el=detailpage&ps=default",
		Timeout:      5 * time.Second,
		MaxRedirects: 1,
		UserAgent:    "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:25.0) Gecko/20100101 Firefox/25.0",
	}.Do()

	if err != nil {
		return "", err
	}

	body, _ := res.Body.ToString()

	getItems := func(query string) map[string]string {
		items := make(map[string]string)
		s := strings.Split(query, "&")
		for _, item := range s {
			m := strings.Split(item, "=")
			v, _ := url.QueryUnescape(m[1])
			items[m[0]] = v
		}
		return items
	}

	getUrl := func(query string) string {
		s := strings.Split(query, "&")
		for _, item := range s {
			m := strings.Split(item, "=")
			if m[0] == "url" {
				v, _ := url.QueryUnescape(m[1])
				return v
			}
		}
		return ""
	}

	items := getItems(body)
	url := getUrl(items["url_encoded_fmt_stream_map"])

	return url, nil
}


func DailyMotion(id string) (string, error) {
	res, err := goreq.Request{
		Uri:          "http://dailymotion.com/embed/video/" + id,
		Timeout:      5 * time.Second,
		MaxRedirects: 1,
		UserAgent:    "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:25.0) Gecko/20100101 Firefox/25.0",
	}.Do()

	if err != nil {
		return "", err
	}

	body, _ := res.Body.ToString()
	reDM := regexp.MustCompile(`(?m)var info = ({.*?}),$`)

	re := reDM.FindAllStringSubmatch(body, -1)
	if len(re) > 0 {
		info := re[0][1]

		var data map[string]interface{}
		err = json.Unmarshal([]byte(info), &data)
		if err != nil {
			return "", err
		}

		url := data["stream_h264_url"].(string)
		return url, nil
	}

	return "", nil
}


func Vimeo(id string) (string, error) {
	res, err := goreq.Request{
		Uri:          "http://player.vimeo.com/video/" + id,
		Timeout:      5 * time.Second,
		MaxRedirects: 1,
		UserAgent:    "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:25.0) Gecko/20100101 Firefox/25.0",
	}.Do()

	if err != nil {
		return "", err
	}

	body, _ := res.Body.ToString()
	reVM := regexp.MustCompile(`,a=({.*?});`)
	re := reVM.FindAllStringSubmatch(body, -1)
	if len(re) > 0 {
		a := re[0][1]

		var data map[string]interface{}
		err = json.Unmarshal([]byte(a), &data)
		if err != nil {
			return "", err
		}

		request, ok := data["request"].(map[string]interface{})
		if ok {
			files := request["files"].(map[string]interface{})
			h264 := files["h264"].(map[string]interface{})
			sd := h264["sd"].(map[string]interface{})
			url := sd["url"].(string)
			return url, nil
		}
	}

	return "", nil
}
