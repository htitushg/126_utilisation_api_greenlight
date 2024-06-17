package utils

import (
	"log"
	"net"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

/*
var (
	_, b, _, _ = runtime.Caller(0)
	Path       = filepath.Dir(filepath.Dir(filepath.Dir(b))) + "/"
)
*/
// durationToString -> just for fun ;)
func durationToString(d time.Duration) string {
	var hours, minutes string
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h < 10 {
		hours += "0"
	}
	if m < 10 {
		minutes += "0"
	}
	hours += strconv.Itoa(h)
	minutes += strconv.Itoa(m)
	return hours + "H" + minutes
}

// SetDailyTimer sets a waiting time to match a certain `hour`.
func SetDailyTimer(hour int) time.Duration {
	hour = hour % 24
	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day(), hour, 0, 0, 0, t.Location())
	d := n.Sub(t)
	if d < 0 {
		n = n.Add(24 * time.Hour)
		d = n.Sub(t)
	}
	log.Println("SetDailyTimer() value: ", durationToString(d), "until", n.Format("02 Jan 15H04")) // verbose
	return d
}

func GetIP(r *http.Request) string {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		// get last IP in list since ELB prepends other user defined IPs, meaning the last one is the actual client IP.
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String()
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Fatalln(err)
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1"
		}
		return ip
	}

	log.Fatalln(err)
	return ""
}

func GetCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func CheckEmail(email string) bool {
	reg := regexp.MustCompile(`^[\w&#$.%+-]+@[\w&#$.%+-]+\.[a-z]{2,6}?$`)
	return reg.MatchString(email)
}

// CheckPasswd
// checks if the password's format is according to the rules.
func CheckPasswd(passwd string) bool {

	// Matches any password containing at least one digit, one lowercase,
	// one uppercase, one symbol and 8 characters in total.
	//regex := regexp.MustCompile(`^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*([^\w\s]|_)).{8,}$`) // Alas not supported by the regexp library
	digit := regexp.MustCompile(`\d+`)
	lower := regexp.MustCompile(`[a-z]+`)
	upper := regexp.MustCompile(`[A-Z]+`)
	symbol := regexp.MustCompile(`([^\w\s]|_)+`)
	minLen := regexp.MustCompile(`.{8,}`)
	return digit.MatchString(passwd) && lower.MatchString(passwd) && upper.MatchString(passwd) && symbol.MatchString(passwd) && minLen.MatchString(passwd)
}
func CheckPseudo(pseudo string) bool {
	reg := regexp.MustCompile(`^[A-Za-z0-9]{5,}$`)
	return reg.MatchString(pseudo)
}

// Fonction qui vérifie si le champ name contient une des valeurs de sl
// fuction to check given string is in array or not
func Contains(sl []string, name string) bool {
	// iterate over the array and compare given string to each element
	for _, value := range sl {
		// log.Printf("Fonction miscellaneous.Contains : name = '%v', value = %v\n", name, value) //Testing
		if value == name {
			return true
		}
	}
	return false
}
