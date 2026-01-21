package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collections := []string{"system_stats", "container_stats"}

		for _, collectionName := range collections {
			collection, err := app.FindCollectionByNameOrId(collectionName)
			if err != nil {
				return err
			}

			// Find the "type" field and add "1d" to its values
			typeField := collection.Fields.GetByName("type")
			if typeField == nil {
				continue
			}

			if selectField, ok := typeField.(*core.SelectField); ok {
				// Check if "1d" already exists
				for _, v := range selectField.Values {
					if v == "1d" {
						goto skipCollection
					}
				}
				selectField.Values = append(selectField.Values, "1d")
			}

			if err := app.Save(collection); err != nil {
				return err
			}

		skipCollection:
		}

		return nil
	}, func(app core.App) error {
		// Rollback: remove "1d" from values
		collections := []string{"system_stats", "container_stats"}

		for _, collectionName := range collections {
			collection, err := app.FindCollectionByNameOrId(collectionName)
			if err != nil {
				return err
			}

			typeField := collection.Fields.GetByName("type")
			if typeField == nil {
				continue
			}

			if selectField, ok := typeField.(*core.SelectField); ok {
				newValues := make([]string, 0, len(selectField.Values))
				for _, v := range selectField.Values {
					if v != "1d" {
						newValues = append(newValues, v)
					}
				}
				selectField.Values = newValues
			}

			if err := app.Save(collection); err != nil {
				return err
			}
		}

		return nil
	})
}
