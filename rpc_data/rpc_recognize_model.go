package rpc_data

import (
	"fmt"
	"sort"
	"strings"

	"code.chenji.com/cj/scan_paper/utils/lists"
	"code.chenji.com/cj/scan_paper/utils/sets"
	"code.chenji.com/pkg/common/tool"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/algo/algo_enum"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/algo/algo_model"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/algo/algo_rpc"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/common"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/resource/resource_enum"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/resource/resource_model"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/scan_paper/scan_paper_enum"
)

type SheetMatchParam struct {
	PaperId        int64                   // 试卷id
	SheetId        int64                   // 页码id
	TemplateImages []*Image                // 模版图片（必传页码）
	ImageGroup     *ImageGroup             // 需要匹配矫正的图片
	TosDir         string                  // 矫正后图像的tos文件目录
	SheetType      resource_enum.SheetType // 试卷类型
	TenantId       int64                   // 租户id
	UserId         int64                   // 用户id
	TraceInfo      *TraceInfo              // 追踪信息
	ExamCategory   common.ExamCategory     // 考试类型
}

func NewSheetMatchReq(param *SheetMatchParam) *algo_rpc.MatchAnsSheetPairReq {

	if param == nil || param.ImageGroup == nil || len(param.ImageGroup.Images) <= 0 {
		return nil
	}
	var firstPage *algo_model.Image
	if len(param.ImageGroup.Images) > 0 {
		firstPage = &algo_model.Image{
			ImageId:  param.ImageGroup.Images[0].ImageId,
			StoreKey: param.ImageGroup.Images[0].OriginalUri,
			PageNum:  tool.Ptr(param.ImageGroup.Images[0].PageNum),
		}
	}
	var secondPage *algo_model.Image
	if len(param.ImageGroup.Images) > 1 {
		secondPage = &algo_model.Image{
			ImageId:  param.ImageGroup.Images[1].ImageId,
			StoreKey: param.ImageGroup.Images[1].OriginalUri,
			PageNum:  tool.Ptr(param.ImageGroup.Images[1].PageNum),
		}
	}
	return &algo_rpc.MatchAnsSheetPairReq{
		PaperId:      param.PaperId,
		SheetId:      param.SheetId,
		First:        firstPage,
		Second:       secondPage,
		RectifiedDir: param.TosDir,
		SheetType:    param.SheetType,
		TenantId:     param.TenantId,
		UserId:       param.UserId,
		ExamCategory: param.ExamCategory,
		TraceInfo:    convTraceInfo(param.TraceInfo),
	}
}

type SheetMatchResult struct {
	ImageGroupId      int64               // 图片组id
	MatchImageResults []*MatchImageResult // 匹配结果
}

func NewSheetMatchResult(res *algo_rpc.MatchAnsSheetPairResp, req *algo_rpc.MatchAnsSheetPairReq) *SheetMatchResult {

	if res == nil || req == nil {
		return nil
	}
	result := &SheetMatchResult{
		ImageGroupId:      0,
		MatchImageResults: nil,
	}
	imageResults := make([]*MatchImageResult, 0, 2)
	if res.First != nil && req.First != nil {
		if res.First.PageNum > 0 && len(res.First.RectifiedKey) > 0 {
			imageResults = append(imageResults, &MatchImageResult{
				ImageId:     req.First.ImageId,
				OriginalUri: req.First.StoreKey,
				RegulateUri: res.First.RectifiedKey,
				PageNum:     res.First.PageNum,
				SheetId:     req.SheetId,
			})
		}
	}
	if res.Second != nil {
		if res.Second.PageNum > 0 && len(res.Second.RectifiedKey) > 0 {
			imageResults = append(imageResults, &MatchImageResult{
				ImageId:     req.Second.ImageId,
				OriginalUri: req.Second.StoreKey,
				RegulateUri: res.Second.RectifiedKey,
				PageNum:     res.Second.PageNum,
				SheetId:     req.SheetId,
			})
		}
	}
	result.MatchImageResults = imageResults
	return result
}

type MatchImageResult struct {
	ImageId     int64  // 图片id
	OriginalUri string // 图片uri
	RegulateUri string // 矫正后的图片uri
	PageNum     int32  // 页码
	SheetId     int64  // 试卷id
}

type Image struct {
	ImageId     int64  // 图片唯一id
	OriginalUri string // 图片uri
	PageNum     *int32 // 页码
}

type RecognizeStudentParam struct {
	ImageInfo    *ImageInfo          // 图片数据
	PaperInfo    *Paper              // 试卷ID
	SheetInfo    *Sheet              // 批分方式
	ExamCategory common.ExamCategory // 考试类型
	TraceInfo    *TraceInfo          // 追踪信息
	TenantId     int64               // 租户id
	UserId       int64               // 用户id
}

