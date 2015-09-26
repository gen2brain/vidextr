package vidextr

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const userAgent string = "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:25.0) Gecko/20100101 Firefox/25.0"

func httpRequest(uri string, method string) (*http.Response, error) {
	timeout := time.Duration(6 * time.Second)

	dialTimeout := func(network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, timeout)
	}

	transport := http.Transport{
		Dial:            dialTimeout,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := http.Client{
		Transport: &transport,
	}

	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}

	req.Close = true
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", userAgent)

	res, err := httpClient.Do(req)
	if err != nil || res == nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Status Code %d received", res.StatusCode))
	}

	return res, nil
}

func YouTube(id string) (string, error) {
	res, err := httpRequest("http://www.youtube.com/get_video_info?video_id="+id+"&el=detailpage&ps=default", "GET")

	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	items, err := url.ParseQuery(string(body[:]))
	if err != nil {
		return "", err
	}

	streamMap, err := url.ParseQuery(items.Get("url_encoded_fmt_stream_map"))
	if err != nil {
		return "", err
	}

	uri := streamMap.Get("url")
	uri, err = url.QueryUnescape(uri)
	if err != nil {
		return "", err
	}

	return uri, nil
}

func DailyMotion(id string) (string, error) {
	res, err := httpRequest("http://dailymotion.com/embed/video/"+id, "GET")

	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	reDM := regexp.MustCompile(`(?mU)mp4","url":"(.*)"`)

	re := reDM.FindAllStringSubmatch(string(body[:]), -1)
	if len(re) > 0 {
		var u string
		if len(re) > 1 {
			u = re[1][1]
		} else {
			u = re[0][1]
		}
		u = strings.Replace(u, "\\", "", -1)

		return u, nil
	}

	return "", nil
}

func Vimeo(id string) (string, error) {
	res, err := httpRequest("http://player.vimeo.com/video/"+id, "GET")

	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	reVM := regexp.MustCompile(`var t=({.*?});`)

	re := reVM.FindAllStringSubmatch(string(body[:]), -1)
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
			u := sd["url"].(string)
			return u, nil
		}
	}

	return "", nil
}
