package request

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ompluscator/dynamic-struct"
	"github.com/xuri/excelize/v2"
	"golang.org/x/exp/slices"
	"log"
	"strings"
	"unicode"
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
	Success SumData `json:"success"`
	Error   string  `json:"error"`
}
type SumData struct {
	Data  []interface{} `json:"data"`
	Limit int           `json:"limit"`
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
		return SumarizeResponse{Error: err.Error()}
	}
	defer file.Close()
	excelFile, err := excelize.OpenReader(file)
	if err != nil {
		return SumarizeResponse{Error: err.Error()}
	}
	data, err := excelFile.GetCols(excelFile.GetSheetName(0))
	if err != nil {
		return SumarizeResponse{Error: err.Error()}
	}
	var colNames []string

	for _, v := range data {
		response.ValidateLowerUpper(v[0])
		if response.Error != "" {
			return response
		}
		colNames = append(colNames, v[0])
	}
	if !slices.Contains(colNames, "Mail") && !slices.Contains(colNames, "Email") && !slices.Contains(colNames, "Correo") {
		return SumarizeResponse{Error: "No existe columna Mail, Email o Correo para enviar a destinatarios"}
	}
	instance := dynamicstruct.NewStruct()
	for _, v := range colNames {
		instance.AddField(v, "", fmt.Sprintf(`json:"%s"`, strings.ToLower(v)))
	}
	dynamicStruct := instance.Build().New()

	data, err = excelFile.GetRows(excelFile.GetSheetName(0))
	if err != nil {
		return SumarizeResponse{Error: err.Error()}
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
		response.Success.Data = append(response.Success.Data, dynamicStruct)
	}
	response.Success.Limit = len(data) - 1
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
func (resp *SumarizeResponse) ValidateLowerUpper(s string) {
	for i, v := range s {
		if i == 0 {
			if !unicode.IsUpper(v) {
				resp.Error = "Verifique que los nombres de columbas empiecen con mayúscula"
			}
		} else {
			if !unicode.IsLower(v) {
				resp.Error = "Verifique que los nombres de columbas empiecen con mayúscula y luego minúscula"
			}
		}
	}
}
