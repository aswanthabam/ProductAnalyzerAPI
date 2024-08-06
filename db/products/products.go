package products_db

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	PRODUCT_ACCESS_KEY_SCOPE_ALL   = "all"
	PRODUCT_ACCESS_KEY_SCOPE_VISIT = "visit"
)

type Product struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"` // primary key
	Name        string               `bson:"name"`          // name of the product
	Description string               `bson:"description"`   // description of the product
	BaseUrl     string               `bson:"base_url"`      // base url of the product, the domain
	ProductID   string               `bson:"product_id"`    // product id of the product, used for identifying the product when requesting
	UserID      primitive.ObjectID   `bson:"user_id"`       // user id of the user who created the product
	AccessKeys  []primitive.ObjectID `bson:"access_keys"`   // access keys created for the product
	CreatedAt   primitive.DateTime   `bson:"created_at"`    // time at which the product was created
	UpdatedAt   primitive.DateTime   `bson:"updated_at"`    // time at which the product was last updated
}

type ProductVisit struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"` // primary key
	ProductID  primitive.ObjectID `bson:"product_id"`    // product id of the product
	SessionID  primitive.ObjectID `bson:"session_id"`    // session id of the user
	Refferer   string             `bson:"refferer"`      // refferer of the user
	Activities []ProductActivity  `bson:"activities"`    // activities of the user
}

type Location struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"` // primary key
	Hash     string             `bson:"hash"`          // hash of the location, used for identifying the location and finding if the location changed
	City     string             `bson:"city"`          // city of the location
	Region   string             `bson:"region"`        // region of the location
	Country  string             `bson:"country"`       // country of the location
	ZipCode  string             `bson:"zip_code"`      // zip code of the location
	TimeZone string             `bson:"time_zone"`     // time zone of the location
}

type ProductUserSession struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // primary key
	Hash      string             `bson:"hash"`          // hash of the session, used for identifying the session and finding if the session changed
	ProductID string             `bson:"product_id"`    // product id of the product
	IPAddress string             `bson:"ip_address"`    // ip address of the user
	Location  primitive.ObjectID `bson:"location_id"`   // location id of the user
	Lat       float64            `bson:"lat"`           // latitude of the location
	Lon       float64            `bson:"lon"`           // longitude of the location
	UserAgent string             `bson:"user_agent"`    // user agent of the user
	Proxy     bool               `bson:"proxy"`         // whether the user is using a proxy
	Isp       string             `bson:"isp"`           // internet service provider of the user
	Device    string             `bson:"device"`        // device type of the user, mobile, tablet, desktop
	Os        string             `bson:"os"`            // operating system of the user
	Browser   string             `bson:"browser"`       // browser of the user
	Bot       bool               `bson:"bot"`           // whether the user is a bot
	CreatedAt primitive.DateTime `bson:"created_at"`    // time at which the session was created
}

type ProductAccessKey struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // primary key
	ProductID primitive.ObjectID `bson:"product_id"`    // product id of the product
	AccessKey string             `bson:"access_key"`    // access key of the product
	Scope     string             `bson:"scope"`         // scope of the access key
	CreatedAt primitive.DateTime `bson:"created_at"`    // time at which the access key was created
}

/* INDIRECT TYPES */

type ProductActivity struct {
	From   string             `json:"from"`   // from where the user accessed the page
	Page   string             `json:"page"`   // page the user accessed
	Method string             `json:"method"` // method used to access the page, GET, POST, etc
	Time   primitive.DateTime `json:"time"`   // time at which the user accessed the page
}
