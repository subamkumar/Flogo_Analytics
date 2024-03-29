package sample

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {

	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	c := &http.Client{}

	//ctx.Logger().Debugf("Setting: %s", s.API_BASE_URL)

	act := &Activity{settings: s, client: c} //add aSetting to instance

	return act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
	settings *Settings
	client   *http.Client
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	urlString := a.settings.API_Gateway
	ctx.Logger().Debugf("API_Gateway: %s", urlString)
	/*if urlString!= "" {
			_,err := url.Parse(urlString)
	 		if err != nil {
	 			return true, err
	 		}

			method := "GET"

			if input.ProcessType != "" {

				if strings.LastIndex(urlString,"/") == len(urlString)-1 {
					urlString = urlString+input.ProcessType
				}else{
					urlString = urlString+"/"+input.ProcessType
				}


				if input.PathParamId != 0 {
					urlString = urlString+"/"+strconv.Itoa(input.PathParamId)
				}else{

					queryString := ""

					if input.CollectionId != 0 {
						queryString = queryString+"collection_id="+strconv.Itoa(input.CollectionId)
					}

					if input.ActivityId != 0 {
						if len(queryString) > 0 {
							queryString = queryString+"&"+"activity_id="+strconv.Itoa(input.ActivityId)
						}else{
							queryString = queryString+"activity_id="+strconv.Itoa(input.ActivityId)
						}
					}

					if input.LocationId != 0 {
						if len(queryString) > 0 {
							queryString = queryString+"&"+"location_id="+strconv.Itoa(input.LocationId)
						}else{
							queryString = queryString+"location_id="+strconv.Itoa(input.LocationId)
						}
					}

					if len(queryString) > 0 {
						urlString = urlString+"?"+queryString
					}
				}

				//ctx.Logger().Debugf("FORMED_URL: %s", urlString)

				req, _ := http.NewRequest(method, urlString, nil)
				resp, err := a.client.Do(req)

				var responseData interface{}

				respContentType := resp.Header.Get("Content-Type")
				switch respContentType {
					case "application/json":
						d := json.NewDecoder(resp.Body)
						d.UseNumber()
						err = d.Decode(&responseData)
						if err != nil {
							switch {
								case err == io.EOF:
								default:
									return false, err
							}
						}
					default:
						b, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							return false, err
						}
					responseData = string(b)
				}

				output := &Output{ResponseCode: resp.StatusCode, ResponseData: responseData}
				err = ctx.SetOutputObject(output)
				if err != nil {
					return true, err
				}

				return true, nil
			}
			return true, activity.NewError("Required Process Type is not provided","",nil)
		}*/

	body := make(map[string]interface{})

	body["collection"] = input.Collection
	body["location"] = input.Location
	body["activity"] = input.Activity
	//body["username"] = input.Username

	jsonData, _ := json.Marshal(body)
	byteData := bytes.NewBuffer(jsonData)

	method := "POST"

	req, _ := http.NewRequest(method, urlString, byteData)
	resp, _ := a.client.Do(req)

	//return true, activity.NewError("API Gateway URL is not provided","",nil)
	ctx.Logger().Debugf("Input: %s", input.Collection)
	ctx.Logger().Debugf("Input: %s", input.Location)
	ctx.Logger().Debugf("Input: %s", input.Activity)
	ctx.Logger().Debugf("Input: %s", input.Username)
	ctx.Logger().Debugf("Input: %s", input.Password)
	ctx.Logger().Debugf("Input: %s", input.RequestId)
	ctx.Logger().Debugf("Input: %s", input.SecretId)

	ctx.Logger().Debugf("Input: %d", resp.StatusCode)

	//ctx.Logger().Debugf("Input: %d", input.ActivityId)
	return true, nil
}
