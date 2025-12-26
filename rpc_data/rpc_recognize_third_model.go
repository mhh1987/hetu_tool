package rpc_data

import (
	"fmt"
	"sort"
	"strconv"

	"code.chenji.com/cj/scan_paper/utils/lists"
	"code.chenji.com/cj/scan_paper/utils/sets"
	thriftcommon "code.chenji.com/jk/thriftidlgen/kitex_gen/edu/ark/common"
	"code.chenji.com/jk/thriftidlgen/kitex_gen/edu/ark/paper_parse"
	"code.chenji.com/jk/thriftidlgen/kitex_gen/edu/ark/paper_parse/answer_sheet_recognize"
	"code.chenji.com/jk/thriftidlgen/kitex_gen/edu/ark/paper_parse/sheet_recognize"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/common"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/resource/resource_enum"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/scan_paper/scan_paper_model"
)

// ===============================
// 识别涂抹区参数与结果的转换
// ===============================

// RecognizeSmearContentParam 识别涂抹区参数
type RecognizeSmearContentParam struct {
	PaperInfo      *Paper                // 试卷ID
	SheetInfo      *Sheet                // 批分方式
	MarkingPoints  []*MarkingPoint       // 判分点
	ImageInfo      *ImageInfo            // 图片数据
	CorrectionType common.CorrectionType // 批改类型
}

type ImageInfo struct {
	ImageId      int64  // 图片Id
	Url          string // 图片URL
	StoreKey     string // 图片TOS存储Key
	PageNum      int32  // 页码
	UnionExamId  int64  // 联合考试Id
	ExamId       int64  // 考试Id
	BatchId      int64  // 批次Id
	ImageGroupId int64  // 图片组Id
	TenantId     int64  // 租户Id
}

// NewSmearRegionRecognizeRequest 创建识别涂抹区参数
func NewSmearRegionRecognizeRequest(param *RecognizeSmearContentParam) *sheet_recognize.SmearRegionRecognizeRequest {
	if param == nil || param.SheetInfo == nil || param.PaperInfo == nil || param.MarkingPoints == nil || len(param.MarkingPoints) == 0 || param.ImageInfo == nil {
		return nil
	}
	imageInfo := param.ImageInfo
	sheetInfo := param.SheetInfo
	markPoints := param.MarkingPoints
	correctionType := param.CorrectionType
	imageIdStr := fmt.Sprintf("%d", imageInfo.ImageId)
	// 获取当前页的position
	pagePositions := make([]*SheetPositioning, 0, len(markPoints))
	for _, position := range sheetInfo.Positions { // 获取当前页的position
		if position.PageNum != imageInfo.PageNum {
			continue
		}
		pagePositions = append(pagePositions, position)
	}

	pageMarkingPoints := make([]*MarkingPoint, 0, len(markPoints))
	// 获取当前页的判分点
	for _, point := range markPoints { // 遍历所有判分点
		for _, position := range pagePositions { // 遍历当前页的position
			if point.Id == position.ItemId && position.ItemType == resource_enum.SheetItemType_SheetItemTypeMarkingPoint { // 判分点在当前页
				//if position.PositionType == resource_enum.PositionType_PositionTypeTitleDivision || position.PositionType == resource_enum.PositionType_PositionTypeChosen { // 选做题图改区
				if position.PositionType == resource_enum.PositionType_PositionTypeTitleDivision || position.PositionType == resource_enum.PositionType_PositionTypeOption {
					pageMarkingPoints = append(pageMarkingPoints, point)
					break
				}
			}
		}
	}

	req := &sheet_recognize.SmearRegionRecognizeRequest{
		ImageUrl:              imageInfo.Url,
		BizType:               convBizType(correctionType),
		StudentNoType:         ConvStudentNoType(sheetInfo.StudentCodeType),
		AbsentRegion:          getAbsentRegion(pagePositions),
		BarcodeRegion:         getBarcodeRegion(pagePositions, sheetInfo.StudentCodeType),
		StudentSmearRegions:   getStudentSmearRegions(pagePositions),
		ChoiceQuestionRegions: getChoiceQuestionRegions(pagePositions, pageMarkingPoints, sheetInfo.ObjectiveDirection),
		FillScoreRegions:      getFillScoreRegions(pagePositions, pageMarkingPoints),
		AnswerScoreRegions:    getAnswerScoreRegions(pagePositions, pageMarkingPoints),
		ChosenGroups:          getChosenGroups(pagePositions, pageMarkingPoints),
		PaperId:               &sheetInfo.PaperId,
		ImageId:               &imageIdStr,
		MonitorInfo: &sheet_recognize.MonitorInfo{
			UnionExamId:  imageInfo.UnionExamId,
			ExamId:       imageInfo.ExamId,
			BatchId:      imageInfo.BatchId,
			ImageGroupId: imageInfo.ImageGroupId,
			SchoolId:     imageInfo.TenantId,
		},
		PageNo: &imageInfo.PageNum,
		//RecognizeType:         nil,   // 不传默认识别全部内容
		//UserId:                nil,   // 暂时不需要
		//DeviceId:              nil,   // 暂时不需要
		//BatchId:               nil,   // 暂时不需要
	}
	recognizeTypes := make([]sheet_recognize.RecognizeType, 0, 10)
	if req.BarcodeRegion != nil || len(req.StudentSmearRegions) > 0 {
		recognizeTypes = append(recognizeTypes, sheet_recognize.RecognizeType_StudentNo)
	}
	if req.AbsentRegion != nil {
		recognizeTypes = append(recognizeTypes, sheet_recognize.RecognizeType_Absent)
	}
	if len(req.ChoiceQuestionRegions) > 0 {
		recognizeTypes = append(recognizeTypes, sheet_recognize.RecognizeType_Select)
	}
	if len(req.FillScoreRegions) > 0 {
		recognizeTypes = append(recognizeTypes, sheet_recognize.RecognizeType_Fill)
	}
	if len(req.AnswerScoreRegions) > 0 {
		recognizeTypes = append(recognizeTypes, sheet_recognize.RecognizeType_Answer_Score)
	}
	if len(req.ChosenGroups) > 0 {
		recognizeTypes = append(recognizeTypes, sheet_recognize.RecognizeType_Chosen)
	}
	req.RecognizeType = recognizeTypes
	return req
}

