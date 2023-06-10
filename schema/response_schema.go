package schema

type TSBModel struct {
	SelectedYear any      `json:"selected_year"`
	AllYears     []string `json:"all_years"`
	Files        []struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		TaskType   string `json:"task_type"`
		FilePath   string `json:"file_path"`
		MtvYear    any    `json:"mtv_year"`
		KaskoYear  string `json:"kasko_year"`
		KaskoMonth string `json:"kasko_month"`
	} `json:"files"`
	ActualFile struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		TaskType   string `json:"task_type"`
		FilePath   string `json:"file_path"`
		MtvYear    any    `json:"mtv_year"`
		KaskoYear  string `json:"kasko_year"`
		KaskoMonth string `json:"kasko_month"`
	} `json:"actual_file"`
	ArchiveYears []string `json:"archive_years"`
}
