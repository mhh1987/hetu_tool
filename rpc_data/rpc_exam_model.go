package rpc_data

import (
	"code.chenji.com/cj/scan_paper/utils/lists"
	"code.chenji.com/pkg/common/tool"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/common"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/exam/exam"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/exam/exam_enum"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/organization/organization_enum"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/resource/resource_enum"
)

type UnionExam struct {
	UnionExamMeta *UnionExamMeta
	ExamGroupList []*ExamGroup
}

func ConvUnionExam(data *exam.UnionExam) *UnionExam {
	return &UnionExam{
		UnionExamMeta: convUnionExamMeta(data.UnionExamMeta),
		ExamGroupList: lists.Map(data.ExamGroupList, convExamGroup),
	}
}

type UnionExamMeta struct {
	UnionExamId    int64                          // 联考ID
	UnionExamType  common.UnionExamType           // 联考类型
	Creator        int64                          // 创建人ID
	CreatorName    string                         // 创建人姓名
	TenantId       int64                          // 租户ID
	TenantType     organization_enum.TenantType   // 租户类型
	UnionExamName  string                         // 联考名称
	CreateTime     int64                          // 创建时间
	ExamStatus     common.UnionExamStatus         // 联考状态
	ExamMode       common.ExamMode                // 考试模式
	ManagementType common.UnionExamManagementType // 联考管理类型
	ExamSource     exam_enum.ExamSource           // 考试来源：0：历史数据，1：机构端，2：教师端
	ExamCategory   common.ExamCategory            // 考试、作业
	ExamType       common.ExamType                // 类型
	SchoolYear     int32                          // 学年
	ScanType       common.ScanType                // 扫描类型
	Stage          organization_enum.Stage        // 学段
	BaseGrade      common.GradeValue              // 基础年级:
}

func convUnionExamMeta(data *exam.UnionExamMeta) *UnionExamMeta {
	return &UnionExamMeta{
		UnionExamId:    data.UnionExamId,
		UnionExamType:  data.UnionExamType,
		Creator:        data.Creator,
		CreatorName:    data.CreatorName,
		TenantId:       data.InstitutionId,
		TenantType:     data.InstitutionType,
		UnionExamName:  data.UnionExamName,
		CreateTime:     data.CreateTime,
		ExamStatus:     data.UnionExamStatus,
		ExamMode:       data.ExamMode,
		ManagementType: data.ManagementType,
		ExamSource:     data.ExamSource,
		ExamCategory:   data.ExamCategory,
		ExamType:       data.ExamType,
		SchoolYear:     data.SchoolYear,
		ScanType:       data.ScanType,
		Stage:          data.Stage,
		BaseGrade:      data.BaseGrade,
	}
}

type ExamGroup struct {
	ExamGroupMeta *ExamGroupMeta // 考试组元数据
	ExamList      []*Exam        // 考试列表
	ExamClassList []*ExamClass   // 考试班级列表
}

func convExamGroup(data *exam.ExamGroup) *ExamGroup {
	return &ExamGroup{
		ExamGroupMeta: convExamGroupMeta(data.ExamGroupMeta),
		ExamList:      lists.Map(data.ExamList, ConvExam),
		ExamClassList: lists.Map(data.ExamClassList, convExamClass),
	}
}

type ExamGroupMeta struct {
	ExamGroupId   int64                        // 考试组ID
	UnionExamId   int64                        // 联考ID
	ExamGroupName string                       // 考试组名称
	Creator       int64                        // 创建人ID
	CreatorName   string                       // 创建人姓名
	BaseGrade     common.GradeValue            // 基础年级
	ExamType      common.ExamType              // 考试类型
	StartTime     int64                        // 考试开始时间
	EndTime       int64                        // 考试结束时间
	TenantId      int64                        // 租户ID
	TenantType    organization_enum.TenantType // 租户类型
	SchoolYear    int32                        // 年份
	Stage         int32                        // 学段
	CreateTime    int64                        // 创建时间
	ScanType      common.ScanType              // 扫描类型
	ExamStatus    common.UnionExamStatus       // 考试状态
	ExamSource    exam_enum.ExamSource         // 考试来源
	ExamCategory  common.ExamCategory          // 考试类别
}

func convExamGroupMeta(data *exam.ExamGroupMeta) *ExamGroupMeta {
	return &ExamGroupMeta{
		ExamGroupId:   data.ExamGroupId,
		UnionExamId:   data.UnionExamId,
		ExamGroupName: data.ExamGroupName,
		Creator:       data.Creator,
		CreatorName:   data.CreatorName,
		BaseGrade:     data.BaseGrade,
		ExamType:      data.ExamType,
		StartTime:     data.StartTime,
		EndTime:       data.EndTime,
		TenantId:      data.InstitutionId,
		TenantType:    data.InstitutionType,
		SchoolYear:    data.SchoolYear,
		Stage:         data.Stage,
		CreateTime:    data.CreateTime,
		ScanType:      data.ScanType,
		ExamStatus:    data.ExamStatus,
		ExamSource:    data.ExamSource,
		ExamCategory:  data.ExamCategory,
	}
}

type Exam struct {
	ExamMeta               *ExamMeta       // 考试元数据
	ExamStepStatus         *ExamStepStatus // 考试步骤状态
	Creators               []*ExamTeacher  // 创建人
	Teachers               []*ExamTeacher  // 教师
	CorrectionTeachers     []*ExamTeacher  // 批改老师
	ArbitrationTeachers    []*ExamTeacher  // 仲裁老师
	ExceptionTeachers      []*ExamTeacher  // 异常处理老师
	ScanTeachers           []*ExamTeacher  // 扫描老师
	ClassList              []*ExamClass    // 班级列表
	SchoolAdminTeacher     []*ExamTeacher  // 学校管理员
	MakePaperTeachers      []*ExamTeacher  // 出卷老师
	ModifyScoreTeachers    []*ExamTeacher  // 修改分数老师
	ExamGroupAdminTeachers []*ExamTeacher  // 考试组管理员
	AreaAdminTeachers      []*ExamTeacher  // 区域管理员
}