// 转换批改类型
func convBizType(correctionType common.CorrectionType) sheet_recognize.BizType {

	switch correctionType {
	case common.CorrectionType_CTWebCorrect:
		return sheet_recognize.BizType_Mark
	case common.CorrectionType_CTHandCorrect:
		return sheet_recognize.BizType_Handing
	default:
		return sheet_recognize.BizType_Unknown
	}
}

// ConvStudentNoType 转换学生考号类型
func ConvStudentNoType(codeType resource_enum.StudentCodeType) *sheet_recognize.StudentNoType {

	switch codeType {
	case resource_enum.StudentCodeType_StudentCodeTypeQrCode, resource_enum.StudentCodeType_StudentCodeTypeBarCode:
		return sheet_recognize.StudentNoTypePtr(sheet_recognize.StudentNoType_Barcode)
	case resource_enum.StudentCodeType_StudentCodeTypeStuNo, resource_enum.StudentCodeType_StudentCodeTypeSmallStuNo:
		return sheet_recognize.StudentNoTypePtr(sheet_recognize.StudentNoType_Smear)
	default:
		return sheet_recognize.StudentNoTypePtr(sheet_recognize.StudentNoType_Unknown)
	}
}

// 获取缺考区域
func getAbsentRegion(positions []*SheetPositioning) *sheet_recognize.Region {

	if len(positions) == 0 {
		return nil
	}

	var absentRegion *sheet_recognize.Region
	for _, position := range positions {
		if position.PositionType == resource_enum.PositionType_PositionTypeAbsent {
			if len(position.Rects) == 0 {
				continue
			}
			rect := position.Rects[0]
			absentRegion = &sheet_recognize.Region{
				LeftPoint: int32(rect.X),
				TopPoint:  int32(rect.Y),
				Width:     int32(rect.W),
				Height:    int32(rect.H),
			}
		}
	}
	return absentRegion
}

