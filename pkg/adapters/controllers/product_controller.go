package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/internal/utils"
	repositories "github.com/shayja/go-template-api/pkg/adapters/repositories/product"
	"github.com/shayja/go-template-api/pkg/entities"
)

type ProductController struct {
	Db *sql.DB
}

func CreateProductController(db *sql.DB) ProductControllerInterface {
	return &ProductController{Db: db}
}

// GetAll implements ProductControllerInterface
func (m *ProductController) GetAll(c *gin.Context) {
	AddRequestHeader(c)
	DB := m.Db
	repositories := repositories.NewProductRepository(DB)
	page, err := strconv.Atoi(c.Query("page"));

	if (err != nil) {
		panic(err)
    }

	res, err := repositories.GetAll(page)
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

	var uri entities.ProductUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}
	repositories := repositories.NewProductRepository(DB)
	
	res, err := repositories.GetSingle(uri.Id)
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
	var post *entities.ProductRequest
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}
	repositories := repositories.NewProductRepository(DB)
	insertedId, err := repositories.Create(post)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	if utils.IsValidUUID(insertedId) {
		c.JSON(http.StatusCreated, gin.H{"status": "success", "msg": nil, "id": insertedId})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "msg": "insert product failed"})
	}
}

// Update implements ProductControllerInterface
func (m *ProductController) Update(c *gin.Context) {
	AddRequestHeader(c)
	DB := m.Db
	var product *entities.ProductRequest
	
	if err := c.ShouldBind(&product); err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}
   
	var uri entities.ProductUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}
	repositories := repositories.NewProductRepository(DB)
	res, err := repositories.Update(uri.Id, product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	if utils.IsValidUUID(res.Id) {
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": res, "msg": nil})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "data": nil, "msg": "update product failed"})
	}

}

//change a specific product price
func (m *ProductController) UpdatePrice(c *gin.Context){
	AddRequestHeader(c)
    var uri entities.ProductUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	var price *entities.ProductPriceRequest
	if err := c.ShouldBind(&price); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	if price.Price <= 0 {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid price value."}) //return custom request for bad request or book not found
        return
    }

	DB := m.Db
	repositories := repositories.NewProductRepository(DB)
	product, err := repositories.GetSingle(uri.Id)

    if err != nil || !utils.IsValidUUID(product.Id) {
		 //return custom request for bad request or item not found
        c.JSON(http.StatusNotFound, gin.H{"message": "Product not found."})
        return
    }

    res, err := repositories.UpdatePrice(uri.Id, price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	if utils.IsValidUUID(res.Id)  {
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": res, "msg": nil})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "data": nil, "msg": "update product failed"})
	}
}

func (m *ProductController) UpdateImage(c *gin.Context){
	AddRequestHeader(c)
    var uri entities.ProductUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	file, err := c.FormFile("image")
   
	if err != nil {
	  log.Println("Error in uploading Image : ", err)
	  c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "data": nil, "msg": "update product failed"})
	}
	
	 uniqueId := utils.CreateNewUUID()

	 filename := strings.Replace(uniqueId.String(), "-", "", -1)
	
	 fileExt := strings.Split(file.Filename, ".")[1]
	
	 imageName := fmt.Sprintf("%s.%s", filename, fileExt)
	
	//todo: Add file size/dimentions and file ext. validation here.
	// Just for the demo, Do not use this for any real world solution. 
	// Do not store uploaded files on your web server, use AWS/GCP cloud bucket instead.
	 err = c.SaveUploadedFile(file, fmt.Sprintf("./images/%s", imageName))
	
	 if err != nil {
	  log.Println("Error in saving Image :", err)
	  c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "msg": "Error in saving Image"})
	 }
	
	 
	 imageUrl := fmt.Sprintf("http://YOUR_DOMAIN_HERE.com/images/%s", imageName)
	
	 data := map[string]interface{}{
	
	  "imageName": imageName,
	  "imageUrl":  imageUrl,
	  "header":    file.Header,
	  "size":      file.Size,
	 }
	
	 var image *entities.ProductImageRequest
	 if err := c.ShouldBind(&image); err != nil {
		 c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		 return
	 }
 
	 DB := m.Db
	 repositories := repositories.NewProductRepository(DB)
	 res, err := repositories.UpdateImage(uri.Id, image)
	 if err != nil {
		 c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		 return
	 }
	 log.Println("UpdateImage res :", res)


	 c.JSON(http.StatusCreated, gin.H{"status": "success", "data": "Image uploaded successfully", "msg": data})
   }

// Delete implements ProductController Interface
func (m *ProductController) Delete(c *gin.Context) {
	AddRequestHeader(c)
	DB := m.Db
	var uri entities.ProductUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}
	repositories := repositories.NewProductRepository(DB)
	res := repositories.Delete(uri.Id)
	if res {
		c.JSON(http.StatusOK, gin.H{"status": "success", "msg": nil})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "msg": "delete product failed"})
		return
	}
}