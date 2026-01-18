package assist

import (
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"

	"bitbucket.org/lewington/erosai-server/globals"
	"bitbucket.org/lewington/erosai-server/lg"
)

var AEST = time.FixedZone("AEST", -60*60*4)
var TimeLocation, _ = time.LoadLocation("Australia/Melbourne")
var BetfairLocation, _ = time.LoadLocation("GMT")

func IPIsValid(IP string) bool {
	parts := strings.Split(IP, ".")

	if len(parts) != 4 {
		return false
	}

	for _, x := range parts {
		if i, err := strconv.Atoi(x); err == nil {
			if i < 0 || i > 255 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// Time related

func Wait(milliseconds int) {
	if milliseconds > 0 {
		time.Sleep(time.Millisecond * time.Duration(milliseconds))
	}
}

func MedWait() {
	Wait(500)
}

// Typecasters

func ToInt(str string) int {
	i, err := strconv.Atoi(str)
	Check(err)
	return i
}

func IsNullValue(str string) bool {
	return str == "" || str == "-"
}

func ToIntOrNegOne(str string) int {
	if IsNullValue(str) {
		return -1
	}

	return ToInt(str)
}

func FloatToIntOrNegOne(str string) int {
	if IsNullValue(str) {
		return -1
	}

	return FloatStringToInt(str)
}

func ToBool(str string) bool {
	if str == "1" {
		return true
	}
	if str != "0" {
		panic(fmt.Sprintf("Unexpected String Input: %v", str))
	}
	return false
}

// Because floats are bad we convert them all into integers
// preserving 2 decimal places
func FloatStringToInt(str string) int {
	f, err := strconv.ParseFloat(str, 64)
	Check(err)

	return int(math.Round(f * 100))
}

func FloatStringToIntSafe(str string) (int, error) {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}

	return FloatToIntRep(f), nil
}

func FloatToIntRep(value float64) int {
	return int(math.Round(value * 100))
}

func PercentageToInt(percentage string) int {
	return ToIntOrNegOne(strings.Trim(percentage, "%"))
}

func IntToFloat(value int) float64 {
	return float64(float64(value) / 100.0)
}

func IntToFloatString(value int) string {
	return TwoDecimalPlaces(IntToFloat(value))
}
func IntToNoDecimalNumber(value int) string {
	return strconv.Itoa(value / 100)
}

func FloatToString(value float64) string {
	return TwoDecimalPlaces(value / 100)
}

func NamesMatch(a string, b string) bool {
	return TrimmedLower(a) == TrimmedLower(b)
}

func TwoDecimalPlaces(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func ToFloat(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	Check(err)
	return f
}

func RoundedInt(value int) int {
	return int(math.Round((float64(value) / 100.0)) * 100.0)
}

func ToFloatSafe(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func ToTwoDigits(number int) string {
	return fmt.Sprintf("%02d", number)
}

// Identifiers

func MeetingID(meetingNumber int, date time.Time) int {
	val, err := strconv.Atoi(strconv.Itoa(TimeInt(date)) + strconv.Itoa(meetingNumber))
	Check(err)
	return val
}

func ContestantIDFromEvent(horseName string, EventID string) string {
	return fmt.Sprintf("%v-%v", horseName, EventID)
}

func ContestantIdentifier(horseName string, meetingID int, raceNumber int) string {
	return fmt.Sprintf("%v-%v", horseName, EventID(meetingID, raceNumber))
}

func EventID(meetingID int, raceNumber int) string {
	return fmt.Sprintf("%v_%02d", meetingID, raceNumber)
}

func Flattened(values [][]interface{}) []interface{} {
	final := make([]interface{}, 10000)
	for _, list := range values {
		for _, value := range list {
			final = append(final, value)
		}
	}

	return final
}

func FlattenedInt(values [][]int) []int {
	final := make([]int, 10000)
	for _, list := range values {
		for _, value := range list {
			final = append(final, value)
		}
	}

	return final
}

// Slice manipulation
func Contains(haystack []string, needle string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}

	return false
}

func MinValue(haystack []int) int {
	min := haystack[0]
	for _, v := range haystack {
		if v < min {
			min = v
		}
	}

	return min
}

func IntSliceContains(haystack []int, needle int) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}

	return false
}

// NZVenue returns whether the given venue looks like
// it's a venue in New Zeeland.
func NZVenue(venueName string) bool {
	return Contains(globals.NZVenues, venueName)
}

func Timestamp() time.Time {
	return time.Now().In(TimeLocation)
}

func ToBetfair(localTime time.Time) time.Time {
	return localTime.In(BetfairLocation)
}

func TimeInt(date time.Time) int {
	year, month, day := date.Date()
	val, err := strconv.Atoi(strconv.Itoa(year) + strconv.Itoa(int(month)) + strconv.Itoa(day))
	Check(err)
	return val
}

func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, TimeLocation)
}

