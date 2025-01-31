package model

/* ---- Requests ---- */
// type Response struct {
// 	Success bool        `json:"success"`
// 	Msg     string      `json:"msg"`
// 	Obj     interface{} `json:"obj"`
// }

/* ---- Login ---- */
const LoginPath = "/panel/api/auth/login"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginObj struct {
}

/* ---- Inbound ---- */
const GetInboundPath = "/panel/api/inbounds/get/{{.inboundId}}"

type GetInboundRequest struct {
}

type GetInboundObj struct {
	Inbound
}

const AddInboundPath = "/panel/api/inbounds/add"

type AddInboundRequest struct {
	Inbound
}

type AddInboundObj struct {
	Inbound
}

const InboundsPath = "/panel/api/inbounds/list"

type InboundsRequest struct {
}

type InboundsObj struct {
	Inbounds []Inbound `json:"inbounds"`
}
