package drivers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/jorgerojas26/lazysql/models"

	"net/http"
)

type WebProxyClient struct {
	endpoint         string
	connectionString string
	provider         string
}

func postRequest[T any](cli *WebProxyClient, actionId string, keyval ...string) (T, error) {
	var result T
	m := map[string]string{
		"connectionString": cli.connectionString,
	}
	count := len(keyval)
	dat, err := json.Marshal(m)
	if err != nil {
		return result, err
	}
	b := bytes.NewBuffer(dat)
	for i := 0; i < count; i += 2 {
		m[keyval[i]] = m[keyval[i+1]]
	}
	r, err := http.Post(fmt.Sprintf("%s/%s", cli.endpoint, actionId), "application/json", b)
	if err != nil {
		return result, err
	}
	out := bytes.NewBufferString("")
	if _, err := io.Copy(out, r.Body); err != nil {
		return result, err
	}
	if err := json.NewDecoder(out).Decode(&result); err != nil {
		return result, err
	}

	return result, err

}
func (cli *WebProxyClient) TestConnection(urlstr string) (err error) {
	_, err = postRequest[any](cli, "TestConnection", "connectionString", urlstr)
	return err
}

func (cli *WebProxyClient) Connect(urlstr string) (err error) {
	_, err = postRequest[any](cli, "Connect", "connectionString", urlstr)
	return err
}

func (cli *WebProxyClient) GetDatabases() ([]string, error) {
	return postRequest[[]string](cli, "GetDatabases")
}

func (cli *WebProxyClient) GetTables(database string) (map[string][]string, error) {
	return postRequest[map[string][]string](cli, "GetTables", "database", database)
}

func (cli *WebProxyClient) GetTableColumns(database, table string) ([][]string, error) {
	return postRequest[[][]string](cli, "GetTableColumns", "database", database, "table", table)

}

func (cli *WebProxyClient) GetConstraints(table string) ([][]string, error) {
	return postRequest[[][]string](cli, "GetConstraints", "table", table)

}

func (cli *WebProxyClient) GetForeignKeys(table string) ([][]string, error) {
	return postRequest[[][]string](cli, "GetForeignKeys", "table", table)

}

func (cli *WebProxyClient) GetIndexes(table string) ([][]string, error) {
	return postRequest[[][]string](cli, "GetIndexes", "table", table)
}

func (cli *WebProxyClient) GetRecords(table, where, sort string, offset, limit int) ([][]string, int, error) {
	params := []string{"table", table, "sort", sort, "where", where,
		"offset", strconv.FormatInt(int64(offset), 10),
		"limit", strconv.FormatInt(int64(limit), 10),
	}
	result, err := postRequest[[][]string](cli, "GetRecords", params...)
	if err != nil {
		return nil, -1, err
	}
	n, err := postRequest[int](cli, "GetRecordCount", params...)
	if err != nil {
		return nil, -1, err
	}
	return result, n, err

}

func (cli *WebProxyClient) ExecuteQuery(query string) ([][]string, error) {
	return postRequest[[][]string](cli, "ExecuteQuery", "query", query)
}

func (cli *WebProxyClient) UpdateRecord(table, column, value, primaryKeyColumnName, primaryKeyValue string) error {
	_, err := postRequest[any](cli, "UpdateRecord", "table", table, "column", column, "value", value,
		"primaryKeyColumnName", primaryKeyColumnName, "primaryKeyValue", primaryKeyValue)
	return err
}

// TODO: Rewrites this logic to use the primary key instead of the id
func (cli *WebProxyClient) DeleteRecord(table, primaryKeyColumnName, primaryKeyValue string) error {
	_, err := postRequest[any](cli, "DeleteRecord", "table", table,
		"primaryKeyColumnName", primaryKeyColumnName, "primaryKeyValue", primaryKeyValue)
	return err
}

func (cli *WebProxyClient) ExecuteDMLStatement(query string) (string, error) {
	return postRequest[string](cli, "ExecuteDMLStatement", "query", query)
}

func (cli *WebProxyClient) ExecutePendingChanges(changes []models.DbDmlChange, inserts []models.DbInsert) error {
	panic("todo")
}

func (cli *WebProxyClient) SetProvider(provider string) {
	cli.provider = provider
}

func (cli *WebProxyClient) GetProvider() string {
	return cli.provider
}
