package services

import (
	"errors"
	"math"
	"time"

	"github.com/JZ23-2/splitbill-backend/database"
	"github.com/JZ23-2/splitbill-backend/dtos"
	"github.com/JZ23-2/splitbill-backend/models"
	"github.com/google/uuid"
)

func CreateBillWithoutParticipant(req dtos.CreateBillWithoutParticipantRequest) (*dtos.CreateBillWithoutParticipantResponse, error) {
	bill := models.Bill{
		BillTitle:   req.BillTitle,
		StoreName:   req.StoreName,
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
	totalAmountAfterTax := 0

	for _, item := range req.Items {
		itemID := uuid.NewString()
		itemSubtotal := float64(item.Price) * float64(item.Quantity)

		taxPortion := itemSubtotal * float64(req.Tax) / 100
		servicePortion := itemSubtotal * float64(req.Service) / 100

		priceAfterTax := int(math.Round(itemSubtotal + taxPortion + servicePortion))
		totalAmountAfterTax += priceAfterTax

		newItem := models.Item{
			ItemID:        itemID,
			BillID:        bill.BillID,
			Name:          item.Name,
			Price:         item.Price,
			Quantity:      item.Quantity,
			PriceAfterTax: priceAfterTax,
		}

		if err := database.DB.Create(&newItem).Error; err != nil {
			return nil, errors.New("failed to create item: " + err.Error())
		}

		itemResponses = append(itemResponses, dtos.CreateBillWithoutParticipantItemResponse{
			ItemID:        itemID,
			Name:          item.Name,
			Quantity:      item.Quantity,
			UnitPrice:     item.Price,
			PriceAfterTax: priceAfterTax,
		})
	}

	bill.TotalAmountAfterTax = totalAmountAfterTax
	if err := database.DB.Save(&bill).Error; err != nil {
		return nil, errors.New("failed to update bill with totalAmountAfterTax: " + err.Error())
	}

	resp := &dtos.CreateBillWithoutParticipantResponse{
		BillID:              bill.BillID,
		BillTitle:           bill.BillTitle,
		StoreName:           bill.StoreName,
		Date:                req.Date,
		Tax:                 req.Tax,
		Service:             req.Service,
		TotalAmount:         req.TotalAmount,
		TotalAmountAfterTax: totalAmountAfterTax,
		CreatedAt:           bill.CreatedAt.Format(time.RFC3339),
		CreatorID:           bill.CreatorID,
		Items:               itemResponses,
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

func AssignParticipantsToItem(req dtos.AssignParticipantsRequest) (*dtos.AssignedParticipantResponse, error) {
	var item models.Item

	if err := database.DB.First(&item, "item_id = ?", req.ItemID).Error; err != nil {
		return nil, errors.New("item not found")
	}

	if len(req.ParticipantID) == 0 {
		return nil, errors.New("no participant IDs provided")
	}

	share := item.PriceAfterTax / len(req.ParticipantID)
	var assigned []dtos.AssignedParticipant

	for _, pid := range req.ParticipantID {
		participant := models.Participant{
			ParticipantID: pid,
			ItemID:        item.ItemID,
			AmountOwed:    share,
			IsPaid:        false,
		}

		if err := database.DB.Create(&participant).Error; err != nil {
			return nil, errors.New("failed to assign participant")
		}

		assigned = append(assigned, dtos.AssignedParticipant{
			ParticipantID: pid,
			ItemID:        item.ItemID,
			AmountOwed:    share,
			IsPaid:        false,
		})

	}

	resp := &dtos.AssignedParticipantResponse{
		ItemID:       item.ItemID,
		ItemName:     item.Name,
		Participants: assigned,
	}

	return resp, nil
}

func GetBillsByParticipantID(participantID string) ([]dtos.ParticipantBillResponse, error) {
	if participantID == "" {
		return nil, errors.New("participantId is required")
	}

	var participantRecords []models.Participant
	if err := database.DB.Where("participant_id = ?", participantID).Find(&participantRecords).Error; err != nil {
		return nil, errors.New("failed to fetch participant data: " + err.Error())
	}

	billMap := make(map[string]bool)
	itemToBill := make(map[string]string)

	for _, p := range participantRecords {
		var item models.Item
		if err := database.DB.Select("item_id", "bill_id").First(&item, "item_id = ?", p.ItemID).Error; err != nil {
			continue
		}
		billMap[item.BillID] = true
		itemToBill[item.ItemID] = item.BillID
	}

	var responses []dtos.ParticipantBillResponse

	for billID := range billMap {
		var bill models.Bill
		if err := database.DB.First(&bill, "bill_id = ?", billID).Error; err != nil {
			continue
		}

		var items []models.Item
		if err := database.DB.Where("bill_id = ?", bill.BillID).Find(&items).Error; err != nil {
			continue
		}

		var itemDTOs []dtos.ParticipantItemResponse
		var allParticipants []dtos.ParticipantListResponse
		participantSet := make(map[string]bool)

		for _, item := range items {
			itemDTOs = append(itemDTOs, dtos.ParticipantItemResponse{
				ItemID:        item.ItemID,
				Name:          item.Name,
				Quantity:      item.Quantity,
				Price:         item.Price,
				PriceAfterTax: item.PriceAfterTax,
			})

			var participants []models.Participant
			if err := database.DB.Where("item_id = ?", item.ItemID).Find(&participants).Error; err != nil {
				continue
			}

			for _, p := range participants {
				if !participantSet[p.ParticipantID] {
					participantSet[p.ParticipantID] = true
					allParticipants = append(allParticipants, dtos.ParticipantListResponse{
						ParticipantID: p.ParticipantID,
						AmountOwed:    p.AmountOwed,
						IsPaid:        p.IsPaid,
					})
				}
			}
		}

		response := dtos.ParticipantBillResponse{
			BillID:      bill.BillID,
			BillTitle:   bill.BillTitle,
			TotalAmount: bill.TotalAmount,
			CreatorID:   bill.CreatorID,
			Tax:         bill.Tax,
			Service:     bill.Service,
			Item:        itemDTOs,
			Participant: allParticipants,
		}

		responses = append(responses, response)
	}

	return responses, nil
}
