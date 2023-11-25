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
		var value interface{}
		for i, key := range keys {
			if nested, ok := input[key]; ok {
				if i == len(keys)-1 {
					value = nested
				} else {
					input = nested.(map[string]interface{})
				}
			} else {
				// Se a chave não for encontrada, retornar nil
				return nil
			}
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

func resolveValue(value interface{}, cfg map[string]interface{}) interface{} {
	// Verificar se a expressão "$" está presente em "value"
	if valueExpr, ok := cfg["value"].(string); ok && valueExpr != "" {
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
		if value != nil {
			// Processar mapeamento recursivamente
			result := make(map[string]interface{})
			for _, mapCfg := range mapping {
				mapCfgMap := mapCfg.(map[string]interface{})
				mapKeyRef := mapCfgMap["key_ref"].(string)
				mapKeyResult := mapCfgMap["key_result"].(string)
				mapValue := applyFiltersAndMapping(value.(map[string]interface{})[mapKeyRef], mapCfgMap)
				result[mapKeyResult] = mapValue
			}
			return result
		}
	}

	// Retornar valor final
	return value
}