func NewRecognizeStudentReq(param *RecognizeStudentParam) *algo_rpc.RecognizeIdInfoReq {

	if param == nil || param.ImageInfo == nil || param.PaperInfo == nil || param.SheetInfo == nil {
		return nil
	}

	recognizeTypes := make([]algo_enum.IdType, 0, 4)
	recognizeTypes = append(recognizeTypes, algo_enum.IdType_HandStuntNo)
	recognizeTypes = append(recognizeTypes, algo_enum.IdType_HandName)
	req := &algo_rpc.RecognizeIdInfoReq{
		Image: &algo_model.Image{
			ImageId:  param.ImageInfo.ImageId,
			StoreKey: param.ImageInfo.StoreKey,
			PageNum:  tool.Ptr(param.ImageInfo.PageNum),
		},
		BarcodeRegion:  convBarcodeRegion(param.SheetInfo, param.ImageInfo.PageNum),
		SmearRegions:   convRegions(param.SheetInfo, param.ImageInfo.PageNum),
		RecognizeTypes: recognizeTypes,
		TenantId:       param.TenantId,
		UserId:         param.UserId,
		ExamCategory:   param.ExamCategory,
		TraceInfo:      convTraceInfo(param.TraceInfo),
	}
	if req.BarcodeRegion != nil {
		req.RecognizeTypes = append(req.RecognizeTypes, algo_enum.IdType_Barcode)
	}
	if len(req.SmearRegions) > 0 {
		req.RecognizeTypes = append(req.RecognizeTypes, algo_enum.IdType_Smear)
	}
	if param.SheetInfo.SheetType == resource_enum.SheetType_SheetTypeContentSheetExt {
		req.StudentInfoRegion = convBarcodeRegion(param.SheetInfo, param.ImageInfo.PageNum) // 三方卷的时候，这个区域就是学生信息区域
	}
	return req
}

func convBarcodeRegion(sheet *Sheet, pageNum int32) *resource_model.PositioningRect {

	if sheet == nil || len(sheet.Positions) <= 0 {
		return nil
	}
	for _, position := range sheet.Positions {
		if position == nil || position.PositionType != resource_enum.PositionType_PositionTypeCandidateNo {
			continue
		}
		if position.PageNum != pageNum {
			continue
		}
		if len(position.Rects) <= 0 {
			continue
		}
		return convRectToResource(position.Rects[0])
	}
	return nil
}

func convRegions(sheet *Sheet, pageNum int32) []*algo_model.StudentNoSmearRegion {

	if sheet == nil || len(sheet.Positions) <= 0 {
		return nil
	}

	var studentSmearRegions []*algo_model.StudentNoSmearRegion
	groupMap := make(map[int32][]*SheetPositioning, 20)
	groupSet := sets.NewHashSet[int32]()
	for _, position := range sheet.Positions {
		if position.PageNum != pageNum {
			continue
		}
		if position.PositionType == resource_enum.PositionType_PositionTypeTestNumberFillItem {
			groupMap[position.Group] = append(groupMap[position.Group], position)
			groupSet.Add(position.Group)
		}
	}
	groups := groupSet.ToSlice()
	// 按照组排序
	sort.Slice(groups, func(i, j int) bool {
		return groups[i] < groups[j]
	})
	for _, group := range groups {
		items, ok := groupMap[group]
		if !ok {
			continue
		}
		if len(items) == 0 {
			continue
		}
		sort.Slice(items, func(i, j int) bool { // 组内按照顺序排序
			return items[i].Sequence < items[j].Sequence
		})
		var smearRegions []*resource_model.PositioningRect
		for _, position := range items {
			if len(position.Rects) == 0 {
				continue
			}
			rect := position.Rects[0]
			smearRegions = append(smearRegions, convRectToResource(rect))

		}
		studentSmearRegions = append(studentSmearRegions, &algo_model.StudentNoSmearRegion{
			GroupId: fmt.Sprintf("%d", group),
			Rects:   smearRegions,
		})
	}
	return studentSmearRegions
}

type RecognizeStudentResult struct {
	IdInfos []*IdInfo // 识别结果
}

func NewRecognizeStudentResult(resp *algo_rpc.RecognizeIdInfoResp) *RecognizeStudentResult {

	if resp == nil || len(resp.IdInfos) <= 0 {
		return nil
	}
	result := &RecognizeStudentResult{
		IdInfos: lists.Map(resp.IdInfos, convIdInfo),
	}
	return result
}

type IdInfo struct {
	IdType  algo_enum.IdType
	Rect    *Rect
	Content string
}

func convIdInfo(info *algo_model.IdInfo) *IdInfo {

	if info == nil {
		return nil
	}
	result := &IdInfo{
		IdType:  info.IdType,
		Content: info.Content,
	}
	if info.Rect != nil {
		result.Rect = &Rect{
			X: int64(info.Rect.X),
			Y: int64(info.Rect.Y),
			W: int64(info.Rect.W),
			H: int64(info.Rect.H),
		}
	}
	return result
}

