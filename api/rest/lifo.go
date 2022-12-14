package rest

import (
	"WebLiFo/dto"
	"WebLiFo/factory"
	"WebLiFo/managers"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/wissance/stringFormatter"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (webApiContext *WebApiContext) GetAllLifo(respWriter http.ResponseWriter, request *http.Request) {
	webApiContext.beforeHandle(&respWriter)
	var result interface{}
	status := http.StatusOK
	items, err := managers.GetLifoList(webApiContext.DbContext, webApiContext.Logger)
	if err != nil {
		webApiContext.Logger.Error(stringFormatter.Format("An error occurred during getting Lifo data: {0}", err.Error()))
		status = http.StatusInternalServerError
	} else {
		dtoItems := make([]dto.LifoInfo, len(items))
		for index, item := range items {
			dtoItems[index] = factory.CreateLifoInfo(&item)
		}
		result = dtoItems
	}
	webApiContext.afterHandle(&respWriter, status, result)
}

func (webApiContext *WebApiContext) GetLifoById(respWriter http.ResponseWriter, request *http.Request) {
	webApiContext.beforeHandle(&respWriter)
	var result interface{}
	status := http.StatusOK
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		status = http.StatusBadRequest
	} else {
		lifo, err := managers.GetLifoById(uint(id), webApiContext.DbContext, webApiContext.Logger)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				webApiContext.Logger.Warn(stringFormatter.Format("A lifo with id: \"{0}\" was not found ", id))
				status = http.StatusNotFound
			} else {
				webApiContext.Logger.Error(stringFormatter.Format("An unexpected error occurred during getting lifo by id: \"{0}\", error: {1}", id, err.Error()))
				status = http.StatusInternalServerError
			}
		}
		lifoDto := factory.CreateLifoInfo(&lifo)
		result = lifoDto
	}
	webApiContext.afterHandle(&respWriter, status, result)
}

func (webApiContext *WebApiContext) GetLifoByIdWithItems(respWriter http.ResponseWriter, request *http.Request) {
	webApiContext.beforeHandle(&respWriter)
	var result interface{}
	status := http.StatusOK
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		status = http.StatusBadRequest
	} else {
		lifo, err := managers.GetLifoByIdWithItems(uint(id), webApiContext.DbContext, webApiContext.Logger)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				webApiContext.Logger.Warn(stringFormatter.Format("A lifo with id: \"{0}\" was not found ", id))
				status = http.StatusNotFound
			} else {
				webApiContext.Logger.Error(stringFormatter.Format("An unexpected error occurred during getting lifo by id: \"{0}\", error: {1}", id, err.Error()))
				status = http.StatusInternalServerError
			}
		}
		lifoDto := factory.CreateLifoWithItems(&lifo)
		result = lifoDto
	}
	webApiContext.afterHandle(&respWriter, status, result)
}

func (webApiContext *WebApiContext) CreateLifo(respWriter http.ResponseWriter, request *http.Request) {
	webApiContext.beforeHandle(&respWriter)
	var result interface{}
	var lifoInfo dto.LifoInfo
	status := http.StatusCreated
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&lifoInfo)
	if err != nil {
		status = http.StatusBadRequest
	} else {
		lifo, err := managers.CreateLifo(&lifoInfo, webApiContext.DbContext, webApiContext.Logger)
		if err != nil {
			if errors.Is(err, managers.BadLifoName) {
				status = http.StatusBadRequest
			} else {
				status = http.StatusInternalServerError
			}
		} else {
			result = factory.CreateLifoInfo(&lifo)
		}
	}
	webApiContext.afterHandle(&respWriter, status, result)
}

