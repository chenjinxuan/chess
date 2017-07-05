package main

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/vmihailenco/msgpack.v2"
)

import (
	"rank/dos"
	"rank/ss"
)

const (
	UPPER_THRESHOLD = 1024 // storage changed to tree when elements exceeds this
	LOWER_THRESHOLD = 512  // storage changed to sortedset when elements below this
)

const (
	SORTEDSET = iota
	RBTREE
)

// a ranking set
type RankSet struct {
	R    dos.Tree        // rbtree
	S    ss.SortedSet    // sorted-set
	M    map[int32]int32 // ID  => SCORE
	Type int
	sync.RWMutex
}

func NewRankSet() *RankSet {
	r := new(RankSet)
	r.M = make(map[int32]int32)
	r.Type = SORTEDSET // default in sortedset
	return r
}

// toggle storage base on Type
func (r *RankSet) toggle() {
	switch r.Type {
	case SORTEDSET:
		for k, v := range r.M {
			r.R.Insert(v, k)
		}
		r.S.Clear()
		r.Type = RBTREE
		log.Debugf("convert sortedset to rbtree %v", len(r.M))
	case RBTREE:
		for k, v := range r.M {
			r.S.Insert(v, k)
		}
		r.R.Clear()
		r.Type = SORTEDSET
		log.Debugf("convert rbtree to sortedset %v", len(r.M))
	}
}

func (r *RankSet) Update(id, newscore int32) {
	r.Lock()
	defer r.Unlock()

	oldscore, ok := r.M[id]
	if !ok { // new element
		if r.Type == SORTEDSET && len(r.M) > UPPER_THRESHOLD { // do convert
			r.toggle()
		}

		switch r.Type {
		case SORTEDSET:
			r.S.Insert(id, newscore)
		case RBTREE:
			r.R.Insert(newscore, id)
		}
	} else {
		switch r.Type {
		case SORTEDSET:
			r.S.Update(id, newscore)
		case RBTREE:
			_, n := r.R.Locate(oldscore, id)
			r.R.Delete(id, n)
			r.R.Insert(newscore, id)
		}
	}
	r.M[id] = newscore
	return
}

func (r *RankSet) Delete(userid int32) {
	r.Lock()
	defer r.Unlock()
	if r.Type == RBTREE && len(r.M) < LOWER_THRESHOLD { // do convert
		r.toggle()
	}

	switch r.Type {
	case SORTEDSET:
		r.S.Delete(userid)
	case RBTREE:
		score := r.M[userid]
		_, n := r.R.Locate(score, userid)
		r.R.Delete(userid, n)
	}
	delete(r.M, userid)
}

func (r *RankSet) Count() int32 {
	r.RLock()
	defer r.RUnlock()
	return int32(len(r.M))
}

// range [A,B]
func (r *RankSet) GetList(A, B int) (ids []int32, scores []int32) {
	if A < 1 || A > B {
		return
	}
	r.RLock()
	defer r.RUnlock()

	if A > len(r.M) {
		return
	}

	if B > len(r.M) {
		B = len(r.M)
	}

	switch r.Type {
	case SORTEDSET:
		ids, scores = r.S.GetList(A, B)
	case RBTREE:
		ids, scores = r.R.GetList(A, B)
	}
	return
}

// rank of a user
func (r *RankSet) Rank(userid int32) (rank int32, score int32) {
	r.RLock()
	defer r.RUnlock()

	switch r.Type {
	case SORTEDSET:
		rankno := r.S.Locate(userid)
		return int32(rankno), r.M[userid]
	case RBTREE:
		rankno, _ := r.R.Locate(r.M[userid], userid)
		return int32(rankno), r.M[userid]
	}
	return
}

// serialization
func (r *RankSet) Marshal() ([]byte, error) {
	r.RLock()
	defer r.RUnlock()
	return msgpack.Marshal(r.M)
}

func (r *RankSet) Unmarshal(bin []byte) error {
	m := make(map[int32]int32)
	r.Lock()
	defer r.Unlock()
	err := msgpack.Unmarshal(bin, &m)
	if err != nil {
		return err
	}

	r.M = m
	if len(r.M) > UPPER_THRESHOLD {
		for id, score := range m {
			r.R.Insert(score, id)
		}
		r.Type = RBTREE
		log.Debugf("rank restored into rbtree %v", len(r.M))
	} else {
		for id, score := range m {
			r.S.Insert(id, score)
		}
		r.Type = SORTEDSET
		log.Debugf("rank restored into sortedset %v", len(r.M))
	}

	return nil
}
