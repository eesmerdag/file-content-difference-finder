package diff_finder

type Range struct {
	start int
	end   int
}

type ChangeTypes string

const (
	ChangeTypesUpdated ChangeTypes = "UPDATED"
	ChangeTypesRemoved ChangeTypes = "REMOVED"
	ChangeTypesAdded   ChangeTypes = "ADDED"
)

type UpdatedIndex struct {
	OldValue string
	NewValue string
	Index    int
	Type     ChangeTypes
}
