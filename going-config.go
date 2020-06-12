package GoingConfig

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"encoding/json"
)

type ConfigSettingStruct struct {
	ConfigAddress	string
	ConfigPort	string
	ConfigUserName	string
	ConfigPassword	string
	ConfigID		string
}

type ConfigTokenStruct struct {
	Access_token	string
	Token_type		string
	Expires_in		int
	Refresh_token	string
}


type ConfigStruct struct {
	ID	string
	Name 	string
	Description	string
	Type	string
	Detail	interface{}
}


func getConfig(address, port, user, password, config_id string) (string, error) {
	var config_str string
	var err error

	config_token, err := getConfigToken(address, port, user, password)
	if err != nil || config_token == nil {
		err = fmt.Errorf("func getConfigToken(): %v, config_token: %v", err, config_token)
	} else {

		url := "http://" + address + ":" + port + "/api/workers/configurations/" + config_id
		method := "GET"
	  
		client := &http.Client {
		}
		req, err := http.NewRequest(method, url, nil)
	  
		if err != nil {
		  fmt.Println(err)
		}
		req.Header.Add("accept", "application/json")
		req.Header.Add("Authorization", config_token.Token_type + " " + config_token.Access_token)
	  
		res, err := client.Do(req)
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
	  

		var config ConfigStruct
		if err = json.Unmarshal(body, &config); err != nil {
			return config_str, fmt.Errorf("Unable to parse config: %v, error: %v, %v", string(body), err)
		}
		
		jsonBtye, err := json.Marshal(config.Detail)
		if err != nil {
			return config_str, fmt.Errorf("Unable to Marshal config.Detail: %v, error: %v", config.Detail, err)
		}
		config_str = string(jsonBtye)
	}
	return config_str, err
}

func getConfigToken(address, port, user, password string) (*ConfigTokenStruct, error) {

	url := "http://" + address + ":" + port + "/api/users/login"
 	method := "POST"
	
	payload := strings.NewReader("username=" + user + "&password=" + password)

	client := &http.Client {
	}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		err = fmt.Errorf("Unable to get config token with address: %s, user: %s, password: %s, err: %v", address, user, password, err)
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	
	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("Unable to get config token with req: %v, err: %v", req, err)
		return nil, err
	}
	
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("Unable to get config token with body: %v, err: %v", body, err)
		return nil, err
	}
	

	var configToken	ConfigTokenStruct
	if err = json.Unmarshal(body, &configToken); err != nil {
		return nil, fmt.Errorf("Unable to parse configToken: %v, error: %v, %v", string(body), err)
	}
	
	return &configToken, err
}
