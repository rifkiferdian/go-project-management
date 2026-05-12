package controllers

import "time"

const roadmapTimezone = "Asia/Jakarta"

func roadmapNow() time.Time {
	location, err := time.LoadLocation(roadmapTimezone)
	if err != nil {
		return time.Now()
	}
	return time.Now().In(location)
}
