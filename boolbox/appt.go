// this is where butch keeps his appointment system.
package boolbox

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
)

// errors
var InvalidDateError = errors.New("Invalid date!")
var InvalidTimeError = errors.New("Invalid time!")

// stuct to hold a date
type Date struct {
	Day   int        `json:"day,omitempty"`
	Month time.Month `json:"month,omitempty"`
	Year  int        `json:"year,omitempty"`
	Ud    bool       `json:"undetermined,omitempty"` // whether the date is undetermined.
}

// struct to hold a time
type Time struct {
	Time [2][2]int `json:"time,omitempty"`
	Ud   bool      `json:"undetermined,omitempty"` // whether the time is undetermined
}

// struct to hold an rsvp
type Rsvp struct {
	User     discord.User `json:"user,omitempty"`
	Time     `json:"pu_time,omitempty"`
	Pickedup bool `json:"picked_up,omitempty"`
}

// struct to hold an appointment
type Appointment struct {
	Name  string `json:"name,omitempty"`
	Date  Date   `json:"date,omitempty"`
	*Time `json:"time,omitempty"`
	Desc  string `json:"desc,omitempty"`
	Resv  []Rsvp `json:"resv,omitempty"`
}

// appointment wrapper for json
type Appointments struct {
	Appts []Appointment `json:"appts,omitempty"`
}

// remove an item from an array of Rsvp
func RemoveRsvp(s []Rsvp, i int) []Rsvp {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// remove an item from an array of Appointment
func RemoveAppointment(s []Appointment, i int) []Appointment {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// make a Date out of xx/xx(/xx)
func MakeDate(str string) (Date, error) {
	date := Date{}
	pstr := strings.Split(str, "/")
	if len(pstr) == 2 || len(pstr) == 3 {
		m, err1 := strconv.Atoi(pstr[0])
		d, err2 := strconv.Atoi(pstr[1])
		if err1 != nil || err2 != nil || m < 1 || m > 12 || d < 1 || d > 31 {
			return Date{}, InvalidDateError
		}

		if len(pstr) == 3 {
			y, err := strconv.Atoi(pstr[2])
			if err != nil {
				return Date{}, InvalidDateError
			}

			date.Year = y
		} else {
			y, _, _ := time.Now().Clock()
			date.Year = y
		}

		date.Month = time.Month(m)
		date.Day = d

		return date, nil
	}

	return Date{}, InvalidTimeError
}

// make a Time out of xx:xx
func MakeTime(str string) (Time, error) {
	time := [2][2]int{}
	pstr := strings.Split(str, ":")
	if len(pstr) == 2 {
		spstr1 := []byte(pstr[0])
		spstr2 := []byte(pstr[1])

		if (len(spstr1) == 1 || len(spstr1) == 2) && len(spstr2) == 2 {
			if len(spstr1) == 2 {
				tsp1, err1 := strconv.Atoi(string(spstr1[0]))
				tsp2, err2 := strconv.Atoi(string(spstr1[1]))
				if err1 != nil || err2 != nil {
					return Time{}, InvalidTimeError
				}

				time[0][0] = tsp1
				time[0][1] = tsp2
			} else {
				tsp1, err1 := strconv.Atoi(string(spstr1[0]))
				if err1 != nil {
					return Time{}, InvalidTimeError
				}

				time[0][1] = tsp1
			}

			tsp3, err3 := strconv.Atoi(string(spstr2[0]))
			tsp4, err4 := strconv.Atoi(string(spstr2[1]))
			if err3 != nil || err4 != nil {
				return Time{}, InvalidTimeError
			}

			time[1][0] = tsp3
			time[1][1] = tsp4

			return Time{Time: time}, nil
		}

	}

	return Time{}, InvalidTimeError
}

func BuildDate(date Date) string {
	if date.Ud {
		return "n/a"
	}

	return date.Month.String() + " " + strconv.Itoa(date.Day) + ", " + strconv.Itoa(date.Year)
}

func BuildTime(time *Time) string {
	if time.Ud {
		return "n/a"
	}

	return strconv.Itoa(time.Time[0][0]) + strconv.Itoa(time.Time[0][1]) + ":" + strconv.Itoa(time.Time[1][0]) + strconv.Itoa(time.Time[1][1])
}

func BuildApptDesc(appointment Appointment) string {
	return "**Time: " + BuildTime(appointment.Time) + "\nDate: " + BuildDate(appointment.Date) + "**\n" + HelpDivider + appointment.Desc + "\n"
}