type RecognizeContentParam struct {
	PaperInfo      *Paper                // 试卷ID
	SheetInfo      *Sheet                // 批分方式
	MarkingPoints  []*MarkingPoint       // 判分点
	ImageInfo      *ImageInfo            // 图片数据
	CorrectionType common.CorrectionType // 批改类型
	TraceInfo      *TraceInfo            // 追踪信息
	ExamCategory   common.ExamCategory   // 考试类型
	TenantId       int64                 // 租户id
	UserId         int64                 // 用户id
}

func NewRecognizeContentReq(param *RecognizeContentParam) *algo_rpc.RecognizeHandwritingContentReq {
	if param == nil || param.SheetInfo == nil || param.PaperInfo == nil || param.ImageInfo == nil {
		return nil
	}
	isExam := param.ExamCategory == common.ExamCategory_ECExam || param.ExamCategory == common.ExamCategory_ECUnknown
	pageMarkPoints := make([]*MarkingPoint, 0, 20)
	if isExam {
		pageMarkPoints = GetPageNumMarkingPointForExam(param.MarkingPoints, param.SheetInfo.Positions, param.ImageInfo.PageNum)
	} else {
		pageMarkPoints = GetPageNumMarkingPoint(param.MarkingPoints, param.SheetInfo.Positions, param.ImageInfo.PageNum)
	}

	req := &algo_rpc.RecognizeHandwritingContentReq{
		PaperId:               param.PaperInfo.Id,
		SheetId:               param.SheetInfo.Id,
		PageNum:               param.ImageInfo.PageNum,
		Page:                  convPageFromImage(param.ImageInfo),
		AbsentRegion:          convAbsentRegion(param.SheetInfo.Positions, param.ImageInfo.PageNum),
		ObjectiveSmearRegions: convObjectiveSmearRegions(param.SheetInfo.Positions, pageMarkPoints, param.SheetInfo.ObjectiveDirection),
		ObjectiveRegions:      convObjectiveAnswerAreaRegions(param.SheetInfo.Positions, pageMarkPoints),
		FillinRegions:         convFillRegions(param.SheetInfo.Positions, pageMarkPoints),
		EassyRegions:          convEssayScoreRegions(param.SheetInfo.Positions, pageMarkPoints),
		TenantId:              param.TenantId,
		UserId:                param.UserId,
		ExamCategory:          param.ExamCategory,
		TraceInfo:             convTraceInfo(param.TraceInfo),
	}
	if param.SheetInfo.SheetType == resource_enum.SheetType_SheetTypeContentSheet || param.SheetInfo.SheetType == resource_enum.SheetType_SheetTypeContentSheetMix { // 如果当前为三方卷时，则不需要传选项涂抹区域
		req.ObjectiveSmearRegions = nil
	}
	return req
}

func convPageFromImage(imageInfo *ImageInfo) *algo_model.Image {

	if imageInfo == nil {
		return nil
	}
	return &algo_model.Image{
		ImageId:  imageInfo.ImageId,
		StoreKey: imageInfo.StoreKey,
		PageNum:  tool.Ptr(imageInfo.PageNum),
	}
}

func convAbsentRegion(positions []*SheetPositioning, pageNum int32) *resource_model.PositioningRect {

	if len(positions) == 0 {
		return nil
	}
	for _, position := range positions {
		if position == nil || position.PositionType != resource_enum.PositionType_PositionTypeAbsent {
			continue
		}
		if position.PageNum != pageNum {
			continue
		}
		if len(position.Rects) <= 0 {
			continue
		}
		rt := position.Rects[0]
		return convRectToResource(rt)
	}
	return nil
}

