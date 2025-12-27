package excel_template

import (
	"bytes"
	"fmt"
	"time"

	"code.chenji.com/cj/scan_paper/excel_tool"
	"code.chenji.com/cj/scan_paper/utils"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

const (
	ScanResultFilePath      = "scan_paper/excel/"
	AbsentSheetName         = "缺考学生名单"
	AbsentFileNamePrefix    = "缺考学生名单"
	MissPaperSheetName      = "缺少试卷学生名单"
	MissPaperFileNamePrefix = "缺少试卷学生名单"
)

type ScanResultTemplate struct {
}

type ScanResultDataRow struct {
	No          string `sequenceExcel:"1"` // 序号
	SchoolName  string `sequenceExcel:"2"` // 学校名称
	CampusName  string `sequenceExcel:"3"` // 校区名称
	GradeName   string `sequenceExcel:"4"` // 年级名称
	ClassName   string `sequenceExcel:"5"` // 班级名称
	ExamSubject string `sequenceExcel:"6"` // 考试科目
	StudentId   string `sequenceExcel:"7"` // 学生学号
	StudentName string `sequenceExcel:"8"` // 学生姓名
	ExamCode    string `sequenceExcel:"9"` // 学生考号
}

type ScanResultSheetData struct {
	SheetName string
	Rows      []*ScanResultDataRow
}

func NewScanResultTemplate() *ScanResultTemplate {
	return &ScanResultTemplate{}
}

func (t *ScanResultTemplate) getTitleRow() []*ScanResultDataRow {

	titles := make([]*ScanResultDataRow, 0, 2)
	titles = append(titles, &ScanResultDataRow{
		No: fmt.Sprintf("%s%s", "数据生成时间:", utils.TimeFormatToStr(time.Now(), utils.TimeFormatHyphen)),
	})
	titles = append(titles, &ScanResultDataRow{
		No:          "序号",
		SchoolName:  "学校",
		CampusName:  "校区",
		GradeName:   "年级",
		ClassName:   "班级",
		ExamSubject: "学科",
		StudentId:   "学生ID",
		StudentName: "学生姓名",
		ExamCode:    "考号",
	})
	return titles
}

func (t *ScanResultTemplate) CreateScanResultExcel(sheetName string, data []*ScanResultDataRow) ([]byte, error) {

	// 创建excel
	contentData := excel_tool.ConvToExcelData(data)
	titleData := excel_tool.ConvToExcelData(t.getTitleRow())
	excelData := append(titleData, contentData...)
	excelByte, err := excel_tool.WriteExcelByte(sheetName, excelData)
	if err != nil {
		return nil, err
	}
	return excelByte, nil
}

func (t *ScanResultTemplate) CreateScanResultExcelMultiSheet(sheetDataList []*ScanResultSheetData) ([]byte, error) {

	sheetCells := make([]*SheetCellData, 0, len(sheetDataList))
	// 创建excel
	for _, item := range sheetDataList {
		if item == nil {
			continue
		}
		contentData := excel_tool.ConvToExcelData(item.Rows)
		titleData := excel_tool.ConvToExcelData(t.getTitleRow())
		excelData := append(titleData, contentData...)
		sheetCells = append(sheetCells, &SheetCellData{
			SheetName: item.SheetName,
			Cells:     excelData,
		})
	}
	excelByte, err := WriteExcelByteMultiSheet(sheetCells)
	if err != nil {
		return nil, err
	}
	return excelByte, nil
}

type ScanItemData struct {
	StudentId       int64  `json:"student_id"`
	StudentName     string `json:"student_name"`
	SchoolId        int64  `json:"school_id"`
	SchoolName      string `json:"school_name"`
	CampusId        int64  `json:"campus_id"`
	CampusName      string `json:"campus_name"`
	BaseGrade       int32  `json:"base_grade"`
	GradeName       string `json:"grade_name"`
	ClassId         int64  `json:"class_id"`
	ClassName       string `json:"class_name"`
	ExamSubject     int32  `json:"exam_subject"`
	ExamSubjectName string `json:"exam_subject_name"`
	ExamCode        string `json:"exam_code"`
}

func (t *ScanResultTemplate) GetAbsentFileName() string {

	return fmt.Sprintf("%s%s_%s.xlsx", ScanResultFilePath, AbsentFileNamePrefix, utils.TimeFormatToStr(time.Now(), utils.TimeFormatHyphen))
}

func (t *ScanResultTemplate) GetMissPaperFileName() string {
	return fmt.Sprintf("%s%s_%s.xlsx", ScanResultFilePath, MissPaperFileNamePrefix, utils.TimeFormatToStr(time.Now(), utils.TimeFormatHyphen))
}

type SheetCellData struct {
	SheetName string
	Cells     [][]interface{}
}

// WriteExcelByteMultiSheet 写入excel到字节流
// sheetName: 工作表名
// data: 要写入的数据
// 返回字节流
func WriteExcelByteMultiSheet(sheetCells []*SheetCellData) ([]byte, error) {

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	activeSheetIndex := 0
	for _, item := range sheetCells {
		if item == nil {
			continue
		}
		// Create a new sheet.
		index, err := f.NewSheet(item.SheetName)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if activeSheetIndex <= 0 {
			activeSheetIndex = index
		}
		for idx, row := range item.Cells {
			cell, err := excelize.CoordinatesToCellName(1, idx+1)
			if err != nil {
				fmt.Println(err)
				break
			}
			err = f.SetSheetRow(item.SheetName, cell, &row)
			if err != nil {
				return nil, err
			}

		}
	}
	// 获取默认工作表的名称
	defaultSheetName := f.GetSheetName(0)
	// 删除默认工作表
	if err := f.DeleteSheet(defaultSheetName); err != nil {
		return nil, errors.WithStack(err)
	}
	if activeSheetIndex > 0 {
		activeSheetIndex = 0
	}
	f.SetActiveSheet(activeSheetIndex)

	var buf bytes.Buffer
	err := f.Write(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
