package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/kdevb0x/or-maybe-more/or-maybe-more/app"
)

// node acts as an intermediary between client streams. Nodes are ment to be
// deployed in many geographical locations. Durring session negotiation of a
// stream, a node is chosen by geographic proximity, and lowest latency.
type node struct {
	ipaddr   net.IPAddr
	location LocInfo
}

// LocInfo contains the geographical location information of a node.
// Because of the ephemeral nature of location data, the contained information
// is only valid until the time indicated by the ValidUntil field, which may be
// checked using Expired().
type LocInfo struct {
	ValidUntil time.Time
	lifetime   time.Duration
	// Coordinates[0] == longitute; Coordinates[1] == latitude.
	Coordinates [2]float64
	// additional info
	GeoIP *app.GeoIP
}

// Expired returns true if loc has outlived its lifetime and is considered
// invalid.
func (loc *LocInfo) Expired() bool {
	if time.Now().Before(loc.ValidUntil) {
		return false
	}
	return true
}

// SetTTL set the time-to-live (aka lifetime) of loc. After duration t loc is
// invalidated.
func (loc *LocInfo) SetTTL(t time.Duration) {
	loc.lifetime = t
	loc.ValidUntil = time.Now().Add(t)

	/*
		if !loc.Expired() {
			// [kdev]: loc should never live longer than its lifetime, so
			// drop the remaining duration, and start over.

		}
		loc.lifetime = t
	*/

}

// GeoIP holds values for interacting with the freegeoip.net API.
type GeoIP struct {
	// The right side is the name of the JSON variable
	Ip          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name""`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	Zipcode     string  `json:"zipcode"`
	Lat         float32 `json:"latitude"`
	Lon         float32 `json:"longitude"`
	MetroCode   int     `json:"metro_code"`
	AreaCode    int     `json:"area_code"`
}

func (geo *GeoIP) String() string {
	s := `==== IP Geolocation Info ====\n
	IP address:\t%s\n
	Country Code:\t%s\n
	Country Name:\t%s\n
	City:\t%s\n
	Zip Code:\t%s\n
	Latitude:\t%s\n
	Longitude:\t%s\n
	Metro Code:\t%s\n
	Area Code:\t%s\n`

	_, err := fmt.Scanf(s,
		geo.Ip,
		geo.CountryCode,
		geo.CountryName,
		geo.Zipcode,
		geo.City,
		geo.Lat,
		geo.Lon,
		geo.MetroCode,
		geo.AreaCode)

	if err != nil {
		return err.Error()
	}
	return s

}

// GetGeoIP fetches the geographical location of addr using the freegeoip.net api.
func GetGeoIP(addr string) (*GeoIP, error) {
	ip, err := net.ResolveIPAddr("ipv4", addr)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get("https://freegeoip.net/json/" + ip.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var g = new(GeoIP)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, g)
	if err != nil {
		return nil, err
	}
	return g, nil

}
