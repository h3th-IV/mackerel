package models

type User struct {
	Email     string `json:"email"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
	Location  GeoLocation
	IpAddress string `json:"ip_address"`
}

type GeoLocation struct {
	City         string `json:"city"`
	Country      string `json:"country"`
	IpAddress    string `json:"ip"`
	Region       string `json:"region"`
	LatLong      string `json:"loc"` // latitude and longitude
	Organization string `json:"org"` //Netwrok Provide Organization
	TimeZone     string `json:"timezone"`
}

type AttackPayload struct {
	Email         string
	MaliciousLink string
}
