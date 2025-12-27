package excel_template

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/mhh1987/hetu_tool/json_tool"
)

func TestWriteExcelByteMultiSheet(t *testing.T) {

	sheetCells := make([]*ScanResultSheetData, 0, 10)
	sheetCells = append(sheetCells, &ScanResultSheetData{
		SheetName: "11111",
		Rows: []*ScanResultDataRow{
			{
				No:          "1",
				SchoolName:  "11111",
				CampusName:  "a1111",
				GradeName:   "b1111",
				ClassName:   "c1111",
				ExamSubject: "语文",
				StudentId:   "1111111111",
				StudentName: "a1111111111",
				ExamCode:    "1111",
			},
			{
				No:          "2",
				SchoolName:  "11111",
				CampusName:  "a1111",
				GradeName:   "b1111",
				ClassName:   "c1111",
				ExamSubject: "语文",
				StudentId:   "2222222222",
				StudentName: "a2222222222",
				ExamCode:    "2222",
			},
			{
				No:          "3",
				SchoolName:  "11111",
				CampusName:  "a1111",
				GradeName:   "b1111",
				ClassName:   "c1111",
				ExamSubject: "语文",
				StudentId:   "3333333333",
				StudentName: "a3333333333",
				ExamCode:    "3333",
			},
		},
	})

	sheetCells = append(sheetCells, &ScanResultSheetData{
		SheetName: "22222",
		Rows: []*ScanResultDataRow{
			{
				No:          "1",
				SchoolName:  "22222",
				CampusName:  "a2222",
				GradeName:   "b2222",
				ClassName:   "c2222",
				ExamSubject: "语文",
				StudentId:   "4444444444",
				StudentName: "a4444444444",
				ExamCode:    "4444",
			},
			{
				No:          "2",
				SchoolName:  "22222",
				CampusName:  "a2222",
				GradeName:   "b2222",
				ClassName:   "c2222",
				ExamSubject: "语文",
				StudentId:   "5555555555",
				StudentName: "a5555555555",
				ExamCode:    "5555",
			},
			{
				No:          "3",
				SchoolName:  "22222",
				CampusName:  "a2222",
				GradeName:   "b2222",
				ClassName:   "c2222",
				ExamSubject: "语文",
				StudentId:   "6666666666",
				StudentName: "a6666666666",
				ExamCode:    "6666",
			},
		},
	})
	resultTemplate := NewScanResultTemplate()
	excelBytes, err := resultTemplate.CreateScanResultExcelMultiSheet(sheetCells)
	if err != nil {
		t.Error(err)
		return
	}
	outDir := "下载文件"
	outFileName := fmt.Sprintf("%s_%d.xlsx", outDir, time.Now().Unix())
	// 保存图片到文件
	outFile, err := os.Create(outFileName)
	if err != nil {
		t.Error(err)
		return

	}
	_, err = io.Copy(outFile, bytes.NewReader(excelBytes))
	if err != nil {
		t.Error(err)
		return
	}
	err = outFile.Close()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("保存成功")
}

func TestScanResultDataRowGroupBySchool(t *testing.T) {

	rows := []*ScanResultDataRow{
		{
			No:          "1",
			SchoolName:  "桂伟的数据同步单校区测试学校",
			CampusName:  "桂伟的数据同步单校区测试学校-默认校区",
			GradeName:   "高一",
			ClassName:   "1班x",
			ExamSubject: "语文",
			StudentId:   "1373242626001651804",
			StudentName: "学生34",
			//ExamCode:    "3333",
		},
		{
			No:          "2",
			SchoolName:  "桂伟的数据同步单校区测试学校",
			CampusName:  "桂伟的数据同步单校区测试学校-默认校区",
			GradeName:   "高一",
			ClassName:   "1班x",
			ExamSubject: "语文",
			StudentId:   "1373242627754870823",
			StudentName: "学生5",
			//ExamCode:    "2222",
		},
		{
			No:          "3",
			SchoolName:  "桂伟的数据同步单校区测试学校",
			CampusName:  "桂伟的数据同步单校区测试学校-默认校区",
			GradeName:   "高一",
			ClassName:   "1班x",
			ExamSubject: "语文",
			StudentId:   "1373242629571004451",
			StudentName: "学生36",
			//ExamCode:    "3333",
		},
		//{
		//	No:          "1",
		//	SchoolName:  "22222",
		//	CampusName:  "a2222",
		//	GradeName:   "b2222",
		//	ClassName:   "c2222",
		//	ExamSubject: "语文",
		//	StudentId:   "4444444444",
		//	StudentName: "a4444444444",
		//	//ExamCode:    "4444",
		//},
		//{
		//	No:          "2",
		//	SchoolName:  "22222",
		//	CampusName:  "a2222",
		//	GradeName:   "b2222",
		//	ClassName:   "c2222",
		//	ExamSubject: "语文",
		//	StudentId:   "5555555555",
		//	StudentName: "a5555555555",
		//	//ExamCode:    "5555",
		//},
		//{
		//	No:          "3",
		//	SchoolName:  "22222",
		//	CampusName:  "a2222",
		//	GradeName:   "b2222",
		//	ClassName:   "c2222",
		//	ExamSubject: "语文",
		//	StudentId:   "6666666666",
		//	StudentName: "a6666666666",
		//	//ExamCode:    "6666",
		//},
	}
	resultData := ScanResultDataRowGroupBySchool(rows)
	t.Log(fmt.Sprintf("resultData:%s", json_tool.ToJson(resultData)))
}

func ScanResultDataRowGroupBySchool(itemList []*ScanResultDataRow) []*ScanResultSheetData {

	result := make([]*ScanResultSheetData, 0, 10)
	currentSchool := ""
	tempList := make([]*ScanResultDataRow, 0, 100)
	for index, item := range itemList {
		if item == nil {
			continue
		}
		if index == 0 {
			currentSchool = item.SchoolName
		}
		if currentSchool != item.SchoolName {
			if len(tempList) > 0 {
				result = append(result, &ScanResultSheetData{
					SheetName: currentSchool,
					Rows:      tempList,
				})
			}
			currentSchool = item.SchoolName
			tempList = make([]*ScanResultDataRow, 0, 100)
		}
		tempList = append(tempList, item)
		if index == len(itemList)-1 { // 最后一个元素
			result = append(result, &ScanResultSheetData{
				SheetName: currentSchool,
				Rows:      tempList,
			})
		}
	}
	return result
}
