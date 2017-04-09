package main 

// Thread safe data copy from one slice to another
func tsdatacopy(copyfrom int, copyto int, list []string) []string {
	//protect memory by copying only what we know we've got
	var res []string
	copyto = len(linkpool)
	if copyto > 0 && copyto > copyfrom {
		res = make([]string, copyto-copyfrom)
		copy(res, list[copyfrom:copyto])
		copyfrom = copyto
	}
	return res
}
