package clause

func _limit(values ...any) (string, []any) {
	return "LIMIT ?", values
}
