package alert

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type Alert struct {
	Name             string
	Message          string
	Source           *Source
	Query            string
	Destinations     []Destination
	Interval         time.Duration
	ReminderInterval time.Duration
	AlertingAt       time.Time
}

type Destination interface {
	SendAlert(a *Alert, msg string) error
	ResolveAlert(a *Alert, msg string) error
	Name() string
}

type Source struct {
	db   *sql.DB
	name string
}

func NewSource(db *sql.DB, name string) *Source {
	return &Source{
		db:   db,
		name: name,
	}
}

func (a *Alert) ExecQuery() (string, error) {
	var sb strings.Builder
	rows, err := a.Source.db.Query(a.Query)
	if err != nil {
		return "", fmt.Errorf("couldn't execute sql query: %w", err)
	}
	cols, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("error during fetching list of columns: %w", err)
	}
	isFirst := true
	for i := range cols {
		if !isFirst {
			sb.WriteString("\t")
		}
		sb.WriteString(cols[i])
		isFirst = false
	}
	sb.WriteString("\n")
	dataStr := make([][]byte, len(cols))
	data := make([]interface{}, len(cols))
	for i := range data {
		data[i] = &dataStr[i]
	}

	var hasRows bool
	for rows.Next() {
		hasRows = true
		err = rows.Scan(data...)
		if err != nil {
			return "", err
		}
		isFirst = true
		for i := range dataStr {
			if !isFirst {
				sb.WriteString("\t")
			}
			sb.Write(dataStr[i])
			isFirst = false
		}
		sb.WriteString("\n")
	}
	if !hasRows {
		return "", nil
	}
	return sb.String(), nil
}

func (a *Alert) Worker() {
	for {
		a.Check()
		time.Sleep(a.Interval)
	}
}

func (a *Alert) Check() {
	msg, err := a.ExecQuery()
	if err != nil {
		LogError(fmt.Errorf("couldn't check conditions for %s: %s", a.Name, err))
		return
	}
	if msg == "" {
		if !a.AlertingAt.IsZero() {
			log.Printf("INFO: %s is resolved", a.Name)
			a.Resolve(fmt.Sprintf("Resolved after %s", humanizeDuration(time.Now().Sub(a.AlertingAt))))
			a.AlertingAt = time.Time{}
		}
		return
	}
	log.Printf("INFO: %s is alerting", a.Name)
	if a.AlertingAt.IsZero() {
		a.AlertingAt = time.Now()
		a.SendAlert(msg)
	} else if a.ReminderInterval > 0 && time.Now().Sub(a.AlertingAt) > a.ReminderInterval {
		a.SendAlert(msg)
	}

}

func (a *Alert) SendAlert(msg string) {
	if !a.AlertingAt.IsZero() {
		a.AlertingAt = time.Now()
	}
	for i := range a.Destinations {
		err := a.Destinations[i].SendAlert(a, msg)
		if err != nil {
			LogError(fmt.Errorf("couldn't send alert %s to %s", a.Name, a.Destinations[i].Name()))
		}
	}
}

func (a *Alert) Resolve(msg string) {
	for i := range a.Destinations {
		err := a.Destinations[i].ResolveAlert(a, msg)
		if err != nil {
			LogError(fmt.Errorf("error during resolver alert %s to %s", a.Name, a.Destinations))
		}
	}
}

func LogError(err error) {
	log.Printf("ERROR: %v", err)
}