package images

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/hetu_tool/json_tool"
	"github.com/pkg/errors"
)

var JoinImageTaskTimeout = 3 * time.Minute

const ImageMergeQuality = 100

// CropImage 切图，key为目标图片的tos key，rect为想要切取的区域
func CropImage(key string, rect *Rect) string {
	if rect == nil {
		return ""
	}
	if strings.Contains(key, "?") {
		return fmt.Sprintf("%s/crop,w_%d,h_%d,x_%d,y_%d", key, rect.W, rect.H, rect.X, rect.Y)
	}
	return fmt.Sprintf("%s?x-tos-process=image/crop,w_%d,h_%d,x_%d,y_%d", key, rect.W, rect.H, rect.X, rect.Y)
}

// ScaleImage 缩放图片，rect为想要缩放的区域，scaleVal为外扩的像素值
func ScaleImage(rect *Rect, scaleVal int64) *Rect {
	if rect == nil {
		return nil
	}
	if scaleVal <= 0 {
		return rect
	}
	if rect.X <= 0 || (rect.X-scaleVal) <= 0 {
		rect.X = 0
	} else {
		rect.X -= scaleVal
	}
	if rect.Y <= 0 || (rect.Y-scaleVal) <= 0 {
		rect.Y = 0
	} else {
		rect.Y -= scaleVal
	}
	rect.W += 2 * scaleVal
	rect.H += 2 * scaleVal
	return rect
}

// JoinImage 图片拼接, 支持操作超时，超时时间为 JoinImageTaskTimeout
func JoinImage(joinImageItems []*JoinImageItem) ([]byte, error) {
	// 手动注册JPEG格式
	image.RegisterFormat("jpeg", "\xff\xd8", jpeg.Decode, jpeg.DecodeConfig)
	// 手动注册PNG格式
	image.RegisterFormat("png", "\x89PNG\r\n\x1a\n", png.Decode, png.DecodeConfig)
	resultDataChan := make(chan []byte)
	errChan := make(chan error)
	// 启动一个 goroutine 处理任务，超时时间为 JoinImageTaskTimeout
	go func(joinImageItems []*JoinImageItem) {
		defer func() {
			if err := recover(); err != nil {
				//logs.CtxError(context.Background(), "【图像工具类-图像拼接任务管理】捕获到panic，错误信息: %v", err)
			}
		}()
		dataBytes, err := joinImageTask(joinImageItems)
		if err != nil {
			errChan <- err
			return
		}
		resultDataChan <- dataBytes
	}(joinImageItems)
	select {
	case err := <-errChan:
		return nil, err
	case data := <-resultDataChan:
		return data, nil
	case <-time.After(JoinImageTaskTimeout):
		return nil, errors.New(fmt.Sprintf("【图像工具类-图像拼接任务管理】图片拼接任务超时, 入参: %s", json_tool.ToJson(joinImageItems)))
	}
}

// 图像拼接的具体实现逻辑
func joinImageTask(joinImageItems []*JoinImageItem) ([]byte, error) {
	if len(joinImageItems) == 0 {
		return nil, errors.New("【图像工具类-图像拼接】拼图函数入参校验，joinImageItems 为空，参数不合法!")
	}
	sort.Slice(joinImageItems, func(i, j int) bool {
		// 优先按照页码排序
		if joinImageItems[i].PageNum < joinImageItems[j].PageNum {
			return true
		}
		if joinImageItems[i].PageNum > joinImageItems[j].PageNum {
			return false
		}
		// 其次按照顺序排序
		return joinImageItems[i].Sequence < joinImageItems[j].Sequence
	})
	// 获取需要拼接的图像列表
	images := make([]image.Image, 0, len(joinImageItems))
	maxWidth := 0
	totalHeight := 0
	for _, item := range joinImageItems {

		tempImage, _, err := image.Decode(bytes.NewReader(item.CropImage))
		if err != nil {
			//logs.CtxError(context.Background(), fmt.Sprintf("【图像工具类-图像拼接】图像解码[image.Decode]出错，错误信息: %v", err))
			return nil, err
		}
		images = append(images, tempImage)
		if tempImage.Bounds().Dx() > maxWidth {
			maxWidth = tempImage.Bounds().Dx()
		}
		totalHeight += tempImage.Bounds().Dy()
	}

	// 新建一个画布
	mergedImg := image.NewRGBA(image.Rect(0, 0, maxWidth, totalHeight))
	// 绘制背景色
	draw.Draw(mergedImg, mergedImg.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	// 绘制图像
	y := 0
	for _, img := range images {
		draw.Draw(mergedImg, image.Rect(0, y, img.Bounds().Dx(), y+img.Bounds().Dy()), img, image.Point{}, draw.Over)
		y += img.Bounds().Dy()
	}
	// 保存图像
	buf := bytes.NewBuffer([]byte{})
	err := png.Encode(buf, mergedImg)
	//err := jpeg.Encode(buf, mergedImg, &jpeg.Options{Quality: ImageMergeQuality})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("【图像工具类-图像拼接】图像编码[Encode]出错，错误信息: %v", errors.WithStack(err)))
	}
	return buf.Bytes(), nil
}

// DownloadImage 下载图片, url为图片的完整url
func DownloadImage(ctx context.Context, url string) ([]byte, error) {

	// 发送 HTTP GET 请求获取图片
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//logs.CtxError(ctx, fmt.Sprintf("【图像工具类-下载图片】关闭响应体失败: %v", err))
		}
	}(resp.Body)

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("【图像工具类-下载图片】下载失败，状态码： %d, url: %s", resp.StatusCode, url))
	}
	// 读取响应内容
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

