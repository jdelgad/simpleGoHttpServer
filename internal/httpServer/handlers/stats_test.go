package handlers

import (
	"testing"
	"time"
)

func TestConvertNsToMicroSeconds(t *testing.T) {
	twentyThreeThousandNs := 23000 * time.Nanosecond
	var twentyThreeMicroSecs int64
	twentyThreeMicroSecs = 23

	convertedNsToMicroSecs := convertNsToMicroSeconds(twentyThreeThousandNs)

	if convertedNsToMicroSecs != twentyThreeMicroSecs {
		t.Fatalf("conversion from 23000ns to microseconds failed. Expected 23us, Received %v", convertedNsToMicroSecs)
	}
}

func TestStats_Update(t *testing.T) {
	s := &Stats{}
	updateTime := 12 * time.Nanosecond
	s.Update(updateTime)

	if s.Total != 1 {
		t.Fatal("total member not updated")
	}

	if s.timeTaken != updateTime {
		t.Fatal("total member not updated")
	}

	s.Update(updateTime)
	if s.Total != 2 {
		t.Fatal("total member not updated on the 2nd call")
	}

	if s.timeTaken != 2*updateTime {
		t.Fatal("total member not updated on the 2nd call")
	}

}

func TestStats_ToJSON(t *testing.T) {
	s := &Stats{}
	updateTime := 11 * time.Microsecond
	s.Update(updateTime)
	updateTime = 20 * time.Microsecond
	s.Update(updateTime)

	js, err := s.ToJSON()
	if err != nil {
		t.Fatalf("marshaling of Stats struct failure: %v", err)
	}

	toString := string(js)
	if toString != `{"total":2,"average":15}` {
		t.Fatalf("%v", toString)
	}

}
