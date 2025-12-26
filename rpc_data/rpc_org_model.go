package rpc_data

import (
	"code.chenji.com/cj/scan_paper/utils/lists"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/organization/organization_enum"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/organization/organization_rpc"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/organization/organization_rpc_model"
)

type StudentDetail struct {
	StudentBaseInfo *StudentBaseInfo // 学生基本信息
	ClassInfos      []*ClassDetail   // 班级信息

}

func ConvStudentDetail(info *organization_rpc.StudentInfo) *StudentDetail {
	if info == nil {
		return nil
	}
	return &StudentDetail{
		StudentBaseInfo: convStudentBaseInfo(info.Student),
		ClassInfos:      lists.Map(info.Classes, convClassDetail),
	}
}

type StudentBaseInfo struct {
	Id            int64                           // 学生表主键id(业务不要用)
	SchoolId      int64                           // 学校id
	Name          string                          // 学生姓名
	StudentNo     string                          // 学生编号
	StudentCardNo string                          // 学生学籍号
	Avatar        string                          // 头像
	Gender        int32                           // 性别 0:未知 1:男 2:女
	Mobile        string                          // 手机号,默认不输出 输出需时脱敏
	Birthday      int64                           // 生日
	Status        organization_enum.StudentStatus // 状态: 1:正常、2:退学、3:转学、4:休学 、5:删除
	TenantId      int64                           // 租户id
	DotMatrixCode int64                           // 点阵码
	Account       string                          // 登陆账号
	ThirdCode     *string                         // 三方code
	UserId        int64                           // 用户id， 业务id使用，这个一定是要和租户id在一起才有意义
	ShortExamCode string                          // 短考号

}

func convStudentBaseInfo(info *organization_rpc.StudentBaseInfo) *StudentBaseInfo {

	if info == nil {
		return nil
	}
	return &StudentBaseInfo{
		Id:            info.Id,
		SchoolId:      info.SchoolId,
		Name:          info.Name,
		StudentNo:     info.StudentNo,
		StudentCardNo: info.StudentCardNo,
		Avatar:        info.Avatar,
		Gender:        info.Gender,
		Mobile:        info.Mobile,
		Birthday:      info.Birthday,
		Status:        info.Status,
		TenantId:      info.TenantId,
		Account:       info.Account,
		ThirdCode:     info.ThirdCode,
		UserId:        info.UserId,
		ShortExamCode: info.ShortExamCode,
		//DotMatrixCode: 0,
	}
}

type ClassDetail struct {
	ClassInfo  *ClassInfo  // 班级信息
	StageInfo  *StageInfo  // 学段信息
	GradeInfo  *GradeInfo  // 年级信息
	SchoolInfo *SchoolInfo // 租户信息
	CampusInfo *CampusInfo // 校区信息
}

func convClassDetail(info *organization_rpc.ClassDetailInfo) *ClassDetail {

	if info == nil {
		return nil
	}
	return &ClassDetail{
		ClassInfo:  convClassInfo(info.ClassInfo),
		StageInfo:  convStageInfo(info.StageInfo),
		GradeInfo:  convGradeInfo(info.GradeInfo),
		SchoolInfo: convSchoolInfo(info.SchoolInfo),
		CampusInfo: conCampusInfo(info.CampusInfo),
	}
}

type StageInfo struct {
	Stage int32  // 学段
	Name  string // 学段名称
}

func convStageInfo(info *organization_rpc_model.StageInfo) *StageInfo {

	if info == nil {
		return nil
	}
	return &StageInfo{
		Stage: info.Stage,
		Name:  info.Name,
	}
}

type ClassInfo struct {
	Id          int64                                   // 班级id
	Name        string                                  // 班级名称(可选)
	ClassType   organization_enum.ClassType             // 班级类型：1-行政班 2-课程班, 3: 虚拟班
	SubjectType organization_enum.SubjectTypeFieldScope // 学科类型范围
}

func convClassInfo(info *organization_rpc.ClassInfo) *ClassInfo {

	if info == nil {
		return nil
	}
	return &ClassInfo{
		Id:          info.Id,
		Name:        info.Name,
		ClassType:   info.ClassType,
		SubjectType: info.SubjectType,
	}
}

