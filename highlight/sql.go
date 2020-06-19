package highlight

// SQLConfig provides highlight config of sql.
var SQLConfig = &Config{
	DataTypes: []string{
		"bigint", "binary", "bit", "blob", "char", "date", "datetime", "decimal", "double", "enum", "float",
		"int", "json", "longblob", "longtext", "mediumblob", "mediumint", "mediumtext", "nchar", "nvarchar", "set", "smallint",
		"text", "time", "timestamp", "tiny", "tinyblob", "tinyint", "tinytext", "varbinary", "varchar", "year",
	},
	Keywords: []string{
		"accessible", "add", "all", "alter", "analyze", "and", "as", "asc", "asensitive",
		"before", "between", "bigint", "binary", "blob", "both", "by",
		"call", "cascade", "case", "change", "char", "character", "charset", "check", "collate", "column", "comment", "condition", "connection", "constraint", "continue", "contributors", "convert", "create", "cross", "current_date", "current_time", "current_timestamp", "current_user", "cursor",
		"database", "databases", "day_hour", "day_microsecond", "day_minute", "day_second", "dec", "decimal", "declare", "default", "delayed", "delete", "desc", "describe", "deterministic", "distinct", "distinctrow", "div", "double", "drop", "dual",
		"each", "else", "elseif", "enclosed", "engine", "escaped", "exists", "exit", "explain",
		"false", "fetch", "float", "float4", "float8", "for", "force", "foreign", "from", "fulltext",
		"grant", "group",
		"having", "high_priority", "hour_microsecond", "hour_minute", "hour_second",
		"if", "ignore", "in", "index", "infile", "inner", "inout", "insensitive", "insert", "int", "int1", "int2", "int3", "int4", "int8", "integer", "interval", "into", "is", "iterate",
		"join",
		"key", "keys", "kill",
		"leading", "leave", "left", "like", "limit", "linear", "lines", "load", "localtime", "localtimestamp", "lock", "long", "longblob", "longtext", "loop", "low_priority",
		"match", "mediumblob", "mediumint", "mediumtext", "middleint", "minute_microsecond", "minute_second", "mod", "modifies",
		"natural", "not", "no_write_to_binlog", "null", "numeric",
		"on", "optimize", "option", "optionally", "or", "order", "out", "outer", "outfile",
		"precision", "primary", "procedure", "purge",
		"range", "read", "reads", "read_only", "read_write", "real", "references", "regexp", "release", "rename", "repeat", "replace", "require", "restrict", "return", "revoke", "right", "rlike",
		"schema", "schemas", "second_microsecond", "select", "sensitive", "separator", "set", "show", "smallint", "spatial", "specific", "sql", "sqlexception", "sqlstate", "sqlwarning", "sql_big_result", "sql_calc_found_rows", "sql_small_result", "ssl", "starting", "straight_join",
		"table", "terminated", "then", "tinyblob", "tinyint", "tinytext", "to", "trailing", "trigger", "true",
		"undo", "union", "unique", "unlock", "unsigned", "update", "upgrade", "usage", "use", "using", "utc_date", "utc_time", "utc_timestamp",
		"values", "varbinary", "varchar", "varcharacter", "varying", "when", "where", "while", "with", "write", "x509", "xor", "year_month", "zerofill",
	},
	KeySymbols: []rune{'(', ')', ',', ';', '='},
	Blocks: TextBlocks{
		{
			BeginIdentifier: "'",
			EndIdentifier:   "'",
			EscapeChar:      '\\',
			TextType:        String,
		},
		{
			BeginIdentifier: `"`,
			EndIdentifier:   `"`,
			EscapeChar:      '\\',
			TextType:        String,
		},
		{
			BeginIdentifier: "#",
			EndIdentifier:   "\n",
			TextType:        Comment,
		},
		{
			BeginIdentifier: "--",
			EndIdentifier:   "\n",
			TextType:        Comment,
		},
		{
			BeginIdentifier: "/*",
			EndIdentifier:   "*/",
			EscapeChar:      '\\',
			TextType:        Comment,
		},
	},
}
