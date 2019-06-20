package common

import "github.com/spf13/pflag"

func AddCommonFlags(flags *pflag.FlagSet, args *Args) {
	flags.BoolVar(
		&args.Debug,
		"debug",
		false,
		"Enable debug mode.",
	)
	flags.StringArrayVar(
		&args.Parameter,
		"parameter",
		nil,
		"Query parameters to add to the request. The value must be the name of the "+
			"parameter, followed by an optional equals sign and then the value "+
			"of the parameter. Can be used multiple times to specify multiple "+
			"parameters or multiple values for the same parameter.",
	)
	flags.StringArrayVar(
		&args.Header,
		"header",
		nil,
		"Headers to add to the request. The value must be the name of the header "+
			"followed by an optional equals sign and then the value of the "+
			"header. Can be used multiple times to specify multiple headers "+
			"or multiple values for the same header.",
	)
}
