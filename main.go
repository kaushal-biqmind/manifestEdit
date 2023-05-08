package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

func main() {

	data := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Service",
		"metadata": map[string]interface{}{
			"kind": "Service",
			"name": "my-service",
		},
		"spec": map[string]interface{}{
			"selector": map[string]interface{}{
				"app": "my-app",
			},
			"ports": []interface{}{
				map[string]interface{}{
					"name":       "http",
					"port":       80,
					"targetPort": 8080,
				},
			},
		},
	}

	manifestList, err := manifest.New(data)
	if err != nil {
		fmt.Println("error encoding manifest: ", err)
	}
	mfList := manifest.List{}
	mfList = append(mfList, manifestList)
	mfList2, err := ReplaceManifestKey(mfList, "name", "parav")
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("after replacing")
		fmt.Println(mfList2.String())
	}

}
func ReplaceManifestKey(manifests manifest.List, key string, newValue string) (manifest.List, error) {
	// Iterate over the manifests
	var manifestList manifest.List
	for _, mf := range manifests {

		// Convert manifest to json
		data, err := json.Marshal(mf)
		if err != nil {
			fmt.Println("cannot marshal manifest")
		}
		var jsonMap map[string]interface{}
		if err := json.Unmarshal(data, &jsonMap); err != nil {
			fmt.Println("109", err)
		}

		replaceKey(jsonMap, key, newValue)
		updatedManifest, err := manifest.New(jsonMap)
		if err != nil {
			fmt.Println("error encoding manifest: ", err)
		}
		manifestList = append(manifestList, updatedManifest)
	}
	return manifestList, nil
}

// replaceKey recursively searches the JSON object for the specified key.
// If the key is found at a particular level, it updates the value.
// If the value at a particular level is a map or an array, the function
// recursively searches for the key in the child objects or elements.
func replaceKey(jsonObj map[string]interface{}, key string, newValue interface{}) {
	for k, v := range jsonObj {
		if k == key {
			jsonObj[k] = newValue
			return
		}
		switch v := v.(type) {
		case map[string]interface{}:
			replaceKey(v, key, newValue)
		case []interface{}:
			for _, item := range v {
				if childObj, ok := item.(map[string]interface{}); ok {
					replaceKey(childObj, key, newValue)
				}
			}
		}
	}
}