// 获取二维码/条形码区域
func getBarcodeRegion(positions []*SheetPositioning, codeType resource_enum.StudentCodeType) *sheet_recognize.Region {
	if len(positions) == 0 {
		return nil
	}
	var barcodeRegion *sheet_recognize.Region
	for _, position := range positions {
		if position.PositionType == resource_enum.PositionType_PositionTypeCandidateNo && (codeType == resource_enum.StudentCodeType_StudentCodeTypeQrCode || codeType == resource_enum.StudentCodeType_StudentCodeTypeBarCode) {
			if len(position.Rects) == 0 {
				continue
			}
			rect := position.Rects[0]
			barcodeRegion = &sheet_recognize.Region{
				LeftPoint: int32(rect.X),
				TopPoint:  int32(rect.Y),
				Width:     int32(rect.W),
				Height:    int32(rect.H),
			}
		}
	}
	return barcodeRegion
}

// 获取考试填涂区域
func getStudentSmearRegions(positions []*SheetPositioning) []*sheet_recognize.StudentNoSmearRegion {

	if len(positions) == 0 {
		return nil
	}

	var studentSmearRegions []*sheet_recognize.StudentNoSmearRegion
	groupMap := make(map[int32][]*SheetPositioning, 20)
	groupSet := sets.NewHashSet[int32]()
	for _, position := range positions {
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
		var smearRegions []*sheet_recognize.Region
		for _, position := range items {
			if len(position.Rects) == 0 {
				continue
			}
			rect := position.Rects[0]
			smearRegions = append(smearRegions, &sheet_recognize.Region{
				LeftPoint: int32(rect.X),
				TopPoint:  int32(rect.Y),
				Width:     int32(rect.W),
				Height:    int32(rect.H),
			})

		}
		studentSmearRegions = append(studentSmearRegions, &sheet_recognize.StudentNoSmearRegion{
			GroupId:      fmt.Sprintf("%d", group),
			SmearRegions: smearRegions,
		})
	}
	return studentSmearRegions
}

// 获取选择题(单选、多选、判断)区域
func getChoiceQuestionRegions(positions []*SheetPositioning, markPoints []*MarkingPoint, objectiveDirection resource_enum.DirectionStatus) []*sheet_recognize.SelectQuestionChoiceRegion {

	if len(positions) == 0 || markPoints == nil || len(markPoints) == 0 {
		return nil
	}
	var choiceQuestionRegions []*sheet_recognize.SelectQuestionChoiceRegion

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
		var choiceRegions []*sheet_recognize.Region
		for _, position := range items {
			if len(position.Rects) == 0 {
				continue
			}
			rect := position.Rects[0]
			choiceRegions = append(choiceRegions, &sheet_recognize.Region{
				LeftPoint: int32(rect.X),
				TopPoint:  int32(rect.Y),
				Width:     int32(rect.W),
				Height:    int32(rect.H),
			})

		}
		choiceQuestionRegions = append(choiceQuestionRegions, &sheet_recognize.SelectQuestionChoiceRegion{
			QuestionId:           markPointId,
			Single:               markPointSingleMap[markPointId],
			Vertical:             convObjectiveDirection(objectiveDirection),
			QuestionChoiceRegion: choiceRegions,
		})
	}

	return choiceQuestionRegions
}

type OptionDetail struct {
	OptionId int64
	Title    string
	Rect     *Rect
}

// 转换客观题方向
func convObjectiveDirection(direction resource_enum.DirectionStatus) bool {

	switch direction {
	case resource_enum.DirectionStatus_DirectionStatusVertical:
		return true
	default:
		return false
	}
}

