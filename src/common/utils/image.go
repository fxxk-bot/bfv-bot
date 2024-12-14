package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	imagedraw "image/draw"
	"image/png"
	"io"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

func SvgToPng(svg io.Reader) (io.Reader, error) {
	// 读取 SVG 图标
	icon, err := oksvg.ReadIconStream(svg)
	if err != nil {
		return nil, err
	}

	// 获取 SVG 的视图框宽度和高度
	width, height := int(icon.ViewBox.W), int(icon.ViewBox.H)
	if width == 0 || height == 0 {
		return nil, fmt.Errorf("invalid SVG dimensions: width=%d, height=%d", width, height)
	}

	// 设置绘制目标区域
	icon.SetTarget(0, 0, icon.ViewBox.W, icon.ViewBox.H)

	// 创建一个新的 RGBA 图像
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充背景为白色
	imagedraw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.White}, image.Point{}, imagedraw.Src)

	// 创建 Rasterx Dasher 以绘制 SVG
	dasher := rasterx.NewDasher(width, height, rasterx.NewScannerGV(width, height, rgba, rgba.Bounds()))
	icon.Draw(dasher, 1)

	// 将 RGBA 图像编码为 PNG
	buf := &bytes.Buffer{}
	if err := png.Encode(buf, rgba); err != nil {
		return nil, err
	}

	return buf, nil
}
