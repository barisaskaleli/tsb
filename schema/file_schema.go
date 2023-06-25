package schema

type TSBFile struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	TaskType   string `json:"task_type"`
	FilePath   string `json:"file_path"`
	MtvYear    any    `json:"mtv_year"`
	KaskoYear  string `json:"kasko_year"`
	KaskoMonth string `json:"kasko_month"`
}
