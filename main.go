package main

import (
	//"fmt"
	"net/http"

	//"bitbucket.org/cloud-platform/uip2sop_platform/integratedkafka"
	"bitbucket.org/cloud-platform/uip2sop_platform/service"
	"github.com/labstack/echo"
)

/*func init() {
	go func() {
		fmt.Println("---------Subscribe Local Topic Start---------")
		////service.SubscribeMQTT("public")
		//config.ReadConfig()
		//url := viper.GetString(`kafka.url`)
		//integratedkafka.KafkaInit(url)
		//producer, _ := transport.InitProducer(url)
		//transport.PublishMessage(producer, "UNKNOW_TOPIC", kafka.PartitionAny, "INIT UNKNOW_TOPIC")
		//service.SubscribeMQTT("public")
		integratedkafka.HandlerConsumer2()
		////fmt.Println("---------Subscribe Local Topic End---------")
	}()
}*/
func main() {
	e := echo.New()

	//e.Validator = &service.CustomValidator{Validator: validator.New()}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "IOT API FOR ORGANIZATION1")
	})

	//USER MANAGEMENT
	//<--------------------------- USER MANAGEMENT ------------------------------>
	e.POST("/api/sso/user/register", service.CreateUser)
	//e.PATCH("/api/sso/v1/user/updatenduser", service.UpdateEndUser)
	//e.PATCH("/api/sso/v1/user/updatecompanyuser", service.UpdateCompanyUser)
	//e.DELETE("/api/sso/v1/user/deletenduser/:username", service.DeleteEndUser)
	//e.GET("/api/sso/v1/user/info/:username", service.GetUserInfo)
	//e.GET("/api/sso/v1/user/childinfo/:username", service.GetUserChildInfo)
	//e.GET("/api/sso/v1/user/infobycomid/:comid", service.GetAllUserByComId)
	//e.DELETE("/api/sso/v1/user/deleteUserRelationship/:username", service.DeleteRelationshipUser)
	//e.PATCH("/api/sso/v1/user/checkparenthavechild", service.CheckParentHaveChild)
	//e.POST("/api/sso/v1/user/paging", service.PagingUserInfo)
	//e.POST("/api/sso/v1/user/pagingandfilter", service.PagingAndFilterUserInfo)*/
	//<--------------------------- THINGS MANAGEMENT ------------------------------>
	e.POST("/api/sso/things/register", service.CreateThings)
	//<--------------------------- CHANELS MANAGEMENT ------------------------------>
	e.POST("/api/sso/chanels/register", service.CreateChanels)
	//<--------------------------- MAP THINGS CHANELS ------------------------------>
	e.POST("/api/sso/mapthingtochanel/register", service.CreateMapThingChanel)
	//<--------------------------- PUSH MESSAGE ------------------------------>
	e.POST("/api/mqtt/pushMessage", service.PushMessage)
	e.POST("/api/mqtt/pushMessageToPublicBroker", service.PushMessageToPublic)
	e.POST("/api/mqtt/pushMessageToPublic", service.PushMessageToPublicBroker)
	//<---------------------------SSO ------------------------------>
	e.POST("/api/sso/requestToken", service.RequestSSOTokenv2)
	e.POST("/api/sso/parseToken", service.ParseSSOToken)
	//<---------------------------SSO ------------------------------>
	e.POST("/api/sso/publicchanel/register", service.CreateOrganizationPublicChanel)
	//e.POST("/api/sso/v1/user/ioToken", service.RequestSSOIoToken)
	//e.POST("/api/sso/v1/user/ioTAuthenKey", service.GetIoTAuthKeyFromRedis)
	//e.POST("/api/sso/v1/user/getTokenByRFToken", service.RefreshSSOToken)
	e.Logger.Fatal(e.Start(":1323"))

}
