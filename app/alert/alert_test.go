package alert

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/r2k1/sql-alert-manager/app/sqltest"
)

//go:generate  mockery -name Destination -inpkg

func TestAlert_ExecQuery_Postgres(t *testing.T) {
	postgres, err := sql.Open("postgres", sqltest.TestPostgresConnection)
	require.NoError(t, err)
	a := Alert{
		Source: NewSource(postgres, "test"),
		Query: `
			SELECT 1 as col1, 'string' as col2, ARRAY[['meeting', 'lunch'], ['training', 'presentation']] as col3
			UNION ALL
			SELECT 2 as col1, 'string2' as col2, ARRAY['test'] as col3`,
	}
	res, err := a.ExecQuery()
	require.NoError(t, err)
	assert.Equal(t, `col1	col2	col3
1	string	{{meeting,lunch},{training,presentation}}
2	string2	{test}
`, res)
}

func TestAlert_ExecQuery_Mysql(t *testing.T) {
	postgres, err := sql.Open("mysql", sqltest.TestMysqlConnection)
	require.NoError(t, err)
	a := Alert{
		Name:   "my_alert",
		Source: NewSource(postgres, "test"),
		Query: `
			SELECT 1 as col1, 'string' as col2
			UNION ALL
			SELECT 2 as col1, 'string2' as col2`,
	}
	res, err := a.ExecQuery()
	require.NoError(t, err)
	assert.Equal(t, `col1	col2
1	string
2	string2
`, res)
}

func TestAlert_Check(t *testing.T) {
	d := new(MockDestination)
	d.On("SendAlert", mock.Anything, mock.Anything).Return(nil)
	postgres, err := sql.Open("mysql", sqltest.TestMysqlConnection)
	require.NoError(t, err)
	a := Alert{
		Name:         "my_alert",
		Source:       NewSource(postgres, "test"),
		Destinations: []Destination{d},
		Query: `
			SELECT 1 as col1, 'string' as col2
			UNION ALL
			SELECT 2 as col1, 'string2' as col2`,
	}
	a.Check()
	d.AssertNumberOfCalls(t, "SendAlert", 1)
	a.Check()
	a.Check()
	d.AssertNumberOfCalls(t, "SendAlert", 1)
}
