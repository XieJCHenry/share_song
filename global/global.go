package global

var objectMap map[string]UsableObject

func Init() {
	objectMap = make(map[string]UsableObject)
}

func GetGlobalObject(key string) UsableObject {
	if obj, ok := objectMap[key]; ok {
		if obj != nil && obj.Usable() {
			return obj
		}
	}
	return nil
}

func SetGlobalObject(obj UsableObject) {
	key := obj.Key()
	if _, ok := objectMap[key]; !ok && obj != nil && obj.Usable() {
		objectMap[key] = obj
	}
}
