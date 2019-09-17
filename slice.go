package ecstaskmonitoring

// IsTaskContains ... tells whether s contains ss.
func IsTaskContains(s []*Task, ss string) bool {
	if s == nil {
		return false
	}

	for _, v := range s {
		if ss == v.Name {
			return true
		}
	}
	return false
}
