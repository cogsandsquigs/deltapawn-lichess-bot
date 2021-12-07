package lichessapi


type User struct {
	Id		    string `json:"id"`
	Name		string `json:"name"`
	Title		string `json:"title"`
	Online		bool   `json:"online"`
	Playing		bool   `json:"playing"`
	Streaming	bool   `json:"streaming"`
	Patron		bool   `json:"patron"`
}