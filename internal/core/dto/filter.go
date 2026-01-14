package dto

// Filter represents base filtering options for queries
type Filter struct {
	Search string `query:"search" form:"search"`
	ID     *uint  `query:"id" form:"id"`
	Page   int    `query:"page" form:"page"`
	Limit  int    `query:"limit" form:"limit"`
	Sort   string `query:"sort" form:"sort"`
	Order  string `query:"order" form:"order"`
}

// ProfileFilter represents filtering options for profiles
type ProfileFilter struct {
	Filter
	WithPermissions *bool `query:"with_permissions" form:"with_permissions"`
	ListRoot        bool  `query:"list_root" form:"list_root"`
}

// UserFilter represents filtering options for users
type UserFilter struct {
	Filter
	ProfileID uint  `query:"profile_id" form:"profile_id"`
	Status    *bool `query:"status" form:"status"`
}

// ApplyPagination returns pagination values
func (f *Filter) ApplyPagination() (enabled bool, offset, limit int) {
	if f.Page > 0 && f.Limit > 0 {
		return true, (f.Page - 1) * f.Limit, f.Limit
	}
	return false, 0, 0
}

// CalcPages calculates total pages based on count
func (f *Filter) CalcPages(count int64) int64 {
	if count == 0 || f.Limit == 0 || f.Page == 0 {
		if count > 0 {
			return 1
		}
		return 0
	}
	return (count + int64(f.Limit) - 1) / int64(f.Limit)
}
