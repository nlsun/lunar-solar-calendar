package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nlsun/lunar-solar-calendar/lunarsolar"
)

// TODO: Take traditional chinese birth date, convert to birthday of this
// Gregorian year.
// TODO: Generate google calendar for a person, notifications configurable.

type lunarBirthdayForYearRequest struct {
	LunarBirthDate time.Time `json:"lunar_birth_date"`
	IsLeapMonth    bool      `json:"is_leap_month"`
	Year           int       `json:"year"`
}

type lunarBirthdayForYearResponse struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	port := os.Getenv("PORT")
	log.Printf("Listening to port: %s", port)

	assetDir := "assets"

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        mkHandler(assetDir),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func mkHandler(assetDir string) *http.ServeMux {
	sv := http.NewServeMux()
	sv.HandleFunc("/", handleStaticFile(assetDir+"/html/index.html"))
	sv.HandleFunc("/assets/js/script.js", handleStaticFile(assetDir+"/js/script.js"))
	sv.HandleFunc("/api/v1/lunar-birthday-for-year/", handlelunarBirthdayForYear)
	return sv
}

func handleStaticFile(file string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			writeHttpErr(w, http.StatusInternalServerError)
			log.Print(err)
			return
		}
		if _, err := w.Write(b); err != nil {
			log.Print(err)
		}
	}
}

func handlelunarBirthdayForYear(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeHttpErr(w, http.StatusMethodNotAllowed)
		return
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		writeHttpErr(w, http.StatusInternalServerError)
		log.Print(err)
		return
	}

	var reqBody lunarBirthdayForYearRequest
	if err := json.Unmarshal(b, &reqBody); err != nil {
		writeHttpErr(w, http.StatusInternalServerError)
		log.Printf("%s: %s", string(b), err)
		return
	}

	lunarBirthday := lunarsolar.NewLunarTime(reqBody.LunarBirthDate, reqBody.IsLeapMonth)
	birthday, err := lunarBirthdayForYear(lunarBirthday, reqBody.Year)
	if err != nil {
		writeHttpErr(w, http.StatusBadRequest)
		log.Print(err)
		return
	}

	resp := lunarBirthdayForYearResponse{
		Year:  birthday.Year(),
		Month: int(birthday.Month()),
		Day:   birthday.Day(),
	}
	b, err = json.Marshal(resp)
	if err != nil {
		writeHttpErr(w, http.StatusInternalServerError)
		log.Print(err)
		return
	}

	if _, err := w.Write(b); err != nil {
		log.Print(err)
	}
}

// Computes the solar birthday given a birth date as a lunar date and a solar
// year to calculate for.
//
// If the birth date is a leap month, and the target year does not leap that
// month, treat it as a non-leap-month birth date.
func lunarBirthdayForYear(birthDate lunarsolar.LunarTime, solarYear int) (time.Time, error) {
	lunarBirthYear := birthDate.Time().Year()
	if lunarBirthYear > solarYear {
		return time.Time{}, fmt.Errorf("birth year %d can't be greater than input year %d", lunarBirthYear, solarYear)
	}

	yearDiff := solarYear - lunarBirthYear
	lunarBirthday := birthDate.AddDate(yearDiff, 0, 0)

	// If born on a leap month, but this year that month is not a leap month,
	// treat it as the normal month.
	if birthDate.IsLeap() && !lunarsolar.IsLunarLeapMonthPossible(lunarBirthday) {
		return lunarsolar.LunarToSolar(lunarBirthday.AsLeap(false)), nil
	}
	return lunarsolar.LunarToSolar(lunarBirthday), nil
}

func writeHttpErr(w http.ResponseWriter, code int) {
	errResp := errorResponse{Error: http.StatusText(code)}
	b, err := json.Marshal(errResp)
	if err != nil {
		http.Error(w, "{}", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	http.Error(w, string(b), code)
}
