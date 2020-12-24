package hmongo

import (
	"encoding/json"
	"testing"
	"time"
)

func TestHConfig_Get(t *testing.T) {
	cfg := NewHConfig(DefaultDB.Collection("hm_config"))
	var wx *struct {
		Appid  string `json:"appid"`
		Secret string `json:"secret"`
	}
	err := cfg.Get("weixin", &wx)
	if err != nil {
		t.Error(err)
		return
	}
	wxb, err := json.Marshal(wx)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(wxb))
}

func TestHConfig_Watch(t *testing.T) {
	cfg := NewHConfig(DefaultDB.Collection("hm_config"))
	cfg.Watch(func(name string, value []byte) {
		t.Log(name, " ", string(value))
	})
	time.Sleep(1 * time.Hour)
}