func ConvExam(data *exam.Exam) *Exam {
	return &Exam{
		ExamMeta:               convExamMeta(data.ExamMeta),
		ExamStepStatus:         convExamStepStatus(data.ExamStepStatus),
		Creators:               lists.Map(data.Creators, convExamTeacher),
		Teachers:               lists.Map(data.Teachers, convExamTeacher),
		CorrectionTeachers:     lists.Map(data.CorrectionTeachers, convExamTeacher),
		ArbitrationTeachers:    lists.Map(data.ArbitrationTeachers, convExamTeacher),
		ExceptionTeachers:      lists.Map(data.ExceptionTeachers, convExamTeacher),
		ScanTeachers:           lists.Map(data.ScanTeachers, convExamTeacher),
		ClassList:              lists.Map(data.ClassList, convExamClass),
		SchoolAdminTeacher:     lists.Map(data.SchoolAdminTeacher, convExamTeacher),
		MakePaperTeachers:      lists.Map(data.MakePaperTeachers, convExamTeacher),
		ModifyScoreTeachers:    lists.Map(data.ModifyScoreTeachers, convExamTeacher),
		ExamGroupAdminTeachers: lists.Map(data.ExamGroupAdminTeachers, convExamTeacher),
		AreaAdminTeachers:      lists.Map(data.AreaAdminTeachers, convExamTeacher),
	}
}

type ExamClass struct {
	TenantId   int64
	TenantType organization_enum.TenantType
	ClassId    int64
	SchoolYear int32
	CampusId   int64
}

func convExamClass(data *exam.ExamClass) *ExamClass {
	return &ExamClass{
		TenantId:   data.InstitutionId,
		TenantType: data.InstitutionType,
		ClassId:    data.ClassId,
		SchoolYear: data.SchoolYear,
		CampusId:   data.TenantId,
	}
}

type ExamTeacher struct {
	TeacherId   int64
	TeacherName string
	TeacherType common.ExamTeacherType
	ExamId      int64
	TenantId    int64
	TenantType  organization_enum.TenantType
}

func convExamTeacher(data *exam.ExamTeacher) *ExamTeacher {
	return &ExamTeacher{
		TeacherId:   data.TeacherId,
		TeacherName: data.TeacherName,
		TeacherType: data.TeacherType,
		ExamId:      data.ExamId,
		TenantId:    data.InstitutionId,
		TenantType:  data.InstitutionType,
	}
}

type ExamStepStatus struct {
	WebPaperStatus               common.MakePaperStatus
	WebDistributionStatus        common.DistributionStatus
	ScanStatus                   common.ScanStatus
	WebCorrectionStatus          common.CorrectionStatus
	WebStatisticsStatus          common.StatisticsStatus
	BureauStatisticsStatus       common.StatisticsStatus
	AreaBureauStatisticsStatus   common.StatisticsStatus
	CanPublishStatistic          common.YesOrNo
	IsObjectPaper                common.YesOrNo
	CanReadCorrectionProgress    common.YesOrNo
	CanReadCorrectionQuality     common.YesOrNo
	IsAllAnswer                  common.YesOrNo
	SyMakePaperStatus            common.SYMakePaperStatus
	SyStatisticStatus            common.SYStatisticsStatus
	ExamStudentStatus            common.ExamStudentStatus
	ThirdMakePaperStatus         common.ThirdMakePaperStatus
	ThirdImportStatus            common.ThirdImportStatus
	ThirdStatisticStatus         common.ThirdStatisticStatus
	ExamPauseStatus              common.ExamPauseStatus
	CanPublishObjectiveStatistic common.YesOrNo
	ObjectiveStatisticsStatus    common.ObjectiveStatisticsStatus
}

func convExamStepStatus(data *exam.ExamStepStatus) *ExamStepStatus {
	return &ExamStepStatus{
		WebPaperStatus:               data.WebPaperStatus,
		WebDistributionStatus:        data.WebDistributionStatus,
		ScanStatus:                   data.ScanStatus,
		WebCorrectionStatus:          data.WebCorrectionStatus,
		WebStatisticsStatus:          data.WebStatisticsStatus,
		BureauStatisticsStatus:       data.BureauStatisticsStatus,
		AreaBureauStatisticsStatus:   data.AreaBureauStatisticsStatus,
		CanPublishStatistic:          data.CanPublishStatistic,
		IsObjectPaper:                data.IsObjectPaper,
		CanReadCorrectionProgress:    data.CanReadCorrectionProgress,
		CanReadCorrectionQuality:     data.CanReadCorrectionQuality,
		IsAllAnswer:                  data.IsAllAnswer,
		SyMakePaperStatus:            data.SyMakePaperStatus,
		SyStatisticStatus:            data.SyStatisticStatus,
		ExamStudentStatus:            data.ExamStudentStatus,
		ThirdMakePaperStatus:         data.ThirdMakePaperStatus,
		ThirdImportStatus:            data.ThirdImportStatus,
		ThirdStatisticStatus:         data.ThirdStatisticStatus,
		ExamPauseStatus:              data.ExamPauseStatus,
		CanPublishObjectiveStatistic: data.CanPublishObjectiveStatistic,
		ObjectiveStatisticsStatus:    data.ObjectiveStatisticsStatus,
		//SkConfirmStatus:              data.SkConfirmStatus,
	}
}

type ExamMeta struct {
	ExamId                 int64
	ExamSubject            int32
	PaperId                int64
	ExamGroupId            int64
	CorrectionType         common.CorrectionType
	ScanType               common.ScanType
	Creator                int64
	CreatorName            string
	UnionExamId            int64
	TenantId               int64
	TenantType             organization_enum.TenantType
	ExamGroupName          string
	Stage                  int32
	BaseGrade              common.GradeValue
	ExamType               common.ExamType
	UnionExamType          common.UnionExamType
	ExamStatus             common.UnionExamStatus
	SchoolYear             int32
	CorrectionMarkType     common.CorrectionMarkType
	ScoringMode            common.ScoringMode
	CreatedAt              int64
	UpdatedAt              int64
	ExamStudentConfig      int64
	ExamStudentSettingType exam_enum.ExamStudentSettingType
	ExamStudentBaseConfig  int64
	StudentRecognition     exam_enum.StudentRecognition
	AreaBureauId           int64
	ManagementType         common.UnionExamManagementType
	ExamSource             exam_enum.ExamSource
	ExamCategory           common.ExamCategory
}

func convExamMeta(data *exam.ExamMeta) *ExamMeta {
	return &ExamMeta{
		ExamId:                 data.ExamId,
		ExamSubject:            data.ExamSubject,
		PaperId:                data.PaperId,
		ExamGroupId:            data.ExamGroupId,
		CorrectionType:         data.CorrectionType,
		ScanType:               data.ScanType,
		Creator:                data.Creator,
		CreatorName:            data.CreatorName,
		UnionExamId:            data.UnionExamId,
		TenantId:               data.InstitutionId,
		TenantType:             data.InstitutionType,
		ExamGroupName:          data.ExamGroupName,
		Stage:                  data.Stage,
		BaseGrade:              data.BaseGrade,
		ExamType:               data.ExamType,
		UnionExamType:          data.UnionExamType,
		ExamStatus:             data.ExamStatus,
		SchoolYear:             data.SchoolYear,
		CorrectionMarkType:     data.CorrectionMarkType,
		ScoringMode:            data.ScoringMode,
		CreatedAt:              data.CreatedAt,
		UpdatedAt:              data.UpdatedAt,
		ExamStudentConfig:      data.ExamStudentConfig,
		ExamStudentSettingType: data.ExamStudentSettingType,
		ExamStudentBaseConfig:  data.ExamStudentBaseConfig,
		StudentRecognition:     data.StudentRecognition,
		AreaBureauId:           data.AreaBureauId,
		ManagementType:         data.ManagementType,
		ExamSource:             data.ExamSource,
		ExamCategory:           data.ExamCategory,
	}
}

