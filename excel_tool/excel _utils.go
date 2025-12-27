package excel_tool

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"path"
	"reflect"
	"sort"
	"strconv"
)

// ConvToExcelData 转换数据到excel数据格式
func ConvToExcelData[T any](dataList []T) [][]interface{} {

	var excelData [][]interface{}
	for _, data := range dataList {
		row := dataToRow(data)
		excelData = append(excelData, row)
	}
	return excelData
}

type columnInfo struct {
	key      string
	value    interface{}
	sequence int
}

var sequenceTag = "sequenceExcel"

func dataToRow(data interface{}) []interface{} {

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		fmt.Println("输入不是结构体")
		return nil
	}
	t := v.Type()

	var columns []*columnInfo
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		temp := &columnInfo{
			key:      fieldType.Name,
			value:    fieldValue.Interface(),
			sequence: 0,
		}
		sequenceStr := fieldType.Tag.Get(sequenceTag)
		if sequenceStr != "" {
			temp.sequence, _ = strconv.Atoi(sequenceStr)
		}
		columns = append(columns, temp)
	}
	sort.Slice(columns, func(i, j int) bool {
		return columns[i].sequence < columns[j].sequence
	})
	var row []interface{}
	for _, column := range columns {
		row = append(row, column.value)
	}
	return row
}

// WriteExcelLocal 写入excel到本地
// outPath: 输出路径
// outFileName: 输出文件名
// sheetName: 工作表名
// data: 要写入的数据
func WriteExcelLocal(outPath string, outFileName string, sheetName string, data [][]interface{}) error {

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create a new sheet.
	index, err := f.NewSheet(sheetName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for idx, row := range data {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			fmt.Println(err)
			break
		}
		err = f.SetSheetRow(sheetName, cell, &row)
		if err != nil {
			return err
		}

	}
	f.SetActiveSheet(index)

	// Save spreadsheet by the given path.
	if err := f.SaveAs(path.Join(outPath, outFileName)); err != nil {
		fmt.Println(err)
	}

	return nil
}

// WriteExcelByte 写入excel到字节流
// sheetName: 工作表名
// data: 要写入的数据
// 返回字节流
func WriteExcelByte(sheetName string, data [][]interface{}) ([]byte, error) {

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create a new sheet.
	index, err := f.NewSheet(sheetName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for idx, row := range data {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			fmt.Println(err)
			break
		}
		err = f.SetSheetRow(sheetName, cell, &row)
		if err != nil {
			return nil, err
		}

	}
	// 获取默认工作表的名称
	defaultSheetName := f.GetSheetName(0)
	// 删除默认工作表
	if err = f.DeleteSheet(defaultSheetName); err != nil {
		return nil, errors.WithStack(err)
	}
	f.SetActiveSheet(index)

	var buf bytes.Buffer
	err = f.Write(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type SheetData struct {
	SheetName string
	Data      [][]interface{}
}

// WriteSheetsToExcelByte 写入excel到字节流
// sheets: 工作表数据
// 返回字节流
func WriteSheetsToExcelByte(sheets []*SheetData) ([]byte, error) {

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for _, sheet := range sheets {
		// Create a new sheet.
		index, err := f.NewSheet(sheet.SheetName)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		for idx, row := range sheet.Data {
			cell, err := excelize.CoordinatesToCellName(1, idx+1)
			if err != nil {
				fmt.Println(err)
				break
			}
			err = f.SetSheetRow(sheet.SheetName, cell, &row)
			if err != nil {
				return nil, err
			}

		}
		f.SetActiveSheet(index)
	}

	var buf bytes.Buffer
	err := f.Write(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