// DownloadImageWithRetry 下载图片, url为图片的完整url
func DownloadImageWithRetry(ctx context.Context, url string) ([]byte, error) {

	retryCount := 3
	for i := 0; i < retryCount; i++ {
		imageData, err := DownloadImage(ctx, url)
		if err == nil {
			return imageData, nil
		} else {
			//logs.CtxError(ctx, fmt.Sprintf("【图像工具类-下载图片】第%d次下载失败，错误信息: %v", i+1, err))
		}
	}
	return nil, errors.New(fmt.Sprintf("【图像工具类-下载图片】下载失败，重试次数: %d, url: %s", retryCount, url))
}

func GetExamImageTosDir(unionExamId, paperId, batchId int64) string {
	return fmt.Sprintf("scan_paper/join_image/%d/%d/%d", unionExamId, paperId, batchId)
}

func GetExamRectifiedTosDir(unionExamId, paperId, batchId int64) string {
	return fmt.Sprintf("scan_paper/scan_image/%d/%d/%d", unionExamId, paperId, batchId)
}

func GetImageStyle(uri string) (key string, process string) {

	splitStr := "?x-tos-process="
	if !strings.Contains(uri, splitStr) {
		return uri, ""
	}
	strs := strings.Split(uri, splitStr)
	if len(strs) < 2 {
		return uri, ""
	}
	return strs[0], fmt.Sprintf("%s", strs[1])
}

// CropJoinItemByRect 切图，单坐标区域切图
func CropJoinItemByRect(imageData image.Image, cropRect *CropRect) (*JoinImageItem, error) {

	if imageData == nil {
		return nil, errors.New(fmt.Sprintf("【切图-单区域】入参校验，imageItem为空，参数不合法!"))
	}
	if cropRect == nil || cropRect.Rect == nil {
		return nil, errors.New(fmt.Sprintf("【切图-单区域】坐标区域缺失，参数不合法!"))
	}
	resultImage, err := CropImageDataByRect(imageData, cropRect.Rect)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// 返回剪裁后的图像数据
	return &JoinImageItem{
		PageNum:   cropRect.PageNum,
		CropImage: resultImage,
		Sequence:  cropRect.Sequence,
	}, nil
}

func CropImageDataByRect(imageData image.Image, rect *Rect) ([]byte, error) {
	if imageData == nil {
		return nil, errors.New(fmt.Sprintf("【切图-单区域-通过坐标】入参校验，imageItem为空，参数不合法!"))
	}
	if rect == nil {
		return nil, errors.New(fmt.Sprintf("【切图-单区域-通过坐标】坐标区域缺失，参数不合法!"))
	}
	// 手动注册JPEG格式
	image.RegisterFormat("jpeg", "\xff\xd8", jpeg.Decode, jpeg.DecodeConfig)
	// 手动注册PNG格式
	image.RegisterFormat("png", "\x89PNG\r\n\x1a\n", png.Decode, png.DecodeConfig)

	// 确保剪裁区域在原始图像范围内
	if err := checkRect(DealPosition(imageData, rect)); err != nil {
		return nil, errors.New(fmt.Sprintf("【切图-单区域-通过坐标】图像坐标不合法(getImageDataByRect), 原始图像尺寸: %s, 切图区域坐标: %s", json_tool.ToJson(imageData.Bounds()), json_tool.ToJson(rect)))
	}

	cropRectImage := image.Rect(int(rect.X), int(rect.Y), int(rect.X+rect.W), int(rect.Y+rect.H))

	// 创建一个新的图像，用于保存剪裁结果
	croppedImg := image.NewRGBA(cropRectImage)

	// 在新图像上绘制剪裁区域
	draw.Draw(croppedImg, croppedImg.Bounds(), imageData, cropRectImage.Min, draw.Src)
	// 创建一个字节缓冲区
	var buf bytes.Buffer
	// 将剪裁后的图像编码为 JPEG 格式并写入缓冲区
	if err := jpeg.Encode(&buf, croppedImg, nil); err != nil {
		return nil, errors.WithStack(err)
	}
	return buf.Bytes(), nil
}

func GetImageByData(imageData []byte) (image.Image, error) {
	if len(imageData) == 0 {
		return nil, errors.New(fmt.Sprintf("【数据转化Image对象】入参校验，imageData为空，参数不合法!"))
	}
	// 手动注册JPEG格式
	image.RegisterFormat("jpeg", "\xff\xd8", jpeg.Decode, jpeg.DecodeConfig)
	// 手动注册PNG格式
	image.RegisterFormat("png", "\x89PNG\r\n\x1a\n", png.Decode, png.DecodeConfig)

	tempImage, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tempImage, nil
}

func DealPosition(imageData image.Image, rect *Rect) *Rect {

	if imageData == nil {
		return rect
	}
	if rect == nil {
		return rect
	}
	if rect.X < 0 {
		rect.X = 0
	}
	if rect.Y < 0 {
		rect.Y = 0
	}
	maxWidth := int64(imageData.Bounds().Max.X)
	maxHeight := int64(imageData.Bounds().Max.Y)
	if rect.W+rect.X > maxWidth {
		rect.W = maxWidth - rect.X
	}
	if rect.H+rect.Y > maxHeight {
		rect.H = maxHeight - rect.Y
	}
	if rect.W <= 0 || rect.H <= 0 {
		return nil
	}
	return rect
}

func checkRect(rect *Rect) error {
	if rect == nil {
		return errors.New(fmt.Sprintf("【坐标区域校验】入参校验，rect为空，参数不合法!"))
	}
	if rect.X < 0 || rect.Y < 0 || rect.W <= 0 || rect.H <= 0 {
		return errors.New(fmt.Sprintf("【坐标区域校验】坐标区域不合法，参数不合法!"))
	}
	return nil
}
