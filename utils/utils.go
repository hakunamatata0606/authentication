package utils

import "authentication/service/token"

func GetRolesFromClaims(claims *token.ClaimMap) []string {
	roles, ok := (*claims)["role"].([]interface{})
	if !ok {
		return []string{}
	}
	var result []string
	for _, role := range roles {
		roleStr, ok := role.(string)
		if !ok {
			continue
		}
		result = append(result, roleStr)
	}
	return result
}
