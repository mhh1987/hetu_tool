package rpc_data

import (
	"sort"

	"code.chenji.com/cj/scan_paper/utils/lists"
	"code.chenji.com/cj/scan_paper/utils/sets"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/resource/resource_enum"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/resource/resource_model"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/resource/resource_rpc"
)

type Paper struct {
	Id              int64                         // 题卷Id
	Name            string                        // 题卷名称
	SubName         string                        // 题卷副标题
	SubHeading      string                        // 三级标题
	Term            int32                         // 学期
	Stage           int32                         // 学段
	BaseGrade       int32                         // 基础年级
	BaseSubjects    []int32                       // 基础学科列表
	Score           string                        // 题卷总分
	UserId          int64                         // 所属用户Id
	TenantId        int64                         // 所属机构Id
	PaperSourceType resource_enum.PaperSourceType // 题卷来源类型
	PaperScoring    []*PaperScoring               // 题目判分配置
	ChoseGroups     []*PaperChoseGroup            // 选做题组
	Sheets          []*Sheet                      // 题卷关联题卡列表
}

type PaperNodeItem struct {
	PaperId           int64                       // 所属题卷Id
	ItemId            int64                       // 项Id, 根据PaperNodeType区分具体Id
	ParentId          int64                       // 项父节点Id
	StandaloneId      int64                       // 独立题Id
	NodeType          resource_enum.PaperNodeType // 题卷节点类型
	Name              string                      // 名称
	Title             string                      // 标题（name的层级拼接）
	Sequence          int32                       // 全局排序
	IsRelated         bool                        // 是否关联题
	IsOptional        bool                        // 是否选做题
	IsAnswerable      bool                        // 是否作答题
	IsMarkingPoint    bool                        // 是否判分点
	Score             string                      // 分数
	ScoreStep         string                      // 步进
	ScoreRule         *resource_model.ScoreRule   // 给分规则
	AnswerContent     string                      // 答案
	AnswerType        resource_enum.AnswerType    // 作答类型, 节点为作答时值有效
	QuestionId        int64                       // 关联的题目Id
	OriginId          int64                       // 原始题卷题目Id
	MarkingPointId    int64                       // 判分点Id, 项是判分点时值有效
	OptionCount       int32                       // 选项数量
	IsPartialMatching bool                        // 是否支持半对(部分给分)
	IsObjective       bool                        // 是否客观题(作答), 项为题或作答时值有效
	IsSubjective      bool                        // 是否主观题(作答), 项为题或作答时值有效
	IsMerged          bool                        // 是否合并判分
	OrderNum          string                      // 编号
}

func ConvPaperNodeItem(paperNodeItem *resource_model.PaperNodeItem) *PaperNodeItem {

	if paperNodeItem == nil {
		return nil
	}
	return &PaperNodeItem{
		PaperId:           paperNodeItem.PaperId,
		ItemId:            paperNodeItem.ItemId,
		ParentId:          paperNodeItem.ParentId,
		StandaloneId:      paperNodeItem.StandaloneId,
		NodeType:          paperNodeItem.NodeType,
		Name:              paperNodeItem.Name,
		Title:             paperNodeItem.Title,
		Sequence:          paperNodeItem.Sequence,
		IsRelated:         paperNodeItem.IsRelated,
		IsOptional:        paperNodeItem.IsOptional,
		IsAnswerable:      paperNodeItem.IsAnswerable,
		IsMarkingPoint:    paperNodeItem.IsMarkingPoint,
		Score:             paperNodeItem.Score,
		ScoreStep:         paperNodeItem.ScoreStep,
		ScoreRule:         paperNodeItem.ScoreRule,
		AnswerContent:     paperNodeItem.AnswerContent,
		AnswerType:        paperNodeItem.AnswerType,
		QuestionId:        paperNodeItem.QuestionId,
		OriginId:          paperNodeItem.OriginId,
		MarkingPointId:    paperNodeItem.MarkingPointId,
		OptionCount:       paperNodeItem.OptionCount,
		IsPartialMatching: paperNodeItem.IsPartialMatching,
		IsObjective:       paperNodeItem.IsObjective,
		IsSubjective:      paperNodeItem.IsSubjective,
		IsMerged:          paperNodeItem.IsMerged,
		OrderNum:          paperNodeItem.OrderNum,
	}
}

