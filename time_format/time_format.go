package time_format

import (
	"time"
)

type TimeFormat string

const (
	TimeFormatHyphen           = "2006-01-02 15:04:05" //YYYY-MM-dd HH:mm:ss
	TimeFormatHyphenWithOutSec = "2006-01-02 15:04"    //YYYY-MM-dd HH:mm
	TimeFormatSlash            = "2006/01/02 15:04:05" //YYYY/MM/dd HH:mm:ss
	TimeFormatSlashWithOutSec  = "2006/01/02 15:04"    //YYYY/MM/dd HH:mm
)

func TimeFormatToStr(t time.Time, format TimeFormat) string {

	switch format {
	case TimeFormatHyphen:
		return t.Format(TimeFormatHyphen)
	case TimeFormatHyphenWithOutSec:
		return t.Format(TimeFormatHyphenWithOutSec)
	case TimeFormatSlash:
		return t.Format(TimeFormatSlash)
	case TimeFormatSlashWithOutSec:
		return t.Format(TimeFormatSlashWithOutSec)
	}
	return t.Format(TimeFormatHyphen)
}