type GradeInfo struct {
	Id         int64  // 年级Id
	Name       string // 年级名称
	SchoolYear int32  // 学年
	BaseGrade  int32  // 基础年级Id
}

func convGradeInfo(info *organization_rpc.GradeInfo) *GradeInfo {

	if info == nil {
		return nil
	}
	return &GradeInfo{
		Id:         info.Id,
		Name:       info.Name,
		SchoolYear: info.SchoolYear,
		BaseGrade:  info.BaseGrade,
	}
}

type SchoolInfo struct {
	SchoolId   int64  // 学校id
	SchoolName string // 学校名称
}

func convSchoolInfo(info *organization_rpc.SchoolInfo) *SchoolInfo {

	if info == nil {
		return nil
	}
	return &SchoolInfo{
		SchoolId:   info.SchoolId,
		SchoolName: info.Name,
	}
}

type CampusInfo struct {
	CampusId        int64  // 校区id
	CampusName      string // 校区名称
	IsDefaultCampus bool   // 是否默认校区
}

func conCampusInfo(info *organization_rpc_model.CampusInfo) *CampusInfo {

	if info == nil {
		return nil
	}
	return &CampusInfo{
		CampusId:        info.CampusTenantId,
		CampusName:      info.Name,
		IsDefaultCampus: info.DefaultCampus,
	}
}

type UserInfo struct {
	Account          *Account                 // 账户, 默认account, 应用全场景，后期扩展，支持不同场景，使用不同账户
	Mobile           string                   // 手机号，安全风险，默认不输出，如需返回，脱敏输出
	Gender           organization_enum.Gender // 性别: 0:未知 1:男、1:女
	Avatar           string                   // 头像
	Name             string                   // 姓名
	UserId           int64                    // 用户id， 业务上使用该字段
	CloseAccountTime int64                    // 账号注销时间

}

func ConvUserInfo(info *organization_rpc.UserInfo) *UserInfo {
	if info == nil {
		return nil
	}
	return &UserInfo{
		Account:          convAccount(info.Account),
		Mobile:           info.Mobile,
		Gender:           info.Gender,
		Avatar:           info.Avatar,
		Name:             info.Name,
		UserId:           info.UserId,
		CloseAccountTime: info.CloseAccountTime,
	}
}

type Account struct {
	Id               int64
	Account          string                          // 登录账号（用户名）
	AppId            int64                           // 应用id, 默认0， 代表通用全场景
	FirstLogin       int32                           // 是否是首次登陆 0:未知 1:首次登录 2:非首次登录
	CloseAccountTime int64                           // 账号注销时间
	LastLoginTime    int64                           // 最后登陆时间
	UserId           int64                           // 用户id， 业务上使用该字段
	Source           int64                           // 账号来源
	Status           organization_enum.AccountStatus // 账号状态
}

func convAccount(info *organization_rpc.Account) *Account {

	if info == nil {
		return nil
	}
	return &Account{
		Id:               info.Id,
		Account:          info.Account,
		AppId:            info.AppId,
		FirstLogin:       info.FirstLogin,
		CloseAccountTime: info.CloseAccountTime,
		LastLoginTime:    info.LastLoginTime,
		UserId:           info.UserId,
		Source:           info.Source,
		Status:           info.Status,
	}
}

type Class struct {
	Id                 int64                                   // ID
	TenantId           int64                                   // 租户ID
	SchoolYear         int32                                   // 学年
	Stage              organization_enum.Stage                 // 学段
	GradeId            int64                                   // 年级ID
	ClassType          organization_enum.ClassType             // 班级类型：1-行政班 2-课程班
	Name               string                                  // 班级名称
	Status             organization_enum.ClassStatus           // 状态: 0-未知 1-已启用 2-升学年中，新学年启用后，新学年才生效 3-毕业, 毕业、结束 都表示结束 4-结束
	SortNo             int64                                   // 排序
	CreatorId          int64                                   // 创建人
	UpdaterId          int64                                   // 修改人
	IncludeSubjectInfo []*ClassSubject                         // 包含学科（走班、虚拟班有学科信息，行政班无）
	SubjectType        organization_enum.SubjectTypeFieldScope // 学科类型
	UserType           organization_enum.UseType               // 使用类型
	SchoolTenantId     int64                                   // 学校租户ID
	BaseGrade          int32                                   // 基础年级
}

