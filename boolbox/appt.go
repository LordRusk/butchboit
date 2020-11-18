// this is where butch keeps his appointment system.
package boolbox

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/discord"
)

var InvalidDateError = errors.New("Invalid date!")
var InvalidTimeError = errors.New("Invalid time!")

type Date struct {
	Day   int
	Month time.Month
	Year  int
}

// rsvp struct used to keep track
// of pickup times and discord.user's
type Rsvp struct {
	User     discord.User `json:"user,omitempty"`
	PuTime   [2]int       `json:"pu_time,omitempty"`
	Pickedup bool         `json:"picked_up,omitempty"`
}

type Appointment struct {
	Name string `json:"name,omitempty"`
	Date Date   `json:"date,omitempty"`
	Time [2]int `json:"time,omitempty"`
	Desc string `json:"desc,omitempty"`
	Resv []Rsvp `json:"resv,omitempty"`
}

// appointment wrapper for json
type Appointments struct {
	Appts []Appointment `json:"appts,omitempty"`
}

// remove an item from an array of Rsvp
func (box *Box) RemoveRsvp(s []Rsvp, i int) []Rsvp {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// remove an item from an array of Appointment
func (box *Box) RemoveAppointment(s []Appointment, i int) []Appointment {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// check if a date is valid, and if it is, return a Date
func (box *Box) CheckMakeDate(input string) (Date, error) {
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
func (box *Box) CheckMakeTime(input string) ([2]int, error) {
	time := [2]int{}

	pTime := strings.Split(input, ":")
	if len(pTime) == 2 {
		tp1, err1 := strconv.Atoi(pTime[0])
		tp2, err2 := strconv.Atoi(pTime[1])
		if err1 != nil || err2 != nil {
			return time, InvalidTimeError
		}

		time[0] = tp1
		time[1] = tp2

		return time, nil
	}

	return time, InvalidTimeError
}

func (box *Box) BuildDate(date Date) string {
	return date.Month.String() + " " + strconv.Itoa(date.Day) + ", " + strconv.Itoa(date.Year)
}

func (box *Box) BuildTime(time [2]int) string {
	return strconv.Itoa(time[0]) + ":" + strconv.Itoa(time[1])
}

// build appointment description
func (box *Box) BuildApptDesc(appointment Appointment) string {
	var desc strings.Builder

	desc.WriteString("**Time: ")
	desc.WriteString(box.BuildTime(appointment.Time))
	desc.WriteString("\nDate: ")
	desc.WriteString(box.BuildDate(appointment.Date))
	desc.WriteString("**\n")
	desc.WriteString(HelpDivider)
	desc.WriteString(appointment.Desc)
	desc.WriteString("\n")

	return desc.String()
}
