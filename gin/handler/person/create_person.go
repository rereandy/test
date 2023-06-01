package person

import (
	"fmt"
	"time"

	"github.com/LearnGin/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CreatePersonHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// parse request
		tic := time.Now()
		req := new(model.CreatePersonRequest)
		err := c.ShouldBindWith(req, binding.JSON)
		if err != nil {
			fmt.Errorf("Can not bind with model.Person, err: %+v\n", err)
			resp := new(model.CreatePersonResponse)
			resp.Elapse = time.Since(tic)
			resp.BaseResp = model.BaseResp{
				Code:    1,
				Message: fmt.Sprintf("create person failed in binding json, err: %s", err.Error()),
			}
			c.JSON(200, resp)
			return
		}

		// process request
		fmt.Printf("Creating Person: %+v\n", req.Person)

		// jsonify response
		resp := new(model.CreatePersonResponse)
		resp.Person = req.Person
		resp.Elapse = time.Since(tic)
		resp.BaseResp = model.BaseResp{
			Code:    0,
			Message: "success",
		}
		c.JSON(200, resp)
	}
}
