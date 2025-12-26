package rpc_data

import (
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/common"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/paper_make/make_enum"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/paper_make/make_model"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/paper_make/make_rpc"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/resource/resource_enum"
)

type BindHomeworkPaperParam struct {
	HomeworkId             int64                   // 作业ID'
	PaperMakeType          make_enum.PaperMakeType // 试卷制作类型
	PaperId                int64                   // 绑定试卷时，传需要绑定的试卷ID
	SheetType              resource_enum.SheetType // 答题卡类型
	TemplateImageStoreKeys []string                // 根据学生图像出卷时，传此参数
	UserId                 int64                   // 用户ID
	CampusId               int64                   // 校区ID
}

func NewCreateMakeRequest(param *BindHomeworkPaperParam) *make_rpc.CreateMakeRequest {

	if param == nil {
		return nil
	}
	return &make_rpc.CreateMakeRequest{
		ExamId:        param.HomeworkId,
		PaperMakeType: param.PaperMakeType,
		Option: &make_rpc.CreateMakeRequest_Option{
			SheetType: param.SheetType,
			//IsSwitchMakeType:  false,
			PaperId:                param.PaperId,
			TemplateImageStoreKeys: param.TemplateImageStoreKeys,
		},
		UserId:   param.UserId,
		TenantId: param.CampusId,
	}
}

type BindHomeworkPaperResult struct {
	PaperMakeId int64                         // 出卷Id
	Overview    *make_model.PaperMakeOverview // 出卷预览
}

func NewBindHomeworkPaperResult(resp *make_rpc.CreateMakeResponse) *BindHomeworkPaperResult {
	if resp == nil {
		return nil
	}
	return &BindHomeworkPaperResult{
		PaperMakeId: resp.PaperMakeId,
		Overview:    resp.Overview,
	}
}

type FetchMakeParam struct {
	PaperMakeId *int64                // 出卷Id
	PaperId     *int64                // 题卷Id
	ExamId      *int64                // 考试Id
	UserId      int64                 // 当前用户
	TenantId    int64                 // 当前租户
	Option      *PaperMakeFetchOption // 数据可选项
}

func NewFetchMakeRequest(param *FetchMakeParam) *make_rpc.FetchMakeRequest {
	if param == nil {
		return nil
	}
	return &make_rpc.FetchMakeRequest{
		PaperMakeId: param.PaperMakeId,
		PaperId:     param.PaperId,
		ExamId:      param.ExamId,
		UserId:      param.UserId,
		TenantId:    param.TenantId,
		Option:      convPaperMakeFetchOption(param.Option),
	}
}

type PaperMakeFetchOption struct {
	WithPaper  bool // 包含题卷内容
	WithSheet  bool // 包含题卡内容
	WithStable bool // 是否最新稳定版本
}

func convPaperMakeFetchOption(opt *PaperMakeFetchOption) *make_rpc.PaperMakeFetchOption {
	if opt == nil {
		return nil
	}
	return &make_rpc.PaperMakeFetchOption{
		WithPaper:  opt.WithPaper,
		WithSheet:  opt.WithSheet,
		WithStable: opt.WithStable,
	}
}

type FetchMakeResult struct {
	Overview *PaperMakeOverview // 出卷概览
	Paper    *PaperPrime        // 题卷
	Sheet    *SheetPrime        // 题卡
}

func NewFetchMakeResult(resp *make_rpc.FetchMakeResponse) *FetchMakeResult {
	if resp == nil {
		return nil
	}
	return &FetchMakeResult{
		Overview: convPaperMakeOverview(resp.Overview),
		Paper:    convPaperPrime(resp.Paper),
		Sheet:    convSheetPrime(resp.Sheet),
	}
}

type PaperMakeOverview struct {
	PaperId          int64                     // 题卷Id
	PaperName        string                    // 试卷名称
	ExamId           int64                     // 考试Id
	ExamSubject      int32                     // 考试学科
	BaseGrade        int32                     // 基础年级
	Stage            int32                     // 学段
	UserId           int64                     // 出卷用户Id
	TenantId         int64                     // 出卷租户Id
	BaseSubjects     []int32                   // 基础学科
	CreateTime       int64                     // 创建时间
	UpdateTime       int64                     // 更新时间
	PaperMakeId      int64                     // 出卷Id
	PaperMakeType    make_enum.PaperMakeType   // 出卷方式
	PaperMakeStatus  make_enum.PaperMakeStatus // 出卷状态
	ReusePaperMakeId int64                     // 复用出卷Id
	CorrectionType   common.CorrectionType     // 批改方式
	SheetType        resource_enum.SheetType   // 题卡类型

}

