package testing

// SampleMessages returns a sample Message.
func SampleMessages(s string) string {
	return map[string]string{
		"emptyMessage": ``,
		"emptyObject":  `{}`,
		"emptyString":  `""`,
		"validIntention": `
{"Who": "Bodie", "What": "To seek the Holy Grail"}`[1:],
		"timedIntention": `
{"Who": "Bodie", "What": "Something neat", "When": "2009-11-10T23:00:00Z"}`[1:],
	}[s]
}