func convObjectiveSmearRegions(positions []*SheetPositioning, markPoints []*MarkingPoint, objectiveDirection resource_enum.DirectionStatus) []*algo_model.ObjectiveSmearRegion {

	if len(positions) == 0 || markPoints == nil || len(markPoints) == 0 {
		return nil
	}
	var regions []*algo_model.ObjectiveSmearRegion

	// 找出客观题的判分点
	objectiveMarkPoints := make([]*MarkingPoint, 0, len(markPoints))
	markPointSingleMap := make(map[int64]bool, 20)
	for _, item := range markPoints {
		if len(item.Answers) == 0 {
			continue
		}
		answer := item.Answers[0]
		if answer.Type == resource_enum.AnswerType_AnswerTypeSingleChoice || answer.Type == resource_enum.AnswerType_AnswerTypeMultiChoice || answer.Type == resource_enum.AnswerType_AnswerTypeTrueFalse {
			objectiveMarkPoints = append(objectiveMarkPoints, item)

			if answer.Type == resource_enum.AnswerType_AnswerTypeSingleChoice {
				markPointSingleMap[item.Id] = true
			} else {
				markPointSingleMap[item.Id] = false
			}
		}

	}
	// 客观题的判分点ID列表
	objectiveMarkPointIds := lists.Map(objectiveMarkPoints, func(item *MarkingPoint) int64 {
		return item.Id
	})
	// 找出当前判分点的选项区域
	markPointRectMap := make(map[int64][]*SheetPositioning, 20) // map[markPointId][]*SheetPositioning
	for _, position := range positions {
		if !lists.IsContain(objectiveMarkPointIds, position.ItemId) {
			continue
		}
		if position.PositionType == resource_enum.PositionType_PositionTypeOption {
			markPointRectMap[position.ItemId] = append(markPointRectMap[position.ItemId], position)
		}
	}

	for markPointId, items := range markPointRectMap {
		// 找出当前组的选项区域
		if len(items) == 0 {
			continue
		}
		sort.Slice(items, func(i, j int) bool { // 组内按照顺序排序
			return items[i].Sequence < items[j].Sequence
		})

		var optionRegions []*resource_model.PositioningRect
		for _, item := range items {
			optionRegions = append(optionRegions, lists.Map(item.Rects, convResourcePositioningRect)...)
		}
		regions = append(regions, &algo_model.ObjectiveSmearRegion{
			Id:            markPointId,
			Single:        markPointSingleMap[markPointId],
			Direction:     objectiveDirection,
			OptionRegions: optionRegions,
		})
	}

	return regions
}

func convResourceSheetPositioning(position *SheetPositioning) *resource_model.SheetPositioning {

	if position == nil {
		return nil
	}

	return &resource_model.SheetPositioning{
		Id:                      position.Id,
		PageId:                  position.PageId,
		PaperQuestionId:         position.PaperQuestionId,
		PaperQuestionAnswerId:   position.PaperQuestionAnswerId,
		ItemId:                  position.ItemId,
		ItemType:                position.ItemType,
		PositionType:            position.PositionType,
		Rects:                   lists.Map(position.Rects, convResourcePositioningRect),
		PageNum:                 position.PageNum,
		Group:                   position.Group,
		Sequence:                position.Sequence,
		PaperQuestionUuid:       position.PaperQuestionUuid,
		PaperQuestionAnswerUuid: position.PaperQuestionAnswerUuid,
		ItemUuid:                position.ItemUuid,
		ScoreExtra:              convResourceScoreExtra(position.ScoreExtra),
	}
}

func convResourcePositioningRect(rect *PositioningRect) *resource_model.PositioningRect {

	if rect == nil {
		return nil
	}
	return &resource_model.PositioningRect{
		X: rect.X,
		Y: rect.Y,
		W: rect.W,
		H: rect.H,
	}
}

func convResourceScoreExtra(extra *PositioningScore) *resource_model.PositioningScore {

	if extra == nil {
		return nil
	}
	return &resource_model.PositioningScore{
		ScoreType: extra.ScoreType,
		Score:     extra.Score,
	}
}

func convObjectiveAnswerAreaRegions(positions []*SheetPositioning, markPoints []*MarkingPoint) []*algo_model.ObjectiveQuestionRegion {

	if len(positions) == 0 || markPoints == nil || len(markPoints) == 0 {
		return nil
	}
	var regions []*algo_model.ObjectiveQuestionRegion
	// 找出客观题的判分点
	objectiveMarkPoints := make([]*MarkingPoint, 0, len(markPoints))
	markPointSingleMap := make(map[int64]bool, 20)
	for _, item := range markPoints {
		if len(item.Answers) == 0 {
			continue
		}
		answer := item.Answers[0]
		if answer.Type == resource_enum.AnswerType_AnswerTypeSingleChoice || answer.Type == resource_enum.AnswerType_AnswerTypeMultiChoice || answer.Type == resource_enum.AnswerType_AnswerTypeTrueFalse {
			objectiveMarkPoints = append(objectiveMarkPoints, item)

			if answer.Type == resource_enum.AnswerType_AnswerTypeSingleChoice {
				markPointSingleMap[item.Id] = true
			} else {
				markPointSingleMap[item.Id] = false
			}
		}

	}
	// 客观题的判分点ID列表
	objectiveMarkPointIds := lists.Map(objectiveMarkPoints, func(item *MarkingPoint) int64 {
		return item.Id
	})
	// 找出当前判分点的选项区域
	markPointRectMap := make(map[int64][]*SheetPositioning, 20) // map[markPointId][]*SheetPositioning
	for _, position := range positions {
		if !lists.IsContain(objectiveMarkPointIds, position.ItemId) {
			continue
		}
		if position.PositionType == resource_enum.PositionType_PositionTypeQuestionAnsweringArea {
			markPointRectMap[position.ItemId] = append(markPointRectMap[position.ItemId], position)
		}
	}
	for markPointId, positionList := range markPointRectMap {
		if len(positionList) <= 0 {
			continue
		}
		if positions[0] == nil || len(positions[0].Rects) <= 0 {
			continue
		}
		regions = append(regions, &algo_model.ObjectiveQuestionRegion{
			Id:   markPointId,
			Rect: convRectToResource(positions[0].Rects[0]),
		})
	}
	return regions
}