type ExamStudentDetail struct {
	StudentMeta *ExamStudentMeta // 学生元数据
	ClassList   []*ExamClassMeta // 班级列表
	Grade       *ExamGradeMeta   // 学生年级；优先行政班所在年级，其次使用走班年级
	Campus      *ExamCampusMeta  // 学生校区信息
	School      *ExamSchoolMeta  // 学生学校信息

}

func ConvExamStudentDetail(studentDetail *exam.ExamStudentDetail) *ExamStudentDetail {

	if studentDetail == nil {
		return nil
	}
	return &ExamStudentDetail{
		StudentMeta: convExamStudentMeta(studentDetail.StudentMeta),
		ClassList:   lists.Map(studentDetail.ClassList, convExamClassMeta),
		Grade:       convExamGrade(studentDetail.Grade),
		Campus:      convExamCampusMeta(studentDetail.Campus),
		School:      convExamSchoolMeta(studentDetail.School),
	}
}

type ExamStudentMeta struct {
	StudentId       int64
	StudentName     string
	StudentCode     string
	ExamCode        string
	TransientStatus int32 // 借读状态
	SchoolId        int64
	ExamHall        *ExamHallMeta // 考场
	ExamSeat        string
	StudentCategory int32
	ShortExamCode   string // 短考号

}

func convExamStudentMeta(data *exam.ExamStudentMeta) *ExamStudentMeta {

	if data == nil {
		return nil
	}
	return &ExamStudentMeta{
		StudentId:       data.StudentId,
		StudentName:     data.StudentName,
		StudentCode:     data.StudentCode,
		ExamCode:        data.ExamCode,
		TransientStatus: data.TransientStatus,
		SchoolId:        data.SchoolId,
		ExamHall:        convExamHall(data.ExamHall),
		ExamSeat:        data.ExamSeat,
		StudentCategory: data.StudentCategory,
		ShortExamCode:   data.ShortExamCode,
	}
}

// ExamHallMeta 考场信息
type ExamHallMeta struct {
	HallId   int64  // 考场ID
	HallName string // 考场名称
}

func convExamHall(data *exam.ExamHallMeta) *ExamHallMeta {
	if data == nil {
		return nil
	}
	return &ExamHallMeta{
		HallId:   data.Id,
		HallName: data.Name,
	}
}

type ExamClassMeta struct {
	ClassId      int64
	ClassName    string
	ClassType    organization_enum.ClassType
	BaseSubjects []int32
	SubjectType  organization_enum.SubjectTypeFieldScope
	IsExamClass  bool
	Sequence     int32
}

func convExamClassMeta(data *exam.ExamClassMeta) *ExamClassMeta {
	if data == nil {
		return nil
	}
	return &ExamClassMeta{
		ClassId:      data.ClassId,
		ClassName:    data.ClassName,
		ClassType:    data.ClassType,
		BaseSubjects: data.BaseSubjects,
		SubjectType:  data.SubjectType,
		IsExamClass:  data.IsExamClass,
	}
}

type ExamGradeMeta struct {
	GradeId   int64
	GradeName string
	BaseGrade int32
}

func convExamGrade(data *exam.ExamGradeMeta) *ExamGradeMeta {
	if data == nil {
		return nil
	}
	return &ExamGradeMeta{
		GradeId:   data.GradeId,
		GradeName: data.GradeName,
		BaseGrade: data.BaseGrade,
	}
}

type ExamCampusMeta struct {
	CampusId   int64  // 校区ID
	CampusName string // 校区名称
}

func convExamCampusMeta(data *exam.ExamCampusMeta) *ExamCampusMeta {
	if data == nil {
		return nil
	}
	return &ExamCampusMeta{
		CampusId:   data.CampusId,
		CampusName: data.CampusName,
	}
}

type ExamSchoolMeta struct {
	SchoolId   int64  // 学校ID
	SchoolName string // 学校名称
	System     int32
}

func convExamSchoolMeta(data *exam.ExamSchoolMeta) *ExamSchoolMeta {
	if data == nil {
		return nil
	}
	return &ExamSchoolMeta{
		SchoolId:   data.SchoolId,
		SchoolName: data.SchoolName,
		System:     data.System,
	}
}

type ExamClassStudentDetail struct {
	ClassMeta   *ExamClassMeta
	StudentMeta []*ExamStudentMeta
}

func ConvExamClassStudentDetail(data *exam.ExamClassStudentDetail) *ExamClassStudentDetail {
	if data == nil {
		return nil
	}
	return &ExamClassStudentDetail{
		ClassMeta:   convExamClassMeta(data.ClassMeta),
		StudentMeta: lists.Map(data.StudentMeta, convExamStudentMeta),
	}
}

type UpdateExamSubjectStatusParam struct {
	ExamId     int64
	TaskType   exam_enum.TaskType
	StepStatus *ExamStepStatus
	//SyScanForm *TTSYScanForm
}

type GetExamStudentCountParam struct {
	UnionExamId int64
	// 不传union_exam_id时，传exam_id，返回指定学科
	// 有值时，忽略exam_subjects
	// 当exam_id > 0， len(school_ids) == 0, 联考自校下学校用户需要限定学校，只统计exam_id对应学校数据
	ExamId int64
	// exam_id == 0时，不传传exam_subjects，返回所有学科
	ExamSubjects []int32
	// 不传时，统计所有参考学校
	Schools []*CountExamStudentRequestSchool
}

type GetEExamStudentCountResultItem struct {
	SchoolId     int64
	CampusId     int64
	ExamSubject  int32
	StudentCount int32
}

func ConvGetExamStudentCountResultItem(data *exam.CountExamStudentResponse_Item) *GetEExamStudentCountResultItem {

	return &GetEExamStudentCountResultItem{
		SchoolId:     data.SchoolId,
		CampusId:     data.CampusId,
		ExamSubject:  data.ExamSubject,
		StudentCount: data.StudentCnt,
	}
}

type CountExamStudentRequestSchool struct {
	SchoolId int64 // 学校ID
	// 不传时，统计所有校区
	CampusIds []int64 // 校区ID
}

