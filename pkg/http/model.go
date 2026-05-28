package xhttp

// type ResponseOption func (*responseOption)

type SuccessResDto struct {
	// Code    string      `json:"code,omitempty" example:"OK"`
	Status  int         `json:"status" example:"200"`
	Message string      `json:"message" example:"call api success"`
	Data    interface{} `json:"data,omitempty"`
	Doc     string      `json:"doc,omitempty"`
}

type ResponseErrDto struct {
	// Code    string `json:"code,omitempty" example:"BAD_REQUEST"`
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"bad request"`
	Doc     string `json:"doc"`
}

type ValidationErrDto struct {
	Code    string `json:"code,omitempty" example:"ERR_REQUIRED"`
	Field   string `json:"field,omitempty" example:"name"`
	Message string `json:"message,omitempty" example:"Name is required"`
}

type ListDataResponseDto struct {
	Rows  interface{} `json:"rows"`
	Total int64       `json:"total,omitempty"`
}

type PaginationOptionsDto struct {
	Page         int  `query:"page" default:"1" validate:"omitempty,min=1"`
	Limit        int  `query:"limit" default:"10" validate:"omitempty,min=1,max=100"`
	ExcludeTotal bool `query:"exclude_total" validate:"omitempty"`
}

type DateRangeOptionsDto struct {
	FromDate string `query:"from_date" validate:"omitempty,datetime=2006-01-02"`
	ToDate   string `query:"to_date" validate:"omitempty,datetime=2006-01-02"`
	RangeBy  string `query:"range_by" default:"created_at" validate:"omitempty,oneof=created_at updated_at"`
}

type SortOptionsDto struct {
	SortBy  string `query:"sort_by" default:"created_at" validate:"omitempty,oneof=created_at updated_at"`
	OrderBy string `query:"order_by" default:"desc" validate:"omitempty,oneof=asc desc"`
}
