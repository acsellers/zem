package store

import (
	"strings"

	"github.com/acsellers/inflections"
	_ "github.com/lib/pq"
)

func (ac AppConfig) SQLTable(table string) string {
	if ac.SpecialTables == nil || ac.SpecialTables[table] == "" {
		return strings.ToLower(inflections.Pluralize(inflections.Underscore(table)))
	}
	return ac.SpecialTables[table]
}

func (ac AppConfig) SQLColumn(table, column string) string {
	return strings.ToLower(inflections.Underscore(column))
}
