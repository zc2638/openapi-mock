package config

/**
 * Created by zc on 2019-11-21.
 */
const ConfigPath = "config.yml"

const ServerPort = "8080"

var OpenTracingHeaders = []string{
	"x-request-id",
	"x-b3-traceid",
	"x-b3-spanid",
	"x-b3-parentspanid",
	"x-b3-sampled",
	"x-b3-flags",
	"x-ot-span-context",
}