// 获取填空题判分区域
func getFillScoreRegions(positions []*SheetPositioning, markPoints []*MarkingPoint) []*sheet_recognize.FillQuestionScoreRegion {

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
	var fillScoreRegions []*sheet_recognize.FillQuestionScoreRegion
	for markPointId, items := range markPointRectMap {
		// 找出当前组的选项区域
		if len(items) == 0 {
			continue
		}
		sort.Slice(items, func(i, j int) bool { // 组内按照顺序排序
			return items[i].Sequence < items[j].Sequence
		})
		var scoreRegions []*sheet_recognize.Region
		var scoreTypes []thriftcommon.PositionScoreType
		for _, position := range items {
			if len(position.Rects) == 0 {
				continue
			}
			rect := position.Rects[0]
			scoreRegions = append(scoreRegions, &sheet_recognize.Region{
				LeftPoint: int32(rect.X),
				TopPoint:  int32(rect.Y),
				Width:     int32(rect.W),
				Height:    int32(rect.H),
			})
			if position.ScoreExtra != nil {
				scoreTypes = append(scoreTypes, thriftcommon.PositionScoreType(position.ScoreExtra.ScoreType))
			}
		}
		fillScoreRegions = append(fillScoreRegions, &sheet_recognize.FillQuestionScoreRegion{
			QuestionId:   markPointId,
			ScoreRegions: scoreRegions,
			ScoreTypes:   scoreTypes,
		})
	}
	return fillScoreRegions
}

// 获取简答题判分区域
func getAnswerScoreRegions(positions []*SheetPositioning, markPoints []*MarkingPoint) []*sheet_recognize.AnswerQuestionScoreRegion {

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
	var answerScoreRegions []*sheet_recognize.AnswerQuestionScoreRegion
	for markPointId, items := range markPointRectMap {
		// 找出当前组的选项区域
		if len(items) == 0 {
			continue
		}
		sort.Slice(items, func(i, j int) bool { // 组内按照顺序排序
			return items[i].Sequence < items[j].Sequence
		})
		var scoreRegions []*sheet_recognize.Region
		var scoreTypes []*thriftcommon.PositionScore
		for _, position := range items {
			if len(position.Rects) == 0 {
				continue
			}
			rect := position.Rects[0]
			scoreRegions = append(scoreRegions, &sheet_recognize.Region{
				LeftPoint: int32(rect.X),
				TopPoint:  int32(rect.Y),
				Width:     int32(rect.W),
				Height:    int32(rect.H),
			})
			if position.ScoreExtra != nil {

				scoreTypes = append(scoreTypes, &thriftcommon.PositionScore{
					PositionScoreType: thriftcommon.PositionScoreType(position.ScoreExtra.ScoreType),
					Score:             position.ScoreExtra.Score,
				})
			}
		}
		answerScoreRegions = append(answerScoreRegions, &sheet_recognize.AnswerQuestionScoreRegion{
			QuestionId:   markPointId,
			ScoreRegions: scoreRegions,
			ScoreTypes:   scoreTypes,
		})
	}
	return answerScoreRegions
}

func getChosenGroups(positions []*SheetPositioning, markPoints []*MarkingPoint) []*sheet_recognize.ChosenGroup {
	if len(positions) == 0 || markPoints == nil || len(markPoints) == 0 {
		return nil
	}

	chosenGroups := make([]*sheet_recognize.ChosenGroup, 0, 20)

	// todo 待完善

	return chosenGroups
}

// RecognizeSmearContentResult 识别涂抹区结果
type RecognizeSmearContentResult struct {
	ImageId               int64                      `json:"image_id"`
	BarcodeResult         string                     `json:"barcode_result"`
	IsAbsent              bool                       `json:"is_absent"`
	ObjectiveResults      []*ObjectiveQuestionResult `json:"objective_results"`
	FillAnswerResults     []*FillQuestionResult      `json:"fill_answer_results"`
	AnswerScoreResults    []*AnswerQuestionResult    `json:"answer_score_results"`
	ChosenQuestionResults []*ChosenQuestionResult    `json:"chosen_question_results"`
	BarcodeText           *string
}

