package rpc_data

import (
	"code.chenji.com/cj/scan_paper/utils/lists"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/scan_paper_crop_image/scan_paper_crop_image_model"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/scan_paper_crop_image/scan_paper_crop_image_rpc"
)

type PageImage struct {
	ImageKey string
	PageNum  int32
}

func convPageImage(pageImage *PageImage) *scan_paper_crop_image_model.PageImage {
	if pageImage == nil {
		return nil
	}
	return &scan_paper_crop_image_model.PageImage{
		ImageKey: pageImage.ImageKey,
		PageNum:  pageImage.PageNum,
	}
}

type CropRect struct {
	PageNum  int32
	Sequence int32
	Rect     *Rect
}

func convRectFromCropRect(rect *Rect) *scan_paper_crop_image_model.Rect {
	if rect == nil {
		return nil
	}
	return &scan_paper_crop_image_model.Rect{
		X:      rect.X,
		Y:      rect.Y,
		Width:  rect.W,
		Height: rect.H,
	}
}

func convCropRect(cropRect *CropRect) *scan_paper_crop_image_model.CropRect {
	if cropRect == nil {
		return nil
	}
	return &scan_paper_crop_image_model.CropRect{
		PageNum:  cropRect.PageNum,
		Sequence: cropRect.Sequence,
		Rect:     convRectFromCropRect(cropRect.Rect),
	}
}

type MarkingPointPosition struct {
	MarkingPointId int64
	AnswerRects    []*CropRect
	FillRects      []*CropRect
}

func convMarkingPointPosition(markingPointPosition *MarkingPointPosition) *scan_paper_crop_image_model.MarkingPointPosition {

	if markingPointPosition == nil {
		return nil
	}
	return &scan_paper_crop_image_model.MarkingPointPosition{
		MarkingPointId: markingPointPosition.MarkingPointId,
		AnswerRects:    lists.Map(markingPointPosition.AnswerRects, convCropRect),
		FillRects:      lists.Map(markingPointPosition.FillRects, convCropRect),
	}
}

type ScanPaperCropImageParam struct {
	PageImages            []*PageImage            // 可供裁剪的图像信息
	MarkingPointPositions []*MarkingPointPosition // 判分点点的位置信息
	OutPath               string                  // 指定输出文件的路径
}

func NewScanPaperCropImageReq(param *ScanPaperCropImageParam) *scan_paper_crop_image_rpc.ScanPaperCropImageReq {
	if param == nil {
		return nil
	}
	return &scan_paper_crop_image_rpc.ScanPaperCropImageReq{
		PageImages:            lists.Map(param.PageImages, convPageImage),
		MarkingPointPositions: lists.Map(param.MarkingPointPositions, convMarkingPointPosition),
		Option: &scan_paper_crop_image_rpc.ScanPaperCropImageReq_Option{
			OutPath: param.OutPath,
		},
	}
}

type ResultImageItem struct {
	Sequence int32
	ImageKey string
	PageNum  int32
}

func convResultImageItem(resultImageItem *scan_paper_crop_image_model.ResultImageItem) *ResultImageItem {
	if resultImageItem == nil {
		return nil
	}
	return &ResultImageItem{
		Sequence: resultImageItem.Sequence,
		ImageKey: resultImageItem.ImageKey,
		PageNum:  resultImageItem.PageNum,
	}
}

type ScanPaperCropImageResult struct {
	MarkingPointCropResults []*MarkingPointCropImageResult // 判分点的裁剪结果
}

func NewScanPaperCropImageResult(resp *scan_paper_crop_image_rpc.ScanPaperCropImageResp) *ScanPaperCropImageResult {
	if resp == nil {
		return nil
	}
	return &ScanPaperCropImageResult{
		MarkingPointCropResults: lists.Map(resp.MarkingPointCropResults, convMarkingPointCropImageResult),
	}
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

type MarkingPointCropImageResult struct {
	MarkingPointId     int64              // 判分点ID
	AnswerUri          string             // 作答区域截图的Tos Key
	FillBlankUri       string             // 填空题完整区域截图的Tos Key（可能为空）
	AnswerSliceUris    []*ResultImageItem // 作答区域截图的Tos
	FillBlankSliceUris []*ResultImageItem // 填空题完整区域截图的Tos（可能为空）
}

func convMarkingPointCropImageResult(markingPointCropImageResult *scan_paper_crop_image_model.MarkingPointCropResult) *MarkingPointCropImageResult {

	if markingPointCropImageResult == nil {
		return nil
	}
	return &MarkingPointCropImageResult{
		MarkingPointId:     markingPointCropImageResult.MarkingPointId,
		AnswerUri:          markingPointCropImageResult.JoinAnswerImageKey,
		FillBlankUri:       markingPointCropImageResult.JoinFillFullImageKey,
		AnswerSliceUris:    lists.Map(markingPointCropImageResult.AnswerImageSlices, convResultImageItem),
		FillBlankSliceUris: lists.Map(markingPointCropImageResult.FillFullImageSlices, convResultImageItem),
	}
}
