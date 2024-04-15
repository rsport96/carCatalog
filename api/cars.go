package api

import (
	"catalog/api/request"
	"catalog/app"
	"catalog/db"
	"catalog/models"
	"catalog/saturator"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofrs/uuid"
	"regexp"
	"strconv"
)

func getCars(app *app.App) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		filter := models.Car{}

		idString := ctx.Query("id", "")
		filter.ID = uuid.FromStringOrNil(idString)
		if idString != "" && filter.ID == uuid.Nil {
			log.Info("id is not valid")
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		filter.RegNum = ctx.Query("regNum", "")
		filter.Model = ctx.Query("model", "")
		filter.Mark = ctx.Query("mark", "")

		yearString := ctx.Query("year", "0")
		var err error
		filter.Year, err = strconv.Atoi(yearString)
		if err != nil {
			filter.Year = 0
		}

		offsetString := ctx.Query("offset", "0")
		offset, err := strconv.Atoi(offsetString)
		if err != nil || offset < 0 {
			log.Info("offset is not valid")
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		limitString := ctx.Query("limit", "100")
		limit, err := strconv.Atoi(limitString)
		if err != nil || limit < 0 {
			log.Info("limit is not valid")
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		cars, err := db.GetCars(app, filter, offset, limit)

		if err != nil {
			log.Debug(err.Error())
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		log.Info("Cars were successfully received")
		return ctx.Status(fiber.StatusOK).JSON(cars)
	}
}

func deleteCar(app *app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idString := c.Params("id", "")
		id := uuid.FromStringOrNil(idString)

		if id == uuid.Nil {
			log.Info("id is not valid")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err := db.DeleteCarById(app, id)
		if errors.Is(err, db.NotFound) {
			log.Info("car not found")
			return c.SendStatus(fiber.StatusNotFound)
		} else if err != nil {
			log.Debug(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Info("Car was successfully deleted")
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func updateCars(app *app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cars []*models.Car
		err := c.BodyParser(&cars)
		if err != nil {
			log.Info("error while parsing body", err.Error())
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err = db.UpdateCars(app, cars)
		if errors.Is(err, db.NotFound) {
			log.Info("car not found")
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		} else if errors.Is(err, db.AlreadyExists) {
			log.Info("car already exists")
			return fiber.NewError(fiber.StatusConflict, err.Error())
		} else if err != nil {
			log.Debug(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Info("Cars were successfully updated")
		return c.Status(fiber.StatusOK).JSON(cars)
	}
}

func addCars(app *app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var regNums request.AddCarsRequest
		if err := c.BodyParser(&regNums); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		var cars []*models.Car
		for _, regNum := range regNums.RegNums {
			// regexp for russian car numbers
			numPattern := regexp.MustCompile("^[А-Я]\\d{3}[А-Я]{2}\\d{3}$")
			if !numPattern.MatchString(regNum) {
				return c.SendStatus(fiber.StatusUnprocessableEntity)
			}
			car, err := app.Saturator().Saturate(regNum)
			if errors.Is(err, saturator.NilClientError) {
				log.Debug("saturator client is nil", err.Error())
				return c.SendStatus(fiber.StatusInternalServerError)
			} else if errors.Is(err, saturator.StatusBadRequestError) {
				log.Info("bad request to saturator", err.Error())
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			} else if errors.Is(err, saturator.StatusInternalServerErrorError) {
				log.Info("internal server error in saturator", err.Error())
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			} else if err != nil {
				log.Debug(err.Error())
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			cars = append(cars, car)
		}

		err := db.AddCars(app, cars)
		if errors.Is(err, db.AlreadyExists) {
			log.Info("car already exists")
			return fiber.NewError(fiber.StatusConflict, err.Error())
		} else if err != nil {
			log.Debug(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Info("Cars were successfully added")
		return c.Status(fiber.StatusCreated).JSON(cars)
	}
}