func convFillRegions(positions []*SheetPositioning, markPoints []*MarkingPoint) []*algo_model.FillinQuestionRegion {
	if len(positions) == 0 || markPoints == nil || len(markPoints) == 0 {
		return nil
	}

	fillMarkPoints := make([]*MarkingPoint, 0, len(markPoints))
	for _, item := range markPoints {
		if len(item.Answers) == 0 {
			continue
		}
		answer := item.Answers[0]
		if answer.Type == resource_enum.AnswerType_AnswerTypeFillin && !item.IsMerged {
			fillMarkPoints = append(fillMarkPoints, item)
		}
	}
	// 判分点ID列表
	markPointIds := lists.Map(fillMarkPoints, func(item *MarkingPoint) int64 {
		return item.Id
	})
	// 找出当前判分点的选项区域
	markPointRectMap := make(map[int64][]*SheetPositioning, 20) // map[markPointId][]*SheetPositioning
	for _, position := range positions {
		if !lists.IsContain(markPointIds, position.ItemId) {
			continue
		}
		if position.PositionType == resource_enum.PositionType_PositionTypeTitleDivision {
			markPointRectMap[position.ItemId] = append(markPointRectMap[position.ItemId], position)
		}
	}
	// 组装填空题识别信息
	var fillScoreRegions []*algo_model.FillinQuestionRegion
	for markPointId, items := range markPointRectMap {
		// 找出当前组的选项区域
		if len(items) == 0 {
			continue
		}
		sort.Slice(items, func(i, j int) bool { // 组内按照顺序排序
			return items[i].Sequence < items[j].Sequence
		})
		var scoreRegions []*resource_model.PositioningRect
		var scoreTypes []resource_enum.PositionScoreType
		for _, position := range items {
			if len(position.Rects) == 0 {
				continue
			}
			rect := position.Rects[0]
			scoreRegions = append(scoreRegions, convRectToResource(rect))
			if position.ScoreExtra != nil {
				scoreTypes = append(scoreTypes, position.ScoreExtra.ScoreType)
			}
		}
		fillScoreRegions = append(fillScoreRegions, &algo_model.FillinQuestionRegion{
			Id:         markPointId,
			Rects:      scoreRegions,
			ScoreTypes: scoreTypes,
		})
	}
	return fillScoreRegions
}

func convEssayScoreRegions(positions []*SheetPositioning, markPoints []*MarkingPoint) []*algo_model.EassyQuestionRegion {

	if len(positions) == 0 || markPoints == nil || len(markPoints) == 0 {
		return nil
	}
	answerMarkPoints := make([]*MarkingPoint, 0, len(markPoints))
	for _, markPoint := range markPoints {

		if len(markPoint.Answers) == 0 {
			continue
		}
		answer := markPoint.Answers[0]
		if answer.Type == resource_enum.AnswerType_AnswerTypeEssay || (answer.Type == resource_enum.AnswerType_AnswerTypeFillin && markPoint.IsMerged) {
			answerMarkPoints = append(answerMarkPoints, markPoint)
		}
	}

	// 判分点ID列表
	markPointIds := lists.Map(answerMarkPoints, func(item *MarkingPoint) int64 {
		return item.Id
	})
	// 找出当前判分点的选项区域
	markPointRectMap := make(map[int64][]*SheetPositioning, 20) // map[markPointId][]*SheetPositioning
	for _, position := range positions {
		if !lists.IsContain(markPointIds, position.ItemId) {
			continue
		}
		if position.PositionType == resource_enum.PositionType_PositionTypeTitleDivision {
			markPointRectMap[position.ItemId] = append(markPointRectMap[position.ItemId], position)
		}
	}
	// 组装答题题识别信息
	var answerScoreRegions []*algo_model.EassyQuestionRegion
	for markPointId, items := range markPointRectMap {
		// 找出当前组的选项区域
		if len(items) == 0 {
			continue
		}
		sort.Slice(items, func(i, j int) bool { // 组内按照顺序排序
			return items[i].Sequence < items[j].Sequence
		})
		var scoreRegions []*resource_model.PositioningRect
		var positionScores []*algo_model.PositionScore
		for _, position := range items {
			if len(position.Rects) == 0 {
				continue
			}
			rect := position.Rects[0]
			scoreRegions = append(scoreRegions, convRectToResource(rect))
			if position.ScoreExtra != nil {
				positionScores = append(positionScores, &algo_model.PositionScore{
					ScoreType: position.ScoreExtra.ScoreType,
					Score:     position.ScoreExtra.Score,
				})
			}
		}
		answerScoreRegions = append(answerScoreRegions, &algo_model.EassyQuestionRegion{
			Id:             markPointId,
			Rects:          scoreRegions,
			PositionScores: positionScores,
		})
	}
	return answerScoreRegions
}

