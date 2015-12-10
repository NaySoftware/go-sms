package telemsg

import "encoding/json"

// import "fmt"
import "net/http"
import "bytes"
import "io/ioutil"
import "errors"

const (
	serverUrl = "https://rest.telemessage.com/rest/message/send"

	AccountClass   = "telemessage.web.services.AuthenticationDetails"
	MsgClass       = "telemessage.web.services.Message"
	RecipientClass = "telemessage.web.services.Recipient"
	MsgType        = "SMS"
	MsgDescription = "-" // an optinal value | read provider docs for more info
)

type Sms struct {
	Username   string         `json:"username,omitempty"`
	Password   string         `json:"password,omitempty"`
	Body       string         `json:"textMessage,omitempty"`
	To         string         `json:"-"`
	Class      string         `json:"class,omitempty"`
	Recipients []SmsRecipient `json:"recipients,omitempty"`
}

type SmsRecipient struct {
	Value       string `json:"value,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Class       string `json:"class,omitempty"`
}

type StatusResponse struct {
	Ok                bool   `json:"-"`
	Class             string `json:"class,omitempty"`
	MsgId             int    `json:"messageID,omitempty"`
	MsgKey            string `json:"messageKey,omitempty"`
	ResultCode        int    `json:"resultCode,omitempty"`
	ResultDescription string `json:"resultDescription,omitempty"`
}

func NewClient(user string, pass string) *Sms {
	sms := new(Sms)

	sms.Username = user
	sms.Password = pass

	return sms
}

func (this *Sms) NewMsg(to string, body string) *Sms {

	this.To = to
	this.Body = body

	return this
}

func (this *Sms) Send() (StatusResponse, error) {

	var statres StatusResponse
	statres.Ok = false

	jsonByte, eror := this.generateJson()
	if eror != nil {

		return statres, eror
	}

	req, err := http.NewRequest("POST", serverUrl, bytes.NewBuffer(jsonByte))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return statres, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var status []StatusResponse
	if err := json.Unmarshal(body, &status); err != nil {
		return statres, err
	}

	if len(status) > 0 {
		s := status[0]
		if s.ResultCode == 100 || s.ResultCode == 0 {
			s.Ok = true
			return s, nil
		} else {
			return s, errors.New("FAILED: " + string(s.ResultCode))
		}
	} else {
		return statres, errors.New("FAILED: no response")
	}

}

func (this *Sms) generateJson() ([]byte, error) {

	var objSms [2]Sms
	var rec SmsRecipient

	rec.Class = RecipientClass
	rec.Description = MsgDescription
	rec.Type = MsgType
	rec.Value = this.To

	this.Recipients = append(this.Recipients, rec)

	objSms[0] = Sms{
		Class:    AccountClass,
		Username: this.Username,
		Password: this.Password,
	}

	objSms[1] = Sms{
		Class:      MsgClass,
		Body:       this.Body,
		Recipients: this.Recipients,
	}

	d, err := json.Marshal(objSms)

	return d, err
}
