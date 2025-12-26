package rpc_data

import "code.chenji.com/pkg/idlgen/kitex_gen/edu/meta/meta_model"

type LabelValue struct {
	Id        int64  // 标签值ID
	ParentId  int64  // 父节点
	RootId    int64  // 根节点
	LabelId   int64  // 所属的标签ID
	EnumValue int32  // 枚举值
	Name      string // 名称
}

func ConvLabelValue(labelValue *meta_model.LabelValue) *LabelValue {
	if labelValue == nil {
		return nil
	}
	return &LabelValue{
		Id:        labelValue.Id,
		ParentId:  labelValue.ParentId,
		RootId:    labelValue.RootId,
		LabelId:   labelValue.LabelId,
		EnumValue: labelValue.EnumValue,
		Name:      labelValue.Name,
	}
}
