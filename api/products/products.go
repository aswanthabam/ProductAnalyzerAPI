package products

import (
	"productanalyzer/api/db"
	db_types "productanalyzer/api/db/types"
	api_error "productanalyzer/api/errors"
	"productanalyzer/api/utils"
	response "productanalyzer/api/utils/response"
	"productanalyzer/api/websockets"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func VisitProduct(c *gin.Context) {
	prod, exists := c.Get("product")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	var params VisitProductRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	var session *db.ProductUserSession
	if params.SessionId != "" {
		sessionId, err := primitive.ObjectIDFromHex(params.SessionId)
		if err != nil {
			response.SendFailureResponse(c, api_error.NewAPIError("Invalid Session ID", 400, "Invalid Session ID"))
			return
		}
		session, err = db.ProductRepository.GetSessionById(sessionId)
		if err != nil {
			response.SendFailureResponse(c, err)
			return
		}

	}
	clientIp, error := utils.GetIPAddress(c)
	if error != nil {
		response.SendFailureResponse(c, error)
		return
	}
	if session != nil && session.IPAddress != clientIp {
		response.SendFailureResponse(c, api_error.NewAPIError("Invalid Session ID", 400, "Invalid Session ID, Session is not of the current user"))
		return
	}
	product := prod.(*db.Product)
	var location db.Location
	productRepository := db.ProductRepository
	if session == nil {
		info, err := utils.GetIPAddressInfo(clientIp)
		if err != nil {
			response.SendFailureResponse(c, err)
			return
		}
		userAgent := c.GetHeader("User-Agent")
		ua := utils.GetUserAgentDetails(userAgent)
		// referer := c.GetHeader("Referer")
		location = db.Location{
			City:     info.City,
			Region:   info.Region,
			Country:  info.Country,
			ZipCode:  info.Zip,
			TimeZone: info.Timezone,
		}
		if location.Hash, err = utils.HashStruct(location); err != nil {
			response.SendFailureResponse(c, err)
			return
		}
		if exists, err := productRepository.ExistsLocationHash(location); err != nil {
			response.SendFailureResponse(c, err)
			return
		} else if exists {
			if err := productRepository.GetLocationByHash(&location); err != nil {
				response.SendFailureResponse(c, err)
				return
			}
		} else {
			if err = productRepository.SaveLocation(&location); err != nil {
				response.SendFailureResponse(c, err)
				return
			}
		}
		referer := c.Request.Header.Get("Referer")
		if referer == "" {
			referer = c.Query("referer")
		}
		session = &db.ProductUserSession{
			ProductID: product.ID,
			IPAddress: clientIp,
			Location:  location.ID,
			Lat:       info.Lat,
			Lon:       info.Lon,
			UserAgent: userAgent,
			Proxy:     info.Proxy,
			Isp:       info.ISP,
			Device:    ua.DeviceType,
			Os:        ua.OS,
			Browser:   ua.Browser,
			Bot:       c.GetBool("isBot"),
			Referer:   referer,
		}
		if err = productRepository.HashProductUserSession(session); err != nil {
			response.SendFailureResponse(c, err)
			return
		}

		if exists, err := productRepository.GetProductUserSessionByHash(session); err != nil {
			response.SendFailureResponse(c, err)
			return
		} else if exists {
			expireTime := session.UpdatedAt.Add(time.Minute * 2).UTC()
			if expireTime.Before(utils.GetUTCTime()) {
				session.ID = primitive.NilObjectID
				if err = productRepository.SaveProductUserSession(session); err != nil {
					response.SendFailureResponse(c, err)
					return
				}
			}
		} else {
			if err = productRepository.SaveProductUserSession(session); err != nil {
				response.SendFailureResponse(c, err)
				return
			}
		}
	}
	activity := db_types.ProductActivity{
		From:   params.From,
		Page:   params.Page,
		Method: params.Method,
		Time:   utils.GetCurrentTime(),
	}
	err := productRepository.VisitProduct(session, activity)
	websockets.SendLog(product.ID.Hex(), websockets.Visit{
		Country: location.Country,
		Referer: session.Referer,
	})
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	response.SendSuccessResponse(c, "LOL", ProductVisitResponse{
		SessionId:  session.ID.Hex(),
		LocationId: session.Location.Hex(),
	}, nil)
}
