package core

import "github.com/bsphere/le_go"

/*
 * Get logger according to rapid7 configurations
 * It will return a Rapid7 logger if LogEntries token is defined.
 * If not, it will return nil
 */
func GetLogger() *le_go.Logger {
	r7Token := GetEnv("LOGENTRIES_TOKEN", "")
	if len(r7Token) > 0 {
		var err error
		le, err := le_go.Connect(r7Token)
		if err != nil {
			panic("unable to setup rapid7 connection with provided token")
		}
		return le
	}
	return nil
}
