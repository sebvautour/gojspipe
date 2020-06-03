package utils

import "github.com/sirupsen/logrus"

func (u *Utils) logger() *logrus.Entry {
	if u.logrusEntry == nil {
		return nil
	}
	return u.logrusEntry.WithField("script", u.ScriptName)
}

// Log adds an info level log entry
func (u *Utils) Log(m interface{}) {
	u.LogInfo(m)
}

// LogInfo add an info level log entry
func (u *Utils) LogInfo(m interface{}) {
	l := u.logger()
	if l == nil {
		return
	}
	l.Info(m)
}

// LogWarning add a warning level log entry
func (u *Utils) LogWarning(m interface{}) {
	l := u.logger()
	if l == nil {
		return
	}
	l.Warn(m)
}

// LogWarn is a shorthand for LogWarning
func (u *Utils) LogWarn(m interface{}) {
	u.LogWarn(m)
}

// LogError add a error level log entry
func (u *Utils) LogError(m interface{}) {
	l := u.logger()
	if l == nil {
		return
	}
	l.Error(m)
}

// LogErr is a shorthand for LogError
func (u *Utils) LogErr(m interface{}) {
	u.LogError(m)
}