func ConvClass(info *organization_rpc_model.Class) *Class {

	if info == nil {
		return nil
	}
	return &Class{
		Id:                 info.Id,
		TenantId:           info.TenantId,
		SchoolYear:         info.SchoolYear,
		Stage:              info.Stage,
		GradeId:            info.GradeId,
		ClassType:          info.ClassType,
		Name:               info.Name,
		Status:             info.Status,
		SortNo:             info.SortNo,
		CreatorId:          info.CreatorId,
		UpdaterId:          info.UpdaterId,
		IncludeSubjectInfo: lists.Map(info.IncludeSubjectInfo, convClassSubject),
		SubjectType:        info.SubjectType,
		UserType:           info.UserType,
		SchoolTenantId:     info.SchoolTenantId,
		BaseGrade:          info.BaseGrade,
	}
}

// ClassSubject 班级学科
type ClassSubject struct {
	SubjectId   int64 // 学科ID
	BaseSubject int32 // 基础学科(基础标签服务)
}

func convClassSubject(info *organization_rpc_model.ClassSubject) *ClassSubject {
	if info == nil {
		return nil
	}
	return &ClassSubject{
		SubjectId:   info.SubjectId,
		BaseSubject: info.BaseSubject,
	}
}

type TenantBaseInfo struct {
	Id              int64                          // 租户id
	ParentId        int64                          // 所属租户id
	RootId          int64                          // 根节点id
	Name            string                         // 租户名称
	TenantType      organization_enum.TenantType   // 租户类型
	AreaId          int64                          // 所在地区
	DetailedAddress string                         // 详细地址
	AdminPhone      string                         // 管理员手机号
	AdminName       string                         // 管理员姓名
	UseType         organization_enum.UseType      // 使用类型
	OpenTime        int64                          // 开通时间
	CloseTime       int64                          // 到期时间
	Status          organization_enum.TenantStatus // 状态

}

func ConvTenantBaseInfo(info *organization_rpc_model.Tenant) *TenantBaseInfo {

	return &TenantBaseInfo{}
}

type TenantDetail struct {
	TenantBaseInfo *TenantBaseInfo // 租户基本信息
	SchoolRef      *SchoolRef      // 学校信息
}

func ConvTenantDetail(info *organization_rpc_model.TenantRef) *TenantDetail {

	if info == nil {
		return nil
	}
	return &TenantDetail{
		TenantBaseInfo: ConvTenantBaseInfo(info.TenantInfo),
		SchoolRef:      ConvSchoolRef(info.School),
	}
}

type SchoolRef struct {
	TenantId        int64                          // 学校租户id
	Name            string                         // 学校名称
	SchoolType      organization_enum.SchoolType   // 学校办学类型
	SchoolSystem    organization_enum.SchoolSystem // 学校学制
	SchoolOffice    organization_enum.SchoolOffice // 学校办别
	Stages          []organization_enum.Stage      // 学段列表
	CampusList      []*CampusRef                   // 校区列表
	DefaultCampusId int64                          // 默认校区租户id
	SchoolYear      int32                          // 学年
}

func ConvSchoolRef(info *organization_rpc_model.SchoolRef) *SchoolRef {

	if info == nil {
		return nil
	}
	return &SchoolRef{
		TenantId:        info.TenantId,
		Name:            info.Name,
		SchoolType:      info.SchoolType,
		SchoolSystem:    info.SchoolSystem,
		SchoolOffice:    info.SchoolOffice,
		Stages:          info.Stages,
		CampusList:      lists.Map(info.CampusList, ConvCampusRef),
		DefaultCampusId: info.DefaultCampusId,
		SchoolYear:      info.SchoolYear,
	}
}

type CampusRef struct {
	CampusId       int64                     // 校区id
	Name           string                    // 校区名称
	SchoolTenantId int64                     // 所属学校
	Stages         []organization_enum.Stage // 学段列表
	DefaultCampus  bool                      // 是否为默认校区 1-是 2-否
	ThirdTenantId  *string                   // 第三方校区code
}

