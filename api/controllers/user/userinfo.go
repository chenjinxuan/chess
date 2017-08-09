package c_user

type UserInfo struct { 
    Id int `json:"id"`
    NickName string `json:"nick_name"`
    MobileNumber string `json:"mobile_number"`
    Gender  int `json:"gender"`
    Avatar string `json:"avatar"`
    Type int `json:"type"`
    Status int `json:"status"`
    IsFresh int `json:"is_fresh"`
    Balance int `json:"balance"`
   // Diamonds
}
