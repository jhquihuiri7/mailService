package request

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ompluscator/dynamic-struct"
	"github.com/xuri/excelize/v2"
	"log"
	"strings"
)

type RequestStandard struct {
	ClientName string `json:"clientName"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Mail       string `json:"mail"`
	Message    string `json:"message"`
}
type RequestBulk struct {
	ClientName string            `json:"clientName"`
	Template   string            `json:"template"`
	Tos        []RequestStandard `json:"tos"`
	Limits     []int             `json:"limits"`
}
type RequestResponse struct {
	Success string `json:"success"`
	Error   string `json:"error"`
}
type SumarizeResponse struct {
	Success []interface{} `json:"success"`
	Error   string        `json:"error"`
}

func (r *RequestStandard) ParseRequestStandardData(c *gin.Context) {
	err := json.NewDecoder(c.Request.Body).Decode(&r)
	if err != nil {
		log.Fatal(err)
	}
}
func (r *RequestBulk) ParseRequestBulkData(c *gin.Context) {
	err := json.NewDecoder(c.Request.Body).Decode(&r)
	if err != nil {
		log.Fatal(err)
	}
}
func (r *RequestBulk) ValidateDataInput(c *gin.Context) SumarizeResponse {
	var response SumarizeResponse
	c.Request.ParseMultipartForm(10 << 20)
	file, _, err := c.Request.FormFile("data")
	if err != nil {
		log.Fatal(err)
		return response
	}
	defer file.Close()
	excelFile, err := excelize.OpenReader(file)
	if err != nil {
		log.Fatal(err)
		return response
	}
	data, err := excelFile.GetCols(excelFile.GetSheetName(0))
	if err != nil {
		log.Fatal(err)
		return response
	}
	var colNames []string
	for _, v := range data {
		colNames = append(colNames, v[0])
	}
	instance := dynamicstruct.NewStruct()
	for _, v := range colNames {
		instance.AddField(v, "", fmt.Sprintf(`json:"%s"`, strings.ToLower(v)))
	}
	dynamicStruct := instance.Build().New()

	data, err = excelFile.GetRows(excelFile.GetSheetName(0))
	if err != nil {
		log.Fatal(err)
		return response
	}

	for i, v := range data {
		if i < 1 || i >= 10 {
			continue
		}
		rawData := "{"
		lenData := len(v)
		for index, val := range v {

			rawData += "\"" + strings.ToLower(colNames[index]) + "\":\"" + val + "\","
			if index == lenData-1 {
				rawData = rawData[:len(rawData)-1]
			}
		}
		rawData += "}"
		err = json.Unmarshal([]byte(rawData), &dynamicStruct)
		if err != nil {
			log.Fatal(err)
		}
		response.Success = append(response.Success, dynamicStruct)
	}
	return response
}
func (resp *SumarizeResponse) Marshal() string {
	JSONresponse, _ := json.Marshal(resp)
	return string(JSONresponse)
}
func (resp *RequestResponse) Marshal() string {
	JSONresponse, _ := json.Marshal(resp)
	return string(JSONresponse)
}
