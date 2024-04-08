package mysqlEntity

type Table struct {
	Name string `db:"TABLE_NAME"`
}

type Column struct {
	Name string `db:"COLUMN_NAME"`
	Type string `db:"DATA_TYPE"`
}
