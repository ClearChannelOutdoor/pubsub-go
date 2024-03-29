package pb

func merge(attrs ...map[string]string) map[string]string {
	// merge attributes
	mgd := make(map[string]string)

	// merge all the dictionaries into a single dictionary
	for _, mp := range attrs {
		for k, v := range mp {
			mgd[k] = v
		}
	}

	return mgd
}