type RecognizeContentResult struct {
	ImageId            int64                          `json:"image_id"`
	BarcodeResult      string                         `json:"barcode_result"`
	IsAbsent           bool                           `json:"is_absent"`
	ObjectiveResults   []*ObjectiveMarkingPointResult `json:"objective_results"`
	FillAnswerResults  []*FillMarkingPointResult      `json:"fill_answer_results"`
	AnswerScoreResults []*AnswerMarkingPointResult    `json:"answer_score_results"`
	BarcodeText        *string                        `json:"barcode_text"`
}

func NewRecognizeContentResult(resp *algo_rpc.RecognizeHandwritingContentResp, param *RecognizeContentParam) *RecognizeContentResult {
	if resp == nil {
		return nil
	}
	return &RecognizeContentResult{
		ImageId:            param.ImageInfo.ImageId,
		IsAbsent:           resp.IsAbsent,
		ObjectiveResults:   convObjectiveMarkingPointResult(resp.ObjectiveResults),
		FillAnswerResults:  lists.Map(resp.FillinResults, convFillAnswerResult),
		AnswerScoreResults: lists.Map(resp.EassyResults, convAnswerMarkingPointResult),
		BarcodeText:        nil,
		BarcodeResult:      "",
	}
}

type ObjectiveMarkingPointResult struct {
	MarkingPointId        int64                 // 判分点Id
	IsAnswered            bool                  // 是否有作答
	Answer                string                // 学生作答
	StuHandwritingResults []*StuHandwritingRes  // 学生作答识别结果
	IsRecognized          bool                  // 是否被识别[客观题的：算法]
	IsMarked              bool                  // 是否有或归属到批改内容
	MarkContents          []*MarkHandwritingRes // 批改内容
}

func convObjectiveMarkingPointResult(choiceHandwritingResults []*algo_model.ObjectiveQuestionRes) []*ObjectiveMarkingPointResult {

	results := make([]*ObjectiveMarkingPointResult, 0, len(choiceHandwritingResults))
	if len(choiceHandwritingResults) > 0 {
		for _, item := range choiceHandwritingResults {
			if item == nil {
				continue
			}
			answerContents := make([]string, 0, len(item.StuRes))
			stuRects := make([]*Rect, 0, len(item.StuRes))
			for _, stu := range item.StuRes {
				if len(stu.Content) <= 0 {
					continue
				}
				answerContents = append(answerContents, stu.Content)
				stuRects = append(stuRects, convRect(stu.Rect))
			}
			result := &ObjectiveMarkingPointResult{
				MarkingPointId:        item.Id,
				IsAnswered:            item.IsAnswered,
				Answer:                AssembleAnswerContent(answerContents),
				IsRecognized:          !(len(answerContents) == 0),
				StuHandwritingResults: lists.Map(item.StuRes, convStuHandwritingRes),
				IsMarked:              item.IsMarked,
				MarkContents:          lists.Map(item.MarkRes, convMarkHandwritingRes),
			}
			results = append(results, result)
		}
	}
	return results
}

type StuHandwritingRes struct {
	Rect    *Rect  // 学生答案区域
	Content string // 学生答案内容
}

func convStuHandwritingRes(stuRes *algo_model.StuHandwritingRes) *StuHandwritingRes {

	if stuRes == nil {
		return nil
	}
	return &StuHandwritingRes{
		Rect:    convRect(stuRes.Rect),
		Content: stuRes.Content,
	}
}

type MarkHandwritingRes struct {
	Rect         *Rect
	Content      string // 当symboltype为数字时为分值；其他情况为空。 根据MarkScore、ComputeScore选择的结果分值
	SymbolType   algo_enum.SymbolType
	MarkScore    string // 算法返回的分值
	ComputeScore string // 算法拆分得到的分数
}

func convMarkHandwritingRes(markRes *algo_model.MarkHandwritingRes) *MarkHandwritingRes {

	if markRes == nil {
		return nil
	}

	content := ""
	if markRes.MarkScore != nil {
		content = markRes.GetMarkScore()
	}
	if len(content) <= 0 && markRes.ComputeScore != nil {
		content = markRes.GetComputeScore()
	}
	if markRes.IsSmearScore && markRes.SmearResType == algo_enum.SmearScoreResType_SRTScore {
		content = markRes.GetMarkScore()
	}
	return &MarkHandwritingRes{
		Rect:         convRect(markRes.Rect),
		Content:      content,
		SymbolType:   markRes.SymbolType,
		MarkScore:    markRes.GetMarkScore(),
		ComputeScore: markRes.GetComputeScore(),
	}
}

