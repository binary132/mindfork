package core_test

func sampleMessages(s string) []byte {
	return map[string][]byte{
		"emptyMessage": []byte(``),
		"emptyObject":  []byte(`{}`),
		"emptyString":  []byte(`""`),
		"validIntention": []byte(`
{"Who": "Bodie", "What": "To seek the Holy Grail"}`[1:]),
		"timedIntention": []byte(`
{"Who": "Bodie", "What": "Something neat", "When": "2009-11-10T23:00:00Z"}`[1:]),
	}[s]
}
