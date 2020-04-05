package exconn

import "github.com/humamfauzi/go-registration/utils"

func ComposeConnectionFromEnv(connection interface{}, vendor string) string {
	switch connection.(type) {
	case map[string]interface{}:
		switch vendor {
		case "mysql":
			parsed := connection.(map[string]interface{})
			composed := utils.InterpretInterfaceString(parsed["username"], "root") + ":"
			composed += utils.InterpretInterfaceString(parsed["password"], "") + "@"
			composed += utils.InterpretInterfaceString(parsed["protocol"], "tcp") + "("
			composed += utils.InterpretInterfaceString(parsed["adress"], "localhost") + ")/"
			composed += utils.InterpretInterfaceString(parsed["dbname"], "try1")
			composed += GetAdditionalDbConnectionParams(parsed)
			return composed
		case "postgres":
			parsed := connection.(map[string]interface{})
			composed := "host=" + utils.InterpretInterfaceString(parsed["address"], "localhost") + " "
			composed += "port=" + utils.InterpretInterfaceString(parsed["port"], "5432") + " "
			composed += "dbname=" + utils.InterpretInterfaceString(parsed["dbname"], "try1") + " "
			composed += "user=" + utils.InterpretInterfaceString(parsed["username"], "root") + " "
			composed += "password=" + utils.InterpretInterfaceString(parsed["password"], "")
			return composed
		case "cassandra":
			//  Because cassandra could work in cluster then the return
			//  value is still array but will be parsed later to array of string
			parsed := connection.(map[string]interface{})
			addressArray := parsed["address"].([]interface{})
			stringArray := ""
			for _, v := range addressArray {
				stringArray += v.(string) + ","
			}
			composed := stringArray[0 : len(stringArray)-1]
			return composed
		case "mongo":
			composed := ""
			return composed
		default:
			panic("FAILED TO PARSE EXTERNAL CONNECTION PROFILE")
		}
	default:
		panic("FAILED TO PARSE EXTERNAL CONNECTION PROFILE")
	}
}

func GetAdditionalDbConnectionParams(connectionParams map[string]interface{}) string {
	reservedKey := utils.StringArray([]string{"username", "password", "protocol", "address", "dbname"})
	var connParams string
	for key, value := range connectionParams {
		if !reservedKey.Includes(key) {
			if key == "parsetime" {
				key = "parseTime"
			}
			connParams += key + "=" + value.(string) + "&"
		}
	}
	if connParams != "" {
		connParams = "?" + connParams
		connParams = connParams[0 : len(connParams)-1]
	}
	return connParams
}