//func convRectFromAlgo(rect *algo_model.Rect) *Rect {
//
//	if rect == nil {
//		return nil
//	}
//	return &Rect{
//		X: int64(rect.Left),
//		Y: int64(rect.Top),
//		W: int64(rect.Width),
//		H: int64(rect.Height),
//	}
//}

type FillMarkingPointResult struct {
	MarkingPointId        int64                       // 判分点ID
	Result                algo_enum.SmearScoreResType // 结果
	StuAnswer             string                      // 学生答案
	StuHandwritingResults []*StuHandwritingRes        // 学生作答识别结果
	IsMarked              bool                        // 是否有批改痕迹
	MarkContents          []*MarkHandwritingRes       // 批改内容
}

func convFillAnswerResult(result *algo_model.FillinQuestionRes) *FillMarkingPointResult {

	if result == nil {
		return nil
	}
	answerContents := make([]string, 0, len(result.StuRes))
	for _, stu := range result.StuRes {
		if len(stu.Content) <= 0 {
			continue
		}
		answerContents = append(answerContents, stu.Content)
	}
	fillResult := algo_enum.SmearScoreResType_SRTUnknown
	for _, markRe := range result.MarkRes {
		if !markRe.IsSmearScore {
			continue
		}
		fillResult = markRe.SmearResType
		if fillResult > 0 {
			break
		}
	}
	return &FillMarkingPointResult{
		MarkingPointId:        result.Id,
		Result:                fillResult,
		StuAnswer:             AssembleAnswerContent(answerContents),
		StuHandwritingResults: lists.Map(result.StuRes, convStuHandwritingRes),
		IsMarked:              result.IsMarked,
		MarkContents:          lists.Map(result.MarkRes, convMarkHandwritingRes),
	}
}

type AnswerMarkingPointResult struct {
	MarkingPointId        int64                 // 判分点ID
	StuAnswer             string                // 学生答案
	StuHandwritingResults []*StuHandwritingRes  // 学生作答识别结果
	IsMarked              bool                  // 是否有批改痕迹
	MarkContents          []*MarkHandwritingRes // 批改内容
}

func convAnswerMarkingPointResult(result *algo_model.EassyQuestionRes) *AnswerMarkingPointResult {

	if result == nil {
		return nil
	}
	answerContents := make([]string, 0, len(result.StuRes))
	for _, stu := range result.StuRes {
		if len(stu.Content) <= 0 {
			continue
		}
		answerContents = append(answerContents, stu.Content)
	}
	return &AnswerMarkingPointResult{
		MarkingPointId:        result.Id,
		StuAnswer:             AssembleAnswerContent(answerContents),
		StuHandwritingResults: lists.Map(result.StuRes, convStuHandwritingRes),
		IsMarked:              result.IsMarked,
		MarkContents:          lists.Map(result.MarkRes, convMarkHandwritingRes),
	}
}

type AlgoRecognizeQrCodeParam struct {
	UnionExamId  int64               // 联合考试ID
	PaperId      int64               // 试卷ID
	ImageGroup   *ImageGroup         // 图片组
	TenantId     int64               // 租户ID
	UserId       int64               // 用户ID
	ExamCategory common.ExamCategory // 考试类型
	TraceInfo    *TraceInfo          // 追踪信息
}

func NewRecognizeQrcodeReq(param *AlgoRecognizeQrCodeParam) *algo_rpc.RecognizeBarcodeReq {

	if param == nil || param.ImageGroup == nil || len(param.ImageGroup.Images) == 0 {
		return nil
	}
	algoImages := make([]*algo_model.Image, 0, 1)
	for _, image := range param.ImageGroup.Images {
		if image == nil {
			continue
		}
		storeKey := image.RegulateUri
		if len(storeKey) <= 0 {
			storeKey = image.OriginalUri
		}
		algoImages = append(algoImages, &algo_model.Image{
			ImageId:  image.ImageId,
			StoreKey: storeKey,
		})
	}
	return &algo_rpc.RecognizeBarcodeReq{
		Images:       algoImages,
		TenantId:     param.TenantId,
		UserId:       param.UserId,
		TraceInfo:    convTraceInfo(param.TraceInfo),
		ExamCategory: param.ExamCategory,
	}
}

type AlgoRecognizeQrcodeResult struct {
	BarcodeResults []*BarcodeResult
}

func NewAlgoRecognizeQrcodeResult(resp *algo_rpc.RecognizeBarcodeResp) *AlgoRecognizeQrcodeResult {

	if resp == nil {
		return nil
	}
	return &AlgoRecognizeQrcodeResult{
		BarcodeResults: lists.Map(resp.BarcodeResults, convBarcodeResult),
	}
}

type BarcodeResult struct {
	ImageId int64
	PageNum int32
	Text    string
}

