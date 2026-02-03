package dto

// RoleOutput represents output data for a role
type RoleOutput struct {
	ID          *string   `json:"id,omitempty"`
	Name        *string   `json:"name,omitempty"`
	Permissions *[]string `json:"permissions,omitempty"`
	Enabled     *bool     `json:"enabled,omitempty"`
}

// UserOutput represents output data for a user
type UserOutput struct {
	ID       *string       `json:"id,omitempty"`
	Name     *string       `json:"name,omitempty"`
	Username *string       `json:"username,omitempty"`
	Email    *string       `json:"email,omitempty"`
	Status   *bool         `json:"status,omitempty"`
	New      *bool         `json:"new,omitempty"`
	Roles    []*RoleOutput `json:"roles,omitempty"`
}

// AuthOutput represents output data for authentication
type AuthOutput struct {
	User         *UserOutput `json:"user,omitempty"`
	AccessToken  string      `json:"accesstoken"`
	RefreshToken string      `json:"refreshtoken"`
}

// ItemOutput represents a simple item output (id + name)
type ItemOutput struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// PaginationOutput represents pagination metadata
type PaginationOutput struct {
	Page       uint `json:"page"`
	Limit      uint `json:"limit"`
	TotalItems uint `json:"total_items"`
	TotalPages uint `json:"total_pages"`
}

// paginableOutput defines which types can be used in PaginatedOutput
type paginableOutput interface {
	RoleOutput | UserOutput | ItemOutput
}

// PaginatedOutput represents a paginated list of items
// T must be one of: RoleOutput, UserOutput, ItemOutput
type PaginatedOutput[T paginableOutput] struct {
	Items      []T              `json:"items"`
	Pagination PaginationOutput `json:"pagination"`
}

// NewPaginatedOutput creates a new paginated output
func NewPaginatedOutput[T paginableOutput](items []T, page, limit int, totalItems int64) *PaginatedOutput[T] {
	totalPages := uint(0)
	if limit > 0 && totalItems > 0 {
		totalPages = uint((totalItems + int64(limit) - 1) / int64(limit))
	} else if totalItems > 0 {
		totalPages = 1
	}

	actualPage := uint(page)
	if actualPage == 0 {
		actualPage = 1
	}

	actualLimit := uint(limit)
	if actualLimit == 0 {
		actualLimit = uint(len(items))
	}

	return &PaginatedOutput[T]{
		Items: items,
		Pagination: PaginationOutput{
			Page:       actualPage,
			Limit:      actualLimit,
			TotalItems: uint(totalItems),
			TotalPages: totalPages,
		},
	}
}
