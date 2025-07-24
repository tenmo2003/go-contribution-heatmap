package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

var CELL_WIDTH = 2

var EMPTY_CELL = strings.Repeat(" ", CELL_WIDTH)

var WEEKDAY_LENGTH = 3
var MONTH_LENGTH = 3

var COLORS = map[string]string{
	"white/gray":    "\033[48;5;255m" + EMPTY_CELL + "\033[0m",
	"light green":   "\033[48;5;120m" + EMPTY_CELL + "\033[0m",
	"medium green":  "\033[48;5;34m" + EMPTY_CELL + "\033[0m",
	"dark green":    "\033[48;5;28m" + EMPTY_CELL + "\033[0m",
	"darkest green": "\033[48;5;22m" + EMPTY_CELL + "\033[0m",
}

var DEFAULT_RELATIVE_TIME = "1 year"

type FirstDay int

const (
	Sunday FirstDay = iota
	Monday
)

var weekdays = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

func generateWeekdayMaps(firstDay FirstDay) (map[string]int, map[int]string) {
	WEEKDAY_INDEX := make(map[string]int)
	INDEX_WEEKDAY := make(map[int]string)

	start := int(firstDay)
	for i := range 7 {
		day := weekdays[(start+i)%7]
		WEEKDAY_INDEX[day] = i
		INDEX_WEEKDAY[i] = day
	}

	return WEEKDAY_INDEX, INDEX_WEEKDAY
}

type Cell struct {
	CommitCount int
	Display     string
	Date        time.Time
}

func (c *Cell) String() string {
	if c.CommitCount < 0 {
		return EMPTY_CELL
	}
	return fmt.Sprintf("%3s %3d %s", c.Date.Format("2006-01-02"), c.CommitCount, c.Display)
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getGitDir(path string) (string, error) {
	gitDir := filepath.Join(path, ".git")
	exists, err := pathExists(gitDir)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("not a git repository")
	}
	return gitDir, nil
}

func expandPath(path string) string {
	if path == "." {
		return os.Getenv("PWD")
	}
	if strings.HasPrefix(path, "~") {
		return strings.Replace(path, "~", os.Getenv("HOME"), 1)
	}
	return path
}

func realignCells(cells [][]Cell) bool {
	changed := false
	for i := range len(cells) - 1 {
		if len(cells[i+1]) > len(cells[i]) {
			for j := range i + 1 {
				cells[j] = slices.Insert(cells[j], 0, Cell{CommitCount: -1, Display: EMPTY_CELL})
				changed = true
			}
		}
	}

	return changed
}

func displayMonths(cells [][]Cell, monthsIndex map[int]string) {
	monthDisplayCells := []string{}
	monthDisplayCells = append(monthDisplayCells, strings.Repeat(" ", WEEKDAY_LENGTH))

	for i := 0; i < len(cells[0]); i++ {
		if _, ok := monthsIndex[i]; ok {
			monthDisplayCells = append(monthDisplayCells, monthsIndex[i][:MONTH_LENGTH])
		} else {
			monthDisplayCells = append(monthDisplayCells, EMPTY_CELL)
		}
	}

	borrowed := 0
	for i := 0; i < len(monthDisplayCells); i++ {
		cellValue := monthDisplayCells[i]

		if cellValue == EMPTY_CELL {
			fmt.Print(" " + strings.Repeat(" ", CELL_WIDTH-borrowed))
			borrowed = 0
		}

		if cellValue != EMPTY_CELL {
			if cellValue != strings.Repeat(" ", WEEKDAY_LENGTH) {
				borrowed += MONTH_LENGTH - CELL_WIDTH
			}
			fmt.Print(cellValue + " ")
		}
	}

	fmt.Println()
}

func displayCells(INDEX_WEEKDAY map[int]string, cells [][]Cell) {
	displayCells := [][]string{}
	for range cells {
		displayCells = append(displayCells, []string{})
	}
	for i := range cells {
		for j := 0; j < len(cells[i]); j++ {
			displayCells[i] = append(displayCells[i], cells[i][j].Display)
		}
	}

	for i := 0; i < len(INDEX_WEEKDAY); i++ {
		displayCells[i] = slices.Insert(displayCells[i], 0, INDEX_WEEKDAY[i][:WEEKDAY_LENGTH])
	}

	for i := 0; i < len(displayCells); i++ {
		fmt.Println(strings.Join(displayCells[i], " "))
		fmt.Println()
	}
}

func parseRelativeTime(relativeTime string) (from time.Time, to time.Time, err error) {
	relativeTimeParts := regexp.MustCompile(`(?i)(\d+)\s*(day|week|month|year)s?`).FindStringSubmatch(relativeTime)
	if len(relativeTimeParts) != 3 {
		return time.Time{}, time.Time{}, fmt.Errorf("Invalid relative time")
	}

	value, err := strconv.Atoi(relativeTimeParts[1])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	switch strings.ToLower(relativeTimeParts[2]) {
	case "day":
		return time.Now().AddDate(0, 0, -1*value).AddDate(0, 0, 1), time.Now(), nil
	case "week":
		return time.Now().AddDate(0, 0, -7*value).AddDate(0, 0, 1), time.Now(), nil
	case "month":
		return time.Now().AddDate(0, -1*value, 0).AddDate(0, 0, 1), time.Now(), nil
	case "year":
		return time.Now().AddDate(-1*value, 0, 0).AddDate(0, 0, 1), time.Now(), nil
	}

	return time.Time{}, time.Time{}, fmt.Errorf("Invalid relative time unit")
}

