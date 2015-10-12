package logging

import (
	"fmt"
	"os"
	"util"
)

type Fields map[string]interface{}

var localhostFields Fields

func init() {
	localhostFields = getLocalhostFields()
}

func (lhs Fields) Update(rhs Fields) Fields {
	for k, v := range rhs {
		lhs[k] = v
	}
	return lhs
}

func (fields Fields) populateStandardFields(lvl level, name string) {
	fields["level"] = levelToString(lvl)
	fields["logger"] = name
	fields.Update(localhostFields)
	/*
		TODO: these are present in python but not yet implemented here:
			path to file (available from runtime.Caller)
			function name (available from runtime.Stack)
			line number (available from runtime.Caller)
			error defaults to filepath:function
			exception (passed in as an arg)
				- type
				- trace
				- message
	*/
}

func (fields Fields) Dupe() Fields {
	dupe := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		dupe[k] = v
	}
	return dupe
}

func (fields Fields) WithError(err error) Fields {
	res := fields.Dupe()
	res["error_message"] = fmt.Sprintf("%v", err)
	return res
}

func getLocalhostFields() Fields {
	fields := make(map[string]interface{})
	fields["pid"] = os.Getpid()
	fields["executable"] = os.Args[0]
	fields["argv"] = os.Args

	localhost, err := os.Hostname()
	if err != nil {
		initError(fmt.Sprintf("Unable to lookup localhost hostname.\n"))
		return fields
	}
	hostInfo := util.GetHostInfo(localhost)
	if hostInfo == nil {
		initError(fmt.Sprintf("Unable to extract host info from %v.\n", localhost))
		return fields
	}

	for k, v := range hostInfo.Map() {
		fields[k] = v
	}
	return fields
}
