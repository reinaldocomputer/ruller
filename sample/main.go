package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/flaviostutz/ruller"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Infof("====Starting Ruller Sample====")
	err := ruller.Add("test", "rule1", func(ctx ruller.Context) (map[string]interface{}, error) {
		output := make(map[string]interface{})
		output["opt1"] = "Some tests rule 1"
		output["rule1-opt2"] = 129.99
		rnd := fmt.Sprintf("v%d", rand.Int())
		if ctx.Input["menu"] == true {
			child := make(map[string]interface{})
			child["rule1-c1"] = "123"
			child["rule1-c2"] = rnd
			output["menu"] = child
		}
		output["rule1"] = true
		return output, nil
	})
	if err != nil {
		panic(err)
	}

	ruller.AddRequiredInput("test", "age", ruller.Float64)
	ruller.AddRequiredInput("test", "children", ruller.Bool)

	err = ruller.AddChild("test", "rule1.1", "rule1", func(ctx ruller.Context) (map[string]interface{}, error) {
		output := make(map[string]interface{})
		output["rule1.1-mydata"] = "myvalue"
		output["rule1.1"] = true
		return output, nil
	})
	if err != nil {
		panic(err)
	}

	err = ruller.Add("test", "rule2", func(ctx ruller.Context) (map[string]interface{}, error) {
		output := make(map[string]interface{})
		output["opt1"] = "Lots of tests rule 2"
		output["rule2"] = true
		// logrus.Debugf("children output from rule 2.1 is %s", ctx.ChildrenOutput)
		return output, nil
	})
	if err != nil {
		panic(err)
	}

	err = ruller.AddChild("test", "rule2.1", "rule2", func(ctx ruller.Context) (map[string]interface{}, error) {
		output := make(map[string]interface{})
		age, ok := ctx.Input["age"].(float64)
		if !ok {
			return nil, fmt.Errorf("Invalid 'age' detected. age=%s", ctx.Input["age"])
		}
		if age > 60 {
			output["category"] = "elder rule2.1"
		} else {
			output["category"] = "young rule2.1"
		}
		output["geo-city"] = ctx.Input["_ip_city"]
		output["geo-state"] = ctx.Input["_ip_state"]
		output["rule2.1"] = true
		return output, nil
	})
	if err != nil {
		panic(err)
	}

	ruller.SetRequestFilter(func(r *http.Request, input map[string]interface{}) error {
		logrus.Debugf("filtering request. input=%s", input)
		input["_something"] = "test"
		return nil
	})

	ruller.SetResponseFilter(func(w http.ResponseWriter, input map[string]interface{}, output map[string]interface{}, outBytes []byte) (bool, error) {
		logrus.Debugf("filtering response. input=%s", input)
		output["filter-attribute"] = "added by sample filter"
		if input["_something"] == "test" {
			w.Write([]byte("{\"a\":"))
			w.Write(outBytes)
			w.Write([]byte("}"))
		}
		return true, nil
	})

	err = ruller.AddChild("test", "rule2.2", "rule2", func(ctx ruller.Context) (map[string]interface{}, error) {
		output := make(map[string]interface{})
		output["rule2.2-type"] = "any"
		output["rule2.2"] = true
		return output, nil
	})
	if err != nil {
		panic(err)
	}

	for a := 0; a < 10; a++ {
		err = ruller.AddChild("test", fmt.Sprintf("rule2.2-%d", a), "rule2.2", func(ctx ruller.Context) (map[string]interface{}, error) {
			output := make(map[string]interface{})
			output["opt1"] = "any1"
			return output, nil
		})
		if err != nil {
			panic(err)
		}
		for b := 0; b < 1; b++ {
			err = ruller.AddChild("test", fmt.Sprintf("rule2.2-%d-%d", a, b), fmt.Sprintf("rule2.2-%d", a), func(ctx ruller.Context) (map[string]interface{}, error) {
				output := make(map[string]interface{})
				output["opt2"] = "any2"
				return output, nil
			})
			if err != nil {
				panic(err)
			}
			for c := 0; c < 1; c++ {
				err = ruller.AddChild("test", fmt.Sprintf("rule2.2-%d-%d-%d", a, b, c), fmt.Sprintf("rule2.2-%d-%d", a, b), func(ctx ruller.Context) (map[string]interface{}, error) {
					output := make(map[string]interface{})
					output["opt3"] = "any3"
					return output, nil
				})
				if err != nil {
					panic(err)
				}
			}
		}
	}

	h := ruller.NewHTTPServer()
	err = h.StartServer()
	if err != nil {
		logrus.Errorf("Error starting server. err=%s", err)
	}

}