func parseYear(year int) (from time.Time, to time.Time, err error) {
	if year < 0 {
		return time.Time{}, time.Time{}, fmt.Errorf("Year must be positive")
	}
	from = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	to = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
	return from, to, nil
}

func main() {
	dirFlag := flag.String("dir", "", "Path to the parent directory of git repositories to scan (depth = 1)")
	authorFlag := flag.String("author-email", "", "Author's email")
	firstWeekdayConfigFlag := flag.String("first-weekday", "sunday", "First day of the week")
	yearFlag := flag.Int("year", -1, "Year to scan")
	relativeTimeFlag := flag.String("relative-time", "", "Relative time to scan")

	flag.Parse()

	if *dirFlag == "" {
		fmt.Println("No parent directory provided")
		return
	}

	if *authorFlag == "" {
		// get the author's email from git config
		cmd := exec.Command("git", "config", "user.email")
		out, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
			return
		}
		*authorFlag = strings.TrimSpace(string(out))
		fmt.Println("Author email not provided, deriving from git config:", *authorFlag)
	}

	if *authorFlag == "" {
		fmt.Println("No author email provided")
		return
	}

	if *relativeTimeFlag != "" && *yearFlag != -1 {
		fmt.Println("Relative time and year are both set, please set only one")
		return
	}

	if *relativeTimeFlag == "" && *yearFlag == -1 {
		*relativeTimeFlag = DEFAULT_RELATIVE_TIME
	}

	dayToFirstDay := map[string]FirstDay{
		"sunday": Sunday,
		"monday": Monday,
	}

	if _, ok := dayToFirstDay[*firstWeekdayConfigFlag]; !ok {
		fmt.Println("Invalid first weekday:", *firstWeekdayConfigFlag)
		return
	}

	firstDay := dayToFirstDay[*firstWeekdayConfigFlag]

	WEEKDAY_INDEX, INDEX_WEEKDAY := generateWeekdayMaps(firstDay)

	root := expandPath(*dirFlag)

	repos := []string{}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if path == root {
			return nil
		}
		if err != nil {
			fmt.Println(err)
			return err
		}
		if d.IsDir() {
			gitDir, err := getGitDir(path)
			if err != nil {
				if err.Error() == "not a git repository" {
					return filepath.SkipDir
				}
				return err
			}
			repos = append(repos, gitDir)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	commitCountByDate := map[string]int{}

	var from time.Time
	var to time.Time

	if *relativeTimeFlag != "" {
		from, to, err = parseRelativeTime(*relativeTimeFlag)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		from, to, err = parseYear(*yearFlag)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	formattedFrom := from.Format("2006-01-02")
	formattedTo := to.Format("2006-01-02")

	fmt.Println("From", formattedFrom, "To", formattedTo)

	for _, repo := range repos {
		cmd := exec.Command(
			"git",
			"--git-dir="+repo,
			"log",
			"--pretty=format:%ae %ad %s",
			"--reverse",
			"--date=format:%Y-%m-%d",
			"--author="+*authorFlag,
			"--since="+formattedFrom,
			"--until="+formattedTo,
		)

		out, err := cmd.Output()
		if err != nil {
			continue
		}

		outputStr := strings.TrimSpace(string(out))
		if outputStr == "" {
			continue
		}
		lines := strings.SplitSeq(outputStr, "\n")

		for line := range lines {
			parts := strings.Split(line, " ")
			date := parts[1]

			commitCountByDate[date]++
		}
	}

	cells := [][]Cell{}

	for i := 0; i < len(WEEKDAY_INDEX); i++ {
		cells = append(cells, []Cell{})
	}

	monthsIndex := map[int]string{}
	months := []string{}
	for d := from; d.After(to) == false; d = d.AddDate(0, 0, 1) {
		commitCount := commitCountByDate[d.Format("2006-01-02")]
		var colorName string
		if commitCount == 0 {
			colorName = "white/gray"
		} else if commitCount < 5 {
			colorName = "light green"
		} else if commitCount < 10 {
			colorName = "medium green"
		} else if commitCount < 20 {
			colorName = "dark green"
		} else {
			colorName = "darkest green"
		}

		cell := Cell{
			CommitCount: commitCount,
			Display:     COLORS[colorName],
			Date:        d,
		}

		cells[WEEKDAY_INDEX[d.Weekday().String()]] = append(cells[WEEKDAY_INDEX[d.Weekday().String()]], cell)

		if (d.Equal(from) || WEEKDAY_INDEX[d.Weekday().String()] == 0) && (len(months) == 0 || months[len(months)-1] != d.Month().String()) {
			months = append(months, d.Month().String())
			monthsIndex[len(cells[WEEKDAY_INDEX[d.Weekday().String()]])-1] = d.Month().String()
		}

	}

	changed := realignCells(cells)

	if changed {
		modifiedMonthsIndex := map[int]string{}
		for index, month := range monthsIndex {
			if index != 0 {
				modifiedMonthsIndex[index+1] = month
			} else {
				modifiedMonthsIndex[index] = month
			}
		}
		monthsIndex = modifiedMonthsIndex
	}

	displayMonths(cells, monthsIndex)

	displayCells(INDEX_WEEKDAY, cells)
}
