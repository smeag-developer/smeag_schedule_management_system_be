package handlers

import (
	"fmt"
	"log/slog"
	"nxt_match_event_manager_api/internal/interfaces"
	config "nxt_match_event_manager_api/internal/models/config"
	"nxt_match_event_manager_api/internal/utils"
	"nxt_match_event_manager_api/internal/utils/grpc"
	"nxt_match_event_manager_api/pb/buff"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
)

type GrpcNotification struct {
	mu           sync.Mutex
	repo         interfaces.NotificationRespositoryInterface
	serviceNotif buff.GroupNotificationServiceClient
	eventNotif   buff.EventNotificationServiceClient
}

func NewGrpcNotificationHandler(
	repo interfaces.NotificationRespositoryInterface,
	groupJoinBuff buff.GroupNotificationServiceClient,
	eventNotif buff.EventNotificationServiceClient,
	host *config.HostConfig) *GrpcNotification {
	return &GrpcNotification{
		repo:         repo,
		serviceNotif: groupJoinBuff,
		eventNotif:   eventNotif,
	}
}

func (s *GrpcNotification) RequestGroupJoinService(ctx *gin.Context, req *buff.CreateGroupJoinRequest) error {

	res, err := s.serviceNotif.RequestJoinGroupService(ctx, req)

	if err != nil {
		grpc.ValidateStatusMessage(err)
		return fmt.Errorf("[grpc_client]: %v ", err)
	}

	slog.Info("is group_created associated :%v", res.IsGroupCreated)

	if !res.IsGroupCreated {
		return fmt.Errorf("[grpc_client]: %v ", err)
	}

	return nil
}

func (s *GrpcNotification) RequestJoinEventService(ctx *gin.Context, req *buff.CreateEventJoinRequest) error {

	res, err := s.eventNotif.RequestUserJoinEvent(ctx, req)

	if err != nil {
		// validate grpc response
		grpc.ValidateStatusMessage(err)
		return fmt.Errorf("[grpc_client]: %v ", err)
	}

	if !bool(res.IsEventCreated) {
		return fmt.Errorf("[grpc_client]: event not created")
	}

	return nil

}

func (s *GrpcNotification) UpdateGroupNotifStatus(ctx *gin.Context, req *buff.GroupStatusRequest) error {

	res, err := s.serviceNotif.RequestUpdateGroupJoinStatus(ctx, req)

	if err != nil {
		// validate grpc response
		grpc.ValidateStatusMessage(err)
	}

	g := req.GroupNotifStatusRequest
	// cast to hex
	notifHexId, err := utils.MustHex(g.NotificationId)

	if err != nil {
		return fmt.Errorf("invalid hex id")
	}

	if res.GroupNotifStatusResponse.StatusCode == codes.OK.String() {
		return s.repo.UpsertNotificationStatus(ctx, g.GroupJoiners[0].Status, notifHexId)
	}

	return nil
}

func (s *GrpcNotification) UpdateEventNotifStatus(ctx *gin.Context, req *buff.EventStatusRequest) error {

	res, err := s.eventNotif.RequestUpdateEventStatus(ctx, req)

	if err != nil {
		// validate grpc response
		grpc.ValidateStatusMessage(err)
	}

	r := req.EventNotifStatusRequest

	// cast to hex
	notifHexId, err := utils.MustHex(r.NotificationId)

	if err != nil {
		return fmt.Errorf("invalid hex id")
	}

	if res.EventNotifStatusResponse.StatusCode == codes.OK.String() {
		return s.repo.UpsertNotificationStatus(ctx, r.EventJoiner[0].Status, notifHexId)
	}
	// return s.repo.UpsertNotificationStatus(ctx, status, notifId)

	return nil
}

func (s *GrpcNotification) BroadCastToUserService(ctx *gin.Context) error {
	return s.repo.BroadCastToUser(ctx)
}