type PaperScoring struct {
	Id                    int64                     // Id
	PaperQuestionId       int64                     // 题卷题目Id
	PaperQuestionAnswerId int64                     // 题卷题目作答Id
	Score                 string                    // 满分
	ScoreStep             string                    // 步进
	ScoreRule             *resource_model.ScoreRule // 给分规则
	IsPartialMatching     bool                      // 是否支持半对(部分给分)
	IsMarkingPoint        bool                      // 是否判分点
	IsMerged              bool                      // 是否合并判分
	MergedStandaloneIds   []int64                   // 合并的题卷独立题Id集合
}

type PaperChoseGroup struct {
	Id               int64   // Id
	ChoseCount       int32   // 选做数量
	PaperQuestionIds []int64 // 选做范围(题卷题目Id列表)
}

func ConvPaper(paper *resource_model.Paper) *Paper {
	if paper == nil {
		return nil
	}
	return &Paper{
		Id:              paper.Id,
		Name:            paper.Name,
		SubName:         paper.SubName,
		SubHeading:      paper.SubHeading,
		Term:            paper.Term,
		Stage:           paper.Stage,
		BaseGrade:       paper.BaseGrade,
		BaseSubjects:    paper.BaseSubjects,
		Score:           paper.Score,
		UserId:          paper.UserId,
		TenantId:        paper.TenantId,
		PaperSourceType: paper.PaperSourceType,
		PaperScoring:    lists.Map(paper.Scorings, convPaperScoring),
		ChoseGroups:     lists.Map(paper.ChoseGroups, convPaperChoseGroup),
		Sheets:          lists.Map(paper.Sheets, ConvSheet),
	}
}

func convPaperScoring(paperScoring *resource_model.PaperScoring) *PaperScoring {

	if paperScoring == nil {
		return nil
	}
	return &PaperScoring{
		Id:                    paperScoring.Id,
		PaperQuestionId:       paperScoring.PaperQuestionId,
		PaperQuestionAnswerId: paperScoring.PaperQuestionAnswerId,
		Score:                 paperScoring.Score,
		ScoreStep:             paperScoring.ScoreStep,
		ScoreRule:             paperScoring.ScoreRule,
		IsPartialMatching:     paperScoring.IsPartialMatching,
		IsMarkingPoint:        paperScoring.IsMarkingPoint,
		IsMerged:              paperScoring.IsMerged,
		MergedStandaloneIds:   paperScoring.MergedStandaloneIds,
	}
}

func convPaperChoseGroup(paperChoseGroup *resource_model.PaperChoseGroup) *PaperChoseGroup {

	if paperChoseGroup == nil {
		return nil
	}
	return &PaperChoseGroup{
		Id:               paperChoseGroup.Id,
		ChoseCount:       paperChoseGroup.ChoseCount,
		PaperQuestionIds: paperChoseGroup.PaperQuestionIds,
	}
}

type MarkingPoint struct {
	Id                  int64                     // 判分点Id
	Title               string                    // 标题(默认同name, 作答时为名称路径，如: 1.(1))
	Name                string                    // 名称
	Score               string                    // 满分
	ScoreStep           string                    // 步进
	ScoreRule           *resource_model.ScoreRule // 给分规则
	IsPartialMatching   bool                      // 是否支持半对(部分给分)
	IsMerged            bool                      // 是否合并判分
	MergedStandaloneIds []int64                   // 合并的题卷独立题Id
	Answers             []*MarkingPointAnswer     // 作答信息
	OrderNum            string                    // 编号
	IsObjective         bool                      // 是否客观题
	IsSubjective        bool                      // 是否主观题
	PaperQuestionId     int64                     // 题卷题目Id
}

