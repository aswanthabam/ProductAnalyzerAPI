package products_db

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	PRODUCT_ACCESS_KEY_SCOPE = "all"
)

type Product struct {
	ID          primitive.ObjectID `json:"id"`          // primary key
	Name        string             `json:"name"`        // name of the product
	Description string             `json:"description"` // description of the product
	BaseUrl     string             `json:"base_url"`    // base url of the product, the domain
	ProductID   string             `json:"product_id"`  // product id of the product, used for identifying the product when requesting
	AccessKeys  []ProductAccessKey `json:"access_keys"` // access keys created for the product
	CreatedAt   primitive.DateTime `json:"created_at"`  // time at which the product was created
	UpdatedAt   primitive.DateTime `json:"updated_at"`  // time at which the product was last updated
}

type ProductVist struct {
	ID         primitive.ObjectID `json:"id"`         // primary key
	ProductID  string             `json:"product_id"` // product id of the product
	SessionID  string             `json:"session_id"` // session id of the user
	Activities []ProductActivity  `json:"activities"` // activities of the user
}

type Location struct {
	ID       primitive.ObjectID `json:"id"`        // primary key
	Hash     string             `json:"hash"`      // hash of the location, used for identifying the location and finding if the location changed
	City     string             `json:"city"`      // city of the location
	Region   string             `json:"region"`    // region of the location
	Country  string             `json:"country"`   // country of the location
	ZipCode  string             `json:"zip_code"`  // zip code of the location
	TimeZone string             `json:"time_zone"` // time zone of the location
}

type ProductUserSession struct {
	ID        primitive.ObjectID `json:"id"`          // primary key
	Hash      string             `json:"hash"`        // hash of the session, used for identifying the session and finding if the session changed
	ProductID string             `json:"product_id"`  // product id of the product
	IPAddress string             `json:"ip_address"`  // ip address of the user
	Location  primitive.ObjectID `json:"location_id"` // location id of the user
	Lat       float64            `json:"lat"`         // latitude of the location
	Lon       float64            `json:"lon"`         // longitude of the location
	UserAgent string             `json:"user_agent"`  // user agent of the user
	Mobile    bool               `json:"mobile"`      // whether the user is using a mobile device
	Bot       bool               `json:"bot"`         // whether the user is a bot
	Proxy     bool               `json:"proxy"`       // whether the user is using a proxy
	Isp       string             `json:"isp"`         // internet service provider of the user
	Refferer  string             `json:"refferer"`    // refferer of the user
	CreatedAt primitive.DateTime `json:"created_at"`  // time at which the session was created
}

/* INDIRECT TYPES */

type ProductActivity struct {
	From   string             `json:"from"`   // from where the user accessed the page
	Page   string             `json:"page"`   // page the user accessed
	Method string             `json:"method"` // method used to access the page, GET, POST, etc
	Time   primitive.DateTime `json:"time"`   // time at which the user accessed the page
}

type ProductAccessKey struct {
	AccessKey string `json:"access_key"` // access key of the product
	Scope     string `json:"scope"`      // scope of the access key
}
