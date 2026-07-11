package models

type Record struct {
	ID          string  `json:"ID"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	CategoryID  *string `json:"categoryID,omitempty"`
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Note        string  `json:"note"`
	Balance     float64 `json:"balance"`
}

type GroupedRecord struct {
	Group string  `json:"group"`
	Total float64 `json:"total"`
	Count int     `json:"count"`
}

type PaginationMetadata struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	TotalPages int  `json:"totalPages,omitempty"`
	TotalCount int  `json:"totalCount,omitempty"`
	HasPrev    bool `json:"hasPrev,omitempty"`
	HasNext    bool `json:"hasNext,omitempty"`
}

type SearchRecordsResult struct {
	Records []Record        `json:"records,omitempty"`
	Groups  []GroupedRecord `json:"groups,omitempty"`
	PaginationMetadata
}

type Category struct {
	ID    string `json:"ID"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

type Budget struct {
	ID         string  `json:"ID"`
	CategoryID string  `json:"categoryID"`
	Category   string  `json:"category"`
	Month      int     `json:"month"`
	Year       int     `json:"year"`
	Amount     float64 `json:"amount"`
}

type BudgetProgress struct {
	Budget
	Spent      float64 `json:"spent"`
	Percentage float64 `json:"percentage"`
}

type Goal struct {
	ID                  string  `json:"ID"`
	Name                string  `json:"name"`
	TargetAmount        float64 `json:"targetAmount"`
	CurrentAmount       float64 `json:"currentAmount"`
	TargetDate          string  `json:"targetDate,omitempty"`
	Category            string  `json:"category,omitempty"`
	CategoryID          *string `json:"categoryID,omitempty"`
	Status              string  `json:"status"`
	Description         string  `json:"description,omitempty"`
	MonthlyContribution float64 `json:"monthlyContribution,omitempty"`
}

type CategoryDetail struct {
	ID         string  `json:"ID"`
	CategoryID string  `json:"categoryID"`
	Category   string  `json:"category"`
	Amount     float64 `json:"amount"`
}

type Summary struct {
	Expenses     []CategoryDetail `json:"expenses"`
	Incomes      []CategoryDetail `json:"incomes"`
	Net          float64          `json:"net"`
	Opening      float64          `json:"opening"`
	Closing      float64          `json:"closing"`
	TotalExpense float64          `json:"totalExpense"`
	TotalIncome  float64          `json:"totalIncome"`
}

type SearchRecordsParams struct {
	FromDate   string
	ToDate     string
	Category   string
	RecordType string
	MinAmount  float64
	MaxAmount  float64
	Search     string
	GroupBy    string
	Page       int
	Limit      int
}

type SummaryParams struct {
	FromDate   string
	ToDate     string
	Category   string
	RecordType string
}
