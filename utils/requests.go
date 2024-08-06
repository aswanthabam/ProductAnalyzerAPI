package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
)

const (
	IP_API_URL    = "http://ip-api.com/json/"
	IP_API_PARAMS = "?fields=status,continent,continentCode,country,countryCode,region,regionName,city,district,zip,lat,lon,timezone,offset,currency,isp,org,as,asname,mobile,proxy,hosting,query"
)

type IPInfoResponse struct {
	Status        string  `json:"status"`
	IP            string  `json:"query"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continentCode"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"countryCode"`
	Region        string  `json:"region"`
	RegionName    string  `json:"regionName"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	Zip           string  `json:"zip"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	Timezone      string  `json:"timezone"`
	Currency      string  `json:"currency"`
	ISP           string  `json:"isp"`
	Org           string  `json:"org"`
	AS            string  `json:"as"`
	ASName        string  `json:"asname"`
	Mobile        bool    `json:"mobile"`
	Proxy         bool    `json:"proxy"`
	Hosting       bool    `json:"hosting"`
}

type UserAgentDetails struct {
	DeviceType string `json:"device_type"`
	OS         string `json:"os"`
	Browser    string `json:"browser"`
}

// GetIPAddress returns the IP address of the client making the request
func GetIPAddress(c *gin.Context) (string, error) {
	ip := c.ClientIP()
	forwardedFor := c.Request.Header.Get("x-forwarded-for")
	if forwardedFor != "" {
		ip = forwardedFor
	}
	if ip == "" {
		return "", errors.New("could not get ip address")
	}
	return ip, nil
}

// GetIPAddressInfo returns information about the given IP address
func GetIPAddressInfo(ip string) (*IPInfoResponse, error) {
	url := IP_API_URL + ip + IP_API_PARAMS
	var response IPInfoResponse
	if err := SendGetRequest(url, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetUserAgentDetails returns information about the given user agent
func GetUserAgentDetails(agent string) UserAgentDetails {
	ua := useragent.New(agent)

	deviceType := "desktop"
	if ua.Mobile() {
		deviceType = "mobile"
	}

	osInfo := ua.OS()
	browserName, browserVersion := ua.Browser()

	return UserAgentDetails{
		DeviceType: deviceType,
		OS:         osInfo,
		Browser:    browserName + " " + browserVersion,
	}
}

// sends a GET request to the given URL and parses the response into the provided struct
func SendGetRequest(url string, result interface{}) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}
	if result.(*IPInfoResponse).Status != "success" {
		return errors.New("error getting IP info")
	}
	return nil
}

var BOT_USER_AGENTS = []string{
	"Googlebot",
	"Bingbot",
	"Slurp",
	"DuckDuckBot",
	"Baiduspider",
	"YandexBot",
	"Sogou",
	"Exabot",
	"facebookexternalhit",
	"ia_archiver",
	"Alexa Crawler",
	"AhrefsBot",
	"Applebot",
	"archive.org_bot",
	"Barkrowler",
	"BLEXBot",
	"BUbiNG",
	"CCBot",
	"Cliqzbot",
	"coccocbot",
	"Daum",
	"Diffbot",
	"DotBot",
	"EveryoneSocialBot",
	"Findxbot",
	"Gluten Free Crawler",
	"Google-Read-Aloud",
	"heritrix",
	"HubSpot",
	"ichiro",
	"LinkpadBot",
	"MJ12bot",
	"MojeekBot",
	"OpenAI GPTBot",
	"OpenAI ChatGPT",
	"PetalBot",
	"ping.blo.gs",
	"Pinterest",
	"pyspider",
	"redditbot",
	"rogerbot",
	"SemrushBot",
	"SeznamBot",
	"Snapchat",
	"Storebot-Google",
	"TelegramBot",
	"Twitterbot",
	"Vagabondo",
	"WhatsApp",
	"WordupInfoSearch",
	"YaK",
	"YandexAccessibilityBot",
	"YandexMetrika",
	"Adsbot",
	"APIs-Google",
}

var BOT_IP_RANGES = []string{
	"66.249.", // Google
	"40.77.",  // Bing
}

func IsBot(r *http.Request) bool {
	userAgent := r.UserAgent()
	for _, botUA := range BOT_USER_AGENTS {
		if strings.Contains(strings.ToLower(userAgent), strings.ToLower(botUA)) {
			return true
		}
	}

	// ip := r.RemoteAddr
	// for _, botIP := range BOT_IP_RANGES {
	// 	if strings.HasPrefix(ip, botIP) {
	// 		return true
	// 	}
	// }

	botPattern := regexp.MustCompile(`(?i)bot|crawler|spider|crawling`)
	return botPattern.MatchString(userAgent)
}
