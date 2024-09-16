package date

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	Sydney, err := time.LoadLocation("Australia/Sydney")
	if err != nil {
		panic(err)
	}

	for i, test := range []struct {
		want string
		got  interface{}
	}{
		{
			"1970-01-01",
			MustFromString("1970-01-01"),
		},
		{
			"2018-03-30",
			FromTime(time.Date(2018, time.March, 30, 12, 01, 45, 0, time.UTC)),
		},
		{
			"2005-04-21",
			FromTime(MustFromString("2005-04-21").Time()),
		},
		{
			"2001-04-11 00:00:00 +1000 AEST",
			MustFromString("2001-04-11").TimeIn(Sydney),
		},
		{
			"2000-02-25",
			MustFromString("1999-12-25").AddMonths(2),
		},
		{
			"1999-10-25",
			MustFromString("1999-12-25").AddMonths(-2),
		},
		{
			"1999-07-01", // May only has 31 days, so normalises to 1st of July.
			MustFromString("1999-05-31").AddMonths(1),
		},
		{
			"2000-12-25",
			MustFromString("1999-12-25").AddYears(1),
		},
		{
			"2017-03-01", // Feb 2016 has 31 days, but Feb 2017 has 28 days, so normalises.
			MustFromString("2016-02-29").AddYears(1),
		},
		{
			"1999-12-25",
			MustFromString("1999-12-26").AddDays(-1),
		},
		{
			"1999-12-25",
			MustFromString("1999-12-24").AddDays(1),
		},
		{
			"2022-11-01",
			MustFromString("2022-11-15").StartOfMonth(),
		},
		{
			"2022-12-01",
			MustFromString("2022-12-15").StartOfMonth(),
		},
		{
			"2022-01-01",
			MustFromString("2022-01-15").StartOfMonth(),
		},
		{
			"2022-11-30",
			MustFromString("2022-11-01").EndOfMonth(),
		},
		{
			"2022-12-31",
			MustFromString("2022-12-15").EndOfMonth(),
		},
		{
			"2022-01-31",
			MustFromString("2022-01-31").EndOfMonth(),
		},
		{
			"30",
			MustFromString("2022-11-01").DaysInMonth(),
		},
		{
			"31",
			MustFromString("2022-12-15").DaysInMonth(),
		},
		{
			"31",
			MustFromString("2022-01-31").DaysInMonth(),
		},
		{
			"2022-01-01",
			MustFromString("2022-01-15").StartOfQuarter(),
		},
		{
			"2022-01-01",
			MustFromString("2022-02-15").StartOfQuarter(),
		},
		{
			"2023-01-01",
			MustFromString("2022-12-15").StartOfNextQuarter(),
		},
		{
			"2022-07-01",
			MustFromString("2022-06-15").StartOfNextQuarter(),
		},
		{
			"2022-07-01",
			MustFromString("2022-04-22").StartOfNextQuarter(),
		},

		{"4", MustFromString("2015-03-04").Day()},
		{"March", MustFromString("2015-03-04").Month()},
		{"2015", MustFromString("2015-03-04").Year()},
		{"63", MustFromString("2015-03-04").YearDay()},

		{
			time.Now().Format("2006-01-02"),
			Today(),
		},
		{
			time.Now().Add(-24 * time.Hour).Format("2006-01-02"),
			Yesterday(),
		},
		{
			time.Now().Add(24 * time.Hour).Format("2006-01-02"),
			Tomorrow(),
		},

		{
			time.Now().In(Sydney).Format("2006-01-02"),
			TodayIn(Sydney),
		},
		{
			time.Now().In(Sydney).Add(-24 * time.Hour).Format("2006-01-02"),
			YesterdayIn(Sydney),
		},
		{
			time.Now().In(Sydney).Add(24 * time.Hour).Format("2006-01-02"),
			TomorrowIn(Sydney),
		},

		{"Wednesday", MustFromString("1989-06-14").Weekday()},
		{"Thursday", MustFromString("2014-12-25").Weekday()},
		{"Saturday", MustFromString("2018-08-18").Weekday()},

		{
			"2018-05-26",
			New(2018, time.May, 26),
		},

		{
			"2018-05-05",
			Max(MustFromString("2018-05-05"), MustFromString("2018-01-01")),
		},
		{
			"2018-05-05",
			Max(MustFromString("2018-01-01"), MustFromString("2018-05-05")),
		},
		{
			"2018-01-01",
			Min(MustFromString("2018-05-05"), MustFromString("2018-01-01")),
		},
		{
			"2018-01-01",
			Min(MustFromString("2018-01-01"), MustFromString("2018-05-05")),
		},

		{
			"366",
			MustFromString("2024-01-01").DaysInYear(),
		},
		{
			"366",
			MustFromString("2000-01-01").DaysInYear(),
		},
		{
			"365",
			MustFromString("2100-01-01").DaysInYear(),
		},
		{
			"365",
			MustFromString("2023-01-01").DaysInYear(),
		},
	} {
		if gotStr := fmt.Sprintf("%v", test.got); gotStr != test.want {
			t.Errorf("i=%d got=%v want=%v", i, gotStr, test.want)
		}
	}
}

