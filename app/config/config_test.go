package config

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/r2k1/sql-alert-manager/app/sqltest"
)

func TestMain(m *testing.M) {
	// set env variable for substitution in config file
	_ = os.Setenv("TEST_POSTGRES_CONNECTION", sqltest.TestPostgresConnection)
	_ = os.Setenv("TEST_MYSQL_CONNECTION", sqltest.TestMysqlConnection)
	os.Exit(m.Run())
}

func Test_loadTOMLConfig(t *testing.T) {
	config, err := loadTOMLConfig("example.toml")
	require.NoError(t, err)
	falseB := false
	trueB := true
	expected := tomlConfig{
		AlertOnError: &trueB,
		ReminderInterval: tomlDuration{Duration: 3 * time.Hour},
		Alerts: map[string]tomlAlert{
			"test-alert-1": {
				Query:            `SELECT * FROM test WHERE something > 2`,
				Message:          "Something is broken",
				DBs:              []string{"my-postgres-db"},
				Destinations:     []string{"slacks.my-slack"},
				Interval:         tomlDuration{Duration: time.Second * 60},
				ReminderInterval: tomlDuration{Duration: time.Hour},
			},
			"test-alert-2": {
				Query:            "SELECT * FROM test",
				Message:          "Something is broken",
				DBs:              []string{"my-mysql-db"},
				Destinations:     []string{"slacks.my-slack"},
				Interval:         tomlDuration{Duration: time.Minute * 90},
				ReminderInterval: tomlDuration{Duration: 0},
				AlertOnError:     &falseB,
			},
		},
		DB: map[string]tomlDB{
			"my-postgres-db": {
				Driver:     "postgres",
				Connection: os.Getenv("TEST_POSTGRES_CONNECTION"),
			},
			"my-mysql-db": {
				Driver:     "mysql",
				Connection: os.Getenv("TEST_MYSQL_CONNECTION"),
			},
		},
		Slack:     map[string]tomlSlack{"my-slack": {WebhookURL: "https://slack.com/something"}},
		PagerDuty: nil,
	}
	assertEqualConfig(t, expected, config)
}

func assertEqualConfig(t *testing.T, expected, actual tomlConfig) {
	stripSpacesFromQueries := func(alerts map[string]tomlAlert) {
		for k, v := range alerts {
			v.Query = stripSpaces(v.Query)
			alerts[k] = v
		}
	}
	stripSpacesFromQueries(expected.Alerts)
	stripSpacesFromQueries(actual.Alerts)
	assert.Equal(t, expected, actual)
}

func TestLOADTOMLConfig(t *testing.T) {
	config, err := loadTOMLConfig("example.toml")
	require.NoError(t, err)
	alerts, err := prepareAlerts(config)
	require.NoError(t, err)
	assert.Equal(t, 2, len(alerts))
}

func stripSpaces(str string) string {
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "\t", "")
	str = strings.ReplaceAll(str, "\n", "")
	return str
}

func TestLoadAlerts(t *testing.T) {
	alerts, err := LoadAlerts("example.toml")
	require.NoError(t, err)
	assert.Equal(t, 2, len(alerts))
	assert.NotNil(t, alerts[0].Source)
}
