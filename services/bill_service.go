package services

import (
	"errors"
	"time"

	"github.com/JZ23-2/splitbill-backend/database"
	"github.com/JZ23-2/splitbill-backend/dtos"
	"github.com/JZ23-2/splitbill-backend/models"
	"github.com/google/uuid"
)

func AssignParticipantToBill(req dtos.AssignParticipantsRequest) (*dtos.AssignParticipantsResponse, error) {
	itemToParticipants := make(map[string][]string)
	for _, p := range req.Participants {
		for _, item := range p.Items {
			itemToParticipants[item.ItemID] = append(itemToParticipants[item.ItemID], p.ParticipantID)
		}
	}

	itemIDs := keys(itemToParticipants)
	var items []models.Item
	if err := database.DB.Where("item_id IN ?", itemIDs).Find(&items).Error; err != nil {
		return nil, errors.New("failed to fetch items: " + err.Error())
	}

	itemPriceMap := make(map[string]int)
	for _, item := range items {
		itemPriceMap[item.ItemID] = item.Price
	}

	var responseParticipants []dtos.AssignParticipantDetailResponse

	for _, p := range req.Participants {
		participant := models.Participant{
			ParticipantID: p.ParticipantID,
			BillID:        req.BillID,
			IsPaid:        p.IsPaid,
		}

		if err := database.DB.Create(&participant).Error; err != nil {
			return nil, errors.New("failed to create participant: " + err.Error())
		}

		totalOwed := 0
		var responseItems []dtos.AssignParticipantItemWithAmount

		for _, item := range p.Items {
			price := itemPriceMap[item.ItemID]
			numParticipants := len(itemToParticipants[item.ItemID])
			if numParticipants == 0 {
				return nil, errors.New("no participants for item: " + item.ItemID)
			}

			share := price / numParticipants
			totalOwed += share

			participantItem := models.ParticipantItem{
				ParticipantID: p.ParticipantID,
				ItemID:        item.ItemID,
				Amount:        share,
			}

			if err := database.DB.Create(&participantItem).Error; err != nil {
				return nil, errors.New("failed to create participant item: " + err.Error())
			}

			responseItems = append(responseItems, dtos.AssignParticipantItemWithAmount{
				ItemID: item.ItemID,
				Amount: share,
			})
		}

		if err := database.DB.Model(&models.Participant{}).
			Where("participant_id = ? AND bill_id = ?", p.ParticipantID, req.BillID).
			Update("amount_owed", totalOwed).Error; err != nil {
			return nil, errors.New("failed to update amount owed: " + err.Error())
		}

		responseParticipants = append(responseParticipants, dtos.AssignParticipantDetailResponse{
			ParticipantID: p.ParticipantID,
			IsPaid:        p.IsPaid,
			AmountOwed:    totalOwed,
			Items:         responseItems,
		})
	}

	resp := &dtos.AssignParticipantsResponse{
		BillID:       req.BillID,
		Participants: responseParticipants,
	}

	return resp, nil
}

func keys(m map[string][]string) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	return out
}

func CreateBillWithoutParticipant(req dtos.CreateBillWithoutParticipantRequest) (*dtos.CreateBillWithoutParticipantResponse, error) {
	bill := models.Bill{
		BillTitle:   req.StoreName,
		TotalAmount: req.TotalAmount,
		Tax:         req.Tax,
		Service:     req.Service,
		CreatorID:   req.CreatorID,
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&bill).Error; err != nil {
		return nil, errors.New("failed to create bill: " + err.Error())
	}

	var itemResponses []dtos.CreateBillWithoutParticipantItemResponse

	for _, item := range req.Items {
		itemID := uuid.NewString()
		newItem := models.Item{
			ItemID:   itemID,
			BillID:   bill.BillID,
			Name:     item.Name,
			Price:    item.UnitPrice,
			Quantity: item.Quantity,
		}

		if err := database.DB.Create(&newItem).Error; err != nil {
			return nil, errors.New("failed to create item: " + err.Error())
		}

		itemResponses = append(itemResponses, dtos.CreateBillWithoutParticipantItemResponse{
			ItemID:        itemID,
			Name:          item.Name,
			Quantity:      item.Quantity,
			UnitPrice:     item.UnitPrice,
			PriceAfterTax: item.PriceAfterTax,
			PriceInHBAR:   item.PriceInHBAR,
		})
	}

	resp := &dtos.CreateBillWithoutParticipantResponse{
		BillID:      bill.BillID,
		StoreName:   req.StoreName,
		Date:        req.Date,
		Tax:         req.Tax,
		Service:     req.Service,
		TotalAmount: req.TotalAmount,
		CreatedAt:   bill.CreatedAt.Format(time.RFC3339),
		CreatorID:   req.CreatorID,
		Items:       itemResponses,
	}

	return resp, nil
}

func GetBillsByCreator(creatorID string, billID string) ([]dtos.GetBillByCreatorResponse, error) {
	if creatorID == "" {
		return nil, errors.New("creatorId is required")
	}

	var bills []models.Bill
	query := database.DB.Preload("Items").Where("creator_id = ?", creatorID)

	if billID != "" {
		query = query.Where("bill_id = ?", billID)
	}

	if err := query.Find(&bills).Error; err != nil {
		return nil, errors.New("failed to fetch bills: " + err.Error())
	}

	var response []dtos.GetBillByCreatorResponse
	for _, bill := range bills {
		var itemResponses []dtos.GetBillByCreatorItemResponse
		for _, item := range bill.Items {
			itemResponses = append(itemResponses, dtos.GetBillByCreatorItemResponse{
				ItemID:   item.ItemID,
				Name:     item.Name,
				Price:    item.Price,
				Quantity: item.Quantity,
			})
		}

		response = append(response, dtos.GetBillByCreatorResponse{
			BillID:      bill.BillID,
			BillTitle:   bill.BillTitle,
			TotalAmount: bill.TotalAmount,
			Tax:         bill.Tax,
			Service:     bill.Service,
			CreatedAt:   bill.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Items:       itemResponses,
		})
	}

	return response, nil
}
