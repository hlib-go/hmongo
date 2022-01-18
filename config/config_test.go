package config

import (
	_ "embed"
	"testing"
)

//go:embed .dconfig_mongo_dsn
var dsn string

func TestNew(t *testing.T) {
	t.Log("........")
	c := New(dsn, "dsys_config_dev")
	t.Log(c.GetCache())
}
