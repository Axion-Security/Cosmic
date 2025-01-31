package executor

var TargetFileArg = "%FILE%"
var TargetFile string

func ReplaceArgs(args []string) []string {
	for i, arg := range args {
		if arg == TargetFileArg {
			args[i] = TargetFile
		}
	}
	return args
}
