package test

import (
	"sort"
	"time_app/app/repository/model"
)

// Equal filter for interval repo pipeline
func filterIntervalsEqualPipeline(rawIntervalList []model.Interval) []model.Interval {
	intervalMap := make(map[string][]model.Interval)
	for _, interval := range rawIntervalList {
		if interval.EndAt != nil {
			groupKey := interval.UserUUID + "_" + interval.CategoryUUID
			intervalMap[groupKey] = append(intervalMap[groupKey], interval)
		}
	}

	var result []model.Interval
	for _, intervals := range intervalMap {
		sort.Slice(intervals, func(i, j int) bool {
			return intervals[i].StartedAt < intervals[j].StartedAt
		})

		result = append(result, intervals[:len(intervals)-1]...)
	}
	return result
}
