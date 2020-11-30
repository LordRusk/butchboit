// this is where butch keeps his appointment system.
package boolbox

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
)

var InvalidDateError = errors.New("Invalid date!")
var InvalidTimeError = errors.New("Invalid time!")

type Date struct {
	Day   int        `json:"day,omitempty"`
	Month time.Month `json:"month,omitempty"`
	Year  int        `json:"year,omitempty"`
	Ud    bool       `json:"undetermined,omitempty"` // whether the date is undetermined.
}

type Time struct {
	Time [2][2]int `json:"time,omitempty"`
	Ud   bool      `json:"undetermined,omitempty"` // whether the time is undetermined
}

// rsvp struct used to keep track
// of pickup times and discord.user's
type Rsvp struct {
	User     discord.User `json:"user,omitempty"`
	*Time    `json:"pu_time,omitempty"`
	Pickedup bool `json:"picked_up,omitempty"`
}

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

// check if a date is valid, and if it is, return a Date
func CheckMakeDate(input string) (Date, error) {
	date := Date{}
	pDate := strings.Split(input, "/")

	if len(pDate) == 2 || len(pDate) == 3 {
		mInt, err1 := strconv.Atoi(pDate[0])
		dInt, err2 := strconv.Atoi(pDate[1])
		if err1 != nil || err2 != nil {
			return date, InvalidDateError
		}

		if mInt < 1 || mInt > 12 || dInt < 1 || dInt > 31 {
			return date, InvalidDateError
		}

		if len(pDate) == 3 {
			yInt, err := strconv.Atoi(pDate[2])
			if err != nil {
				return date, InvalidDateError
			}

			date.Year = yInt
		} else {
			year, _, _ := time.Now().Clock()
			date.Year = year
		}

		date.Month = time.Month(mInt)
		date.Day = dInt

		return date, nil
	}

	return date, InvalidTimeError
}

// check if a time is valid, and if it is, return
// a corrently formatted time.
func CheckMakeTime(input string) (*Time, error) {
	time := [2][2]int{}

	pTime := strings.Split(input, ":")

	if len(pTime) == 2 {
		spTime1 := []byte(pTime[0])
		spTime2 := []byte(pTime[1])

		if (len(spTime1) == 1 || len(spTime1) == 2) && len(spTime2) == 2 {
			if len(spTime1) == 2 {
				tsp1, err1 := strconv.Atoi(string(spTime1[0]))
				tsp2, err2 := strconv.Atoi(string(spTime1[1]))
				if err1 != nil || err2 != nil {
					return &Time{}, InvalidTimeError
				}

				time[0][0] = tsp1
				time[0][1] = tsp2
			} else {
				tsp1, err1 := strconv.Atoi(string(spTime1[0]))
				if err1 != nil {
					return &Time{}, InvalidTimeError
				}

				time[0][1] = tsp1
			}

			tsp3, err3 := strconv.Atoi(string(spTime2[0]))
			tsp4, err4 := strconv.Atoi(string(spTime2[1]))
			if err3 != nil || err4 != nil {
				return &Time{}, InvalidTimeError
			}

			time[1][0] = tsp3
			time[1][1] = tsp4

			return &Time{Time: time}, nil
		}

	}

	return &Time{}, InvalidTimeError
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

// build appointment description
func BuildApptDesc(appointment Appointment) string {
	return "**Time: " + BuildTime(appointment.Time) + "\nDate: " + BuildDate(appointment.Date) + "**\n" + HelpDivider + appointment.Desc + "\n"
}
