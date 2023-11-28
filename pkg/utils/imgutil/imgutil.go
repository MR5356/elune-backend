package imgutil

import (
	"encoding/base64"
	"errors"
	"github.com/gabriel-vasile/mimetype"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

var miniTypeToBasePrefix = map[string]string{
	"image/jpeg":                "data:image/jpeg;base64,",
	"image/png":                 "data:image/png;base64,",
	"image/gif":                 "data:image/gif;base64,",
	"image/webp":                "data:image/webp;base64,",
	"image/bmp":                 "data:image/bmp;base64,",
	"image/svg+xml":             "data:image/svg+xml;base64,",
	"image/x-icon":              "data:image/x-icon;base64,",
	"text/plain; charset=utf-8": "data:text/plain;charset=utf-8;base64,",
}

func ImgLinkToBase64(imgUrl string) (string, error) {
	logrus.Debugf("imgUrl: %s", imgUrl)
	//获取远端图片
	res, err := http.Get(imgUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// 读取获取的[]byte数据
	data, _ := io.ReadAll(res.Body)

	miniType := mimetype.Detect(data).String()
	logrus.Debugf("miniType: %s", miniType)

	if mt, ok := miniTypeToBasePrefix[miniType]; ok {
		return mt + base64.StdEncoding.EncodeToString(data), nil
	} else {
		return "", errors.New("unsupported image type")
	}
}
