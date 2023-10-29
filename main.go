package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Note struct{
	Date string
	Name string
	Phone string
	Message string
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(sendNotes)
}
func sendNotes() {
//	func main() {	
	currentTime := time.Now()
	currentDate := fmt.Sprintf("%d/%d/%d",currentTime.Month(),currentTime.Day(),currentTime.Year())
	file, err := os.Open("Notes.csv")
	if err != nil {
    	fmt.Println(err)
    	return
	}
	defer file.Close()
	reader := csv.NewReader(file) 
	var notes []Note
	
	b, err :=reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range b{
		if (currentDate==record[0]){
			note:=Note{
				Date: record[0],
				Name: record[1],
				Phone: record[2],
				Message: record[3],
			}
			notes=append(notes, note)
		}
	}
	

	
	
	// Your Twilio Account SID and Auth Token
	accountSid := os.Getenv("accountSid")
	authToken := os.Getenv("authToken")
	from := os.Getenv("from")


	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	for _,v :=range notes{
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(v.Phone)
	params.SetFrom(from)
	params.SetBody(v.Message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}
	fmt.Println(v)
}
}