func ConvMarkingPoint(markingPoint *resource_model.MarkingPoint) *MarkingPoint {

	if markingPoint == nil {
		return nil
	}
	return &MarkingPoint{
		Id:                  markingPoint.Id,
		Title:               markingPoint.Title,
		Name:                markingPoint.Name,
		Score:               markingPoint.Score,
		ScoreStep:           markingPoint.ScoreStep,
		ScoreRule:           markingPoint.ScoreRule,
		IsPartialMatching:   markingPoint.IsPartialMatching,
		IsMerged:            markingPoint.IsMerged,
		MergedStandaloneIds: markingPoint.MergedStandaloneIds,
		Answers:             lists.Map(markingPoint.Answers, convMarkingPointAnswer),
		OrderNum:            markingPoint.OrderNum,
		IsObjective:         markingPoint.IsObjective,
		IsSubjective:        markingPoint.IsSubjective,
		PaperQuestionId:     markingPoint.PaperQuestionId,
	}
}

type MarkingPointAnswer struct {
	Type    resource_enum.AnswerType // 作答类型
	Content string                   // 答案
	//Options []*MarkingPointOption    // 当前判分点的选项信息
}

func convMarkingPointAnswer(answer *resource_model.MarkingPoint_Answer) *MarkingPointAnswer {
	if answer == nil {
		return nil
	}
	return &MarkingPointAnswer{
		Type:    answer.Type,
		Content: answer.Content,
		//Options: lists.Map(answer.Options, convMarkingPointOption),
	}
}

type MarkingPointOption struct {
	Id    int64  // 选项Id
	Title string // 选项标题
}

//func convMarkingPointOption(data *resource_model.MarkingPoint_Option) *MarkingPointOption {
//
//	if data == nil {
//		return nil
//	}
//	return &MarkingPointOption{
//		Id:    data.Id,
//		Title: data.Title,
//	}
//}

type Sheet struct {
	Id                 int64                     // 题卡Id
	PaperId            int64                     // 题卷Id
	SheetType          resource_enum.SheetType   // 题卡类型
	PrintFormat        resource_enum.PrintFormat // 打印板式
	MarkingMode        resource_enum.MarkingMode // 批分方式
	Column             int32                     // 栏数
	Version            int32                     // 版本号
	FileId             int64                     // 题卡完整pdf文件Id
	EditConfig         string                    // 端上编辑配置
	IsAvailable        bool                      // 是否可用
	IsSheetPageOneSide bool                      // 单面题卡页
	// 排版类枚举
	StudentCodeType    resource_enum.StudentCodeType // 考号类型
	ObjectiveDirection resource_enum.DirectionStatus // 客观题选项方向
	Pages              []*SheetPage                  // 题卡页列表
	Questions          []*SheetQuestion              // 题卡题目列表
	Positions          []*SheetPositioning           // 题卡定位点列表
}

type SheetPage struct {
	Id          int64                     // Id
	FileId      int64                     // pdf文件Id
	ImageFileId int64                     // 图片文件Id
	PageNum     int32                     // 页码
	Column      int32                     // 栏数
	PrintFormat resource_enum.PrintFormat // 打印板式
	SheetId     int64                     // 题卡Id
	Uuid        string                    // 业务交互UUID
}

type SheetQuestion struct {
	Id                int64  // Id
	PageId            int64  // 题卡页Id
	PaperQuestionId   int64  // 题卷题目Id
	Sequence          int32  // 页内顺序
	EditContent       string // 题卡可编辑内容
	Delta             string // 端上自定义内容
	PageNum           int32  // 题卡页页码
	PaperQuestionUuid string // 业务交互题卷题目UUID
}

type SheetPositioning struct {
	Id                      int64                       // Id
	PageId                  int64                       // 所在题卡页Id
	PaperQuestionId         int64                       // 题卡项关联的题卷题目Id
	PaperQuestionAnswerId   int64                       // 题卡项关联的题卷题目作答Id
	ItemId                  int64                       // 关联题卡项Id
	ItemType                resource_enum.SheetItemType // 关联题卡项类型
	PaperSectionId          int64                       // 关联题节点Id
	PaperScoringId          int64                       // 题卡项关联的题卷给分配置Id
	PositionType            resource_enum.PositionType  // 点位位置类型
	Rects                   []*PositioningRect          // 点位位置列表
	PageNum                 int32                       // 题卡页页码
	Group                   int32                       // 定位分组
	Sequence                int32                       // 全局排序
	PaperQuestionUuid       string                      // 业务交互题卷题目UUID
	PaperQuestionAnswerUuid string                      // 业务交互题卷题目作答UUID
	ItemUuid                string                      // 业务交互题卡项UUID
	ScoreExtra              *PositioningScore           // 打分框扩展信息
}