func convBarcodeResult(barcodeResult *algo_rpc.RecognizeBarcodeResp_BarcodeResult) *BarcodeResult {

	if barcodeResult == nil {
		return nil
	}
	return &BarcodeResult{
		ImageId: barcodeResult.ImageId,
		PageNum: barcodeResult.PageNum,
		Text:    barcodeResult.Text,
	}
}

func ParseAnswerContent(content string) []string {

	return strings.Split(content, ",")
}

func AssembleAnswerContent(answers []string) string {

	if len(answers) == 1 {
		return answers[0]
	}
	return strings.Join(answers, ",")
}

func convRectToResource(rect *PositioningRect) *resource_model.PositioningRect {
	if rect == nil {
		return nil
	}
	return &resource_model.PositioningRect{
		X: rect.X,
		Y: rect.Y,
		W: rect.W,
		H: rect.H,
	}
}

func convRect(rect *resource_model.PositioningRect) *Rect {

	if rect == nil {
		return nil
	}
	return &Rect{
		X: int64(rect.X),
		Y: int64(rect.Y),
		W: int64(rect.W),
		H: int64(rect.H),
	}
}

type TraceInfo struct {
	UnionExamId  *int64
	ImageGroupId *int64
	ExamId       *int64
	PaperId      *int64
	SheetId      *int64
	ImageId      *int64
	QuestionId   *int64
	StudentId    *int64
}

func convTraceInfo(traceInfo *TraceInfo) *algo_model.TraceInfo {

	if traceInfo == nil {
		return nil
	}
	return &algo_model.TraceInfo{
		ExamId:     traceInfo.ExamId,
		PaperId:    traceInfo.PaperId,
		SheetId:    traceInfo.SheetId,
		ImageId:    traceInfo.ImageId,
		QuestionId: traceInfo.QuestionId,
		StudentId:  traceInfo.StudentId,
	}
}

type MatchAnsSheetForceParam struct {
	PaperId      int64
	SheetId      int64
	PageNum      int32
	RawStuUrl    string
	RectifiedDir string
	SheetType    resource_enum.SheetType
	TenantId     int64
	UserId       int64
	TraceInfo    *TraceInfo
	IsHomework   bool
	ExamCategory common.ExamCategory
}

type MatchAnsSheetForceResult struct {
	PageResult *PageResult
}

type PageResult struct {
	RectifiedKey string                     // 矫正图
	PageNum      int32                      // 页码
	Status       algo_enum.PaperMatchStatus // 匹配状态
}

func convPageResult(pageResult *algo_model.PageResult) *PageResult {

	if pageResult == nil {
		return nil
	}
	return &PageResult{
		RectifiedKey: pageResult.RectifiedKey,
		PageNum:      pageResult.PageNum,
		Status:       pageResult.Status,
	}
}

func NewMatchAnsSheetForceParam(req *MatchAnsSheetForceParam) *algo_rpc.MatchAnsSheetForceReq {

	if req == nil {
		return nil
	}
	return &algo_rpc.MatchAnsSheetForceReq{
		PaperId:      req.PaperId,
		SheetId:      req.SheetId,
		PageNum:      req.PageNum,
		RawStuUrl:    req.RawStuUrl,
		RectifiedDir: req.RectifiedDir,
		SheetType:    req.SheetType,
		TenantId:     req.TenantId,
		UserId:       req.UserId,
		TraceInfo:    convTraceInfo(req.TraceInfo),
		ExamCategory: req.ExamCategory,
	}
}

func NewMatchAnsSheetForceResult(resp *algo_rpc.MatchAnsSheetForceResp) *MatchAnsSheetForceResult {

	if resp == nil {
		return nil
	}
	return &MatchAnsSheetForceResult{
		PageResult: convPageResult(resp.PageResult),
	}
}

func ConvFillAnswerResType(fillResult *FillMarkingPointResult) scan_paper_enum.RecordResultType {

	if fillResult == nil {
		return scan_paper_enum.RecordResultType_RT_Unknown
	}
	switch fillResult.Result {
	case algo_enum.SmearScoreResType_SRTRight:
		return scan_paper_enum.RecordResultType_RT_Right
	case algo_enum.SmearScoreResType_SRTHalfRight:
		return scan_paper_enum.RecordResultType_RT_HalfRight
	case algo_enum.SmearScoreResType_SRTWrong:
		return scan_paper_enum.RecordResultType_RT_Wrong
	default:
		if !fillResult.IsMarked { // 没有标记默认为对
			return scan_paper_enum.RecordResultType_RT_Right
		}
		return scan_paper_enum.RecordResultType_RT_Unknown
	}
}

func GetScanTeacherTenantId(examInfo *Exam, userId int64) int64 {
	if examInfo == nil || len(examInfo.ScanTeachers) <= 0 {
		return 0
	}
	for _, teacher := range examInfo.ScanTeachers {
		if teacher == nil || teacher.TeacherId != userId {
			continue
		}
		return teacher.TenantId
	}
	return 0
}