func ConvCampusRef(info *organization_rpc_model.CampusRef) *CampusRef {

	if info == nil {
		return nil
	}
	return &CampusRef{
		CampusId:       info.CampusId,
		Name:           info.Name,
		SchoolTenantId: info.SchoolTenantId,
		Stages:         info.Stages,
		DefaultCampus:  info.DefaultCampus,
		ThirdTenantId:  info.ThirdTenantId,
	}
}

type GetTenantUserInfoParam struct {
	TenantId       int64                                         // 租户id, 当传该值时，只返回一条数据
	SchoolTenantId int64                                         // 父租户id， 可选，如学校
	UserId         int64                                         // 用户id
	Statuses       []organization_enum.InstitutionUserInfoStatus // 默认只返回 normal 状态的数据
	TokenInfo      *TokenInfo                                    // 当前登陆人信息
}

type TokenInfo struct {
	AppId     int64  // 当前登陆的应用id
	AccountId int64  // 当前登陆的帐户id，仅在登陆相关时使用
	UserId    int64  // 当前登陆的用户id，用于记录行为日志等
	UserName  string // 当前登陆的用户 name
}

func convTokenInfo(info *TokenInfo) *organization_rpc.TokenInfo {
	if info == nil {
		return nil
	}
	return &organization_rpc.TokenInfo{
		AppId:     info.AppId,
		AccountId: info.AccountId,
		UserId:    info.UserId,
		UserName:  info.UserName,
	}
}

type TenantUserInfo struct {
	UserInfo   *InstitutionUserInfo // 人员信息
	DeptList   []*DeptInfo          // 人员所属部门(机构人员&老师)
	TenantInfo *TenantInfo          // 租户信息
	SchoolInfo *SchoolInfo          // 学校校区，仅对校区人员有意义
	Roles      []*UserRoleInfo      // 角色信息
}

type InstitutionUserInfo struct {
	Id          int64                                       // 机构人员表主键id(业务别用)
	UserId      int64                                       // 用户id， 业务方用这个id！
	TenantId    int64                                       // 租户ID
	Name        string                                      // 机构人员姓名
	UserNo      string                                      // 机构人员编号
	Avatar      string                                      // 头像
	Phone       string                                      // 手机号, 安全风险，脱敏输出
	Gender      organization_enum.Gender                    // 性别: 0:未知 1:男、1:女
	Status      organization_enum.InstitutionUserInfoStatus // 状态: 正常、离职、退休等
	ThirdUserId *string                                     // 三方老师id
	UseType     organization_enum.UseType                   // 使用类型类型

}

func convInstitutionUserInfo(info *organization_rpc.InstitutionUserInfo) *InstitutionUserInfo {
	if info == nil {
		return nil
	}
	return &InstitutionUserInfo{
		Id:          info.Id,
		UserId:      info.UserId,
		TenantId:    info.TenantId,
		Name:        info.Name,
		UserNo:      info.UserNo,
		Avatar:      info.Avatar,
		Phone:       info.Phone,
		Gender:      info.Gender,
		Status:      info.Status,
		ThirdUserId: info.ThirdUserId,
		UseType:     info.UseType,
	}
}

type DeptInfo struct {
	Id   int64  // 部门id
	Name string // 部门名称
}

func convDeptInfo(info *organization_rpc.DeptInfo) *DeptInfo {
	if info == nil {
		return nil
	}
	return &DeptInfo{
		Id:   info.Id,
		Name: info.Name,
	}
}

type TenantInfo struct {
	Id            int64                               // 班级id
	Name          string                              // 名称
	TenantType    organization_enum.TenantType        // 租户类型
	ThirdTenantId string                              // 三方租户id
	ThirdSource   organization_enum.ThirdTenantSource // 三方租户类型

}

func convTenantInfo(info *organization_rpc.TenantInfo) *TenantInfo {
	if info == nil {
		return nil
	}
	return &TenantInfo{
		Id:            info.Id,
		Name:          info.Name,
		TenantType:    info.TenantType,
		ThirdTenantId: info.ThirdTenantId,
		ThirdSource:   info.ThirdSource,
	}
}

