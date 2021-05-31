package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"bitbucket.org/cloud-platform/uip2sop_platform/config"
	"github.com/spf13/viper"
)

//

/*type Attributes struct {
	ComDepartment string `json:"comdepartment"`
	ComId         string `json:"comid"`
	Lastname      string `json:"lastname"`
	UserCode      string `json:"usercode"`
	UserDate      string `json:"userdate"`
	UserEmail     string `json:"useremail"`
	UserStatus    string `json:"userstatus"`
	UserTel       string `json:"usertel"`
	UserType      string `json:"usertype"`
	UserParentid  string `json:"userparentid"`
}*/
type Tokens struct {
	Accesstoken  string `json:"access_token"`
	Refreshtoken string `json:"refresh_token"`
	Tokentype    string `json:"token_type"`
	Expiresin    string `json:"expires_in"`
	Scope        string `json:"scope"`
}
type RTokens struct {
	Accesstoken string `json:"access_token"`
	Tokentype   string `json:"token_type"`
	Expiresin   string `json:"expires_in"`
	Scope       string `json:"scope"`
}

func BasicAuth(token string) (s string, err error) {
	config.ReadConfig()
	link := viper.GetString(`sso.url`)
	//link := "mongodb://casuser:Mellon@iotsso:27017/cas?authSource=cas"
	var username string = "exampleOauthClient"
	var passwd string = "exampleOauthClientSecret"
	client := &http.Client{}
	url := link + "/profile?access_token=" + token
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
		return "", err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s = string(bodyText)
	return s, nil
}
func RequestTokenv2(client_id string, client_secret string, username string, password string) (string, string) {
	//readConfig()
	config.ReadConfig()
	link := viper.GetString(`sso.url`)
	//link := "mongodb://casuser:Mellon@iotsso:27017/cas?authSource=cas"
	url := link + "/token?grant_type=password&client_id=" + client_id + "&client_secret=" + client_secret + "&username=" + username + "&password=" + password
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		//Sys log
		fmt.Println(err)
		return "", ""
		//Sys log
	}
	body, _ := ioutil.ReadAll(res.Body)
	stringBody := string(body)
	if stringBody != "" {
		var tokens Tokens
		json.Unmarshal([]byte(stringBody), &tokens)
		return tokens.Accesstoken, tokens.Refreshtoken
	} else {
		return "", ""
	}
}