func StartOfDay(date time.Time) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, TimeLocation)
}

func SameDay(a time.Time, b time.Time) bool {
	year, month, day := a.Date()
	yearB, monthB, dayB := b.Date()
	return year == yearB && month == monthB && day == dayB
}
func SameHour(a time.Time, b time.Time) bool {
	return SameDay(a, b) && a.Hour() == b.Hour()
}

func EndOfDay(date time.Time) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, TimeLocation)
}

func UTCUnixDayStart(date time.Time) int {
	return int(StartOfDay(date).Unix()+39600) - 3600
}

func FutureTime() time.Time {
	time, err := time.Parse("02-01-2006", "01-12-2041")
	Check(err)
	return time
}

func SetTimestamp(seconds int) time.Time {
	return time.Date(2009, 11, 17, 20, 34, seconds, 651387237, time.UTC)
}

func RecordTime(startTime time.Time, action string) {
	if Debug() {
		lg.L.Debug("%v took %v", action, Timestamp().Sub(startTime))
	}

}

func PathToPackage() string {
	return os.Getenv("GOPATH") + "/src/bitbucket.org/lewington/erosai-server"
}

func PathToAsset() string {
	return PathToPackage() + "/public/static/assets/"
}

func ToPlaceName(name string) string {
	return fmt.Sprintf("p-%v", name)
}

// Maths

// FloatRepMultiply assumes that both integers are 100X
// larger than they should be so divides their
// product by 100 before returning.
func FloatRepMultiply(a int, b int) int {
	return a * b / 100
}

func Min(a int, b int) int {
	if a > b {
		return b
	}

	return a
}

func FloatRepDivide(a int, b int) int {
	return (a * 100 / b)
}

func Rounded(amount int, increment int) int {
	remainder := amount % increment

	return amount - remainder
}

func Profit(price int) int {
	return price - 100
}

func ToWinString(win int, odd int) string {
	return TwoDecimalPlaces(float64(ToWin(win, odd)) / 100)
}

func ToWin(win int, odd int) int {
	return (win / (odd - 100)) * 100
}

func Max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

// WithoutIndex returns the same slice with the element
// at the given index removed.
func WithoutIndex(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

func Debug() bool {
	return globals.LogLevel == "debug"
}
func Panicf(message string, args ...interface{}) {
	panic(fmt.Sprintf(message, args...))
}

func RaceOpenStatus(status string) bool {
	switch status {
	case "CLOSED", "FINAL", "ABANDONED":
		return false
	default:
		return true
	}
}

// Hashing

func Hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}

// Reading files
func BytesFromFile(filename string) []byte {
	jsonFile, err := os.Open(filename)
	Check(err)
	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	Check(err)
	return bytes
}

func SafeBytes(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func StrictBytes(body io.ReadCloser) []byte {
	defer body.Close()
	bytes, err := ioutil.ReadAll(body)
	Check(err)

	return bytes
}

// String manipulation

func TrimmedLower(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func PunctRemoved(word string) string {
	newStr := ""
	for _, r := range word {
		if !unicode.IsPunct(r) {
			newStr += string(r)
		}
	}

	return newStr
}

func Joined(value string) string {
	return strings.Replace(value, " ", "-", -1)
}

func UsableIPs() []string {
	ips := []string{}
	ifaces, err := net.Interfaces()
	Check(err)
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		Check(err)
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			ipStr := ip.String()

			if strings.Contains(ipStr, "192.168") || strings.Contains(ipStr, "10.0.0") {
				fmt.Println("using internal ip: ", ipStr)
				ips = append(ips, ipStr)
			}
		}
	}

	return ips
}

func PrintGoquery(s *goquery.Selection) {
	fmt.Println(goquery.OuterHtml(s))
}

func ContainsBookie(haystack []globals.BookieName, needle globals.BookieName) bool {
	for _, name := range haystack {
		if name == needle {
			return true
		}
	}

	return false
}
