package services

import (
	"errors"
	"gobase-app/models"
	"gobase-app/repositories"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ManagementService struct {
	Repo *repositories.ManagementRepository
}

func (s *ManagementService) GetTickets() ([]models.TicketListItem, error) {
	return s.Repo.GetTickets()
}

func (s *ManagementService) GetBoardColumns() ([]models.BoardColumn, error) {
	return s.Repo.GetBoardColumns()
}

func (s *ManagementService) GetRoadmapEpics() ([]models.RoadmapEpic, error) {
	items, err := s.Repo.GetRoadmapEpics()
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i].Progress = progressPercent(items[i].DoneCount, items[i].TicketCount)
		items[i].ProgressLabel = progressLabel(items[i].DoneCount, items[i].TicketCount)
	}
	return items, nil
}

func (s *ManagementService) GetRoadmapSprints() ([]models.RoadmapSprint, error) {
	items, err := s.Repo.GetRoadmapSprints()
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i].Progress = progressPercent(items[i].DoneCount, items[i].TicketCount)
		items[i].ProgressLabel = progressLabel(items[i].DoneCount, items[i].TicketCount)
	}
	return items, nil
}

func (s *ManagementService) CountRoadmapProjects() (int, error) {
	return s.Repo.CountRoadmapProjects()
}

func (s *ManagementService) GetRoadmapTickets() ([]models.RoadmapTicket, error) {
	return s.Repo.GetRoadmapTickets()
}

func (s *ManagementService) GetRoadmapProjectOptions() ([]models.ProjectOption, error) {
	return s.Repo.GetRoadmapProjectOptions()
}

func (s *ManagementService) GetRoadmapEpicOptions() ([]models.RoadmapEpicOption, error) {
	return s.Repo.GetRoadmapEpicOptions()
}

func (s *ManagementService) CreateRoadmapEpic(input models.RoadmapEpicCreateInput) error {
	input.Name = strings.TrimSpace(input.Name)
	input.StartsAt = strings.TrimSpace(input.StartsAt)
	input.EndsAt = strings.TrimSpace(input.EndsAt)

	if input.ProjectID <= 0 {
		return errors.New("project wajib dipilih")
	}
	if input.Name == "" {
		return errors.New("nama epic wajib diisi")
	}
	if input.StartsAt == "" || input.EndsAt == "" {
		return errors.New("tanggal mulai dan akhir wajib diisi")
	}

	start, err := time.Parse("2006-01-02", input.StartsAt)
	if err != nil {
		return errors.New("format tanggal mulai tidak valid")
	}
	end, err := time.Parse("2006-01-02", input.EndsAt)
	if err != nil {
		return errors.New("format tanggal akhir tidak valid")
	}
	if end.Before(start) {
		return errors.New("tanggal akhir tidak boleh lebih kecil dari tanggal mulai")
	}

	return s.Repo.CreateRoadmapEpic(input)
}

func (s *ManagementService) CreateRoadmapTicket(input models.RoadmapTicketCreateInput) error {
	input.Name = strings.TrimSpace(input.Name)
	if input.ProjectID <= 0 {
		return errors.New("project wajib dipilih")
	}
	if input.Name == "" {
		return errors.New("nama ticket wajib diisi")
	}
	if input.ResourceUserID <= 0 {
		return errors.New("resource wajib dipilih")
	}
	if input.Estimation < 0 {
		return errors.New("estimasi tidak valid")
	}
	return s.Repo.CreateRoadmapTicket(input)
}

