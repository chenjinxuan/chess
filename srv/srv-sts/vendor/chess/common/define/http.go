package define

var (
	RetFail = 0

	MsgInternalServiceError = "Internal service error!"
	MsgParamsError          = "Params invalid!"
)

type BaseResult struct {
	Ret int    `json:"ret" description:"1成功 0失败"`
	Msg string `json:"msg,omitempty"`
}

type PagingParams struct {
	PageSize int `form:"page_size" description:"每页显示个数"`
	PageNum  int `form:"page_num" description:"第几页"`
}