type UserRoleInfo struct {
	Id         int64                          // 用户&角色 关系 id
	Name       string                         // 角色名称
	RoleType   organization_enum.RoleType     // 0:未知、1:机构、6:学校、3:运营, 7:校区
	Category   organization_enum.RoleCategory // 0:未知、1:内置角色、2:常规角色
	RoleCode   organization_enum.BuiltInRole  // 角色编号
	SchoolYear int32                          // 角色学年
	// 租户名称 ，如学校/校区名
	TenantName string
	// 租户id
	TenantId int64
	RoleId   int64 // 角色id

}

func convUserRoleInfo(info *organization_rpc.UserRoleInfo) *UserRoleInfo {
	if info == nil {
		return nil
	}
	return &UserRoleInfo{
		Id:         info.Id,
		Name:       info.Name,
		RoleType:   info.RoleType,
		Category:   info.Category,
		RoleCode:   info.RoleCode,
		SchoolYear: info.SchoolYear,
		TenantName: info.TenantName,
		TenantId:   info.TenantId,
		RoleId:     info.RoleId,
	}
}

func NewTenantUserInfoReq(param *GetTenantUserInfoParam) *organization_rpc.GetTenantUserInfoRequest {

	if param == nil {
		return nil
	}
	return &organization_rpc.GetTenantUserInfoRequest{
		TenantId:       param.TenantId,
		SchoolTenantId: param.SchoolTenantId,
		UserId:         param.UserId,
		Statuses:       param.Statuses,
		TokenInfo:      convTokenInfo(param.TokenInfo),
	}
}
func NewTenantUserInfo(resp *organization_rpc.GetTenantUserInfoResponse) *TenantUserInfo {

	if resp == nil || resp.UserDetailInfo == nil || len(resp.UserDetailInfo) <= 0 {
		return nil
	}
	userDetail := resp.UserDetailInfo[0]
	if userDetail == nil {
		return nil
	}
	return &TenantUserInfo{
		UserInfo:   convInstitutionUserInfo(userDetail.UserInfo),
		DeptList:   lists.Map(userDetail.DeptList, convDeptInfo),
		TenantInfo: convTenantInfo(userDetail.TenantInfo),
		SchoolInfo: convSchoolInfo(userDetail.SchoolInfo),
		Roles:      lists.Map(userDetail.Roles, convUserRoleInfo),
	}
}

func NewInstitutionUser(info *organization_rpc.InstitutionUserInfo) *InstitutionUser {

	if info == nil {
		return nil
	}
	return &InstitutionUser{
		Id:          info.Id,
		UserId:      info.UserId,
		TenantId:    info.TenantId,
		Name:        info.Name,
		UserNo:      info.UserNo,
		Avatar:      info.Avatar,
		Phone:       info.Phone,
		Gender:      info.Gender,
		Status:      info.Status,
		ThirdUserId: info.ThirdUserId,
		UseType:     info.UseType,
	}
}

type BatchGetInstitutionUserByUserIdParam struct {
	TenantId       int64                                         // 租户id, 当有值时，只返回该租户下的用户数据
	SchoolTenantId int64                                         // 学校id
	UserIds        []int64                                       // 机构人员用户id
	Statuses       []organization_enum.InstitutionUserInfoStatus // 状态，如果不带获取normal状态的
}

type InstitutionUser struct {
	Id          int64                                       // 机构人员表主键id(业务别用)
	UserId      int64                                       // 用户id， 业务方用这个id！
	TenantId    int64                                       // 租户ID
	Name        string                                      // 机构人员姓名
	UserNo      string                                      // 机构人员编号
	Avatar      string                                      // 头像
	Phone       string                                      // 手机号, 安全风险，脱敏输出
	Gender      organization_enum.Gender                    // 性别: 0:未知 1:男、1:女
	Status      organization_enum.InstitutionUserInfoStatus // 状态: 正常、离职、退休等
	ThirdUserId *string                                     // 三方老师id
	UseType     organization_enum.UseType                   // 使用类型类型
}
