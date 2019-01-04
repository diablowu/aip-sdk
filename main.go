package baidu_ai_sdk

const (
	aip_endpoint   = "https://aip.baidubce.com"
	api_token      = aip_endpoint + "/oauth/2.0/token"
	api_unit_skill = aip_endpoint + "/rpc/2.0/unit/bot/chat"
	api_unit_bot   = aip_endpoint + "/rpc/2.0/unit/service/chat"
	api_unit_mgt   = aip_endpoint + "/rpc/2.0/unit"
	api_unit_file  = aip_endpoint + "/file/2.0/unit"
	aip_grant_type = "client_credentials"
)
