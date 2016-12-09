package lib

import (
	"testing"
	"fmt"
)

func TestGetHostStat (t *testing.T) {
	v := new(ServerStat)
	err := v.GetHostStat()
	if err != nil {
		t.Errorf("error %v", err)
	}

	empty := &ServerStat{}
	if v == empty {
		t.Errorf("error : cannot get hostStat")
	}

	if v.HostName == "" {
		t.Errorf("error : cannot get hostName")
	}

	if v.HostID == "" {
		t.Errorf("error : cannot get hostID")
	}

	if v.VirtualizationSystem == "" {
		t.Errorf("error : cannot get virtualizationSystem")
	}
}

func TestMemoryStat (t *testing.T) {
	v := new(ServerStat)
	err := v.GetMemoryStat()
	if err != nil {
		t.Errorf("error %v", err)
	}

	empty := &ServerStat{}
	if v == empty {
		t.Errorf("error : cannot get memoryStat")
	}
}

func TestDiskIOStat (t *testing.T) {
	v := new(ServerStat)
	err := v.GetDiskIOStat()
	if err != nil {
		t.Errorf("error %v", err)
	}

	empty := &ServerStat{}
	if v == empty {
		t.Errorf("error : cannot get serverStat")
	}
}

func TestApacheStat (t *testing.T) {
	v := new(ServerStat)
	err := v.GetApacheStat()
	if err != nil {
		t.Errorf("error %v", err)
	}

	empty := &ServerStat{}
	if v == empty {
		t.Errorf("error : cannot get apacheStat")
	}
}

func TestDiskStat_String (t *testing.T) {
	v := DiskStat {
		Name: "disk",
		IoTime: 100,
		WeightedIO: 100,
	}

	e := `{"name":"disk","ioTime":100,"weightedIO":100}`

	if e != fmt.Sprintf("%v", v) {
		t.Errorf("DiskStat string is invalid: %v", v)
	}
}

func TestServerStat_String (t *testing.T) {
	vd1 := DiskStat {
		Name: "disk1",
		IoTime: 123,
		WeightedIO: 123,
	}

	vd2 := DiskStat {
		Name: "disk2",
		IoTime: 200,
		WeightedIO: 300,
	}

	vs := ServerStat {
		HostName:             "host",
		HostID:               "123",
		VirtualizationSystem: "vbox",
		Total:                123456,
		Available:            123456,
		UsedPercent:          123.456,
		DiskIO: []DiskStat {
			vd1,
			vd2,
		},
		Time: "00:00:00",
		ApacheStat: 123.456,
	}

	e := `{"hostname":"host","hostid":"123","virtualizationSystem":"vbox","total":123456,"available":123456,"usedPercent":123.456,"diskIO":[{"name":"disk1","ioTime":123,"weightedIO":123},{"name":"disk2","ioTime":200,"weightedIO":300}],"time":"00:00:00","apacheStat":123.456,"errorInfo":null}`

	if e != fmt.Sprintf("%v", vs) {
		t.Errorf("ServerStat string is invalid: %v", vs)
	}
}
