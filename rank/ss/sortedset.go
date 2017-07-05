package ss

type sortpair struct {
	id    int32
	score int32
}

type SortedSet struct {
	set []sortpair
}

func (ss *SortedSet) Clear() {
	ss.set = nil
}

func (ss *SortedSet) Insert(id, score int32) {
	p := sortpair{id: id, score: score}
	if len(ss.set) == 0 {
		ss.set = []sortpair{p}
		return
	}

	// grow
	ss.set = append(ss.set, p)
	update_idx := -1
	for k := range ss.set {
		if score > ss.set[k].score {
			update_idx = k
			break
		}
	}

	if update_idx == -1 { // already appended
		return
	}

	ss.rshift(update_idx, len(ss.set)-1, p)
}

func (ss *SortedSet) Delete(id int32) {
	for k := range ss.set {
		if ss.set[k].id == id {
			ss.set = append(ss.set[:k], ss.set[k+1:]...)
			return
		}
	}
}

func (ss *SortedSet) Locate(id int32) int32 {
	for k := range ss.set {
		if ss.set[k].id == id {
			return int32(k + 1)
		}
	}
	return -1
}

func (ss *SortedSet) Update(id, score int32) {
	p := sortpair{id: id, score: score}
	idx := -1
	update_idx := -1
	for k := range ss.set {
		if idx == -1 && ss.set[k].id == id {
			idx = k
		}
		if update_idx == -1 && score > ss.set[k].score {
			update_idx = k // insert point
		}

		if idx != -1 && update_idx != -1 { // both set, break
			break
		}
	}
	if idx == -1 {
		return
	}

	n := len(ss.set) - 1
	// smallest
	if update_idx == -1 {
		update_idx = n
	}

	// shift
	if update_idx > idx {
		ss.lshift(idx, update_idx, p)
	} else if update_idx < idx {
		ss.rshift(update_idx, idx, p)
	} else {
		ss.set[idx] = p
	}
}

// shift left [i,j]
func (ss *SortedSet) lshift(i, j int, p sortpair) {
	copy(ss.set[i:j+1], ss.set[i+1:j+1])
	ss.set[j] = p
}

// shift right [i,j]
func (ss *SortedSet) rshift(i, j int, p sortpair) {
	copy(ss.set[i+1:j+1], ss.set[i:j+1])
	ss.set[i] = p
}

func (ss *SortedSet) GetList(a, b int) (ids []int32, scores []int32) {
	ids, scores = make([]int32, b-a+1), make([]int32, b-a+1)
	for k := a - 1; k <= b-1; k++ {
		ids[k-a+1] = ss.set[k].id
		scores[k-a+1] = ss.set[k].score
	}
	return
}
