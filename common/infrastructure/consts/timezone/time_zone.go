package timezone

import (
	"git-new.okkia.site/crk/decimal-cricket-common/application/utils"
	"strconv"
	"time"
)

var Locations = gen25Locations()

func gen25Locations() []*time.Location {
	ret := make([]*time.Location, 0, 25)
	for i := -12; i <= 12; i++ {
		utcStr := buildUtcStr(i)
		ret = append(ret, utils.GetTimeLocationByTimeZone(utcStr))
	}
	return ret
}

func buildUtcStr(offsetHour int) string {
	var utcStr string
	if offsetHour >= 0 {
		utcStr = utils.UTC + "+" + strconv.Itoa(offsetHour)
	} else {
		utcStr = utils.UTC + strconv.Itoa(offsetHour)
	}
	return utcStr
}
