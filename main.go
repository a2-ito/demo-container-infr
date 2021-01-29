package main

import (
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"io/ioutil"
	yaml "gopkg.in/yaml.v2"
	"fmt"
	"os"
  "demo-container-infr/cloudrun"
)

const confDir = "./config/"

type GC_Key struct {
    ProjectId	string	`json:"project_id"`
}

func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {
    err := createGcpResources(ctx)
    if err != nil {
      return err
    }
    return nil
  })

}

func createGcpResources(ctx *pulumi.Context) error {

	// set environment variables
	GOOGLE_PROJECT					:= os.Getenv("GOOGLE_PROJECT")
	if GOOGLE_PROJECT == "" {
		panic("failed to get SECRET_GOOGLE_BILLING_ACCOUNT_ID.")
	}

	appMode := os.Getenv("APP_MODE")
	if appMode == "" {
		  appMode = "dev"
	//	panic("failed to get application mode, check whether APP_MODE is set.")
	}
	fmt.Println(appMode)

  configbuf, err := ioutil.ReadFile(confDir + appMode + ".yaml")
  if err != nil {
    fmt.Println(err)
    return err
  }

  yamldata := make(map[interface{}]interface{})
  err = yaml.Unmarshal(configbuf, &yamldata)
  if err != nil {
    fmt.Println(err)
    return err
  }

	err = cloudrun.CreateCloudRunService(
		ctx,
		GOOGLE_PROJECT,
		nil,
		yamldata["cloudrun"].(map[interface {}]interface {})["serviceapi"].(string),
		yamldata["cloudrun"].(map[interface {}]interface {})["computeserviceapi"].(string),
		yamldata["cloudrun"].(map[interface {}]interface {})["service"].(map[interface {}]interface {})["name"].(string),
		yamldata["cloudrun"].(map[interface {}]interface {})["service"].(map[interface {}]interface {})["image"].(string),
		yamldata["cloudrun"].(map[interface {}]interface {})["urlmap"].(map[interface {}]interface {})["name"].(string),
		yamldata["cloudrun"].(map[interface {}]interface {})["iam"].(map[interface {}]interface {})["name"].(string))
  if err != nil {
    return err
  }

  return nil
}