func ConvCountExamStudentRequestSchool(data *CountExamStudentRequestSchool) *exam.CountExamStudentRequest_School {

	if data == nil {
		return nil
	}
	return &exam.CountExamStudentRequest_School{
		SchoolId:  data.SchoolId,
		CampusIds: data.CampusIds,
	}
}

type TTSYScanForm struct {
	ExceptionCount int32 // 异常个数
	BatchId        int64 // 批次id
	InstitutionId  int64 // 当前批次老师机构id
	TeacherId      int64 // 当前批次老师id
	PaperId        int64 // 试卷id
}

type GetUnionExamSubjectStudentCountParam struct {
	UnionExamId  int64   // 联考ID
	ExamSubjects []int32 // 考试科目
	SchoolIds    []int64 // 学校ID
}

type GetUnionExamSubjectStudentCountResult struct {
	SchoolId                   int64           // 学校ID
	ExamSubjectStudentCountMap map[int32]int32 // 科目ID-学生数量
}

type GetExamSchoolInfoParam struct {
	ExamId         int64 // 考试ID，当前用户登录的学校对应的考试ID即可，考试管理会根据此考试是否混阅来判断返回数据的范围
	InstitutionIds []int64
	CampusIds      []int64
	GradeIds       []int64
	ClassIds       []int64
	IgnoreSubject  bool
}

type SearchExamStudentParam struct {
	Offset    int32
	Limit     int32
	ExamId    int64
	SchoolIds []int64 // 筛选机构
	CampusIds []int64 // 筛选校区
	ClassIds  []int64
	// 定位学生id;使用此条件忽略offset,通过student位置计算offset
	StudentId   int64
	UnionExamId int64
	// required
	InstitutionId int64 // 当前机构, 学校id、教育局id etc. 非校区id
	// exam_id == 0 时生效， exam_id == 0 && exam_subject == 0  表示查询统一配置下的学生
	ExamSubject int32    // 考试学科
	Keyword     string   // 搜索关键字
	TermKeyword []string // 精确查询关键字, 目前支持student_name
}

func NewSearchExamStudentReq(param *SearchExamStudentParam) *exam.SearchExamStudentRequest {
	if param == nil {
		return nil
	}
	limit := param.Limit
	if limit <= 0 {
		limit = 2000
	}
	return &exam.SearchExamStudentRequest{
		Offset:        param.Offset,
		Limit:         limit,
		ExamId:        param.ExamId,
		SchoolIds:     param.SchoolIds,
		CampusIds:     param.CampusIds,
		ClassIds:      param.ClassIds,
		StudentId:     param.StudentId,
		UnionExamId:   param.UnionExamId,
		InstitutionId: param.InstitutionId,
		ExamSubject:   param.ExamSubject,
		Keyword:       param.Keyword,
		TermKeyword:   param.TermKeyword,
	}
}

type SearchExamStudentResult struct {
	StudentList []*ExamStudentDetail // 学生列表
	Total       int32                // 筛选总数
}

func NewSearchExamStudentResult(data *exam.SearchExamStudentResponse) *SearchExamStudentResult {

	if data == nil {
		return nil
	}
	return &SearchExamStudentResult{
		StudentList: lists.Map(data.StudentList, ConvExamStudentDetail),
		Total:       data.Total,
	}
}

type GetAuthorityOfTeacherParam struct {
	TeacherId   int64
	ExamGroupId int64
	ExamId      int64
	// union_exam_id + institution_id、exam_group_id、exam_id  三选一
	UnionExamId   int64
	InstitutionId int64
	// exam_subject == 0、exam_id 二选一， 都为0时返回所有学科 权限
	ExamSubject int32
}

func NewGetAuthorityOfTeacherParam(param *GetAuthorityOfTeacherParam) *exam.GetAuthorityOfTeacherRequest {
	if param == nil {
		return nil
	}
	return &exam.GetAuthorityOfTeacherRequest{
		TeacherId:     param.TeacherId,
		ExamGroupId:   param.ExamGroupId,
		ExamId:        param.ExamId,
		UnionExamId:   param.UnionExamId,
		InstitutionId: param.InstitutionId,
		ExamSubject:   param.ExamSubject,
	}
}

type ExamSubjectRightInfo struct {
	ExamId            int64
	ExamSubjectRights []exam_enum.ExamSubjectRightType
}

type GetAuthorityOfTeacherResult struct {
	ExamGroupRights []common.ExamGroupRightType
	// key 为 exam_id
	ExamSubjectRightsMap map[int64]*ExamSubjectRightInfo
	TeacherTypes         []common.ExamTeacherType
}

func NewGetAuthorityOfTeacherResult(data *exam.GetAuthorityOfTeacherResponse) *GetAuthorityOfTeacherResult {
	if data == nil {
		return nil
	}
	tempResult := &GetAuthorityOfTeacherResult{
		ExamGroupRights:      data.ExamGroupRights,
		ExamSubjectRightsMap: nil,
		TeacherTypes:         data.TeacherTypes,
	}
	tempResult.ExamSubjectRightsMap = make(map[int64]*ExamSubjectRightInfo, len(data.ExamSubjectRightsMap))
	for id, info := range data.ExamSubjectRightsMap {
		tempResult.ExamSubjectRightsMap[id] = &ExamSubjectRightInfo{
			ExamId:            info.ExamId,
			ExamSubjectRights: info.ExamSubjectRights,
		}
	}
	return tempResult
}

type NoticeExamStudentAbsentParam struct {
	UnionExamId int64   `json:"union_exam_id"`
	ExamSubject int32   `json:"exam_subject"`
	StudentIds  []int64 `json:"student_id"`
	IsAbsent    bool    `json:"is_absent"` // 是否缺考
}

type ExamSchool struct {
	SchoolId       int64
	SchoolName     string
	CampusList     []*ExamCampus
	System         int32  // 学制
	AreaBureauId   int64  // 区域教育局id
	AreaBureauName string // 区域教育局名字

}

func ConvExamSchool(data *exam.ExamSchool) *ExamSchool {
	if data == nil {
		return nil
	}
	campusList := make([]*ExamCampus, 0, 10)
	if len(data.CampusList) > 0 {
		for i, campus := range data.CampusList {
			tempCampus := convExamCampus(campus)
			tempCampus.Sequence = int32(i + 1)
			campusList = append(campusList, tempCampus)
		}
	}
	return &ExamSchool{
		SchoolId:       data.SchoolId,
		SchoolName:     data.SchoolName,
		CampusList:     campusList,
		System:         data.System,
		AreaBureauId:   data.AreaBureauId,
		AreaBureauName: data.AreaBureauName,
	}
}

