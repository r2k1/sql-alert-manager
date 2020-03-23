package config

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/r2k1/sql-alert-manager/app/alert"

	"github.com/BurntSushi/toml"
	"github.com/a8m/envsubst"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func LoadAlerts(configPath string) ([]alert.Alert, error) {
	c, err := loadTOMLConfig(configPath)
	if err != nil {
		return nil, err
	}
	return prepareAlerts(c)
}

type tomlAlert struct {
	Query            string       `toml:"query"`
	DBs              []string     `toml:"dbs"`
	Destinations     []string     `toml:"destinations"`
	Interval         tomlDuration `toml:"interval"`
	ReminderInterval tomlDuration `toml:"reminder_interval"`
	Message          string       `toml:"message"`
	AlertOnError     *bool        `toml:"alert_on_error"` // can hold 3 possible values: true/false/not-specified
}

type tomlDuration struct {
	time.Duration
}

func (d *tomlDuration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

type tomlDB struct {
	Driver     string `toml:"driver"`
	Connection string `toml:"connection"`
}

type tomlSlack struct {
	WebhookURL string `toml:"webhook_url"`
}

type tomlPagerDuty struct {
}

type tomlConfig struct {
	Alerts           map[string]tomlAlert     `toml:"alerts"`
	DB               map[string]tomlDB        `toml:"dbs"`
	Slack            map[string]tomlSlack     `toml:"slacks"`
	PagerDuty        map[string]tomlPagerDuty `toml:"pager_duties"`
	Message          string                   `toml:"message"`
	ReminderInterval tomlDuration             `toml:"reminder_interval"`
	AlertOnError     *bool                     `toml:"alert_on_error"` // can hold 3 possible values: true/false/not-specified
}

func loadTOMLConfig(path string) (tomlConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return tomlConfig{}, fmt.Errorf("couldn't open file %v: %w", path, err)
	}
	input, err := ioutil.ReadAll(f)
	if err != nil {
		return tomlConfig{}, fmt.Errorf("couldn't read file %v: %w", path, err)
	}
	return parseConfig(input)
}

// split from ParseConfigFile for testing purposes
func parseConfig(input []byte) (tomlConfig, error) {
	var c tomlConfig
	inputWithENV, err := envsubst.String(string(input))
	if err != nil {
		return c, fmt.Errorf("error during env substition: %w", err)
	}
	err = toml.Unmarshal([]byte(inputWithENV), &c)
	if err != nil {
		return c, fmt.Errorf("error parsing TOML: %w", err)
	}
	return c, nil
}

func prepareAlerts(config tomlConfig) ([]alert.Alert, error) {
	alerts := make([]alert.Alert, 0)
	dbs, err := prepareDBs(config)
	if err != nil {
		return nil, fmt.Errorf("invalid configration for db section: %w", err)
	}
	destinations, err := prepareDestinations(config)
	if err != nil {
		return nil, fmt.Errorf("invalid configration for destinations section: %w", err)
	}

	for alertName, alertConfig := range config.Alerts {
		var alertOnError bool
		if alertConfig.AlertOnError == nil {
			if config.AlertOnError == nil {
				alertOnError = true
			} else {
				alertOnError = *config.AlertOnError
			}
		} else {
			alertOnError = *alertConfig.AlertOnError
		}
		reminderInterval := alertConfig.ReminderInterval.Duration
		if reminderInterval == 0 {
			reminderInterval = config.ReminderInterval.Duration
		}

		alertDestinations := make([]alert.Destination, 0)
		for _, destName := range alertConfig.Destinations {
			dest, ok := destinations[destName]
			if !ok {
				return nil, fmt.Errorf("%v destination for %v isn't configured", destName, alertName)
			}
			alertDestinations = append(alertDestinations, dest)
		}

		for _, dbName := range alertConfig.DBs {
			db, ok := dbs[dbName]
			if !ok {
				return nil, fmt.Errorf("%v database is not configured (%v)", dbName, alertName)
			}

			alerts = append(alerts, alert.Alert{
				Name:             fmt.Sprintf("%s (%s)", alertName, dbName),
				Message:          alertConfig.Message,
				Source:           alert.NewSource(db, dbName),
				Query:            alertConfig.Query,
				Interval:         alertConfig.Interval.Duration,
				ReminderInterval: reminderInterval,
				Destinations:     alertDestinations,
				AlertOnError:     alertOnError,
			})
		}
	}

	return alerts, nil
}

func prepareDestinations(config tomlConfig) (map[string]alert.Destination, error) {
	destinations := make(map[string]alert.Destination)
	for name, slackConf := range config.Slack {
		destinations["slacks."+name] = alert.NewSlack(name, slackConf.WebhookURL)
	}
	return destinations, nil
}

func prepareDBs(config tomlConfig) (map[string]*sql.DB, error) {
	dbs := map[string]*sql.DB{}
	for key, dbConf := range config.DB {
		var err error
		dbs[key], err = sql.Open(dbConf.Driver, dbConf.Connection)
		if err != nil {
			return nil, fmt.Errorf("invalid connection string for %v: %w", key, err)
		}
	}
	return dbs, nil
}
