// Package auditlog provides structured logging for auditing purposes
package auditlog

import (
	"github.com/fatih/structs"
	"github.com/sirupsen/logrus"
)

var auditFields = logrus.Fields{"auditlog": "true"}

// Entry is used to control what is tracked in the audit log
type Entry struct {
	CurrentUserUUID  string
	CompanyUUID      string
	TeamUUID         string
	Authorization    string
	TargetType       string
	TargetUUID       string
	OriginalContents interface{}
	UpdatedContents  interface{}
}

// Log sends an audit log event based on an entry data structure
func (a *Entry) Log(logger *logrus.Entry, action string) {
	logger = logger.WithFields(auditFields)
	entryMap := structs.Map(a)
	for k, v := range entryMap {
		logger = logger.WithFields(logrus.Fields{k: v})
	}
	logger.Infof(action)
}
