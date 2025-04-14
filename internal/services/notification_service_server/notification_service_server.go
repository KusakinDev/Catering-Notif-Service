package notificationserviceserver

import (
	"context"
	"encoding/json"
	"errors"

	emailmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/email_model"
	menunotificationmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/menu_notif_model"
	resetnotifmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/reset_notif_model"
	pb "github.com/KusakinDev/Catering-Notif-Service/internal/services/notification_service_gen"
	rabbitmq "github.com/KusakinDev/Catering-Notif-Service/internal/utils/RabbitMQ"
	"github.com/sirupsen/logrus"
)

type Server struct {
	pb.UnimplementedNotificationServiceServer
	Rmq *rabbitmq.RabbitMQ
}

func (s *Server) ResetNotif(ctx context.Context, req *pb.ResetRequest) (*pb.Response, error) {
	var resetNotif resetnotifmodel.ResetNotification
	var errCode int
	resetNotif.Code = int(req.ResetCode)
	resetNotif.Email = req.Email

	errCode = resetNotif.GetTemplate()
	if errCode != 200 {
		return &pb.Response{
			Code:    int32(errCode),
			Message: "Error get reset template",
		}, errors.New("error get reset template")
	}

	body, err := json.Marshal(resetNotif)
	if err != nil {
		logrus.Error("Error marshalling notification: ", err)
		return &pb.Response{
			Code:    int32(errCode),
			Message: "Error marshalling reset notification",
		}, errors.New("error marshalling reset notification")
	}

	err = s.Rmq.Publish(body, "reset")
	if err != nil {
		logrus.Error("Failed to publish message to RabbitMQ: ", err)
		return &pb.Response{
			Code:    int32(errCode),
			Message: "Error send reset email",
		}, errors.New("error send reset email")
	}

	return &pb.Response{
		Code:    int32(200),
		Message: "Success send reset email",
	}, nil
}

func (s *Server) NotifNewMenu(ctx context.Context, req *pb.MenuRequest) (*pb.Response, error) {
	message := req.Message

	var email emailmodel.Email
	emails, code := email.GetAllFromTable()
	if code != 200 {
		return &pb.Response{
			Code:    int32(code),
			Message: "Error get all emails from table",
		}, errors.New("error get all emails from table")
	}

	var notif menunotificationmodel.MenuNotification
	notif.GetTemplateByTag("new_menu")
	notif.Message = message
	for _, email := range emails {
		notif.Email = email

		body, err := json.Marshal(notif)
		if err != nil {
			logrus.Error("Error marshalling notification: ", err)
			return &pb.Response{
				Code:    int32(code),
				Message: "Error marshalling notification",
			}, errors.New("error marshalling notification")
		}

		err = s.Rmq.Publish(body, "menu")
		if err != nil {
			logrus.Error("Failed to publish message to RabbitMQ: ", err)
			return &pb.Response{
				Code:    int32(code),
				Message: "Failed to publish message to RabbitMQ",
			}, errors.New("failed to publish message to rabbit_mq")
		}
	}

	return &pb.Response{
		Code:    int32(200),
		Message: "Success broadcast",
	}, nil
}