type ExamCampus struct {
	CampusId   int64
	CampusName string
	GradeList  []*ExamSchoolGrade
	Sequence   int32
}

func convExamCampus(data *exam.ExamCampus) *ExamCampus {
	if data == nil {
		return nil
	}
	gradeList := make([]*ExamSchoolGrade, 0, 10)
	if len(data.GradeList) > 0 {
		for i, grade := range data.GradeList {
			tempGrade := convExamSchoolGrade(grade)
			tempGrade.Sequence = int32(i + 1)
			gradeList = append(gradeList, tempGrade)
		}
	}
	return &ExamCampus{
		CampusId:   data.CampusId,
		CampusName: data.CampusName,
		GradeList:  gradeList,
	}
}

type ExamSchoolGrade struct {
	GradeId   int64
	GradeName string
	BaseGrade int32
	ClassList []*ExamClassMeta
	Sequence  int32
}

func convExamSchoolGrade(data *exam.ExamSchoolGrade) *ExamSchoolGrade {
	if data == nil {
		return nil
	}
	classList := make([]*ExamClassMeta, 0, 10)
	if len(data.ClassList) > 0 {
		for i, meta := range data.ClassList {
			tempClass := convExamClassMeta(meta)
			tempClass.Sequence = int32(i + 1)
			classList = append(classList, tempClass)
		}
	}
	return &ExamSchoolGrade{
		GradeId:   data.GradeId,
		GradeName: data.GradeName,
		BaseGrade: data.BaseGrade,
		ClassList: classList,
	}
}

type GetExamStudentBatchParam struct {
	ExamId     int64
	StudentIds []int64
}

type GetExamStudentBatchResult struct {
	ExamStudents []*ExamStudentDetail
}

type GetExamStudentParam struct {
	ExamId             int64
	ClassIds           []int64
	StudentIds         []int64
	NeedAdministrative bool
	ExamCodes          []string
	ShortExamCodes     []string
	StudentCodes       []string
}
type PaperSourceInfo struct {
	ExamId          int64                  // 考试id
	ExamPaperSource common.PaperSource     // 出卷方式
	Status          common.MakePaperStatus // 出卷状态 1:出卷中 2:已完成
	UpdateTime      int64                  // 更新时间
	ReuseName       string                 // 复用试卷名称
	ReuseId         int64                  // 复用试卷id
	IsReuse         bool                   // 是否使用答题卡

}

func ConvPaperSourceInfo(data *exam.PaperSourceInfo) *PaperSourceInfo {
	if data == nil {
		return nil
	}
	return &PaperSourceInfo{
		ExamId:          data.ExamId,
		ExamPaperSource: data.ExamPaperSource,
		Status:          data.Status,
		UpdateTime:      data.UpdateTime,
		ReuseName:       data.ReuseName,
		ReuseId:         data.ReuseId,
		IsReuse:         data.IsReuse,
	}
}

func ConvPaperSourceToResource(source common.PaperSource) resource_enum.SheetType {

	switch source {
	case common.PaperSource_PSAnswerSheet:
		return resource_enum.SheetType_SheetTypeAnswerSheet
	case common.PaperSource_PSWordAnswerSheet:
		return resource_enum.SheetType_SheetTypeAnswerSheet
	case common.PaperSource_PSWordContentSheet:
		return resource_enum.SheetType_SheetTypeContentSheet
	case common.PaperSource_PSThirdAnswerSheet:
		return resource_enum.SheetType_SheetTypeGenericAnswerSheet
	case common.PaperSource_PSThirdContentSheet:
		return resource_enum.SheetType_SheetTypeContentSheet
	case common.PaperSource_PSThirdImport:
		return resource_enum.SheetType_SheetTypeGenericAnswerSheet
	case common.PaperSource_PSResourceAnswerSheet:
		return resource_enum.SheetType_SheetTypeAnswerSheet
	case common.PaperSource_PSReuseAnswerSheet:
		return resource_enum.SheetType_SheetTypeAnswerSheet
	default:
		return resource_enum.SheetType_SheetTypeUnknown
	}
}

func IsWebCorrectByPaper(unionExam *UnionExam, paperId int64) bool {

	if unionExam == nil || len(unionExam.ExamGroupList) <= 0 {
		return false
	}
	for _, examGroup := range unionExam.ExamGroupList {
		if len(examGroup.ExamList) <= 0 {
			continue
		}
		for _, examInfo := range examGroup.ExamList {
			if examInfo == nil || examInfo.ExamMeta == nil {
				continue
			}
			if examInfo.ExamMeta.PaperId == paperId && (examInfo.ExamMeta.CorrectionType == common.CorrectionType_CTWebCorrect || examInfo.ExamMeta.CorrectionType == common.CorrectionType_CTThirdCorrect) {
				return true
			}
		}
	}
	return false
}

func IsWebCorrectByExam(examInfo *Exam) *bool {

	if examInfo == nil || examInfo.ExamMeta == nil {
		return nil
	}
	isWebCorrect := examInfo.ExamMeta.CorrectionType == common.CorrectionType_CTWebCorrect || examInfo.ExamMeta.CorrectionType == common.CorrectionType_CTThirdCorrect
	return tool.Ptr(isWebCorrect)
}

func GetExamSubjectByPaper(unionExam *UnionExam, paperId int64) int32 {

	if unionExam == nil || len(unionExam.ExamGroupList) <= 0 {
		return 0
	}
	for _, examGroup := range unionExam.ExamGroupList {
		if len(examGroup.ExamList) <= 0 {
			continue
		}
		for _, examInfo := range examGroup.ExamList {
			if examInfo == nil || examInfo.ExamMeta == nil {
				continue
			}
			if examInfo.ExamMeta.PaperId == paperId {
				return examInfo.ExamMeta.ExamSubject
			}
		}
	}
	return 0
}

func GetExamId(examInfo *Exam) int64 {

	if examInfo == nil || examInfo.ExamMeta == nil {
		return 0
	}
	return examInfo.ExamMeta.ExamId
}

func GetExamFromUnionByExamId(unionExam *UnionExam, examId int64) *Exam {
	if unionExam == nil || len(unionExam.ExamGroupList) <= 0 {
		return nil
	}
	for _, examGroup := range unionExam.ExamGroupList {
		if len(examGroup.ExamList) <= 0 {
			continue
		}
		for _, examInfo := range examGroup.ExamList {
			if examInfo == nil || examInfo.ExamMeta == nil {
				continue
			}
			if examInfo.ExamMeta.ExamId == examId {
				return examInfo
			}
		}
	}
	return nil
}

