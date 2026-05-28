package response

type SuccessResDto struct {
    Status  int            `json:"status"          example:"200"`
    Message string         `json:"message"         example:"OK"`
    Data    any            `json:"data,omitempty"`
    Doc     string         `json:"doc,omitempty"`
    Meta    map[string]any `json:"meta,omitempty"`
}

type ErrorResDto struct {
    Status  int            `json:"status"          example:"400"`
    Message string         `json:"message"         example:"Bad Request"`
    Code    string         `json:"code,omitempty"  example:"USER_NOT_FOUND"`
    Data    any            `json:"data,omitempty"`
    Doc     string         `json:"doc,omitempty"`
    Meta    map[string]any `json:"meta,omitempty"`
}

type ListDataDto[T any] struct {
    Rows  []T   `json:"rows"`
    Total int64 `json:"total,omitempty"`
}