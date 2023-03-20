package protocol

import "fmt"

type NotFoundErr struct {
	msg string
}

func (ne *NotFoundErr) Error() string {
	return ne.msg
}

type Protocol struct {
	Path string                 `json:"path,omitempty"`
	Body map[string]interface{} `json:"body,omitempty"`
}

func (p *Protocol) GetField(key string) (interface{}, error) {
	if val, ok := p.Body[key]; ok {
		return val, nil
	} else {
		return nil, &NotFoundErr{
			msg: fmt.Sprintf("key %s not found", key),
		}
	}
}

func (p *Protocol) GetFieldOrDefault(key string, def interface{}) interface{} {
	if val, ok := p.Body[key]; ok {
		return val
	} else {
		return def
	}
}

func (p *Protocol) ValidAndGetFields(fields ...string) (map[string]interface{}, error) {
	var err error
	result := make(map[string]interface{})
	for _, key := range fields {
		if val, ok := p.Body[key]; ok {
			result[key] = val
		} else {
			err = &NotFoundErr{
				msg: fmt.Sprintf("key %s not found", key),
			}
			break
		}
	}

	return result, err
}

func WrapperError(err error) (map[string]interface{}, error) {
	return map[string]interface{}{
		"error": err.Error(),
	}, nil
}
