package api

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Application Configuraton
type configuration struct {
	Host       string
	DBName     string
	DBUsername string
	DBPassword string

	BCryptCost int
	JWTkey     string
}

var Config = configuration{}

var defaultBcryptCost = 3

func LoadEnvVars() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Error("Couldnt load .env: ", err)
		return err
	}
	var present bool

	Config.DBName, present = os.LookupEnv("dbname")
	if !present {
		log.Error("dbname not found in .env file")
		return errors.New("read .env:unsuccessfull")
	}
	Config.DBUsername, present = os.LookupEnv("dbusername")
	if !present {
		log.Error("dbusername not found in .env file")
		return errors.New("read .env:unsuccessfull")
	}
	Config.DBPassword, present = os.LookupEnv("dbpassword")
	if !present {
		log.Error("dbpassword not found in .env file")
		return errors.New("read .env:unsuccessfull")
	}
	Config.Host, present = os.LookupEnv("host")
	if !present {
		log.Error("host not found in .env file")
		return errors.New("read .env:unsuccessfull")
	}

	Config.JWTkey, present = os.LookupEnv("jwt_key")
	if !present {
		log.Error("jwtkey not found in .env file")
		return errors.New("read .env:unsuccessfull")
	}

	bcoststr, present := os.LookupEnv("bcrypt_cost")
	if !present {
		//default value 3
		Config.BCryptCost = defaultBcryptCost
	} else {

		bcost, err := strconv.Atoi(bcoststr)
		if err != nil {
			log.Error("Unable to convert bcryptcost(string) to int")
			return err
		}
		Config.BCryptCost = bcost
	}

	return nil
}