func convPaperMakeOverview(overview *make_model.PaperMakeOverview) *PaperMakeOverview {
	if overview == nil {
		return nil
	}
	return &PaperMakeOverview{
		PaperId:          overview.PaperId,
		PaperName:        overview.PaperName,
		ExamId:           overview.ExamId,
		ExamSubject:      overview.ExamSubject,
		BaseGrade:        overview.BaseGrade,
		Stage:            overview.Stage,
		UserId:           overview.UserId,
		TenantId:         overview.TenantId,
		BaseSubjects:     overview.BaseSubjects,
		CreateTime:       overview.CreateTime,
		UpdateTime:       overview.UpdateTime,
		PaperMakeId:      overview.PaperMakeId,
		PaperMakeType:    overview.PaperMakeType,
		PaperMakeStatus:  overview.PaperMakeStatus,
		ReusePaperMakeId: overview.ReusePaperMakeId,
		CorrectionType:   overview.CorrectionType,
		SheetType:        overview.SheetType,
	}
}

type PaperPrime struct {
	Id              int64                          // 题卷Id
	Name            *string                        // 题卷名称
	SubName         *string                        // 题卷副标题
	SubHeading      *string                        // 三级标题
	Category        *string                        // 题卷分类
	Term            *int32                         // 学期
	Stage           *int32                         // 学段
	BaseGrade       *int32                         // 基础年级
	BaseSubjects    []int32                        // 基础学科列表
	Score           *string                        // 题卷总分
	PaperSourceType *resource_enum.PaperSourceType // 题卷来源类型
	PaperExamType   *int32                         // 题卷考试类型
	Year            *int32                         // 年份
	ThirdId         *string                        // 三方Id
	Supplier        *resource_enum.Supplier        // 数据提供方
	SubSupplier     *resource_enum.SubSupplier     // 数据次级提供方
}

func convPaperPrime(prime *make_model.PaperPrime) *PaperPrime {
	if prime == nil {
		return nil
	}
	return &PaperPrime{
		Id:              prime.Id,
		Name:            prime.Name,
		SubName:         prime.SubName,
		SubHeading:      prime.SubHeading,
		Category:        prime.Category,
		Term:            prime.Term,
		Stage:           prime.Stage,
		BaseGrade:       prime.BaseGrade,
		BaseSubjects:    prime.BaseSubjects,
		Score:           prime.Score,
		PaperSourceType: prime.PaperSourceType,
		PaperExamType:   prime.PaperExamType,
		Year:            prime.Year,
		ThirdId:         prime.ThirdId,
		Supplier:        prime.Supplier,
		SubSupplier:     prime.SubSupplier,
	}
}

type SheetPrime struct {
	Id                 int64                           // 题卡Id
	SheetType          *resource_enum.SheetType        // 题卡类型
	PrintFormat        *resource_enum.PrintFormat      // 打印板式
	Column             *int32                          // 栏数
	Version            *int32                          // 版本号
	FileId             *int64                          // 题卡完整pdf文件Id
	MarkingMode        *resource_enum.MarkingMode      // 批分方式
	StudentCodeType    *resource_enum.StudentCodeType  // 考号类型
	ObjectiveDirection *resource_enum.DirectionStatus  // 客观题选项方向
	IsSheetPageOneSide *bool                           // 单面题卡页
	CorrectType        *resource_enum.SheetCorrectType // 批改类型
	EditConfig         *string                         // 端上编辑配置
	TosKey             *string                         // 生产资源文件用模版key
}

func convSheetPrime(prime *make_model.SheetPrime) *SheetPrime {
	if prime == nil {
		return nil
	}
	return &SheetPrime{
		Id:                 prime.Id,
		SheetType:          prime.SheetType,
		PrintFormat:        prime.PrintFormat,
		Column:             prime.Column,
		Version:            prime.Version,
		FileId:             prime.FileId,
		MarkingMode:        prime.MarkingMode,
		StudentCodeType:    prime.StudentCodeType,
		ObjectiveDirection: prime.ObjectiveDirection,
		IsSheetPageOneSide: prime.IsSheetPageOneSide,
		CorrectType:        prime.CorrectType,
		EditConfig:         prime.EditConfig,
		TosKey:             prime.TosKey,
	}
}

func GetSheetType(makeResult *FetchMakeResult) resource_enum.SheetType {

	if makeResult == nil {
		return 0
	}
	if makeResult.Overview == nil {
		return 0
	}
	return makeResult.Overview.SheetType
}
