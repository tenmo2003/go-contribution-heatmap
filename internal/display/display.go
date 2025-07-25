package display

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

var WEEKDAY_LENGTH = 3
var MONTH_LENGTH = 3

var CELL_WIDTH = 2

var EMPTY_CELL = strings.Repeat(" ", CELL_WIDTH)

var COLORS = map[string]string{
	"white/gray":    "\033[48;5;255m" + EMPTY_CELL + "\033[0m",
	"light green":   "\033[48;5;120m" + EMPTY_CELL + "\033[0m",
	"medium green":  "\033[48;5;34m" + EMPTY_CELL + "\033[0m",
	"dark green":    "\033[48;5;28m" + EMPTY_CELL + "\033[0m",
	"darkest green": "\033[48;5;22m" + EMPTY_CELL + "\033[0m",
}

type Cell struct {
	CommitCount int
	Display     string
	Date        time.Time
}

func DisplayHeatmap(WEEKDAY_INDEX map[string]int, INDEX_WEEKDAY map[int]string, commitCountByDate map[string]int, from time.Time, to time.Time) {
	cells, monthsIndex := generateHeatmap(WEEKDAY_INDEX, commitCountByDate, from, to)

	changed := realignCells(cells)

	if changed {
		shiftMonthCells(&monthsIndex)
	}

	displayMonths(cells, monthsIndex)

	displayCells(INDEX_WEEKDAY, cells)
}

func generateHeatmap(WEEKDAY_INDEX map[string]int, commitCountByDate map[string]int, from time.Time, to time.Time) ([][]Cell, map[int]string) {
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
	return cells, monthsIndex
}

func realignCells(cells [][]Cell) bool {
	changed := false
	for i := range len(cells) - 1 {
		if len(cells[i]) == 0 {
			cells[i] = append(cells[i], Cell{CommitCount: -1, Display: EMPTY_CELL})
			continue
		}
		if len(cells[i+1]) == 0 {
			cells[i+1] = append(cells[i+1], Cell{CommitCount: -1, Display: EMPTY_CELL})
			continue
		}
		former := cells[i][0]
		latter := cells[i+1][0]
		if former.Date.After(latter.Date) {
			for j := range i + 1 {
				cells[j] = slices.Insert(cells[j], 0, Cell{CommitCount: -1, Display: EMPTY_CELL})
				changed = true
			}
		}
	}

	return changed
}

func shiftMonthCells(monthsIndex *map[int]string) {
	modifiedMonthsIndex := map[int]string{}
	for index, month := range *monthsIndex {
		if index != 0 {
			modifiedMonthsIndex[index+1] = month
		} else {
			modifiedMonthsIndex[index] = month
		}
	}
	*monthsIndex = modifiedMonthsIndex
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
