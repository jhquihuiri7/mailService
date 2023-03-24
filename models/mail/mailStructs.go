package mail

type Clients struct {
	List []Client
}
type Client struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Answer int    `json:"answer"`
	Templates
	Auth
}
type Templates struct {
	TemplateReceive string `json:"templateReceive"`
	TemplateSend    string `json:"templateSend"`
}
type Auth struct {
	Sender   string `json:"from"`
	Alias    string `json:"alias"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}