type PositioningScore struct {
	ScoreType resource_enum.PositionScoreType // 打分框类型
	Score     string                          // 打分框分数
}

type PositioningRect struct {
	X float32 // left
	Y float32 // top
	W float32 // width
	H float32 // height
}

func ConvSheet(sheet *resource_model.Sheet) *Sheet {

	if sheet == nil {
		return nil
	}
	return &Sheet{
		Id:                 sheet.Id,
		PaperId:            sheet.PaperId,
		SheetType:          sheet.SheetType,
		PrintFormat:        sheet.PrintFormat,
		MarkingMode:        sheet.MarkingMode,
		Column:             sheet.Column,
		Version:            sheet.Version,
		FileId:             sheet.FileId,
		EditConfig:         sheet.EditConfig,
		IsAvailable:        sheet.IsAvailable,
		IsSheetPageOneSide: sheet.IsSheetPageOneSide,
		StudentCodeType:    sheet.StudentCodeType,
		ObjectiveDirection: sheet.ObjectiveDirection,
		Pages:              lists.Map(sheet.Pages, convPages),
		Questions:          lists.Map(sheet.Questions, convSheetQuestion),
		Positions:          lists.Map(sheet.Positions, convSheetPositioning),
	}
}

func convPages(page *resource_model.SheetPage) *SheetPage {

	if page == nil {
		return nil
	}
	return &SheetPage{
		Id:          page.Id,
		FileId:      page.FileId,
		ImageFileId: page.ImageFileId,
		PageNum:     page.PageNum,
		Column:      page.Column,
		PrintFormat: page.PrintFormat,
		SheetId:     page.SheetId,
		Uuid:        page.Uuid,
	}
}

func convSheetQuestion(sheetQuestion *resource_model.SheetQuestion) *SheetQuestion {

	if sheetQuestion == nil {
		return nil
	}
	return &SheetQuestion{
		Id:                sheetQuestion.Id,
		PageId:            sheetQuestion.PageId,
		PaperQuestionId:   sheetQuestion.PaperQuestionId,
		Sequence:          sheetQuestion.Sequence,
		EditContent:       sheetQuestion.EditContent,
		Delta:             sheetQuestion.Delta,
		PageNum:           sheetQuestion.PageNum,
		PaperQuestionUuid: sheetQuestion.PaperQuestionUuid,
	}
}

func convSheetPositioning(position *resource_model.SheetPositioning) *SheetPositioning {

	if position == nil {
		return nil
	}
	return &SheetPositioning{
		Id:                      position.Id,
		PageId:                  position.PageId,
		PaperQuestionId:         position.PaperQuestionId,
		PaperQuestionAnswerId:   position.PaperQuestionAnswerId,
		ItemId:                  position.ItemId,
		ItemType:                position.ItemType,
		PaperSectionId:          position.PaperSectionId,
		PaperScoringId:          position.PaperScoringId,
		PositionType:            position.PositionType,
		Rects:                   lists.Map(position.Rects, convPositioningRect),
		PageNum:                 position.PageNum,
		Group:                   position.Group,
		Sequence:                position.Sequence,
		PaperQuestionUuid:       position.PaperQuestionUuid,
		PaperQuestionAnswerUuid: position.PaperQuestionAnswerUuid,
		ItemUuid:                position.ItemUuid,
		ScoreExtra:              convPositioningScore(position.ScoreExtra),
	}
}
func convPositioningScore(scoreExtra *resource_model.PositioningScore) *PositioningScore {
	if scoreExtra == nil {
		return nil
	}
	return &PositioningScore{
		ScoreType: scoreExtra.ScoreType,
		Score:     scoreExtra.Score,
	}
}

func convPositioningRect(rect *resource_model.PositioningRect) *PositioningRect {

	if rect == nil {
		return nil
	}
	w := rect.W
	if w <= 0 {
		w = 1
	}
	h := rect.H
	if h <= 0 {
		h = 1
	}
	return &PositioningRect{
		X: rect.X,
		Y: rect.Y,
		W: w,
		H: h,
	}
}

