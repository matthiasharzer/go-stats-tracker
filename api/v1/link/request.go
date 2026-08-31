package link

type RequestBody struct {
	SheetID       string `json:"sheet_id"`
	UserAccessKey string `json:"user_access_key"`
	State         string `json:"state"`
}
