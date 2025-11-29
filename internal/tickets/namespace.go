package tickets

func Namespaced(agentID, key string) string {
	if agentID == "" {
		agentID = "_"
	}
	if key == "" {
		key = "_"
	}
	return "agent." + agentID + "." + key
}
