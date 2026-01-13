package data_splitting

import (
	"errors"
)

var (
	// delimiterRE tokenizes the following string into walkable with extracted delimiter + key.
	// string:
	// ` %{key}, %{key/2}`
	// into:
	// [["", "key" ], [", ", "key/2"]]
	ordinalIndicator     = "/"
	fixedLengthIndicator = "#"

	skipFieldPrefix      = "?"
	appendFieldPrefix    = "+"
	indirectFieldPrefix  = "&"
	appendIndirectPrefix = "+&"
	indirectAppendPrefix = "&+"
	greedySuffix         = "->"
	pointerFieldPrefix   = "*"
	dataTypeIndicator    = "|"
	dataTypeSeparator    = "\\|" // Needed for regexp

	numberRE = "\\d{1,2}"
	alphaRE  = "[[:alpha:]]*"

	defaultJoinString = " "

	errParsingFailure            = errors.New("parsing failure")
	errEmpty                     = errors.New("empty string provided")
	errMixedPrefixIndirectAppend = errors.New("mixed prefix `&+`")
	errMixedPrefixAppendIndirect = errors.New("mixed prefix `&+`")
	errEmptyKey                  = errors.New("empty key")
	errInvalidDatatype           = errors.New("invalid data type")
	errMissingDatatype           = errors.New("missing data type")
)
