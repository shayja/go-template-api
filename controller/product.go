package controller

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/model"
	"github.com/shayja/go-template-api/repository"
	"github.com/shayja/go-template-api/repository/utils"
)

type ProductController struct {
	Db *sql.DB
}

func CreateProductController(db *sql.DB) ProductControllerInterface {
	return &ProductController{Db: db}
}

func AddRequestHeader(c *gin.Context) {
	c.Header("Content-Type", "application/json")
}

// GetAll implements ProductControllerInterface
func (m *ProductController) GetAll(c *gin.Context) {
	AddRequestHeader(c)
	DB := m.Db
	repository := repository.NewProductRepository(DB)
	page, err := strconv.Atoi(c.Query("page"));

	if (err != nil) {
		panic(err)
    }

	res, err := repository.GetAll(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}
	
	
	if (res != nil) {
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": res, "msg": nil})
		return
	} else { 
		c.JSON(404, gin.H{"status": "failed", "data": nil, "msg": "products not found for this page."})
		return
	}
}

// GetSingle implements ProductControllerInterface
func (m *ProductController) GetSingle(c *gin.Context) {
	AddRequestHeader(c)
	DB := m.Db

	var uri model.ProductUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}
	repository := repository.NewProductRepository(DB)
	
	res, err := repository.GetSingle(uri.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}
	
	if (utils.IsValidUUID(res.Id)) {
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": res, "msg": nil})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"status": "success", "data": nil, "msg": "product not found"})
	}
}

// Create implements ProductControllerInterface
func (m *ProductController) Create(c *gin.Context) {
	AddRequestHeader(c)
	DB := m.Db
	var post model.ValidateProduct
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}
	repository := repository.NewProductRepository(DB)
	insertedId, err := repository.Create(post)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	if (utils.IsValidUUID(insertedId)) {
		c.JSON(http.StatusCreated, gin.H{"status": "success", "msg": nil})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "msg": "insert product failed"})
	}
}

// Update implements ProductControllerInterface
func (m *ProductController) Update(c *gin.Context) {
	AddRequestHeader(c)
	DB := m.Db
	var product model.ValidateProduct
	

	if err := c.ShouldBind(&product); err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}
   
	var uri model.ProductUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}
	repository := repository.NewProductRepository(DB)
	res, err := repository.Update(uri.ID, product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	if (utils.IsValidUUID(res.Id))  {
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": res, "msg": nil})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "data": nil, "msg": "update product failed"})
	}

}


//change a specific product price
func (m *ProductController) UpdatePrice(c *gin.Context){
	AddRequestHeader(c)
    var uri model.ProductUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	var price model.ValidateProductPrice
	if err := c.ShouldBind(&price); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	if price.Price <= 0{
        c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid price value."}) //return custom request for bad request or book not found
        return
    }

	DB := m.Db
	repository := repository.NewProductRepository(DB)
	product, err := repository.GetSingle(uri.ID)

    if err != nil || (!utils.IsValidUUID(product.Id))  {
		 //return custom request for bad request or item not found
        c.JSON(http.StatusNotFound, gin.H{"message": "Product not found."})
        return
    }

    res, err := repository.UpdatePrice(uri.ID, price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	if (utils.IsValidUUID(res.Id))  {
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": res, "msg": nil})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "data": nil, "msg": "update product failed"})
	}
}



// Delete implements ProductControllerInterface
func (m *ProductController) Delete(c *gin.Context) {
	AddRequestHeader(c)
	DB := m.Db
	var uri model.ProductUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}
	repository := repository.NewProductRepository(DB)
	res := repository.Delete(uri.ID)
	if res {
		c.JSON(http.StatusOK, gin.H{"status": "success", "msg": nil})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "msg": "delete product failed"})
		return
	}
}