func ConvQuestionTypeByAnswerType(answerType resource_enum.AnswerType) resource_enum.QuestionType {

	switch answerType {
	case resource_enum.AnswerType_AnswerTypeSingleChoice:
		return resource_enum.QuestionType_QuestionTypeSingleChoice
	case resource_enum.AnswerType_AnswerTypeFillin:
		return resource_enum.QuestionType_QuestionTypeFillin
	case resource_enum.AnswerType_AnswerTypeEssay:
		return resource_enum.QuestionType_QuestionTypeEssay
	case resource_enum.AnswerType_AnswerTypeMultiChoice:
		return resource_enum.QuestionType_QuestionTypeMultiChoice
	case resource_enum.AnswerType_AnswerTypeTrueFalse:
		return resource_enum.QuestionType_QuestionTypeTrueFalse
	default:
		return resource_enum.QuestionType_QuestionTypeUnknown
	}
}

type PaperReferencedFileResult struct {
	IsReferenced bool                    // 文件是否为题卷所引用
	Page         *SheetPage              // 如果有引用关系, 返回文件对应的题卡页信息
	SheetType    resource_enum.SheetType // 题卡类型
}

func ConPaperReferencedFileResult(req *resource_rpc.IsPaperReferencedFileResponse) *PaperReferencedFileResult {
	if req == nil {
		return nil
	}
	return &PaperReferencedFileResult{
		IsReferenced: req.IsReferenced,
		Page:         convPages(req.Page),
		SheetType:    req.SheetType,
	}
}

type SheetBaseInfo struct {
	Id          int64                     // 题卡Id
	PaperId     int64                     // 题卷Id
	SheetType   resource_enum.SheetType   // 题卡类型
	PrintFormat resource_enum.PrintFormat // 打印板式
	MarkingMode resource_enum.MarkingMode // 批分方式
	Column      int32                     // 栏数
	Version     int32                     // 版本号
	FileId      int64                     // 题卡完整pdf文件Id
	EditConfig  string                    // 端上编辑配置
	IsAvailable bool                      // 是否可用
	// 排版类枚举
	StudentCodeType    resource_enum.StudentCodeType // 考号类型
	ObjectiveDirection resource_enum.DirectionStatus // 客观题选项方向
}

func ConvSheetBaseInfo(sheet *resource_model.Sheet) *SheetBaseInfo {
	if sheet == nil {
		return nil
	}
	return &SheetBaseInfo{
		Id:                 sheet.Id,
		PaperId:            sheet.PaperId,
		SheetType:          sheet.SheetType,
		PrintFormat:        sheet.PrintFormat,
		MarkingMode:        sheet.MarkingMode,
		Column:             sheet.Column,
		Version:            sheet.Version,
		FileId:             sheet.FileId,
		EditConfig:         sheet.EditConfig,
		IsAvailable:        sheet.IsAvailable,
		StudentCodeType:    sheet.StudentCodeType,
		ObjectiveDirection: sheet.ObjectiveDirection,
	}
}

// GetPageNumMarkingPoint 获取当前页涉及到的判分点
func GetPageNumMarkingPoint(markingPoints []*MarkingPoint, positions []*SheetPositioning, pageNum int32) []*MarkingPoint {
	if pageNum == 0 || len(markingPoints) <= 0 || len(positions) <= 0 {
		return nil
	}
	markingPointMap := make(map[int64]int32, 20) // map[判分点ID]页码
	for _, position := range positions {
		if position.PageNum != pageNum { // 跳过非当前页
			continue
		}
		if position.ItemType != resource_enum.SheetItemType_SheetItemTypeMarkingPoint { // 跳过非判分点的题目
			continue
		}
		if position.PositionType == resource_enum.PositionType_PositionTypeOption || position.PositionType == resource_enum.PositionType_PositionTypeTitleDivision || position.PositionType == resource_enum.PositionType_PositionTypeQuestionAnsweringArea {
			v, ok := markingPointMap[position.ItemId]
			if !ok { // 没有收集，则直接收集即可
				markingPointMap[position.ItemId] = position.PageNum
				continue
			}
			if position.PageNum < v { // 跨页则当前判分点归属到小的页码，即为归属到题目的起始页码
				markingPointMap[position.ItemId] = position.PageNum
			}
		}
	}
	//logs.Info(fmt.Sprintf("getPageNumMarkingPoint -> 页码:%d, 对应的判分点:%s", pageNum, tool.ToJson(markingPointIds)))
	results := make([]*MarkingPoint, 0, len(markingPointMap))
	for _, point := range markingPoints {
		if _, ok := markingPointMap[point.Id]; ok {
			results = append(results, point)
		}
	}
	return results
}

