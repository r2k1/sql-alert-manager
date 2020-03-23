# SQL Alert Manager

Note: it's a beta version and likely to change in the future. Use at your own risk.

SQL Alert Manager is alerting system based on sql queries.
The idea is simple. You define DB connections and SQL queries. The alert manager occasionally executes the queries and if the result is not empty it will trigger an alert.

## Run

You need to set `CONFIG_PATH` environment variable to specify path to a configuration file.

Example of usage with docker, assuming you have `config.toml` in the working directory.
```docker run --env CONFIG_PATH=/mnt/config.toml --mount type=bind,source=$(pwd)/config.toml,target='/mnt/config.toml' r2k1/sql-alert-manager:latest```


## Configuration

Example:

```toml
reminder_interval = "3h"
alert_on_error = false

[alerts]
    [alerts.test-alert-1]
    query = """
        SELECT *
        FROM test
        WHERE something > 2
    """
    message = "Something is broken"
    dbs = ["my-postgres-db"]
    destinations = ["slacks.my-slack"]
    interval = "60s"
    reminder_interval = "1h"

    [alerts.test-alert-2]
    query = "SELECT * FROM test"
    message = "Something is broken"
    dbs = ["my-mysql-db"]
    destinations = ["slacks.my-slack"]
    interval = "1h30m"
    alert_on_error = true


[dbs]
    [dbs.my-postgres-db]
    driver = "postgres"
    connection = "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"


    [dbs.my-mysql-db]
    driver = "mysql"
    connection="${DB_USER}:${DB_PASSWORD}@tcp(example.com:3306)/mysql_db"


[slacks]
    [slacks.my-slack]
    webhook_url = "https://hooks.slack.com/services/secret/secret/supersecret"

```


Configuration is defined in [TOML](https://github.com/toml-lang/toml) format.
- `reminder_interval` - Optional. Default interval for all alerts after which alert will be triggered again. Set to 0 if you don't need reminders. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h". Examples: "300s", "1.5h" or "2h45m". Default value is 0.
- `alert_on_error` - Option. Defines behavior during error encounters (connection errors, sql syntax errors, timeouts etc ). If set to true any error will trigger an alert. If set to false then error will be logged and ignored. Default value is true.
- `alerts.{alert-name}.reminder_interval` - Optional. Same as above, but for individual alert.
- `alerts.{alert-name}.query` - Required. SQL query to execute at regular interval.
- `alerts.{alert-name}.message` - Optional. Message to pass with alert.
- `alerts.{alert-name}.dbs` - Required. List of database references. Provided query will be executed against each database and trigger an individual alert for each database. All databases must be defined in `dbs` section of the configuration. Example: `["my-postgres-db", "my-mysql-db"]`.
- `alerts.{alert-name}.destinations` - Required. List of destination references to report the alert. Each destination should be defined in related section. Example: `["slacks.channel-1-webhook", "slacks.channel-2-webhook"]`
- `alerts.{alert-name}.interval` - Required. An interval between consecutive query execution. For simplicity it does not take in account time required to execute the query. For example if query execution time is 5s and interval is 10s then interval between two consecutive queries will be 15s. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h". Examples: "300s", "1.5h" or "2h45m"
- `dbs.{db-name}.driver` - Required. SQL driver. Supported drivers: mysql, postgres.
- `dbs.{db-name}.connection` - Required. Connection string for the database. Documentation for databases: [postgres](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING), [mysql](https://github.com/go-sql-driver/mysql#dsn-data-source-name).
- `slacks.{destination-name}.webhook_url` - Required. Webhook URL for slack integration. More information how to set it up in [official slack documentation](https://api.slack.com/messaging/webhooks).

### Environment substitution.

Configuration is supporting environment variable substitution.\
You can put `${MY_ENV_VARIABLE}` in configuration and it will be replaced with the value.\
More information about available options is [here](https://github.com/a8m/envsubst#docs).\
It's a recommended way to provide secrets if you don't want them to be exposed in the configuration file.\
Note: `$` in configuration is treated as an expression. You need to escape it with `$$`.