func GetSchoolIdFromExam(examInfo *Exam) int64 {
	if examInfo == nil || examInfo.ExamMeta == nil {
		return 0
	}
	return examInfo.ExamMeta.TenantId
}

func GetStudentId(student *ExamStudentDetail) int64 {

	if student == nil || student.StudentMeta == nil {
		return 0
	}
	return student.StudentMeta.StudentId
}

func GetStudentNo(student *ExamStudentDetail) string {
	if student == nil || student.StudentMeta == nil {
		return ""
	}
	return student.StudentMeta.StudentCode
}

func GetStudentName(student *ExamStudentDetail) string {

	if student == nil || student.StudentMeta == nil {
		return ""
	}
	return student.StudentMeta.StudentName
}

func GetSchoolIdByStudent(student *ExamStudentDetail) int64 {

	if student == nil || student.School == nil {
		return 0
	}
	return student.School.SchoolId
}

func GetSchoolNameByStudent(student *ExamStudentDetail) string {

	if student == nil || student.School == nil {
		return ""
	}
	return student.School.SchoolName
}

func GetCampusIdByStudent(student *ExamStudentDetail) int64 {
	if student == nil || student.Campus == nil {
		return 0
	}
	return student.Campus.CampusId
}

func GetCampusNameByStudent(student *ExamStudentDetail) string {
	if student == nil || student.Campus == nil {
		return ""
	}
	return student.Campus.CampusName
}

func GetAdministrativeClassIdByStudent(student *ExamStudentDetail) int64 {

	if student == nil || len(student.ClassList) <= 0 {
		return 0
	}
	for _, classMeta := range student.ClassList {
		if classMeta == nil {
			continue
		}
		if classMeta.ClassType == organization_enum.ClassType_ClassTypeAdministrativeClass {
			return classMeta.ClassId
		}
	}
	return 0
}

func GetAdministrativeClassNameByStudent(student *ExamStudentDetail) string {

	if student == nil || len(student.ClassList) <= 0 {
		return ""
	}
	for _, classMeta := range student.ClassList {
		if classMeta == nil {
			continue
		}
		if classMeta.ClassType == organization_enum.ClassType_ClassTypeAdministrativeClass {
			return classMeta.ClassName
		}
	}
	return ""
}

func GetCourseClassNameByStudent(student *ExamStudentDetail) string {

	if student == nil || len(student.ClassList) <= 0 {
		return ""
	}
	for _, classMeta := range student.ClassList {
		if classMeta == nil {
			continue
		}
		if classMeta.ClassType == organization_enum.ClassType_ClassTypeCourseClass {
			return classMeta.ClassName
		}
	}
	return ""
}

func GetGradeIdByStudent(student *ExamStudentDetail) int64 {
	if student == nil || student.Grade == nil {
		return 0
	}
	return student.Grade.GradeId
}

func GetGradeNameByStudent(student *ExamStudentDetail) string {
	if student == nil || student.Grade == nil {
		return ""
	}
	return student.Grade.GradeName
}

func GetExamNoByStudent(student *ExamStudentDetail) string {
	if student == nil || student.StudentMeta == nil {
		return ""
	}
	return student.StudentMeta.ExamCode
}

func GetAdministrativeClassByStudent(student *ExamStudentDetail) int64 {
	if student == nil || len(student.ClassList) <= 0 {
		return 0
	}
	for _, classDetail := range student.ClassList {
		if classDetail == nil {
			continue
		}
		if classDetail.ClassType == organization_enum.ClassType_ClassTypeAdministrativeClass {
			return classDetail.ClassId
		}
	}
	return 0
}

func GetExamIdByPaper(unionExam *UnionExam, paperId, schoolId int64) int64 {

	examInfo := GetExamInfoByPaperSchool(unionExam, paperId, schoolId)
	if examInfo == nil {
		return 0
	}
	if examInfo.ExamMeta == nil {
		return 0
	}
	return examInfo.ExamMeta.ExamId
}

func GetScoringModeByExam(exam *Exam) common.ScoringMode {
	if exam == nil || exam.ExamMeta == nil {
		return common.ScoringMode_SMUnknown
	}
	return exam.ExamMeta.ScoringMode
}

func IsScanMixed(unionExam *UnionExam, paperId int64) bool {

	if unionExam == nil || len(unionExam.ExamGroupList) <= 0 {
		return false
	}
	for _, examGroup := range unionExam.ExamGroupList {
		if len(examGroup.ExamList) <= 0 {
			continue
		}
		for _, examInfo := range examGroup.ExamList {
			if examInfo == nil || examInfo.ExamMeta == nil {
				continue
			}
			if examInfo.ExamMeta.PaperId == paperId {
				return examInfo.ExamMeta.ScanType == common.ScanType_ESTMixScan
			}
		}
	}
	return false
}

func IsScanMixedByExamSubject(unionExam *UnionExam, examSubject int32) bool {

	if unionExam == nil || len(unionExam.ExamGroupList) <= 0 {
		return false
	}
	for _, examGroup := range unionExam.ExamGroupList {
		if len(examGroup.ExamList) <= 0 {
			continue
		}
		for _, examInfo := range examGroup.ExamList {
			if examInfo == nil || examInfo.ExamMeta == nil {
				continue
			}
			if examInfo.ExamMeta.ExamSubject == examSubject {
				return examInfo.ExamMeta.ScanType == common.ScanType_ESTMixScan
			}
		}
	}
	return false
}

func IsHomeworkByUnion(unionExam *UnionExam) bool {

	if unionExam == nil || len(unionExam.ExamGroupList) <= 0 {
		return false
	}
	for _, examGroup := range unionExam.ExamGroupList {
		if len(examGroup.ExamList) <= 0 {
			continue
		}
		for _, examInfo := range examGroup.ExamList {
			if examInfo == nil || examInfo.ExamMeta == nil {
				continue
			}
			if examInfo.ExamMeta.ExamCategory == common.ExamCategory_ECHomeWork {
				return true
			}
		}
	}
	return false
}

func IsHomeworkByExam(examInfo *Exam) bool {

	if examInfo == nil || examInfo.ExamMeta == nil {
		return false
	}
	if examInfo.ExamMeta.ExamCategory == common.ExamCategory_ECHomeWork {
		return true
	}
	return false
}

func GetExamCategoryExam(examInfo *Exam) common.ExamCategory {

	if examInfo == nil || examInfo.ExamMeta == nil {
		return common.ExamCategory_ECUnknown
	}
	return examInfo.ExamMeta.ExamCategory
}

