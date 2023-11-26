package engine

import (
	"regexp"
	"strings"
)

func processData(input map[string]interface{}, config []map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Iterar sobre as configurações
	for _, cfg := range config {
		keyRef := cfg["key_ref"].(string)
		keyResult := cfg["key_result"].(string)

		// Separar as chaves aninhadas usando "."
		keys := strings.Split(keyRef, ".")

		// Iterar sobre as chaves para obter o valor final
		value, found := navigateKeys(input, keys)
		if !found {
			return nil
		}

		// Aplicar filtros e mapeamentos
		value = applyFiltersAndMapping(value, cfg)

		// Resolver valor, levando em consideração expressões "$"
		value = resolveValue(value, cfg)

		// Inserir no resultadoß
		result[keyResult] = value
	}

	return result
}

func navigateKeys(input map[string]interface{}, keys []string) (interface{}, bool) {
	for i, key := range keys {
		// Se estiver na última iteração, retornar o valor
		if i == len(keys)-1 {
			return input[key], true
		}
		if nested, ok := input[key].(map[string]interface{}); ok {
			// Se estiver na última iteração, retornar o valor
			if i == len(keys)-1 {
				return nested, true
			}
			// Se não estiver na última iteração, continuar a navegar
			input = nested
		} else {
			// Se a chave não for encontrada, retornar nil
			return nil, false
		}
	}

	return nil, false
}

func resolveValue(value interface{}, cfg map[string]interface{}) interface{} {
	// Verificar se a expressão "$" está presente em "value"
	if valueExpr, ok := cfg["value"].(string); ok && valueExpr != "" {

		// Se não houver expressão "$", retornar o valor original
		if valueExpr[0] != '$' {
			return valueExpr
		}

		// Extrair a referência da expressão "$"
		refKey := valueExpr[1:]

		// Dividir as referências aninhadas usando "."
		keys := strings.Split(refKey, ".")

		// Iterar sobre as chaves para obter o valor final
		for _, key := range keys {
			if nested, ok := value.(map[string]interface{}); ok {
				value = nested[key]
			} else {
				// Se não for um mapa, retornar nil
				return nil
			}
		}

		return value
	}

	// Se não houver expressão "$", retornar o valor original
	return value
}

func applyMappingRecursively(value interface{}, cfg map[string]interface{}) interface{} {
	if mapping, ok := cfg["mapping"].([]interface{}); ok {
		switch typedValue := value.(type) {
		case map[string]interface{}:
			result := make(map[string]interface{})
			for _, mapCfg := range mapping {
				mapCfgMap := mapCfg.(map[string]interface{})
				mapKeyRef := mapCfgMap["key_ref"].(string)
				mapKeyResult := mapCfgMap["key_result"].(string)
				mapValue := applyMappingRecursively(typedValue[mapKeyRef], mapCfgMap)
				result[mapKeyResult] = mapValue
			}
			return result
		case []interface{}:
			// Se value for um array, aplicar recursivamente para cada elemento
			var resultSlice []interface{}
			for _, item := range typedValue {
				mapValue := applyMappingRecursively(item, cfg)
				resultSlice = append(resultSlice, mapValue)
			}
			return resultSlice
		default:
			// Se value não for nem mapa nem array, retornar o valor original
			return value
		}
	}

	// Se não houver "mapping", retornar o valor original
	return value
}

func applyFiltersAndMapping(value interface{}, cfg map[string]interface{}) interface{} {
	// Aplicar filtros
	if filters, ok := cfg["filters"].([]interface{}); ok {
		for _, filter := range filters {
			filterMap := filter.(map[string]interface{})
			operator := filterMap["operator"].(string)

			switch operator {
			case "not_nil":
				if value == nil {
					return nil
				}
				break
			case "not_empty":
				if value == "" {
					return nil
				}
				break
			case "regex":
				pattern := filterMap["pattern"].(string)
				r, _ := regexp.Compile(pattern)
				if !r.MatchString(value.(string)) {
					return nil
				}
				break
			}
		}
	}

	// Mapear valor, se existir
	if mapping, ok := cfg["mapping"].([]interface{}); ok {
		switch typedValue := value.(type) {
		case map[string]interface{}:
			result := make(map[string]interface{})
			for _, mapCfg := range mapping {
				mapCfgMap := mapCfg.(map[string]interface{})
				mapKeyRef := mapCfgMap["key_ref"].(string)
				mapKeyResult := mapCfgMap["key_result"].(string)
				mapValue := applyFiltersAndMapping(typedValue[mapKeyRef], mapCfgMap)
				result[mapKeyResult] = mapValue
			}
			return result
		case []interface{}:
			// Se value for um array, aplicar recursivamente para cada elemento
			var resultSlice []interface{}
			for _, item := range typedValue {
				for _, mapCfg := range mapping {
					mapCfgMap := mapCfg.(map[string]interface{})
					mapValue := applyFiltersAndMapping(item, mapCfgMap)
					resultSlice = append(resultSlice, mapValue)
				}
			}
			return resultSlice
		default:
			// Se value não for nem mapa nem array, retornar o valor original
			return value
		}
	}

	// Retornar valor final
	return value
}
