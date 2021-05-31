package model

//body request token
type TokenRequestBodyv2 struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Clientid     string `json:"clientid"`
	Clientsecret string `json:"clientsecret"`
}
type AuthenRequestBody struct {
	Token string `json:"token"`
}
type UserSSO struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Userid       string `json:"lastname"`
	Usermail     string `json:"usermail"`
	Userstatus   string `json:"userstatus"`
	Userparentid string `json:"userparentid"`
	Usertype     string `json:"usertype"`
}
type UserSSOs []UserSSO
type IoThing struct {
	Thingname    string `json:"thingname"`
	Thingid      string `json:"thingid"`
	Thingkey     string `json:"thingkey"`
	Thingstatus  string `json:"thingstatus"`
	Userparentid string `json:"userparentid"`
}
type IoTCreateRequest struct {
	Token       string `json:"token"`
	Thingname   string `json:"thingname"`
	Thingstatus string `json:"Thingstatus"`
}
type IoTChanel struct {
	Chanelname   string `json:"chanelname"`
	Chanelid     string `json:"chanelid"`
	Userparentid string `json:"userparentid"`
}
type IoTChanelRequest struct {
	Token      string `json:"token"`
	Chanelname string `json:"chanelname"`
}
type MapThingChanel struct {
	Thingid      string `json:"thingid"`
	Chanelid     string `json:"chanelid"`
	Mapstatus    string `json:"mapstatus"`
	Userparentid string `json:"userparentid"`
}
type MapThingChanelRequest struct {
	Token     string `json:"token"`
	Thingid   string `json:"thingid"`
	Chanelid  string `json:"chanelid"`
	Mapstatus string `json:"mapstatus"`
}

//push
type PushMessage struct {
	Token    string `json:"token"`
	Thingid  string `json:"thingid"`
	Thingkey string `json:"thingkey"`
	Chanelid string `json:"chanelid"`
	Message  string `json:"message"`
}
type PushMessagePublic struct {
	Token       string `json:"token"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Message     string `json:"message"`
}
type DataPublic struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Message     string `json:"message"`
}

//public chanel
type OrganizationPublicChanelRequest struct {
	Token string `json:"token"`
	Name  string `json:"name"`
}
type OrganizationPublicChanel struct {
	Publichanelid   string `json:"publichanelid"`
	Name            string `json:"name"`
	OrganizatioName string `json:"organizationame"`
	OrganizationId  string `json:"organizationid"`
}
