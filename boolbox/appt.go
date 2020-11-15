// this is where butch keeps his appointment system.
package boolbox

import (
	"errors"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
)

// rsvp struct used to keep track
// of pickup times and discord.user's
type Rsvp struct {
	User     discord.User `json:"user,omitempty"`
	PuTime   string       `json:"puTime,omitempty"`
	Pickedup bool         `json:"pickedup,omitempty"`
}

// appointment struct
type Appointment struct {
	Name string `json:"name,omitempty"`
	Date string `json:"date,omitempty"`
	Time string `json:"time,omitempty"`
	Desc string `json:"desc,omitempty"`
	Resv []Rsvp `json:"resv,omitempty"`
}

// appointment wrapper for json
type Appointments struct {
	Appts []Appointment `json:"appts,omitempty"`
}

// remove an item from an array of interface{}
func (box *Box) RemoveRsvp(s []Rsvp, i int) []Rsvp {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// remove an item from an array of Appointment
func (box *Box) RemoveAppointment(s []Appointment, i int) []Appointment {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// check if a date is valid
func (box *Box) CheckDate(input string) error {
	pDate := strings.Split(input, "/")
	if len(pDate) >= 2 {
		_, firstErr := strconv.Atoi(pDate[0])
		_, secondErr := strconv.Atoi(pDate[1])
		if firstErr == nil || secondErr == nil {
			return nil
		}
	}
	if len(pDate) == 3 {
		_, thirdErr := strconv.Atoi(pDate[2])
		if thirdErr != nil {
			return nil
		}
	}

	return errors.New("Invalid date")
}

// check if a time is valid.
func (box *Box) CheckTime(input string) error {
	pTime := strings.Split(input, ":")
	if len(pTime) == 2 {
		_, firstErr := strconv.Atoi(pTime[0])
		_, secondErr := strconv.Atoi(pTime[1])
		if firstErr == nil || secondErr == nil {
			return nil
		}
	}

	return errors.New("Invalid time")
}

// build appointment description
func (box *Box) BuildApptDesc(appointment Appointment) string {
	var desc strings.Builder

	desc.WriteString("**Time: ")
	desc.WriteString(appointment.Time)
	desc.WriteString("\nDate: ")
	desc.WriteString(appointment.Date)
	desc.WriteString("**\n")
	desc.WriteString(HelpDivider)
	desc.WriteString(appointment.Desc)
	desc.WriteString("\n")

	return desc.String()
}
