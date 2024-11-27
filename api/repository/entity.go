package repository

type CustomerInfo struct {
	FirstName      *string  `json:"firstName" pg:",use_zero"`
	LastName       *string  `json:"lastName,omitempty" pg:",use_zero"`
	PatronymicName *string  `json:"patronymicName,omitempty" pg:",use_zero"`
	Phone          *string  `json:"phone" pg:",use_zero"`
	Email          *string  `json:"email" pg:",use_zero"`
	Id             int      `json:"id" pg:",pk"`
	tableName      struct{} `json:"-" pg:"customer"`
}

type DeleteInfo struct {
	Count int   `json:"count"`
	Ids   []int `json:"ids"`
}