// GetPageNumMarkingPointForExam 获取当前页涉及到的判分点
func GetPageNumMarkingPointForExam(markingPoints []*MarkingPoint, positions []*SheetPositioning, pageNum int32) []*MarkingPoint {
	if pageNum == 0 || len(markingPoints) <= 0 || len(positions) <= 0 {
		return nil
	}
	markingPointMap := make(map[int64]int32, 20) // map[判分点ID]页码
	for _, position := range positions {
		if position.PageNum != pageNum { // 跳过非当前页
			continue
		}
		if position.ItemType != resource_enum.SheetItemType_SheetItemTypeMarkingPoint { // 跳过非判分点的题目
			continue
		}
		if position.PositionType == resource_enum.PositionType_PositionTypeOption || position.PositionType == resource_enum.PositionType_PositionTypeTitleDivision {
			v, ok := markingPointMap[position.ItemId]
			if !ok { // 没有收集，则直接收集即可
				markingPointMap[position.ItemId] = position.PageNum
				continue
			}
			if position.PageNum < v { // 跨页则当前判分点归属到小的页码，即为归属到题目的起始页码
				markingPointMap[position.ItemId] = position.PageNum
			}
		}
	}
	//logs.Info(fmt.Sprintf("getPageNumMarkingPoint -> 页码:%d, 对应的判分点:%s", pageNum, tool.ToJson(markingPointIds)))
	results := make([]*MarkingPoint, 0, len(markingPointMap))
	for _, point := range markingPoints {
		if _, ok := markingPointMap[point.Id]; ok {
			results = append(results, point)
		}
	}
	return results
}

// GetPageCorrectMarkingPoint 获取当前页有判分需要的判分点
func GetPageCorrectMarkingPoint(markingPoints []*MarkingPoint, positions []*SheetPositioning, pageNum int32) []*MarkingPoint {

	if pageNum == 0 || len(markingPoints) <= 0 || len(positions) <= 0 {
		return nil
	}

	resultMarkingPointIdMap := make(map[int64]int32, 50) // map[判分点ID]页码
	for _, position := range positions {
		if position.PageNum != pageNum {
			continue
		}
		if position.PositionType == resource_enum.PositionType_PositionTypeOption || position.PositionType == resource_enum.PositionType_PositionTypeTitleDivision { // 选项】、题目判分区区域
			resultMarkingPointIdMap[position.ItemId] = pageNum
		}
	}

	results := make([]*MarkingPoint, 0, len(resultMarkingPointIdMap))
	for _, point := range markingPoints {
		if _, ok := resultMarkingPointIdMap[point.Id]; ok {
			results = append(results, point)
		}
	}
	return results
}

func GetMarkingPointIdByPaperQuestionId(markingPoints []*MarkingPoint, paperQuestionId int64) []int64 {

	if markingPoints == nil {
		return nil
	}
	result := make([]int64, 0, 4)
	for _, markingPoint := range markingPoints {
		if markingPoint.PaperQuestionId == paperQuestionId {
			result = append(result, markingPoint.Id)
		}
	}
	return result
}

