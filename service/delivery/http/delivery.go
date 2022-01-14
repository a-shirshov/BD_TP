
package delivery

import (
	"net/http"
	"bd_tp/response"
	serviceUsecase "bd_tp/service/usecase"
)

type ServiceDelivery struct {
	useCase *serviceUsecase.Usecase
}

func NewServiceDelivery(useCase *serviceUsecase.Usecase) *ServiceDelivery {
	return &ServiceDelivery{
		useCase: useCase,
	}
}

func (d *ServiceDelivery) Clear(w http.ResponseWriter, r *http.Request) {
	err := d.useCase.Clear()
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError,err)
		return
	}
	response.SendResponse(w, http.StatusOK, nil)
}

func (d *ServiceDelivery) GetStatus(w http.ResponseWriter, r *http.Request) {
	status, err := d.useCase.GetStatus()
	if err != nil {
		response.SendResponse(w, http.StatusInternalServerError,err)
		return
	}
	response.SendResponse(w, http.StatusOK, status)
}