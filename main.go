package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	xj "github.com/basgys/goxml2json"
	sj "github.com/bitly/go-simplejson"
)

var rawXml = `<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Date>2020-12-13T20:12:41.6228299</Date>
    <Author>BUIHOANG\hoang</Author>
    <Description>test_trigger description</Description>
    <URI>\test_trigger</URI>
  </RegistrationInfo>
  <Triggers>
    <EventTrigger>
      <Enabled>true</Enabled>
      <Subscription>&lt;QueryList&gt;&lt;Query Id="0" Path="Microsoft-Windows-Audio/CaptureMonitor"&gt;&lt;Select Path="Microsoft-Windows-Audio/CaptureMonitor"&gt;*[System[Provider[@Name='Microsoft-Windows-Audio'] and EventID=13]]&lt;/Select&gt;&lt;/Query&gt;&lt;/QueryList&gt;</Subscription>
    </EventTrigger>
  </Triggers>
  <Principals>
    <Principal id="Author">
      <UserId>S-1-5-21-1869200887-1913775759-1203460732-1001</UserId>
      <LogonType>InteractiveToken</LogonType>
      <RunLevel>LeastPrivilege</RunLevel>
    </Principal>
  </Principals>
  <Settings>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>true</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>true</StopIfGoingOnBatteries>
    <AllowHardTerminate>true</AllowHardTerminate>
    <StartWhenAvailable>false</StartWhenAvailable>
    <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
    <IdleSettings>
      <Duration>PT10M</Duration>
      <WaitTimeout>PT1H</WaitTimeout>
      <StopOnIdleEnd>true</StopOnIdleEnd>
      <RestartOnIdle>false</RestartOnIdle>
    </IdleSettings>
    <AllowStartOnDemand>true</AllowStartOnDemand>
    <Enabled>true</Enabled>
    <Hidden>false</Hidden>
    <RunOnlyIfIdle>false</RunOnlyIfIdle>
    <WakeToRun>false</WakeToRun>
    <ExecutionTimeLimit>PT72H</ExecutionTimeLimit>
    <Priority>7</Priority>
  </Settings>
  <Actions Context="Author">	
    <Exec>
      <Command>C:\Users\hoang\Desktop\trigger\User_Feed_Synchronization-{6FA533E6-FB3A-43F7-8612-EA982FB6D214}.xml</Command>
    </Exec>
	  <Exec>
      <Command>C:\Users\hoang\Desktop\trigger\abc.xml</Command>
    </Exec>
	<SendEmail>
		<To>To 1</To>
		<ReplyTo>ReplyTo 1</ReplyTo>
		<From>From 1</From>
		<Body>Body 1</Body>
		<Subject>Subject 1</Subject>
		<Bcc>Bcc 1</Bcc>
		<Attachments> 
			<File>abcdef</File>
		</Attachments>
	</SendEmail>
  </Actions>
</Task>`

type Scheduler struct {
	RegistrationInfo RegistrationInfo `json:"RegistrationInfo"`
	Principals       []Principal      `json:"Principal"`
	Actions          Actions          `json:"Actions"`
}

type Principal struct {
	UserId      string `json:"UserId"`
	LogonType   string `json:"LogonType"`
	GroupId     string `json:"GroupId"`
	DisplayName string `json:"DisplayName"`
	RunLevel    string `json:"RunLevel"`
}

type Actions struct {
	Exec      []Exec      `json:"Exec"`
	SendEmail []SendEmail `json:"SendEmail"`
}

type Exec struct {
	Command          string `json:"Command"`
	Arguments        string `json:"Arguments"`
	WorkingDirectory string `json:"WorkingDirectory"`
}
type RegistrationInfo struct {
	URI    string `json:"URI"`
	Author string `json:"Author"`
	Date   string `json:"Date"`
	Source string `json:"Source"`
}
type Attachments struct {
	File string `json:"File"`
}

type SendEmail struct {
	Server      string      `json:"Server"`
	Subject     string      `json:"Subject"`
	Cc          string      `json:"Cc"`
	To          string      `json:"To"`
	Bcc         string      `json:"Bcc"`
	ReplyTo     string      `json:"ReplyTo"`
	From        string      `json:"From"`
	Body        string      `json:"Body"`
	Attachments Attachments `json:"Attachments"`
}

func main() {
	// xml is an io.Reader
	rawXml = strings.Replace(rawXml, "UTF-16", "UTF-8", -1)
	xml := strings.NewReader(rawXml)
	jsonEncoder, err := xj.Convert(xml)
	if err != nil {
		panic("That's embarrassing...")
	}
	jsonStr := jsonEncoder.String()
	// fmt.Println(jsonStr)
	// {"hello": "world"}
	jsonMapper, err := sj.NewJson([]byte(jsonStr))
	if err != nil {
		fmt.Print(err)
	}
	var scheduler = Scheduler{}
	keys := []string{"RegistrationInfo", "Principals", "Actions"}
	for key, value := range jsonMapper.Get("Task").MustMap() {
		isExist := findElement(keys, key)
		if isExist {
			filter(&scheduler, value, key)
		}
	}
	fmt.Printf("scheduler: %+v ", scheduler)
}

func findElement(list []string, element string) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}
	return false
}

func filter(scheduler *Scheduler, mapper interface{}, childKey string) {

	isUpdated, _ := updateScheduler(scheduler, mapper, childKey)
	if isUpdated {
		return
	}
	m := mapper.(map[string]interface{})
	for key, value := range m {
		kind := reflect.TypeOf(value).Kind()
		switch kind {
		case reflect.Slice:
			arr := value.([]interface{})
			for _, element := range arr {
				updateScheduler(scheduler, element, key)
			}
		default:
			updateScheduler(scheduler, value, key)
		}

	}
}

func updateScheduler(scheduler *Scheduler, child interface{}, childKey string) (isUpdated bool, err error) {
	jsonBody, err := json.Marshal(child)
	if err != nil {
		return isUpdated, err
	}
	isUpdated = true
	var childInstance interface{}
	switch childKey {
	case "Exec":
		fmt.Println("Exec")
		childInstance, _ := childInstance.(Exec)
		if err := json.Unmarshal(jsonBody, &childInstance); err != nil {
			fmt.Println("err Exec: ", err)
			isUpdated = false
			break
		}
		scheduler.Actions.Exec = append(scheduler.Actions.Exec, childInstance)
	case "SendEmail":
		fmt.Println("SendEmail")

		childInstance, _ := childInstance.(SendEmail)
		if err := json.Unmarshal(jsonBody, &childInstance); err != nil {
			fmt.Println("err SendEmail: ", err)
			isUpdated = false
			break
		}
		scheduler.Actions.SendEmail = append(scheduler.Actions.SendEmail, childInstance)
	case "Principal":
		fmt.Println("Principal")

		childInstance, _ := childInstance.(Principal)
		if err := json.Unmarshal(jsonBody, &childInstance); err != nil {
			fmt.Println("err Principal: ", err)
			isUpdated = false
			break
		}
		scheduler.Principals = append(scheduler.Principals, childInstance)
	case "RegistrationInfo":
		childInstance, _ := childInstance.(RegistrationInfo)
		if err := json.Unmarshal(jsonBody, &childInstance); err != nil {
			fmt.Println("err: ", err)
			isUpdated = false
			break
		}
		scheduler.RegistrationInfo = childInstance
	default:
		isUpdated = false
	}
	return isUpdated, err
}
