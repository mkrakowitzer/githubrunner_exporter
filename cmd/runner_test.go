package cmd

import "testing"

func TestSetRunnerStatusGauge(t *testing.T) {

	result, _ := setRunnerStatusGauge("online")
	if result != 1 {
		t.Errorf("setRunnerStatusGauge(\"online\") failed, expected %v, got %v", 1, result)
	}

	result, _ = setRunnerStatusGauge("offline")
	if result != 1 {
		t.Errorf("setRunnerStatusGauge(\"offline\") failed, expected %v, got %v", 0, result)
	}

}
