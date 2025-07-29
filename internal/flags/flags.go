package flags

import (
	"contribution-heatmap/internal/calendar"
	"contribution-heatmap/internal/git"
	"flag"
	"fmt"
)

const DEFAULT_RELATIVE_TIME = "1 year"

type Config struct {
	AuthorEmail  string
	Year         int
	RelativeTime string
	FirstWeekday calendar.FirstWeekday
	Dirs         []string
}

func ParseFlags() (*Config, error) {
	authorFlag := flag.String("author-email", "", "Author's email")
	yearFlag := flag.Int("year", -1, "Year to scan")
	relativeTimeFlag := flag.String("relative-time", "", "Relative time to scan")
	var firstWeekday calendar.FirstWeekday
	flag.Var(&firstWeekday, "first-weekday", "First day of the week (sunday or monday)")

	flag.Parse()
	dirs := flag.Args()

	if len(dirs) == 0 {
		return nil, fmt.Errorf("No parent directory provided through arguments")
	}

	if *authorFlag == "" {
		email, err := git.GetUserEmailFromGitConfig()
		if err != nil {
			return nil, fmt.Errorf("No author email provided through --author-email flag or git config: %w", err)
		}
		*authorFlag = email
		fmt.Println("Author email not provided, deriving from git config:", *authorFlag)
	}

	if *authorFlag == "" {
		return nil, fmt.Errorf("No author email provided from --author-email flag or git config")
	}

	if *relativeTimeFlag != "" && *yearFlag != -1 {
		return nil, fmt.Errorf("Relative time and year are both set, please set only one")
	}

	if *relativeTimeFlag == "" && *yearFlag == -1 {
		*relativeTimeFlag = DEFAULT_RELATIVE_TIME
	}

	return &Config{
		AuthorEmail:  *authorFlag,
		Year:         *yearFlag,
		RelativeTime: *relativeTimeFlag,
		FirstWeekday: firstWeekday,
		Dirs:         dirs,
	}, nil
}