func NewRecognizeSmearContentResult(result *sheet_recognize.SmearRegionRecognizeResponse, params *RecognizeSmearContentParam) *RecognizeSmearContentResult {

	if result == nil {
		return nil
	}
	imageId := int64(0)
	if params != nil && params.ImageInfo != nil {
		imageId = params.ImageInfo.ImageId
	}

	return &RecognizeSmearContentResult{
		ImageId:               imageId,
		BarcodeResult:         result.BarcodeResult_,
		IsAbsent:              result.IsAbsent,
		ObjectiveResults:      lists.Map(result.ChoiceResults, convObjectiveQuestionResult),
		FillAnswerResults:     lists.Map(result.FillAnswerResults, convFillQuestionResult),
		AnswerScoreResults:    lists.Map(result.AnswerScoreResults, convAnswerQuestionResult),
		ChosenQuestionResults: lists.Map(result.ChosenQuestionResults, convChosenQuestionResult),
		BarcodeText:           result.BarcodeText,
	}
}

func convRegion(region *sheet_recognize.Region) *Rect {
	if region == nil {
		return nil
	}
	return &Rect{
		X: int64(region.LeftPoint),
		Y: int64(region.TopPoint),
		W: int64(region.Width),
		H: int64(region.Height),
	}
}

type ObjectiveQuestionResult struct {
	MarkingPointId int64  // 判分点Id
	Answer         string // 答案
	IsRecognized   bool   // 是否被识别
	Rect           *Rect  // 区域
}

func convObjectiveQuestionResult(objectiveQuestion *sheet_recognize.SelectQuestionAnswerRes) *ObjectiveQuestionResult {

	if objectiveQuestion == nil {
		return nil
	}
	return &ObjectiveQuestionResult{
		MarkingPointId: objectiveQuestion.QuestionId,
		Answer:         objectiveQuestion.PaintResults,
		IsRecognized:   objectiveQuestion.IsRecognized,
		Rect:           convRegion(objectiveQuestion.QuestionRegion),
	}
}

type FillQuestionResult struct {
	MarkingPointId int64                             // 判分点ID
	Result         sheet_recognize.FillAnswerResType // 结果
	IsMarked       bool                              // 是否有批改痕迹
}

func convFillQuestionResult(fillQuestion *sheet_recognize.FillQuestionAnswerRes) *FillQuestionResult {

	if fillQuestion == nil {
		return nil
	}
	return &FillQuestionResult{
		MarkingPointId: fillQuestion.QuestionId,
		Result:         fillQuestion.ResType,
		IsMarked:       fillQuestion.IsMarked,
	}
}

type AnswerQuestionResult struct {
	MarkingPointId int64                             // 判分点ID
	Score          string                            // 得分
	IsMarked       bool                              // 是否有批改痕迹
	Result         sheet_recognize.FillAnswerResType // 不批分模式返回的结果：正确、半对、错误
}

func convAnswerQuestionResult(answerQuestion *sheet_recognize.AnswerQuestionScoreRes) *AnswerQuestionResult {

	if answerQuestion == nil {
		return nil
	}
	return &AnswerQuestionResult{
		MarkingPointId: answerQuestion.QuestionId,
		Score:          answerQuestion.Score,
		IsMarked:       answerQuestion.IsMarked,
		Result:         answerQuestion.ResType,
	}
}

type ChosenQuestionResult struct {
	MarkingPointId int64 // 判分点ID
	IsMarked       bool  // 是否涂改
}

func convChosenQuestionResult(chosenQuestion *sheet_recognize.ChosenQuestionRes) *ChosenQuestionResult {

	if chosenQuestion == nil {
		return nil
	}
	return &ChosenQuestionResult{
		MarkingPointId: chosenQuestion.QuestionId,
		IsMarked:       chosenQuestion.IsMarked,
	}
}

// StudentQrCode 内容识别接口算法返回学生二维码的json字符串对应的数据结构
type StudentQrCode struct {
	Text []string `json:"text"`
}

// ==============================
// 答题卡匹配参数和结果转换
// ==============================

type MatchPaperSheetParam struct {
	PaperId    int64 // 试卷ID
	SheetId    int64 // 答题卡ID
	ImageGroup *ImageGroup
	TosDir     string // 矫正后学生作答页的 tos 上传 文件夹
}

type ImageGroup struct {
	Id     int64        // 图片组ID
	Images []*ScanImage // 图像列表
}

