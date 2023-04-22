package request

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
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
	ClientName string              `json:"clientName"`
	Template   string              `json:"template"`
	Tos        []map[string]string `json:"tos"`
	Limits     []int               `json:"limits"`
}
type RequestTemplate struct {
	ClientName string   `json:"clientName"`
	Template   string   `json:"template"`
	Columns    []string `json:"columns"`
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
	Data    []map[string]string `json:"data"`
	Limit   int                 `json:"limit"`
	Columns []string            `json:"columns"`
}

type ListBulk struct {
	List []RequestBulk
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
		colNames = append(colNames, strings.ToLower(v[0]))
	}
	if !slices.Contains(colNames, "mail") && !slices.Contains(colNames, "email") && !slices.Contains(colNames, "correo") {
		return SumarizeResponse{Error: "No existe columna Mail, Email o Correo para enviar a destinatarios"}
	}
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

			rawData += "\"" + strings.ToUpper(colNames[index][:1]) + strings.ToLower(colNames[index][1:]) + "\":\"" + val + "\","
			if index == lenData-1 {
				rawData = rawData[:len(rawData)-1]
			}
		}
		rawData += "}"
		var result map[string]string
		err = json.Unmarshal([]byte(rawData), &result)
		if err != nil {
			log.Fatal(err)
		}
		response.Success.Data = append(response.Success.Data, result)
		response.Success.Columns = colNames
		r.Tos = response.Success.Data
	}
	response.Success.Limit = len(data) - 1
	return response
}

func (r *RequestTemplate) ParseRequestBulkTemplate(c *gin.Context) {
	err := json.NewDecoder(c.Request.Body).Decode(&r)
	if err != nil {
		log.Fatal(err)
	}
}
func (r *RequestTemplate) ValidateTemplate() RequestResponse {
	var response RequestResponse
	response.Error = "Columna "
	for _, v := range r.Columns {
		v = strings.ToLower(v)
		if v == "mail" || v == "email" || v == "correo" {
			continue
		}
		v = strings.ToUpper(v[:1]) + v[1:]
		if !strings.Contains(r.Template, "{{."+v+"}}") {
			response.Error += strings.ToLower(v) + ","
		}
	}
	if strings.Contains(response.Error, ",") {
		response.Error = response.Error[:len(response.Error)-1] + " no existen en template ingresado"
	} else {
		response.Error = ""
		response.Success = "Template validado con éxito"
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
func (list *ListBulk) GetRequestItemTmp(request RequestTemplate) RequestResponse {
	for i, v := range list.List {
		if v.ClientName == request.ClientName {
			v.Template = request.Template
			list.List[i] = v
			return RequestResponse{Success: "Template validado con éxito"}
		}
	}
	return RequestResponse{Error: "No se encontró cliente"}
}
func (list *ListBulk) GetRequestItemLimits(request RequestBulk) (RequestBulk, RequestResponse) {
	for i, v := range list.List {
		if v.ClientName == request.ClientName {
			v.Limits = request.Limits
			list.List[i] = v
			return list.List[i], RequestResponse{Success: "Template validado con éxito"}
		}
	}
	return RequestBulk{}, RequestResponse{Error: "No se encontró cliente"}
}
