package rpc_data

import (
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/mark_paper/mark_paper_enum"
	"code.chenji.com/pkg/idlgen/kitex_gen/edu/mark_paper/mark_rpc_model"
)

type StudentResult struct {
	UnionExamId int64                             // 联考id
	ExamId      int64                             // 考试id
	PaperId     int64                             // 试卷id
	Type        mark_paper_enum.ScanRecordType    // 扫描类型：1-保存 2-更新客观题 3-更新主观题
	Records     []*mark_rpc_model.MarkPaperRecord // 扫描记录
}

func GetRecordType(isObjective bool) mark_paper_enum.RecordType {

	if isObjective {
		return mark_paper_enum.RecordType_RecordTypeObjective
	}
	return mark_paper_enum.RecordType_RecordTypeSubjective
}