func (s *ManagementService) BuildRoadmapTimeline(epics []models.RoadmapEpic, tickets []models.RoadmapTicket, now time.Time, format string) ([]models.RoadmapWeek, []models.RoadmapTimelineRow, int, int, int, int) {
	rangeStart, rangeEnd := timelineBounds(epics, now, format)
	columns, columnWidth := buildRoadmapColumns(rangeStart, rangeEnd, format)
	timelineWidth := len(columns) * columnWidth

	ticketsByEpic := map[int][]models.RoadmapTicket{}
	for _, ticket := range tickets {
		ticketsByEpic[ticket.EpicID] = append(ticketsByEpic[ticket.EpicID], ticket)
	}

	rows := make([]models.RoadmapTimelineRow, 0, len(epics)+len(tickets))
	for _, epic := range epics {
		epicTickets := ticketsByEpic[epic.ID]
		if len(epicTickets) > 0 {
			totalProgress := 0
			for _, ticket := range epicTickets {
				totalProgress += ticket.Progress
			}
			epic.Progress = totalProgress / len(epicTickets)
			epic.ProgressLabel = strconv.Itoa(epic.Progress) + "%"
		}

		rows = append(rows, makeTimelineRow(
			epic.Name,
			"",
			epic.Progress,
			epic.ProgressLabel,
			epic.StartsAt,
			epic.EndsAt,
			epic.StartsAtISO,
			epic.EndsAtISO,
			"#4f88eb",
			"#35c65a",
			false,
			rangeStart,
			format,
			columnWidth,
			true,
			"",
			epic.ProjectName,
		))

		for index, ticket := range epicTickets {
			rowTone := ""
			if index%2 == 0 {
				rowTone = "bg-amber-50/50"
			}
			rows = append(rows, models.RoadmapTimelineRow{
				Name:           ticket.Name,
				Resource:       ticket.ResourceName,
				Progress:       ticket.Progress,
				ProgressLabel:  strconv.Itoa(ticket.Progress) + "%",
				StartDateLabel: "",
				EndDateLabel:   "",
				BarColor:       "#4f88eb",
				BarAccentColor: "#35c65a",
				BarProgressPct: ticket.Progress,
				ShowBar:        false,
				IsChild:        true,
				SearchText:     strings.ToLower(ticket.Name + " " + ticket.ProjectName + " " + ticket.ResourceName),
				RowTone:        rowTone,
			})
		}
	}

	currentMarkerLeft, currentMarkerWidth := currentMarkerMetrics(now, rangeStart, format, columnWidth)

	return columns, rows, timelineWidth, currentMarkerLeft, currentMarkerWidth, columnWidth
}

func progressPercent(done, total int) int {
	if total <= 0 {
		return 0
	}
	value := (done * 100) / total
	if value > 100 {
		return 100
	}
	if value < 0 {
		return 0
	}
	return value
}

func progressLabel(done, total int) string {
	return strconv.Itoa(done) + "/" + strconv.Itoa(total) + " tickets"
}

func timelineBounds(epics []models.RoadmapEpic, now time.Time, format string) (time.Time, time.Time) {
	var starts []time.Time
	var ends []time.Time

	for _, epic := range epics {
		if start, err := time.Parse("2006-01-02", epic.StartsAtISO); err == nil {
			starts = append(starts, start)
		}
		if end, err := time.Parse("2006-01-02", epic.EndsAtISO); err == nil {
			ends = append(ends, end)
		}
	}
	if len(starts) == 0 || len(ends) == 0 {
		start := startOfWeek(now)
		return start, start.AddDate(0, 0, 83)
	}

	sort.Slice(starts, func(i, j int) bool { return starts[i].Before(starts[j]) })
	sort.Slice(ends, func(i, j int) bool { return ends[i].Before(ends[j]) })

	baseStart := starts[0].AddDate(0, 0, -7)
	baseEnd := ends[len(ends)-1].AddDate(0, 0, 28)

	switch format {
	case "day":
		return startOfDay(baseStart), startOfDay(baseEnd)
	case "month":
		return firstOfMonth(baseStart), endOfMonth(baseEnd)
	default:
		return startOfWeek(baseStart), endOfWeek(baseEnd)
	}
}

func buildRoadmapColumns(start, end time.Time, format string) ([]models.RoadmapWeek, int) {
	switch format {
	case "day":
		var columns []models.RoadmapWeek
		for cursor := start; !cursor.After(end); cursor = cursor.AddDate(0, 0, 1) {
			columns = append(columns, models.RoadmapWeek{
				YearLabel: cursor.Format("2006"),
				DateLabel: cursor.Format("02 Jan"),
			})
		}
		return columns, 62
	case "month":
		var columns []models.RoadmapWeek
		for cursor := firstOfMonth(start); !cursor.After(end); cursor = cursor.AddDate(0, 1, 0) {
			columns = append(columns, models.RoadmapWeek{
				YearLabel: cursor.Format("2006"),
				DateLabel: cursor.Format("Jan"),
			})
		}
		return columns, 88
	default:
		var columns []models.RoadmapWeek
		for cursor := start; !cursor.After(end); cursor = cursor.AddDate(0, 0, 7) {
			columns = append(columns, models.RoadmapWeek{
				YearLabel: cursor.Format("2006"),
				DateLabel: cursor.Format("02 Jan"),
			})
		}
		return columns, 55
	}
}

