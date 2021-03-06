package main

// Thread safe data copy from one slice to another
func tsdatacopy(copyfrom *int, copyto *int, list []string) []string {
	//protect memory by copying only what we know we've got
	*copyto = len(list)
	if *copyto > 0 && *copyto > *copyfrom {
		var res []string
		res = make([]string, *copyto-*copyfrom)
		copy(res, list[*copyfrom:*copyto])
		*copyfrom = *copyto
		return res
	}
	return []string{}
}

// Thread safe data copy from one slice to another
func pldatacopy(copyfrom *int, copyto *int, list []processLog) []processLog {
	//protect memory by copying only what we know we've got
	*copyto = len(list)
	if *copyto > 0 && *copyto > *copyfrom {
		var res []processLog
		res = make([]processLog, *copyto-*copyfrom)
		copy(res, list[*copyfrom:*copyto])
		*copyfrom = *copyto
		return res
	}
	return []processLog{}
}

// Thread safe data copy from one slice to another
// This method allows us to specify a length which may be safer for
// us in the long run...
func pldatacopylen(copyfrom *int, copyto *int, list []processLog, copylen int) []processLog {
	//protect memory by copying only what we know we've got
	*copyto = *copyto + 1
	if *copyto > 0 && *copyto > *copyfrom && *copyto <= len(list) {
		var res []processLog
		res = make([]processLog, *copyto-*copyfrom)
		copy(res, list[*copyfrom:*copyto])
		*copyfrom = *copyto
		return res
	}
	return []processLog{}
}