func ConvImageGroup(group *ImageGroup) *scan_paper_model.ImageGroup {

	if group == nil {
		return nil
	}
	images := make([]*scan_paper_model.ScanImage, 0, len(group.Images))
	for _, image := range group.Images {
		images = append(images, &scan_paper_model.ScanImage{
			ImageId:     image.ImageId,
			PageNum:     image.PageNum,
			OriginalUri: image.OriginalUri,
			RegulateUri: image.RegulateUri,
		})
	}

	return &scan_paper_model.ImageGroup{
		Id:     group.Id,
		Images: images,
	}
}

type ScanImage struct {
	ImageId     int64  // 图片ID
	PageNum     int32  // 页码
	OriginalUri string // 原始图片URL
	RegulateUri string // 矫正后图像URL
}

func NewAnsSheetMatchPairRequest(param *MatchPaperSheetParam) *paper_parse.AnsSheetMatchPairRequest {
	if param == nil || param.ImageGroup == nil || len(param.ImageGroup.Images) != 2 {
		return nil
	}

	images := param.ImageGroup.Images
	return &paper_parse.AnsSheetMatchPairRequest{
		PaperId:       param.PaperId,
		SheetId:       param.SheetId,
		First:         convPageInfo(images[0]),
		Second:        convPageInfo(images[1]),
		TosDir:        param.TosDir,
		InstitutionId: nil, // 暂不传值
	}
}

func convPageInfo(pageImage *ScanImage) *paper_parse.PageInfo {

	if pageImage == nil {
		return nil
	}
	imageIdStr := fmt.Sprintf("%d", pageImage.ImageId)
	return &paper_parse.PageInfo{
		Uri:     pageImage.OriginalUri,
		PageNo:  pageImage.PageNum,
		ImageId: &imageIdStr,
	}
}

type MatchPaperSheetResult struct {
	PaperId    int64
	ImageGroup *ImageGroup
}

func NewMatchPaperSheetResult(result *paper_parse.AnsSheetMatchPairResponse, req *paper_parse.AnsSheetMatchPairRequest, param *MatchPaperSheetParam) *MatchPaperSheetResult {

	if result == nil || req == nil || param == nil {
		return nil
	}
	images := make([]*ScanImage, 0, 2)
	images = append(images, convScanImage(result.First, req.First))   // 第一页
	images = append(images, convScanImage(result.Second, req.Second)) // 第二页
	return &MatchPaperSheetResult{
		PaperId: req.PaperId,
		ImageGroup: &ImageGroup{
			Id:     param.ImageGroup.Id,
			Images: images,
		},
	}
}

func convScanImage(pageResult *paper_parse.PageResult_, paramPage *paper_parse.PageInfo) *ScanImage {

	if pageResult == nil || paramPage == nil {
		return nil
	}
	imageId := int64(0)
	if paramPage.ImageId != nil {
		imageId, _ = strconv.ParseInt(paramPage.GetImageId(), 10, 64) // 转换成int64，err忽略，失败则转为0
	}
	return &ScanImage{
		ImageId:     imageId,
		PageNum:     pageResult.PageNo,
		OriginalUri: paramPage.Uri,
		RegulateUri: pageResult.RectifiedUri,
	}
}

// ==============================
// 答题卡匹配参数和结果转换
// ==============================

type MatchPaperSheetForceParam struct {
	PaperId   int64  // 试卷ID
	SheetId   int64  // 答题卡ID
	PageNum   int32  // 页码
	RawStuUri string // 学生作答页的原始图像URL
	TosDir    string // 矫正后学生作答页的 tos 上传 文件夹
	ImageId   int64  // 图片ID
}

func NewRecognizeAnswerSheetForceReq(param *MatchPaperSheetForceParam) *paper_parse.AnsSheetForceMatchRequest {
	if param == nil {
		return nil
	}
	imageIdStr := fmt.Sprintf("%d", param.ImageId)
	return &paper_parse.AnsSheetForceMatchRequest{
		PaperId:       param.PaperId,
		SheetId:       param.SheetId,
		PageNo:        param.PageNum,
		RawStuUri:     param.RawStuUri,
		TosDir:        param.TosDir,
		ImageId:       &imageIdStr,
		InstitutionId: nil, // 暂不传值
	}
}

