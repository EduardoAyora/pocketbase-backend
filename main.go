package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func main() {
	app := pocketbase.New()

	app.OnRecordBeforeCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		if e.Record.Collection().Name == "material_movements" {
			startingPointId := e.Record.GetString("starting_point")
			truckId := e.Record.GetString("truck")
			startingPoint, err := app.Dao().FindRecordById("places", startingPointId)
			if err != nil {
				return err
			}
			truck, err := app.Dao().FindRecordById("trucks", truckId)
			if err != nil {
				return err
			}
			if e.Record.GetString("movement_type") == "income" {
				e.Record.Set("arrival", "Cantera")
				e.Record.Set("starting_point", startingPoint.GetString("name"))
				pricePerMeter := startingPoint.GetFloat("price_per_meter")
				amount := e.Record.GetInt("number_of_travels") * truck.GetInt("capacity")
				movementCost := pricePerMeter * float64(amount)

				collection, err := app.Dao().FindCollectionByNameOrId("material_movement_details")
				if err != nil {
					return err
				}
				record := models.NewRecord(collection)
				record.Set("amount", amount)
				record.Set("unit_price", pricePerMeter)
				record.Set("description", "Ingreso de material")
				record.Set("movement_cost", movementCost)

				if err := app.Dao().SaveRecord(record); err != nil {
					return err
				}

				material_movement_details := []string{record.GetId()}
				e.Record.Set("material_movement_details", material_movement_details)
				e.Record.Set("total", movementCost)
			} else {
				return apis.NewBadRequestError("El tipo de movimiento no es válido", nil)
			}
		}

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