func makeTimelineRow(name, resource string, progress int, progressLabel, startLabel, endLabel, startISO, endISO, barColor, accentColor string, isChild bool, rangeStart time.Time, format string, columnWidth int, showBar bool, rowTone string, projectName string) models.RoadmapTimelineRow {
	start, errStart := time.Parse("2006-01-02", startISO)
	end, errEnd := time.Parse("2006-01-02", endISO)
	if errStart != nil || errEnd != nil {
		return models.RoadmapTimelineRow{
			Name:           name,
			Resource:       resource,
			Progress:       progress,
			ProgressLabel:  progressLabel,
			StartDateLabel: startLabel,
			EndDateLabel:   endLabel,
			BarLeftPx:      0,
			BarWidthPx:     80,
			BarColor:       barColor,
			BarAccentColor: accentColor,
			BarProgressPct: progress,
			ShowBar:        showBar,
			IsChild:        isChild,
			SearchText:     strings.ToLower(name + " " + resource + " " + projectName),
			RowTone:        rowTone,
		}
	}

	left, width := barMetrics(start, end, rangeStart, format, columnWidth)
	if width < 18 {
		width = 18
	}

	return models.RoadmapTimelineRow{
		Name:           name,
		Resource:       resource,
		Progress:       progress,
		ProgressLabel:  progressLabel,
		StartDateLabel: startLabel,
		EndDateLabel:   endLabel,
		BarLeftPx:      left,
		BarWidthPx:     width,
		BarColor:       barColor,
		BarAccentColor: accentColor,
		BarProgressPct: progress,
		ShowBar:        showBar,
		IsChild:        isChild,
		SearchText:     strings.ToLower(name + " " + resource + " " + projectName),
		RowTone:        rowTone,
	}
}

func barMetrics(start, end, rangeStart time.Time, format string, columnWidth int) (int, int) {
	switch format {
	case "day":
		offset := daysBetween(rangeStart, startOfDay(start))
		span := daysBetween(startOfDay(start), startOfDay(end)) + 1
		return offset * columnWidth, max(span, 1) * columnWidth
	case "month":
		offset := monthsBetween(firstOfMonth(rangeStart), firstOfMonth(start))
		span := monthsBetween(firstOfMonth(start), firstOfMonth(end)) + 1
		return offset * columnWidth, max(span, 1) * columnWidth
	default:
		offset := daysBetween(startOfWeek(rangeStart), startOfWeek(start)) / 7
		span := (daysBetween(startOfWeek(start), startOfWeek(end)) / 7) + 1
		return offset * columnWidth, max(span, 1) * columnWidth
	}
}

func currentMarkerMetrics(now, rangeStart time.Time, format string, columnWidth int) (int, int) {
	switch format {
	case "day":
		return daysBetween(rangeStart, startOfDay(now)) * columnWidth, columnWidth
	case "month":
		return monthsBetween(firstOfMonth(rangeStart), firstOfMonth(now)) * columnWidth, columnWidth
	default:
		return (daysBetween(startOfWeek(rangeStart), startOfWeek(now)) / 7) * columnWidth, columnWidth
	}
}

func startOfWeek(value time.Time) time.Time {
	normalized := time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
	weekday := int(normalized.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return normalized.AddDate(0, 0, -(weekday - 1))
}

func endOfWeek(value time.Time) time.Time {
	return startOfWeek(value).AddDate(0, 0, 6)
}

func startOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

func firstOfMonth(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, value.Location())
}

func endOfMonth(value time.Time) time.Time {
	return firstOfMonth(value).AddDate(0, 1, -1)
}

func daysBetween(start, end time.Time) int {
	return int(startOfDay(end).Sub(startOfDay(start)).Hours() / 24)
}

func monthsBetween(start, end time.Time) int {
	return (end.Year()-start.Year())*12 + int(end.Month()-start.Month())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
