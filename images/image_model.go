package images

import "image"

type Rect struct {
	X int64
	Y int64
	W int64
	H int64
}

type ImageInfo struct {
	PageNum int32
	Width   int64
	Height  int64
	Image   image.Image
}

type JoinImageItem struct {
	PageNum   int32  // 页码
	CropImage []byte // 切面图片对象
	Sequence  int32  // 顺序
}

type CropImageItem struct {
	MarkingPointId int64             // 判分点ID
	AnswerRects    map[int32][]*Rect // 当前判分点的作答区域截图数据 map[页码][]*Rect
	FillBlankFull  map[int32][]*Rect // 当前空所属的填空题完整区域截图数据(可能为空) map[页码][]*Rect
}

type CropImageResultItem struct {
	MarkingPointId int64  // 判分点ID
	AnswerUri      string // 作答区域截图的Tos Key
	FillBlankUri   string // 填空题完整区域截图的Tos Key（可能为空）
}

type CropRect struct {
	PageNum  int32
	Rect     *Rect
	Sequence int32
}

type ImageDimension struct {
	Width  int64
	Height int64
}
