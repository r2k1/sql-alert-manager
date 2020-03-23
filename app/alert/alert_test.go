package alert

import (
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/r2k1/sql-alert-manager/app/sqltest"
)

//go:generate  mockery -name Destination -inpkg

func TestAlert_ExecQuery_Postgres(t *testing.T) {
	postgres := sqltest.GetTestPostgres()
	a := Alert{
		Source: NewSource(postgres, "test"),
		Query: `
			SELECT 1 as col1, 'string' as col2, ARRAY[['meeting', 'lunch'], ['training', 'presentation']] as col3
			UNION ALL
			SELECT 2 as col1, 'string2' as col2, ARRAY['test'] as col3`,
	}
	res, err := a.ExecQuery()
	require.NoError(t, err)
	assertTables(t, `
+------+---------+-------------------------------------------+
| COL1 |  COL2   |                   COL3                    |
+------+---------+-------------------------------------------+
|    1 | string  | {{meeting,lunch},{training,presentation}} |
|    2 | string2 | {test}                                    |
+------+---------+-------------------------------------------+
`,
		res)
}

func TestAlert_ExecQuery_Mysql(t *testing.T) {
	mysql := sqltest.GetTestMysql()
	a := Alert{
		Name:   "my_alert",
		Source: NewSource(mysql, "test"),
		Query: `
			SELECT 1 as col1, 'string' as col2
			UNION ALL
			SELECT 2 as col1, 'string2' as col2`,
	}
	res, err := a.ExecQuery()
	require.NoError(t, err)
	assertTables(t, `
+------+---------+
| COL1 |  COL2   |
+------+---------+
|    1 | string  |
|    2 | string2 |
+------+---------+
`,
		res)
}

func TestAlert_Check(t *testing.T) {
	d := new(MockDestination)
	d.On("SendAlert", mock.Anything, mock.Anything).Return(nil)
	mysql := sqltest.GetTestPostgres()
	a := Alert{
		Name:         "my_alert",
		Source:       NewSource(mysql, "test"),
		Destinations: []Destination{d},
		Query: `
			SELECT 1 as col1, 'string' as col2
			UNION ALL
			SELECT 2 as col1, 'string2' as col2`,
	}
	a.Check()
	assert.NotEqual(t, a.LastAlertSentAt, time.Time{})
	assert.NotEqual(t, a.AlertingAt, time.Time{})
	d.AssertNumberOfCalls(t, "SendAlert", 1)
	a.Check()
	a.Check()
	d.AssertNumberOfCalls(t, "SendAlert", 1)
}

func TestAlert_Check_Error(t *testing.T) {
	d := new(MockDestination)
	d.On("SendAlert", mock.Anything, mock.Anything).Return(nil)
	mysql := sqltest.GetTestPostgres()
	a := Alert{
		Name:         "my_alert",
		Source:       NewSource(mysql, "test"),
		Destinations: []Destination{d},
		AlertOnError: false,
		Query: `
			SELECT 1 as col1, 'string' as col2 FROM`, // invalid query
	}
	a.Check()
	d.AssertNumberOfCalls(t, "SendAlert", 0)
	a.AlertOnError = true
	a.Check()
	assert.NotEqual(t, a.LastAlertSentAt, time.Time{})
	assert.NotEqual(t, a.AlertingAt, time.Time{})
	d.AssertNumberOfCalls(t, "SendAlert", 1)
}


func TestAlert_Reminder(t *testing.T) {
	d := new(MockDestination)
	d.On("SendAlert", mock.Anything, mock.Anything).Return(nil)
	mysql := sqltest.GetTestPostgres()
	a := Alert{
		Name:         "my_alert",
		Source:       NewSource(mysql, "test"),
		Destinations: []Destination{d},
		ReminderInterval: time.Millisecond,
		Query: `
			SELECT 1 as col1, 'string' as col2`,
	}
	a.Check()
	d.AssertNumberOfCalls(t, "SendAlert", 1)
	a.Check()
	d.AssertNumberOfCalls(t, "SendAlert", 1)
	time.Sleep(time.Millisecond)
	a.Check()
	d.AssertNumberOfCalls(t, "SendAlert", 2)
}

func assertTables(t *testing.T, actual string, expected string) {
	actual = strings.TrimPrefix(actual, " ")
	expected = strings.TrimPrefix(actual, " ")
	assert.Equal(t, actual, expected)
}