func GetCropAnswerImageItem(positions []*SheetPositioning, markingPoints []*MarkingPoint) ([]*MarkingPointPosition, error) {

	resultMap := make(map[int64]*MarkingPointPosition, 30)
	for _, position := range positions {
		if position.PositionType == resource_enum.PositionType_PositionTypeQuestionAnsweringArea { // 判分点作答区
			result := resultMap[position.PaperScoringId]
			if result == nil {
				result = &MarkingPointPosition{
					MarkingPointId: position.PaperScoringId,
					AnswerRects:    nil,
					FillRects:      nil,
				}
			}
			for i, rect := range position.Rects {
				result.AnswerRects = append(result.AnswerRects, &CropRect{
					PageNum:  position.PageNum,
					Sequence: int32(i + 1),
					Rect: &Rect{
						X: int64(rect.X),
						Y: int64(rect.Y),
						W: int64(rect.W),
						H: int64(rect.H),
					},
				})
			}
			resultMap[position.PaperScoringId] = result
		}
		if position.PositionType == resource_enum.PositionType_PositionTypeCompletionFullAnswerArea { // 填空题取填空题公共区域

			markingPointIds := GetMarkingPointIdByPaperQuestionId(markingPoints, position.PaperQuestionId) // 这个转化的原因是，当前填空题公共区域的坐标是挂在paperQuestionId上的
			for _, markingPointId := range markingPointIds {
				result := resultMap[markingPointId]
				if result == nil {
					result = &MarkingPointPosition{
						MarkingPointId: markingPointId,
						AnswerRects:    nil,
						FillRects:      nil,
					}
				}
				for i, rect := range position.Rects {
					result.FillRects = append(result.FillRects, &CropRect{
						PageNum:  position.PageNum,
						Sequence: int32(i + 1),
						Rect: &Rect{
							X: int64(rect.X),
							Y: int64(rect.Y),
							W: int64(rect.W),
							H: int64(rect.H),
						},
					})
				}
				resultMap[markingPointId] = result
			}
		}
	}

	return lists.ToValueList(resultMap), nil
}

func CropAnswerImageItemRemoveObjective(items []*MarkingPointPosition, markingPoints []*MarkingPoint) []*MarkingPointPosition {

	if items == nil || len(items) <= 0 {
		return items
	}
	if markingPoints == nil || len(markingPoints) <= 0 {
		return items
	}
	markingPointMap := make(map[int64]*MarkingPoint, len(markingPoints))
	for _, markingPoint := range markingPoints {
		if markingPoint == nil {
			continue
		}
		markingPointMap[markingPoint.Id] = markingPoint
	}

	items = lists.Filter(items, func(item *MarkingPointPosition) bool {
		markingPoint := markingPointMap[item.MarkingPointId]
		if markingPoint == nil { // 判分点不存在, 过滤
			return false
		}
		if markingPoint.IsObjective { // 客观题过滤
			return false
		}
		return true
	})
	return items
}

// GetAnswerMarkingPoints 获取解答题的判分点，包含合并判分
func GetAnswerMarkingPoints(markingPoints []*MarkingPoint) []*MarkingPoint {

	result := make([]*MarkingPoint, 0, 10)
	for _, markingPoint := range markingPoints {
		if markingPoint.IsObjective {
			continue
		}
		if markingPoint.IsMerged {
			result = append(result, markingPoint)
			continue
		}
		if len(markingPoint.Answers) > 0 && markingPoint.Answers[0].Type == resource_enum.AnswerType_AnswerTypeEssay {
			result = append(result, markingPoint)
		}
	}
	return result
}

// GetSheetAllContainMarkingPointPageNums 考试试卷所有包含判分点的页码
func GetSheetAllContainMarkingPointPageNums(positions []*SheetPositioning) *sets.HashSet[int32] {

	sheetPageNums := sets.NewHashSet[int32]()
	if len(positions) <= 0 {
		return sheetPageNums
	}
	for _, position := range positions { // 总页码
		if position == nil || position.PageNum == 0 {
			continue
		}
		if position.PositionType != resource_enum.PositionType_PositionTypeTitleDivision &&
			position.PositionType != resource_enum.PositionType_PositionTypeQuestionAnsweringArea &&
			position.PositionType != resource_enum.PositionType_PositionTypeCompletionFullAnswerArea &&
			position.PositionType != resource_enum.PositionType_PositionTypeOption { // 以上区域均不包含的情况，认为是可以跳过的页码
			continue
		}
		sheetPageNums.Add(position.PageNum)
	}
	return sheetPageNums
}

