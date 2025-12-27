package excel_tool

import (
	"encoding/json"
	"testing"
)

func TestConvToExcelData(t *testing.T) {
	type args[T any] struct {
		dataList []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want [][]interface{}
	}
	tests := []testCase[*TestStruct]{
		{
			name: "test1",
			args: args[*TestStruct]{
				dataList: []*TestStruct{
					{
						Name: "Alice",
						Age:  25,
						City: "New York",
						Sex:  1,
					},
					{
						Name: "Bob",
						Age:  30,
						City: "Los Angeles",
						Sex:  0,
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvToExcelData(tt.args.dataList)
			gotStr, _ := json.Marshal(got)
			t.Log(string(gotStr))
		})
	}
}

// TestStruct 注意！sequenceExcel:"1"为自定义扩展信息，用来表示当前字段在Excel中的列的顺序
type TestStruct struct {
	Name string `sequenceExcel:"1"`
	Age  int    `sequenceExcel:"3"`
	City string `sequenceExcel:"2"`
	Sex  int32  `sequenceExcel:"4"`
}
