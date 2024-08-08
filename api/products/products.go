package products

import (
	products_db "productanalyzer/api/db/products"
	api_error "productanalyzer/api/errors"
	"productanalyzer/api/utils"
	response "productanalyzer/api/utils/response"
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
	var session *products_db.ProductUserSession
	if params.SessionId != "" {
		sessionId, err := primitive.ObjectIDFromHex(params.SessionId)
		if err != nil {
			response.SendFailureResponse(c, api_error.NewAPIError("Invalid Session ID", 400, "Invalid Session ID"))
			return
		}
		session, err = products_db.GetSessionById(sessionId)
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
	product := prod.(*products_db.Product)
	if session == nil {
		info, err := utils.GetIPAddressInfo(clientIp)
		if err != nil {
			response.SendFailureResponse(c, err)
			return
		}
		userAgent := c.GetHeader("User-Agent")
		ua := utils.GetUserAgentDetails(userAgent)
		// referer := c.GetHeader("Referer")
		location := products_db.Location{
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
		if exists, err := location.ExistsHash(); err != nil {
			response.SendFailureResponse(c, err)
			return
		} else if exists {
			if err := location.GetByHash(); err != nil {
				response.SendFailureResponse(c, err)
				return
			}
		} else {
			if err = location.Save(); err != nil {
				response.SendFailureResponse(c, err)
				return
			}
		}
		session = &products_db.ProductUserSession{
			ProductID: product.ID.Hex(),
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
		}
		if err = session.HashSession(); err != nil {
			response.SendFailureResponse(c, err)
			return
		}

		if exists, err := session.GetByHash(); err != nil {
			response.SendFailureResponse(c, err)
			return
		} else if exists {
			expireTime := session.UpdatedAt.Time().Add(time.Minute * 2).UTC()
			if expireTime.Before(utils.GetUTCTime()) {
				session.ID = primitive.NilObjectID
				if err = session.Save(); err != nil {
					response.SendFailureResponse(c, err)
					return
				}
			}
		} else {
			if err = session.Save(); err != nil {
				response.SendFailureResponse(c, err)
				return
			}
		}
	}
	activity := products_db.ProductActivity{
		From:   params.From,
		Page:   params.Page,
		Method: params.Method,
		Time:   utils.GetCurrentTime(),
	}
	err := session.VisitProduct(activity)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	response.SendSuccessResponse(c, "LOL", ProductVisitResponse{
		SessionId:  session.ID.Hex(),
		LocationId: session.Location.Hex(),
	}, nil)
}