func GetAnswerTypeByMarkingPoint(markingPoint *MarkingPoint) resource_enum.AnswerType {

	result := resource_enum.AnswerType_AnswerTypeUnknown
	if markingPoint == nil {
		return result
	}
	if len(markingPoint.Answers) <= 0 {
		return result
	}
	for _, answer := range markingPoint.Answers {
		if answer.Type != resource_enum.AnswerType_AnswerTypeUnknown {
			return answer.Type
		}
	}
	return result
}

// IsThirdSheet 是否是三方卡
func IsThirdSheet(sheet *Sheet) bool {
	if sheet == nil {
		return false
	}
	return sheet.SheetType == resource_enum.SheetType_SheetTypeGenericAnswerSheet
}

func GetMarkingPointScore(points []*MarkingPoint, markingPointId int64) string {

	if len(points) <= 0 || markingPointId <= 0 {
		return ""
	}
	for _, point := range points {
		if point.Id == markingPointId {
			return point.Score
		}
	}
	return ""
}

func GetAllPageNum(sheet *Sheet) *sets.HashSet[int32] {

	if sheet == nil {
		return nil
	}
	result := sets.NewHashSet[int32]()
	for _, page := range sheet.Pages {
		result.Add(page.PageNum)
	}
	return result
}

func IsSheetPageOneSide(sheet *Sheet) bool {
	if sheet == nil {
		return false
	}
	return sheet.IsSheetPageOneSide
}

func IsObjective(AnswerType resource_enum.AnswerType) bool {
	return AnswerType == resource_enum.AnswerType_AnswerTypeSingleChoice ||
		AnswerType == resource_enum.AnswerType_AnswerTypeMultiChoice ||
		AnswerType == resource_enum.AnswerType_AnswerTypeTrueFalse
}

func IsFillBlank(AnswerType resource_enum.AnswerType) bool {
	return AnswerType == resource_enum.AnswerType_AnswerTypeFillin
}

// GetScoreModeMarkingPointIds 获取受scoreMode影响的判分点ID集合
func GetScoreModeMarkingPointIds(markingPoints []*MarkingPoint) []int64 {
	if markingPoints == nil {
		return nil
	}
	scoreModeMarkingPointIds := make([]int64, 0, 10) // 受scoreMode影响的判分点ID集合
	for _, markingPoint := range markingPoints {
		if markingPoint == nil {
			continue
		}
		if markingPoint.IsMerged {
			scoreModeMarkingPointIds = append(scoreModeMarkingPointIds, markingPoint.Id)
			continue
		}
		for _, answer := range markingPoint.Answers {
			if answer == nil {
				continue
			}
			if answer.Type == resource_enum.AnswerType_AnswerTypeEssay {
				scoreModeMarkingPointIds = append(scoreModeMarkingPointIds, markingPoint.Id)
				break
			}
		}
	}
	return scoreModeMarkingPointIds
}

func GetObjectMarkingPointSet(markingPoints []*MarkingPoint) *sets.HashSet[int64] {

	if markingPoints == nil {
		return nil
	}
	result := sets.NewHashSet[int64]()
	for _, markingPoint := range markingPoints {
		if markingPoint == nil || !markingPoint.IsObjective {
			continue
		}
		result.Add(markingPoint.Id)
	}
	return result
}

func GetObjectMarkingPointSetByNodeItem(nodeItems []*PaperNodeItem) *sets.HashSet[int64] {

	if nodeItems == nil {
		return nil
	}
	result := sets.NewHashSet[int64]()
	for _, item := range nodeItems {
		if item == nil || !item.IsMarkingPoint || !item.IsObjective {
			continue
		}
		if item.MarkingPointId <= 0 {
			continue
		}
		result.Add(item.MarkingPointId)
	}
	return result
}

func GetFrontPageNum(pageNums []int32) int32 {

	if len(pageNums) <= 0 {
		return 0
	}
	noZeroPageNums := make([]int32, 0, len(pageNums))
	for _, pageNum := range pageNums {
		if pageNum > 0 {
			noZeroPageNums = append(noZeroPageNums, pageNum)
		}
	}
	if len(noZeroPageNums) <= 0 {
		return 0
	}
	sort.Slice(noZeroPageNums, func(i, j int) bool {
		return noZeroPageNums[i] < noZeroPageNums[j]
	})
	return noZeroPageNums[0]
}
