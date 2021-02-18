/* ****************************************************************************
Copyright © 2020 by Sagacity. All rights reserved.
Filename    : helper.go
File-type   : golang source code.
Compiler    : go version go1.13.1 linux/amd64

Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Date        : 24-Jan-2020
Description :
- Uses as helper functions.

Version History
Version     : 1.0
Author      : om kumar (om.kumar@sagacitysoftware.co.in)
Description : Initial version
**************************************************************************** */

package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	gomail "gopkg.in/gomail.v2"
	gModels "nfon.com/models"

	"os"
	"path/filepath"
	"reflect"
	"unicode"

	"github.com/go-sql-driver/mysql"
	"nfon.com/logger"
)

func ReadFileContent(filePath string) (bool, []byte) {
	rootDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileFullPath := filepath.Join(rootDir, filePath)
	logger.Log(MODULENAME, logger.DEBUG, "fileFullPath: %s", fileFullPath)
	byteData, readError := ioutil.ReadFile(fileFullPath)
	if readError != nil {
		logger.Log(MODULENAME, logger.ERROR, "byteData read error from file: %s, error: %s", fileFullPath, readError.Error())
		return false, []byte{}
	}

	return true, byteData
}

func ConvertToJSON(dataStruct interface{}) (bool, string) {
	dataBytes, err := json.Marshal(dataStruct)

	if err != nil {
		return false, ""
	}
	jsonData := string(dataBytes)

	return true, jsonData
}

func ConvertFromJSON(jsonData string, pConvertType interface{}) bool {
	err := json.Unmarshal([]byte(jsonData), pConvertType)

	if err != nil {
		fmt.Println("failed to convert from json", err.Error())
		return false
	}

	return true
}

func GetDBTagFromPropName(user interface{}, propName string) string {

	dbTagName := "db"

	t := reflect.TypeOf(user)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Name == propName {
			dbTagValue := field.Tag.Get(dbTagName)
			return dbTagValue
		}
	}

	return ""
}

func GetDBTagWithDataTypeFromJSONTag(user interface{}, jsonTag string) (string, string) {

	dbTagName := "db"
	jsonTagName := "json"

	t := reflect.TypeOf(user)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag := field.Tag.Get(jsonTagName)

		if tag == jsonTag {
			dbTagValue := field.Tag.Get(dbTagName)
			return dbTagValue, fmt.Sprintf("%s", field.Type)
		}
	}
	return "", ""
}

func VerifyPasswordStrength(userPassword string) (sevenOrMore, number, upper, lower, special bool) {
	letters := 0
	for _, password := range userPassword {
		switch {
		case unicode.IsNumber(password):
			number = true
			letters++
		case unicode.IsUpper(password):
			upper = true
			letters++
		case unicode.IsPunct(password) || unicode.IsSymbol(password):
			special = true
			letters++
		case unicode.IsLetter(password) || password == ' ':
			lower = true
			letters++
		default:
			return false, false, false, false, false
		}
	}
	if letters >= 8 {
		sevenOrMore = true
	}
	return sevenOrMore, number, upper, lower, special
}

func MaskMobile(s string) string {
	rs := []rune(s)
	for i := 0; i < len(rs)-4; i++ {
		rs[i] = 'X'
	}
	return string(rs)
}

func GetFileContentType(out *os.File) (string, error) {

	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func GetApplicationErrorCodeFromDBError(dbErr error) (errorHandled bool, errorCode int) {

	if err, ok := dbErr.(*mysql.MySQLError); ok {

		switch err.Number {
		case 1062: //Unique key constrain failed
			return true, MOD_OPER_DUPLICATE_RECORD_FOUND
		default:
			return false, int(err.Number)
		}
	}
	return false, MOD_OPER_ERR_DATABASE
}

func SendEmailNotification(smtpDetails gModels.SmtpConfigDetails, emailData gModels.DBEmailConfigRowModel, toEmailAdd string) {
	//fmt.Println("Function invoked to send Email...")
	m := gomail.NewMessage()
	m.SetAddressHeader("From", smtpDetails.SmtpUserName, EMAIL_DISPLAY_NAME)

	m.SetAddressHeader("To", toEmailAdd, toEmailAdd)
	m.SetHeader("Subject", emailData.Subject)
	m.SetBody("text/html", emailData.Body)
	d := gomail.NewPlainDialer(smtpDetails.SmtpHostName, smtpDetails.SmtpHostPort, smtpDetails.SmtpUserName, smtpDetails.SmtpHostPassword)

	if err := d.DialAndSend(m); err != nil {
		logger.Log(MODULENAME, logger.ERROR, "Error in sending Email...%s", err)
	}
}
