package define

const ()

type BaseResult struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg,omitempty"`
}

type PagingParams struct {
	PageSize int `form:"page_size"`
	PageNum  int `form:"page_num"`
}
