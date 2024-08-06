package products

import (
	"fmt"
	products_db "productanalyzer/api/db/products"
	api_error "productanalyzer/api/errors"
	"productanalyzer/api/utils"
	response "productanalyzer/api/utils/response"

	"github.com/gin-gonic/gin"
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
	clientIp, error := utils.GetIPAddress(c)
	if error != nil {
		response.SendFailureResponse(c, error)
		return
	}
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
		fmt.Println("Location does not exist")
		if location.ID, err = location.Save(); err != nil {
			response.SendFailureResponse(c, err)
			return
		}
	}

	product := prod.(*products_db.Product)
	session := products_db.ProductUserSession{
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
	if session.Hash, err = utils.HashStruct(session); err != nil {
		response.SendFailureResponse(c, err)
		return
	}

	if exists, err := session.ExistsHash(); err != nil {
		response.SendFailureResponse(c, err)
		return
	} else if exists {
		if err := session.GetByHash(); err != nil {
			response.SendFailureResponse(c, err)
			return
		}

	} else {
		if session.ID, err = session.Save(); err != nil {
			response.SendFailureResponse(c, err)
			return
		}
	}
	activity := products_db.ProductActivity{
		From:   params.From,
		Page:   params.Page,
		Method: params.Method,
		Time:   utils.GetCurrentTime(),
	}
	visitId, err2 := products_db.VisitProduct(product.ID, session.ID, activity, c.Request.Referer())
	if err2 != nil {
		response.SendFailureResponse(c, err)
		return
	}
	response.SendSuccessResponse(c, "LOL", ProductVisitResponse{
		VisitId:    visitId.Hex(),
		SessionId:  session.ID.Hex(),
		LocationId: location.ID.Hex(),
	}, nil)
}
