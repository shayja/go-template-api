// internal/adapters/controllers/product_controller.go
package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/internal/entities"
	"github.com/shayja/go-template-api/internal/usecases"
	"github.com/shayja/go-template-api/internal/utils"
)

type ProductController struct {
	ProductInteractor usecases.ProductInteractor
}

// GetAll godoc
// @Summary      Get all products
// @Description  Retrieve a paginated list of all products
// @Tags         Products
// @Param        page  query     int  true  "Page number"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Router       /product [get]
func (uc *ProductController) GetAll(c *gin.Context) {
	AddRequestHeader(c)
	page, err := strconv.Atoi(c.Query("page"))

	if err != nil {
		panic(err)
	}

	res, err := uc.ProductInteractor.GetAll(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	if res != nil {
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": res, "msg": nil})
		return
	} else {
		c.JSON(404, gin.H{"status": "failed", "data": nil, "msg": "products not found for this page."})
		return
	}
}

// GetById godoc
// @Summary      Get a product by ID
// @Description  Retrieve product details by product ID
// @Tags         Products
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /product/{id} [get]
func (uc *ProductController) GetById(c *gin.Context) {
	AddRequestHeader(c)

	var uri entities.ProductUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	res, err := uc.ProductInteractor.GetById(uri.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	if utils.IsValidUUID(res.Id) {
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": res, "msg": nil})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"status": "success", "data": nil, "msg": "product not found"})
	}
}

// Create godoc
// @Summary      Create a new product
// @Description  Add a new product to the inventory
// @Tags         Products
// @Param        product  body      entities.ProductRequest  true  "Product data"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /product [post]
func (uc *ProductController) Create(c *gin.Context) {
	AddRequestHeader(c)

	var post *entities.ProductRequest
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	insertedId, err := uc.ProductInteractor.Create(post)

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

// Update godoc
// @Summary      Update product details
// @Description  Update an existing product's details by ID
// @Tags         Products
// @Param        id       path      string                  true  "Product ID"
// @Param        product  body      entities.ProductRequest true  "Updated product data"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /product/{id} [put]
func (uc *ProductController) Update(c *gin.Context) {
	AddRequestHeader(c)

	var product *entities.ProductRequest

	if err := c.ShouldBind(&product); err != nil {
		fmt.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	var uri entities.ProductUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	res, err := uc.ProductInteractor.Update(uri.Id, product)
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
// @Summary Update Product Price
// @Description Update the price of a specific product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param price body entities.ProductPriceRequest true "Product Price"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /products/{id}/price [put]
func (uc *ProductController) UpdatePrice(c *gin.Context){
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

	product, err := uc.ProductInteractor.GetById(uri.Id)

    if err != nil || !utils.IsValidUUID(product.Id) {
		 //return custom request for bad request or item not found
        c.JSON(http.StatusNotFound, gin.H{"message": "Product not found."})
        return
    }

    res, err := uc.ProductInteractor.UpdatePrice(uri.Id, price)
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

// @Summary Update Product Image
// @Description Upload and update the image of a specific product by ID
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Product ID"
// @Param image formData file true "Product Image"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{id}/image [put]
func (uc *ProductController) UpdateImage(c *gin.Context){
	AddRequestHeader(c)
    var uri entities.ProductUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	file, err := c.FormFile("image")
   
	if err != nil {
		fmt.Print("Error in uploading Image : ", err)
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
		fmt.Print("Error in saving Image :", err)
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
 
	 
	 
	 res, err := uc.ProductInteractor.UpdateImage(uri.Id, image)
	 if err != nil {
		 c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		 return
	 }
	fmt.Print("UpdateImage res :", res)


	 c.JSON(http.StatusCreated, gin.H{"status": "success", "data": "Image uploaded successfully", "msg": data})
   }

// Delete implements ProductController Interface
// @Summary Delete Product
// @Description Delete a specific product by ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{id} [delete]
func (uc *ProductController) Delete(c *gin.Context) {
	AddRequestHeader(c)
	
	var uri entities.ProductUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}
	
	res, err := uc.ProductInteractor.Delete(uri.Id)

	if err != nil {
		fmt.Print("Error in Delete product: ", err)
	  	c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "data": nil, "msg": "Delete product failed"})
	}
	
	if res {
		c.JSON(http.StatusOK, gin.H{"status": "success", "msg": nil})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "msg": "delete product failed"})
		return
	}
}
