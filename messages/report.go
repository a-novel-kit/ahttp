package ahttpmessages

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"github.com/a-novel-kit/quicklog"
)

type Metrics struct {
	Latency   time.Duration
	StartedAt time.Time
}

type reportMessage struct {
	metrics   *Metrics
	projectID string
	ginC      *gin.Context

	quicklog.Message
}

func (report *reportMessage) RenderTerminal() string {
	errorMessage := ""
	for _, err := range report.ginC.Errors.Errors() {
		errorMessage += "\n" + lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("9")).Render("- "+err)
	}

	latencyMessage := ""
	if report.metrics != nil {
		latencyMessage = lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf(" (%s)", report.metrics.Latency))
	}

	query := report.ginC.Request.URL.Query()
	queryMessage := ""
	if len(query) > 0 {
		queryTable := table.New().
			Width(quicklog.TermWidth).
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("220"))).
			StyleFunc(func(_, col int) lipgloss.Style {
				if col == 0 {
					return lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
				}

				return lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
			})

		keys := lo.Keys(query)
		sort.Strings(keys)

		for _, key := range keys {
			values := query[key]
			sort.Strings(values)
			queryTable.Row("  "+key, "  "+strings.Join(values, "\n  "))
		}

		queryMessage = "\n" + queryTable.Render()
	}

	statusCode := report.ginC.Writer.Status()
	color := lipgloss.Color("33")
	prefix := "✅ "

	if statusCode > 499 {
		color = "9"
		prefix = "👶🔪🩸 "
	} else if statusCode > 399 {
		color = "202"
		prefix = "⚠ "
	}

	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Render(fmt.Sprintf("%s%v", prefix, statusCode)) +
		lipgloss.NewStyle().
			Foreground(color).
			Render(fmt.Sprintf(" [%s %s]", report.ginC.Request.Method, report.ginC.FullPath())) +
		latencyMessage +
		queryMessage +
		errorMessage +
		"\n\n"
}

func (report *reportMessage) RenderJSON() map[string]interface{} {
	statusCode := report.ginC.Writer.Status()

	severity := "INFO"
	if statusCode > 499 {
		severity = "ERROR"
	} else if statusCode > 399 {
		severity = "WARNING"
	}

	httpRequest := map[string]interface{}{
		"requestMethod": report.ginC.Request.Method,
		"requestUrl":    report.ginC.FullPath(),
		"status":        statusCode,
		"userAgent":     report.ginC.Request.UserAgent(),
		"remoteIp":      report.ginC.ClientIP(),
		"protocol":      report.ginC.Request.Proto,
	}

	output := map[string]interface{}{
		"httpRequest": httpRequest,
		"severity":    severity,
		"ip":          report.ginC.ClientIP(),
		"contentType": report.ginC.ContentType(),
		"errors":      report.ginC.Errors.Errors(),
		"query":       report.ginC.Request.URL.Query(),
	}

	if report.metrics != nil {
		output["start"] = report.metrics.StartedAt
		httpRequest["latency"] = report.metrics.Latency.String()
	}

	if report.projectID != "" {
		traceHeader := report.ginC.GetHeader("X-Cloud-Trace-Context")
		traceParts := strings.Split(traceHeader, "/")

		if len(traceParts) > 0 && len(traceParts[0]) > 0 {
			output["logging.googleapis.com/trace"] = fmt.Sprintf(
				"projects/%s/traces/%s",
				report.projectID, traceParts[0],
			)
		}
	}

	return output
}

func NewReport(metrics *Metrics, projectID string, ginC *gin.Context) quicklog.Message {
	return &reportMessage{
		metrics:   metrics,
		projectID: projectID,
		ginC:      ginC,
	}
}