func GetExamCategoryByUnion(unionExam *UnionExam, paperId int64) common.ExamCategory {

	if unionExam == nil || len(unionExam.ExamGroupList) <= 0 {
		return common.ExamCategory_ECUnknown
	}
	for _, examGroup := range unionExam.ExamGroupList {
		if len(examGroup.ExamList) <= 0 {
			continue
		}
		for _, examInfo := range examGroup.ExamList {
			if examInfo == nil || examInfo.ExamMeta == nil {
				continue
			}
			if examInfo.ExamMeta.PaperId == paperId {
				return examInfo.ExamMeta.ExamCategory
			}
		}
	}
	return common.ExamCategory_ECUnknown
}

type SearchTeacherHomeworkParam struct {
	InstitutionId int64 // 机构id
	// 内部实现，根据用户权限过滤, 高职务教师，根据职务范围，取数据列表，比如 校长、学科主任、年级主任
	TenantId       int64                   // 租户ID
	UserId         int64                   // 用户id
	QueryCondition *QueryHomeworkCondition // 查询条件
	Limit          int32                   // 分页参数
	Offset         int32                   // 分页参数
}

func NewSearchTeacherHomeworkReq(param *SearchTeacherHomeworkParam) *exam.SearchTeacherHomeworkRequest {

	if param == nil {
		return nil
	}
	return &exam.SearchTeacherHomeworkRequest{
		InstitutionId: param.InstitutionId,
		TenantId:      param.TenantId,
		UserId:        param.UserId,
		//QueryCondition: convQueryHomeworkCondition(param.QueryCondition),
		Limit:  param.Limit,
		Offset: param.Offset,
	}
}

type QueryHomeworkCondition struct {
	ExamName          *string                     // 考试名称
	DurationStartTime *int64                      // 时间范围开始时间
	DurationEndTime   *int64                      // 时间范围结束时间
	CorrectionType    []common.CorrectionType     // 筛选考试类型列表
	TeacherTaskType   []exam_enum.TeacherTaskType // 根据老师任务类型筛选
	ExamStatus        []common.ExamStatus         // 根据考试状态筛选
	SchoolYear        *int32                      // 学年
	UnionExamId       *int64                      // 联考id
	UnionExamStatus   *common.UnionExamStatus     // 联考状态
	ExamTeacherTypes  []common.ExamTeacherType    // 老师角色类型
	InstitutionIds    []int64                     // 查询机构列表；若为空则默认基于用户搜索（即用户权限可见考试）
	ExamSubject       *int32                      // 考试科目
	BaseGrade         *int32                      // 参考年级
	ExamCreatorName   *string                     // 考试创建人名
	ExamSources       []exam_enum.ExamSource      //  1：机构端，2：教师端
	ExamCategory      common.ExamCategory         // 考试、作业
	CampusIds         []int64                     // 参考校区id
	ClassIds          []int64                     // 参考班级id
}

func convQueryHomeworkCondition(data *QueryHomeworkCondition) *exam.QueryCondition {
	if data == nil {
		return nil
	}
	return &exam.QueryCondition{
		ExamName:          data.ExamName,
		DurationStartTime: data.DurationStartTime,
		DurationEndTime:   data.DurationEndTime,
		CorrectionType:    data.CorrectionType,
		TeacherTaskType:   data.TeacherTaskType,
		ExamStatus:        data.ExamStatus,
		SchoolYear:        data.SchoolYear,
		UnionExamId:       data.UnionExamId,
		UnionExamStatus:   data.UnionExamStatus,
		ExamTeacherTypes:  data.ExamTeacherTypes,
		InstitutionIds:    data.InstitutionIds,
		ExamSubject:       data.ExamSubject,
		BaseGrade:         data.BaseGrade,
		ExamCreatorName:   data.ExamCreatorName,
		ExamSources:       data.ExamSources,
		ExamCategory:      data.ExamCategory,
		CampusIds:         data.CampusIds,
		ClassIds:          data.ClassIds,
	}
}

type SearchTeacherHomeworkResult struct {
	Items []*Homework `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	Total int32       `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
}

type Homework struct {
	ExamId          int64                     // 考试id(作业ID)
	Name            string                    // 作业名称
	ExamType        common.ExamType           // 考试类型
	InstitutionId   int64                     // 机构id
	BaseGrade       common.GradeValue         // 参考年级
	PaperId         int64                     // 试卷id
	ExamSubject     common.ExamSubjectValue   // 考试科目
	ExamSubjectName string                    // 考试科目名称
	GradeInfo       *ExamSchoolGrade          // 年级、班级名称
	Creator         int64                     // 创建人id
	CreatorName     string                    // 创建人名称
	UnionExamId     int64                     // 联考id
	ExamGroupId     int64                     // 联考组id
	StatisticStatus common.StatisticsStatus   // 学情发布状态
	MarkingMode     common.CorrectionMarkType // 批分模式
	SchoolYear      int32                     // 学年
	CreateAt        int64                     // 创建时间戳（秒）
	UpdateAt        int64                     // 更新时间戳（秒）
	Stage           organization_enum.Stage   // 学段
	CorrectionType  common.CorrectionType     // 阅卷类型， 固定值：5
	ExamCategory    common.ExamCategory       // 考试、作业
	TenantId        int64                     // 租户id
	BaseSubjects    []int32                   // 基础学科
}

func NewSearchTeacherHomeworkResult(data *exam.SearchTeacherHomeworkResponse) *SearchTeacherHomeworkResult {
	if data == nil {
		return nil
	}
	return &SearchTeacherHomeworkResult{
		Items: lists.Map(data.Items, convHomework),
		Total: data.Total,
	}
}

func convHomework(data *exam.Homework) *Homework {
	if data == nil {
		return nil
	}
	return &Homework{
		ExamId:          data.ExamId,
		Name:            data.Name,
		ExamType:        data.ExamType,
		InstitutionId:   data.InstitutionId,
		BaseGrade:       data.BaseGrade,
		PaperId:         data.PaperId,
		ExamSubject:     data.ExamSubject,
		ExamSubjectName: data.ExamSubjectName,
		GradeInfo:       convExamSchoolGrade(data.GradeInfo),
		Creator:         data.Creator,
		CreatorName:     data.CreatorName,
		UnionExamId:     data.UnionExamId,
		ExamGroupId:     data.ExamGroupId,
		StatisticStatus: data.StatisticStatus,
		MarkingMode:     data.MarkingMode,
		SchoolYear:      data.SchoolYear,
		CreateAt:        data.CreateAt,
		UpdateAt:        data.UpdateAt,
		Stage:           data.Stage,
		CorrectionType:  data.CorrectionType,
		ExamCategory:    data.ExamCategory,
		TenantId:        data.TenantId,
		BaseSubjects:    data.BaseSubjects,
	}
}