type MatchPaperSheetForceResult struct {
	PaperId     int64  // 试卷ID
	ImageId     int64  // 图片ID
	OriginalUrl string // 原始图像URL
	RegulateUrl string // 矫正后图像URL
	PageNum     int32  // 页码
	Match       bool   // 是否匹配
}

func NewMatchPaperSheetForceResult(result *paper_parse.AnsSheetForceMatchResponse, param *MatchPaperSheetForceParam) *MatchPaperSheetForceResult {
	if result == nil || param == nil {
		return nil
	}
	return &MatchPaperSheetForceResult{
		PaperId:     param.PaperId,
		ImageId:     param.ImageId,
		OriginalUrl: param.RawStuUri,
		RegulateUrl: result.RectifiedStuUri,
		PageNum:     param.PageNum,
		Match:       result.Match,
	}

}

// ==============================
// 二维码识别参数和结果转换
// ==============================

type RecognizeQrCodeParam struct {
	Image   *ImageData // 图像数据
	ImageId int64      // 图片ID
}

// NewBarcodeRecognizeReq 转换二维码识别请求
func NewBarcodeRecognizeReq(param *RecognizeQrCodeParam) *paper_parse.BarcodeRecognizeRequest {

	if param == nil || param.Image == nil {
		return nil
	}
	imageIdStr := fmt.Sprintf("%d", param.ImageId)
	return &paper_parse.BarcodeRecognizeRequest{
		Image:   convImageData(param.Image),
		ImageId: &imageIdStr,
	}
}

type ImageData struct {
	DataUrl    *string // 图像数据URL
	DataUri    *string // 图像数据URI
	DataBase64 *string // 图像数据Base64
	DataBinary []byte  // 图像数据二进制
}

func convImageData(data *ImageData) *answer_sheet_recognize.Data {

	if data == nil {
		return nil
	}
	return &answer_sheet_recognize.Data{
		DataUrl:    data.DataUrl,
		DataUri:    data.DataUri,
		DataBase64: data.DataBase64,
		DataBinary: data.DataBinary,
	}
}

type RecognizeQrCodeResult struct {
	ImageId int64     // 图片Id
	QrCodes []*QrCode // 二维码列表
}

// NewRecognizeQrCodeResult 转换二维码识别结果
func NewRecognizeQrCodeResult(result *paper_parse.BarcodeRecognizeResponse, param *RecognizeQrCodeParam) *RecognizeQrCodeResult {
	if result == nil {
		return nil
	}
	return &RecognizeQrCodeResult{
		ImageId: param.ImageId,
		QrCodes: lists.Map(result.Items, convQrCode),
	}
}

// QrCode 二维码结构
type QrCode struct {
	Text string // 二维码内容
	Rect *Rect  // 二维码位置
}

func convQrCode(data *answer_sheet_recognize.Item) *QrCode {

	if data == nil {
		return nil
	}
	return &QrCode{
		Text: data.GetText(),
		Rect: &Rect{
			X: int64(data.GetRect().GetX()),
			Y: int64(data.GetRect().GetY()),
			W: int64(data.GetRect().GetW()),
			H: int64(data.GetRect().GetH()),
		},
	}
}

// CodeData 试卷二维码数据结构
type CodeData struct {
	PaperId string `json:"paper_id"`
	FileId  string `json:"file_id"`
	PageNum string `json:"page_num"`
}

func ConvToRpcDataImageGroup(group *scan_paper_model.ImageGroup) *ImageGroup {

	if group == nil {
		return nil
	}
	images := make([]*ScanImage, 0, len(group.Images))
	for _, image := range group.Images {
		images = append(images, &ScanImage{
			ImageId:     image.ImageId,
			PageNum:     image.PageNum,
			OriginalUri: image.OriginalUri,
			RegulateUri: image.RegulateUri,
		})
	}

	return &ImageGroup{
		Id:     group.Id,
		Images: images,
	}
}