func TestFromStringErr(t *testing.T) {
	if _, err := FromString("not a date"); err == nil {
		t.Error("expected error")
	}
}

func TestJSON(t *testing.T) {
	type J struct {
		D Date `json:"d"`
	}
	j := J{MustFromString("2015-05-21")}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(j)
	if err != nil {
		t.Fatalf("Could not encode: %v", err)
	}
	if want := `{"d":"2015-05-21"}`; strings.TrimSpace(buf.String()) != want {
		t.Fatalf("Did not encode correctly, want=%q got=%q",
			want, strings.TrimSpace(buf.String()))
	}

	j = J{}
	err = json.NewDecoder(&buf).Decode(&j)
	if err != nil {
		t.Fatalf("Could not decode: %v", err)
	}
	if want := MustFromString("2015-05-21"); j.D != want {
		t.Fatalf("Did not decode correctly, want=%q got=%q", want, j.D)
	}
}

func TestSQLScan(t *testing.T) {
	var d Date
	if err := d.Scan(time.Date(2013, time.July, 13, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	if d != MustFromString("2013-07-13") {
		t.Fatalf("Got=%v", d)
	}
}

func TestSQLValue(t *testing.T) {
	v, err := MustFromString("2013-07-13").Value()
	if err != nil {
		t.Fatal(err)
	}
	if v != "2013-07-13" {
		t.Fatalf("Got=%v", v)
	}
}

func TestNullDate(t *testing.T) {
	d := MustFromString("2013-07-13")
	nd := NewNullDate(d)
	if d != nd.Date {
		t.Fatalf("Got=%v", nd.Date)
	}
	if !nd.Valid {
		t.Fatalf("Got=%v", nd.Valid)
	}
}

func TestNullDateSQLScanDate(t *testing.T) {
	var nd NullDate
	if err := nd.Scan(time.Date(2013, time.July, 13, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	if nd.Date != MustFromString("2013-07-13") {
		t.Fatalf("Got=%v", nd.Date)
	}
	if !nd.Valid {
		t.Fatalf("Got=%v", nd.Valid)
	}
}

func TestNullDateSQLScanNull(t *testing.T) {
	var nd NullDate
	if err := nd.Scan(nil); err != nil {
		t.Fatal(err)
	}
	if nd.Valid {
		t.Fatalf("Got=%v", nd.Valid)
	}
}

func TestNullDateSQLValueDate(t *testing.T) {
	v, err := NewNullDate(MustFromString("2013-07-13")).Value()
	if err != nil {
		t.Fatal(err)
	}
	if v != "2013-07-13" {
		t.Fatalf("Got=%v", v)
	}
}

func TestNullDateSQLValueNull(t *testing.T) {
	v, err := NullDate{}.Value()
	if err != nil {
		t.Fatal(err)
	}
	if v != nil {
		t.Fatalf("Got=%v", v)
	}
}

func TestNullDateJSONDate(t *testing.T) {
	type J struct {
		Nd NullDate `json:"d"`
	}
	j := J{NewNullDate(MustFromString("2015-05-21"))}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(j)
	if err != nil {
		t.Fatalf("Could not encode: %v", err)
	}
	if want := `{"d":"2015-05-21"}`; strings.TrimSpace(buf.String()) != want {
		t.Fatalf("Did not encode correctly, want=%q got=%q",
			want, strings.TrimSpace(buf.String()))
	}

	j = J{}
	err = json.NewDecoder(&buf).Decode(&j)
	if err != nil {
		t.Fatalf("Could not decode: %v", err)
	}
	if want := MustFromString("2015-05-21"); j.Nd.Date != want {
		t.Fatalf("Did not decode correctly, want=%q got=%q", want, j.Nd.Date)
	}
	if !j.Nd.Valid {
		t.Fatalf("Did not decode correctly, should be valid")
	}
}

func TestNullDateJSONNull(t *testing.T) {
	type J struct {
		Nd NullDate `json:"d"`
	}
	j := J{NullDate{}}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(j)
	if err != nil {
		t.Fatalf("Could not encode: %v", err)
	}
	if want := `{"d":null}`; strings.TrimSpace(buf.String()) != want {
		t.Fatalf("Did not encode correctly, want=%q got=%q",
			want, strings.TrimSpace(buf.String()))
	}

	j = J{}
	err = json.NewDecoder(&buf).Decode(&j)
	if err != nil {
		t.Fatalf("Could not decode: %v", err)
	}
	if j.Nd.Valid {
		t.Fatalf("Did not decode correctly, should be not valid")
	}
}

func TestDiff(t *testing.T) {

	type want struct {
		year  int
		month int
		day   int
	}

	tests := map[string]struct {
		d1   Date
		d2   Date
		want want
	}{
		"1 year diff": {
			d1: MustFromString("2001-01-01"),
			d2: MustFromString("2002-01-01"),
			want: want{
				year:  1,
				month: 12,
				day:   365,
			},
		},
		"1 year diff leap year": {
			d1: MustFromString("2000-01-01"),
			d2: MustFromString("2001-01-01"),
			want: want{
				year:  1,
				month: 12,
				day:   366,
			},
		},
		"1 month diff": {
			d1: MustFromString("2000-01-01"),
			d2: MustFromString("2000-02-01"),
			want: want{
				year:  0,
				month: 1,
				day:   31,
			},
		},
		"1 day diff": {
			d1: MustFromString("2000-01-01"),
			d2: MustFromString("2000-01-02"),
			want: want{
				year:  0,
				month: 0,
				day:   1,
			},
		},
		"1 year, 1 month, 1 day": {
			d1: MustFromString("2001-02-01"),
			d2: MustFromString("2002-03-01"),
			want: want{
				year:  1,
				month: 13,
				day:   393,
			},
		},
		"negative diff": {
			d1: MustFromString("2002-01-01"),
			d2: MustFromString("2001-01-01"),
			want: want{
				year:  1,
				month: 12,
				day:   365,
			},
		},
		"equal dates": {
			d1: MustFromString("2000-01-01"),
			d2: MustFromString("2000-01-01"),
			want: want{
				year:  0,
				month: 0,
				day:   0,
			},
		},
		"february": {
			d1: MustFromString("2001-02-01"),
			d2: MustFromString("2001-03-01"),
			want: want{
				year:  0,
				month: 1,
				day:   28,
			},
		},
		"february leap year": {
			d1: MustFromString("2000-02-01"),
			d2: MustFromString("2000-03-01"),
			want: want{
				year:  0,
				month: 1,
				day:   29,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			y, m, d := Diff(tt.d1, tt.d2)
			if y != tt.want.year || m != tt.want.month || d != tt.want.day {
				t.Fatalf("Got Year: %d, Month: %d, Day: %d, want Year: %d, Month: %d, Day: %d", y, m, d, tt.want.year, tt.want.month, tt.want.day)
			}

		})
	}
}