func (webApiContext *WebApiContext) UpdateLifo(respWriter http.ResponseWriter, request *http.Request) {
	webApiContext.beforeHandle(&respWriter)
	var result interface{}
	var lifoInfo dto.LifoInfo
	status := http.StatusOK
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	decoder := json.NewDecoder(request.Body)
	decodeErr := decoder.Decode(&lifoInfo)
	if err != nil || decodeErr != nil {
		status = http.StatusBadRequest
	} else {
		lifo, err := managers.UpdateLifo(&lifoInfo, uint(id), webApiContext.DbContext, webApiContext.Logger)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				status = http.StatusNotFound
			} else if errors.Is(err, managers.BadSizeError) {
				status = http.StatusBadRequest
			} else {
				webApiContext.Logger.Error(stringFormatter.Format("An unexpected error occurred during lifo with id \"{0}\" update, error: {1}", id, err.Error()))
				status = http.StatusInternalServerError
			}
		} else {
			result = factory.CreateLifoInfo(&lifo)
		}
	}
	webApiContext.afterHandle(&respWriter, status, result)
}

func (webApiContext *WebApiContext) DeleteLifo(respWriter http.ResponseWriter, request *http.Request) {
	status := http.StatusOK
	webApiContext.beforeHandle(&respWriter)
	vars := mux.Vars(request)
	var result interface{}
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		status = http.StatusBadRequest
	} else {
		success, err := managers.DeleteLifo(uint(id), webApiContext.DbContext, webApiContext.Logger)
		if success {
			status = http.StatusNoContent
		} else {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				status = http.StatusNotFound
			} else {
				status = http.StatusInternalServerError
			}
			webApiContext.Logger.Error(err.Error())
		}
	}
	webApiContext.afterHandle(&respWriter, status, result)
}

func (webApiContext *WebApiContext) PushLifo(respWriter http.ResponseWriter, request *http.Request) {
	webApiContext.beforeHandle(&respWriter)
	var result interface{}
	var lifoInfo dto.LifoItem
	status := http.StatusOK
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	decoder := json.NewDecoder(request.Body)
	decodeErr := decoder.Decode(&lifoInfo)
	if err != nil || decodeErr != nil {
		status = http.StatusBadRequest
	} else {
		msg := ""
		newTopItem, err := managers.PushToLifo(uint(id), &lifoInfo, webApiContext.DbContext, webApiContext.Logger)
		if err != nil {
			status = http.StatusInternalServerError
			msg = err.Error()
		}
		result = factory.CreateLifoOperation(dto.Push, err == nil, msg, &newTopItem)
	}
	webApiContext.afterHandle(&respWriter, status, result)
}

func (webApiContext *WebApiContext) PopLifo(respWriter http.ResponseWriter, request *http.Request) {
	webApiContext.beforeHandle(&respWriter)
	var result interface{}
	status := http.StatusOK
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		status = http.StatusBadRequest
	} else {
		poppedItem, err := managers.PopFromLifo(uint(id), webApiContext.DbContext, webApiContext.Logger)
		msg := ""
		if err != nil {
			if err != managers.LifoIsEmpty {
				status = http.StatusInternalServerError
			}
			msg = err.Error()
		}
		result = factory.CreateLifoOperation(dto.Pop, err == nil, msg, &poppedItem)
	}
	webApiContext.afterHandle(&respWriter, status, result)
}

func (webApiContext *WebApiContext) FlushLifo(respWriter http.ResponseWriter, request *http.Request) {

}

// BeforeHandle
/* This function prepare response headers prior to response handle. It sets content-type and CORS headers.
 * Parameters:
 *     - respWriter - gorilla/mux response writer
 * Returns nothing
 */
func (webApiContext *WebApiContext) beforeHandle(respWriter *http.ResponseWriter) {
	(*respWriter).Header().Set("Content-Type", "application/json")
	(*respWriter).Header().Set("Accept", "application/json")
}

// AfterHandle
/* This function finalize response handle: serialize (json) and write object and set status code. If error occur during object serialization status code sets to 500
 * Parameters:
 *     - respWriter - gorilla/mux response writer
 *     - statusCode - http response status
 *     - data - object (json) could be empty
 * Returns nothing
 */
func (webApiContext *WebApiContext) afterHandle(respWriter *http.ResponseWriter, statusCode int, data interface{}) {
	(*respWriter).WriteHeader(statusCode)
	if data != nil {
		err := json.NewEncoder(*respWriter).Encode(data)
		if err != nil {
			(*respWriter).WriteHeader(500)
		}
	}
}
