package main

import (
	"contribution-heatmap/internal/calendar"
	"contribution-heatmap/internal/display"
	"contribution-heatmap/internal/flags"
	"contribution-heatmap/internal/git"
	internalTime "contribution-heatmap/internal/time"
	"contribution-heatmap/internal/utils"
	"fmt"
	"time"
)

func main() {
	config, err := flags.ParseFlags()
	if err != nil {
		fmt.Println(err)
		return
	}

	WEEKDAY_INDEX, INDEX_WEEKDAY := calendar.GenerateWeekdayMaps(config.FirstWeekday)

	roots := utils.ExpandPaths(config.Dirs)

	repos, err := git.GetReposGitDirs(roots)
	if err != nil {
		fmt.Println(err)
		return
	}

	var from time.Time
	var to time.Time

	if config.RelativeTime != "" {
		from, to, err = internalTime.ParseRelativeTime(config.RelativeTime)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		from, to, err = internalTime.ParseYear(config.Year)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	formattedFrom := from.Format("2006-01-02")
	formattedTo := to.Format("2006-01-02")

	fmt.Println("From", formattedFrom, "To", formattedTo)

	commitCountByDate := git.GetCommitCountByDate(config.AuthorEmail, formattedFrom, formattedTo, repos)

	display.DisplayHeatmap(WEEKDAY_INDEX, INDEX_WEEKDAY, commitCountByDate, from, to)
}
