package rpc_data

type QrCodeType int32

const (
	QrCodeType_Unknown QrCodeType = 0 // 未知类型
	QrCodeType_Student QrCodeType = 1 // 学生码
	QrCodeType_Sheet   QrCodeType = 2 // 试卷码
)
