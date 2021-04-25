package test

import (
	"errors"
	"testing"

	"github.com/yscsky/yu"
)

func TestGetCaller(t *testing.T) {
	cases := []struct {
		depth    int
		funcName string
		line     int
	}{
		{
			depth:    1,
			funcName: "test.TestGetCaller",
			line:     23,
		},
	}
	for _, ca := range cases {
		f, line := yu.GetCaller(ca.depth)
		if f != ca.funcName {
			t.Errorf("func wrong should be %s, but %s", ca.funcName, f)
		}
		if line != ca.line {
			t.Errorf("line wrong should be %d, but %d", ca.line, line)
		}
	}
}

func TestLogf(t *testing.T) {
	yu.Logf("test log 1")
	yu.Logf("test log %d", 2)
	yu.Logf("%s log %d", "test", 2)
}

func TestLogErr(t *testing.T) {
	yu.LogErr(errors.New("test err1"), "test err1")
	yu.LogErr(errors.New("test err2"), "test err2")
}

func TestSaveToml(t *testing.T) {
	var data = struct {
		Name      string
		Count     int
		Open      bool
		Databases []string
		Servers   map[string]string
	}{
		Name:      "test",
		Count:     10,
		Open:      true,
		Databases: []string{"mysql", "mariadb", "postgresql", "mongodb", "redis"},
		Servers: map[string]string{
			"server1": "127.0.0.1:8001",
			"server2": "127.0.0.1:8002",
			"server3": "127.0.0.1:8003",
			"server4": "127.0.0.1:8004",
			"server5": "127.0.0.1:8005",
		},
	}
	if err := yu.SaveToml("test.toml", data); err != nil {
		t.Error(err)
	}
}

func TestLoadToml(t *testing.T) {
	var data struct {
		Name      string
		Count     int
		Open      bool
		Databases []string
		Servers   map[string]string
	}
	if err := yu.LoadToml("test.toml", &data); err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestLoadOrSaveToml(t *testing.T) {
	type Data struct {
		Name      string
		Count     int
		Open      bool
		Databases []string
		Servers   map[string]string
	}
	data := &Data{}
	if err := yu.LoadOrSaveToml("test.toml", data, func() interface{} {
		data = &Data{
			Name:      "test",
			Count:     10,
			Open:      true,
			Databases: []string{"mysql", "mariadb", "postgresql", "mongodb", "redis"},
			Servers: map[string]string{
				"server1": "127.0.0.1:8001",
				"server2": "127.0.0.1:8002",
				"server3": "127.0.0.1:8003",
				"server4": "127.0.0.1:8004",
				"server5": "127.0.0.1:8005",
			},
		}
		return data
	}); err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestSaveJSON(t *testing.T) {
	var data = struct {
		Name      string
		Count     int
		Open      bool
		Databases []string
		Servers   map[string]string
	}{
		Name:      "test",
		Count:     10,
		Open:      true,
		Databases: []string{"mysql", "mariadb", "postgresql", "mongodb", "redis"},
		Servers: map[string]string{
			"server1": "127.0.0.1:8001",
			"server2": "127.0.0.1:8002",
			"server3": "127.0.0.1:8003",
			"server4": "127.0.0.1:8004",
			"server5": "127.0.0.1:8005",
		},
	}
	if err := yu.SaveJSON("test.json", data); err != nil {
		t.Error(err)
	}
}

func TestLoadJSON(t *testing.T) {
	var data struct {
		Name      string
		Count     int
		Open      bool
		Databases []string
		Servers   map[string]string
	}
	if err := yu.LoadJSON("test.json", &data); err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestLoadOrSaveJSON(t *testing.T) {
	type Data struct {
		Name      string
		Count     int
		Open      bool
		Databases []string
		Servers   map[string]string
	}
	data := &Data{}
	if err := yu.LoadOrSaveJSON("test.json", data, func() interface{} {
		data = &Data{
			Name:      "test",
			Count:     10,
			Open:      true,
			Databases: []string{"mysql", "mariadb", "postgresql", "mongodb", "redis"},
			Servers: map[string]string{
				"server1": "127.0.0.1:8001",
				"server2": "127.0.0.1:8002",
				"server3": "127.0.0.1:8003",
				"server4": "127.0.0.1:8004",
				"server5": "127.0.0.1:8005",
			},
		}
		return data
	}); err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}
