package main

import (
	"contribution-heatmap/internal/calendar"
	"contribution-heatmap/internal/display"
	"contribution-heatmap/internal/git"
	internalTime "contribution-heatmap/internal/time"
	"contribution-heatmap/internal/utils"
	"flag"
	"fmt"
	"time"
)

var DEFAULT_RELATIVE_TIME = "1 year"

func main() {
	authorFlag := flag.String("author-email", "", "Author's email")
	yearFlag := flag.Int("year", -1, "Year to scan")
	relativeTimeFlag := flag.String("relative-time", "", "Relative time to scan")
	var firstWeekday calendar.FirstWeekday
	flag.Var(&firstWeekday, "first-weekday", "First day of the week (sunday or monday)")

	flag.Parse()
	dirs := flag.Args()

	if len(dirs) == 0 {
		fmt.Println("No parent directory provided through arguments")
		return
	}

	if *authorFlag == "" {
		email, err := git.GetUserEmailFromGitConfig()
		if err != nil {
			fmt.Println("Error getting author email:", err)
			return
		}
		*authorFlag = email
		fmt.Println("Author email not provided, deriving from git config:", *authorFlag)
	}

	if *authorFlag == "" {
		fmt.Println("No author email provided from --author-email flag or git config")
		return
	}

	if *relativeTimeFlag != "" && *yearFlag != -1 {
		fmt.Println("Relative time and year are both set, please set only one")
		return
	}

	if *relativeTimeFlag == "" && *yearFlag == -1 {
		*relativeTimeFlag = DEFAULT_RELATIVE_TIME
	}

	WEEKDAY_INDEX, INDEX_WEEKDAY := calendar.GenerateWeekdayMaps(firstWeekday)

	root := utils.ExpandPath(dirs[0])

	repos, err := git.GetReposGitDirs(root)
	if err != nil {
		fmt.Println(err)
		return
	}

	var from time.Time
	var to time.Time

	if *relativeTimeFlag != "" {
		from, to, err = internalTime.ParseRelativeTime(*relativeTimeFlag)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		from, to, err = internalTime.ParseYear(*yearFlag)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	formattedFrom := from.Format("2006-01-02")
	formattedTo := to.Format("2006-01-02")

	fmt.Println("From", formattedFrom, "To", formattedTo)

	commitCountByDate := git.GetCommitCountByDate(*authorFlag, formattedFrom, formattedTo, repos)

	display.DisplayHeatmap(WEEKDAY_INDEX, INDEX_WEEKDAY, commitCountByDate, from, to)
}