//type GetHomeworkDetailParam struct {
//	HomeworkId int64 // 作业id
//}
//
//func NewGetHomeworkDetailReq(param *GetHomeworkDetailParam) *exam.GetHomeworkDetailRequest {
//	if param == nil {
//		return nil
//	}
//	return &exam.GetHomeworkDetailRequest{
//		ExamId: param.HomeworkId,
//	}
//}
//
//type GetHomeworkDetailResult struct {
//	Homework *Homework // 作业详情
//}
//
//func NewGetHomeworkDetailResult(data *exam.GetHomeworkDetailResponse) *GetHomeworkDetailResult {
//	if data == nil {
//		return nil
//	}
//	return &GetHomeworkDetailResult{
//		Homework: convHomework(data.Homework),
//	}
//}

func GetExamInfoByPaperSchool(unionExam *UnionExam, paperId int64, schoolId int64) *Exam {

	if unionExam == nil || len(unionExam.ExamGroupList) <= 0 {
		return nil
	}
	for _, group := range unionExam.ExamGroupList {

		if group == nil || len(group.ExamList) <= 0 {
			continue
		}
		for _, examInfo := range group.ExamList {
			if examInfo == nil || examInfo.ExamMeta == nil {
				continue
			}
			if examInfo.ExamMeta.PaperId == paperId && examInfo.ExamMeta.TenantId == schoolId {
				return examInfo
			}
		}
	}
	return nil
}

func GetExamIdBySubjectSchool(unionExam *UnionExam, schoolId int64, subject int32) *Exam {

	if unionExam == nil || len(unionExam.ExamGroupList) <= 0 {
		return nil
	}
	for _, group := range unionExam.ExamGroupList {

		if group == nil || len(group.ExamList) <= 0 {
			continue
		}
		for _, examInfo := range group.ExamList {
			if examInfo == nil || examInfo.ExamMeta == nil {
				continue
			}
			if examInfo.ExamMeta.ExamSubject == subject && examInfo.ExamMeta.TenantId == schoolId {
				return examInfo
			}
		}
	}
	return nil
}

type GetCorrectionRuleParam struct {
	ExamId int64 // 作业id
}

type GetCorrectionRuleResult struct {
	Rule       *CorrectionRuleInfo // 批改规则
	ExampleUrl string              // 示例图
}
type CorrectionRuleInfo struct {
	Id        int64
	ExamId    int64                          // 作业 id
	Right     common.CorrectionRule          // 判对
	Wrong     common.CorrectionRule          // 判错
	HalfRight common.CorrectionRule          // 判半对
	NoCorrect common.CorrectionRule          // 不批改
	Status    exam_enum.CorrectionRuleStatus // 批改规则状态
	UpdaterId int64                          // 更新人
}

func convCorrectionRule(data *exam.CorrectionRuleInfo) *CorrectionRuleInfo {
	if data == nil {
		return nil
	}
	return &CorrectionRuleInfo{
		Id:        data.Id,
		ExamId:    data.ExamId,
		Right:     data.Right,
		Wrong:     data.Wrong,
		HalfRight: data.HalfRight,
		NoCorrect: data.NoCorrect,
		Status:    data.Status,
		UpdaterId: data.UpdaterId,
	}
}

func NewGetCorrectionRuleReq(param *GetCorrectionRuleParam) *exam.GetCorrectionRuleRequest {
	if param == nil {
		return nil
	}
	return &exam.GetCorrectionRuleRequest{
		ExamId: param.ExamId,
	}
}

func NewGetCorrectionRuleResult(data *exam.GetCorrectionRuleResponse) *GetCorrectionRuleResult {
	if data == nil {
		return nil
	}
	return &GetCorrectionRuleResult{
		Rule:       convCorrectionRule(data.Rule),
		ExampleUrl: data.ExampleUrl,
	}
}

func GetHomeworkIdFromUnionExam(unionExam *UnionExam) int64 {

	if unionExam == nil || len(unionExam.ExamGroupList) <= 0 {
		return 0
	}
	for _, group := range unionExam.ExamGroupList {
		if len(group.ExamList) <= 0 {
			continue
		}
		for _, e := range group.ExamList {
			if e == nil || e.ExamMeta == nil {
				continue
			}
			if e.ExamMeta.ExamCategory != common.ExamCategory_ECHomeWork {
				return 0
			}
			if e.ExamMeta.ExamId > 0 {
				return e.ExamMeta.ExamId
			}
		}
	}
	return 0
}

func IsInExamByClassId(examInfo *Exam, schoolId int64, classId int64) bool {

	if examInfo == nil || len(examInfo.ClassList) <= 0 {
		return false
	}
	for _, classInfo := range examInfo.ClassList {
		if classInfo == nil {
			continue
		}
		if classInfo.TenantId == schoolId && classInfo.ClassId == classId {
			return true
		}
	}
	return false
}

func FilterExamClassStudents(students []*ExamStudentDetail, examClassId int64) []*ExamStudentDetail {

	resultStudents := make([]*ExamStudentDetail, 0, len(students))
	for _, student := range students {
		if student == nil {
			continue
		}
		if len(student.ClassList) <= 0 {
			continue
		}
		for _, classMeta := range student.ClassList {
			if classMeta == nil {
				continue
			}
			if !classMeta.IsExamClass {
				continue
			}
			if classMeta.ClassId == examClassId {
				resultStudents = append(resultStudents, student)
				break
			}
		}
	}
	return resultStudents
}

func GetStudentExamClassMeta(studentInfo *ExamStudentDetail) *ExamClassMeta {

	if studentInfo == nil {
		return nil
	}
	if len(studentInfo.ClassList) <= 0 {
		return nil
	}
	var courseClass *ExamClassMeta
	var adminClass *ExamClassMeta

	for _, classInfo := range studentInfo.ClassList {
		if classInfo == nil {
			continue
		}
		if !classInfo.IsExamClass {
			continue
		}
		if classInfo.ClassType == organization_enum.ClassType_ClassTypeCourseClass {
			courseClass = classInfo
		}
		if classInfo.ClassType == organization_enum.ClassType_ClassTypeAdministrativeClass {
			adminClass = classInfo
		}
	}
	if courseClass != nil { // 优先选择走班
		return courseClass
	}
	return adminClass
}

func IsScanTeacher(examId int64, authInfo *GetAuthorityOfTeacherResult) bool {

	if examId <= 0 {
		return false
	}
	if authInfo == nil || authInfo.ExamSubjectRightsMap == nil {
		return false
	}
	rightInfo, ok := authInfo.ExamSubjectRightsMap[examId]
	if !ok || rightInfo == nil {
		return false
	}
	for _, right := range rightInfo.ExamSubjectRights {
		if right == exam_enum.ExamSubjectRightType_ESRTScan {
			return true
		}
	}
	return false
}
