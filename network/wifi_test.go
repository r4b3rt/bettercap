package network

import "testing"

func buildExampleWiFi() *WiFi {
	return NewWiFi(buildExampleEndpoint(), func(ap *AccessPoint) {}, func(ap *AccessPoint) {})
}

func TestDot11Freq2Chan(t *testing.T) {
	exampleFreq := 2472
	exp := 13
	got := Dot11Freq2Chan(exampleFreq)
	if got != exp {
		t.Fatalf("expected '%v', got '%v'", exp, got)
	}
}

func TestDot11Chan2Freq(t *testing.T) {
	exampleChan := 13
	exp := 2472
	got := Dot11Chan2Freq(exampleChan)
	if got != exp {
		t.Fatalf("expected '%v', got '%v'", exp, got)
	}
}

func TestNewWiFi(t *testing.T) {
	exampleWiFi := NewWiFi(buildExampleEndpoint(), func(ap *AccessPoint) {}, func(ap *AccessPoint) {})
	if exampleWiFi == nil {
		t.Error("unable to build net wifi struct")
	}
}

func TestWiFiMarshalJSON(t *testing.T) {
	exampleWiFi := buildExampleWiFi()
	json, err := exampleWiFi.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	if len(json) <= 0 {
		t.Error("unable to marshal JSON from WiFi struct")
	}
}

func TestEachAccessPoint(t *testing.T) {
	exampleWiFi := buildExampleWiFi()
	exampleAP := NewAccessPoint("my_wifi", "ff:ff:ff:ff:ff:ff", 2472, int8(0))
	exampleWiFi.aps["ff:ff:ff:ff:ff:f1"] = exampleAP
	exampleWiFi.aps["ff:ff:ff:ff:ff:f2"] = exampleAP
	count := 0
	exampleCB := func(mac string, ap *AccessPoint) {
		count++
	}
	exampleWiFi.EachAccessPoint(exampleCB)
	if count != 2 {
		t.Error("unable to perform callback function for each access point")
	}
}

func TestStations(t *testing.T) {
	exampleWiFi := buildExampleWiFi()
	exampleAP := NewAccessPoint("my_wifi", "ff:ff:ff:ff:ff:ff", 2472, int8(0))
	exampleWiFi.aps["ff:ff:ff:ff:ff:f1"] = exampleAP
	exampleWiFi.aps["ff:ff:ff:ff:ff:f2"] = exampleAP
	exp := 2
	got := len(exampleWiFi.Stations())
	if got != exp {
		t.Fatalf("expected '%v', got '%v'", exp, got)
	}
}
